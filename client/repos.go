package client

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/gosuri/uiprogress"
)

func (c *Client) GetOrgRepos(ctx context.Context, orgName string) ([]*github.Repository, error) {
	org, _, err := c.ghClient.Organizations.Get(ctx, orgName)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("github: hit rate limit")
		}

		return nil, fmt.Errorf("get org: %v", err.Error())
	}

	count := org.GetPublicRepos() + org.GetTotalPrivateRepos()

	repoBar := uiprogress.AddBar(count).
		AppendCompleted().
		PrependElapsed().
		PrependFunc(func(b *uiprogress.Bar) string {
			return fmt.Sprintf("Fetching (%d/%d)", b.Current(), count)
		})

	opts := &github.RepositoryListByOrgOptions{
		Type: "all",
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: 100,
		},
	}

	var repos []*github.Repository
	for {
		rs, resp, err := c.ghClient.Repositories.ListByOrg(ctx, orgName, opts)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return nil, fmt.Errorf("github: hit rate limit")
			}

			return nil, fmt.Errorf("list repos: %v", err.Error())
		}

		for range rs {
			repoBar.Incr()
		}

		repos = append(repos, rs...)

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return repos, nil
}

func (c *Client) GetUserRepos(ctx context.Context, username string) ([]*github.Repository, error) {
	u, _, err := c.ghClient.Users.Get(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get user: %v", err.Error())
	}

	count := u.GetPublicRepos() + u.GetTotalPrivateRepos()

	bar := uiprogress.AddBar(count).
		AppendCompleted().
		PrependElapsed().
		PrependFunc(func(b *uiprogress.Bar) string {
			return fmt.Sprintf("Fetching (%d/%d)", b.Current(), count)
		})

	var repos []*github.Repository

	opts := &github.RepositoryListByOrgOptions{
		Type: "all",
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: 100,
		},
	}

	for {
		rs, resp, err := c.ghClient.Repositories.List(ctx, username, nil)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return nil, fmt.Errorf("github: hit rate limit")
			}

			return nil, fmt.Errorf("list repos: %v", err.Error())
		}

		for range rs {
			bar.Incr()
		}

		repos = append(repos, rs...)

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return repos, nil
}

func (c *Client) ProcessRepos(ctx context.Context, repos []*github.Repository, dryRun bool) ([]string, error) {
	count := len(repos)
	name := repos[0].GetName()
	owner := repos[0].GetOwner().GetLogin()
	appendStr := fmt.Sprintf("\nCurrent Repo: %v/%v", owner, name)

	bar := uiprogress.AddBar(count).
		AppendCompleted().
		PrependElapsed().
		PrependFunc(func(b *uiprogress.Bar) string {
			return fmt.Sprintf("Processing (%d/%d)", b.Current(), count)
		}).
		AppendFunc(func(b *uiprogress.Bar) string {
			return appendStr
		})

	urls := []string{}
	for _, repo := range repos {
		name = repo.GetName()
		owner = repo.GetOwner().GetLogin()
		appendStr = fmt.Sprintf("\nCurrent Repo: %v/%v", owner, name)

		url, err := c.processRepo(ctx, repo, dryRun)
		if err != nil {
			if strings.HasPrefix(err.Error(), "get branch: ") || strings.HasPrefix(err.Error(), "no commits") {
				bar.Incr()
				continue
			}

			fmt.Printf("process repo: %v\n", err.Error())
			os.Exit(1)
		}

		urls = append(urls, url)
		bar.Incr()
	}

	appendStr = ""

	return urls, nil
}

func (c *Client) processRepo(ctx context.Context, repo *github.Repository, dryRun bool) (string, error) {
	name := repo.GetName()
	owner := repo.GetOwner().GetLogin()
	head := repo.GetDefaultBranch()

	_, _, err := c.ghClient.Repositories.GetBranch(ctx, owner, name, c.base)
	if err != nil {
		return "", fmt.Errorf("get branch: %v", err.Error())
	}

	opts := &github.PullRequestListOptions{
		Head: head,
		Base: c.base,
	}

	prs, _, err := c.ghClient.PullRequests.List(ctx, owner, name, opts)
	if err != nil {
		return "", fmt.Errorf("list prs: %v", err.Error())
	}

	if len(prs) > 0 {
		pr := prs[0]

		pr.Title = github.String("Release")

		changes, _ := c.createChangeLog(ctx, owner, name, c.base, head)
		body := prBody(prBodyTemplate, changes)

		pr.Body = github.String(body)

		if !dryRun {
			pr, _, _ = c.ghClient.PullRequests.Edit(ctx, owner, name, pr.GetNumber(), pr)
		}

		return pr.GetHTMLURL(), nil
	}

	changes, err := c.createChangeLog(ctx, owner, name, c.base, head)
	if err != nil {
		return "", err
	}

	body := prBody(prBodyTemplate, changes)

	newPR := &github.NewPullRequest{
		Title:               github.String("Release"),
		Head:                &head,
		Base:                &c.base,
		Body:                github.String(body),
		MaintainerCanModify: github.Bool(true),
	}

	if !dryRun {
		pr, _, err := c.ghClient.PullRequests.Create(ctx, owner, name, newPR)
		if err != nil {
			return "", fmt.Errorf("create pr: %v", err.Error())
		}

		return pr.GetHTMLURL(), nil
	}

	return fmt.Sprintf("https://github.com/%v/%v/compare/%v...%v", owner, name, c.base, head), nil
}

func (c *Client) ReleaseRepos(ctx context.Context, repos []*github.Repository, dryRun bool) ([]string, error) {
	releases, err := c.getReleases(ctx, repos)
	if err != nil {
		return nil, fmt.Errorf("releases: %v\n", err.Error())
	}

	if len(releases) < 1 {
		return nil, nil
	}

	count := len(releases)
	repo := releases[0].GetBase().GetRepo()
	name := repo.GetName()
	owner := repo.GetOwner().GetLogin()
	appendStr := fmt.Sprintf("\nCurrent Repo: %v/%v", owner, name)

	bar := uiprogress.AddBar(count).
		AppendCompleted().
		PrependElapsed().
		PrependFunc(func(b *uiprogress.Bar) string {
			return fmt.Sprintf("Processing Releases (%d/%d)", b.Current(), count)
		}).
		AppendFunc(func(b *uiprogress.Bar) string {
			return appendStr
		})

	var released []string
	for _, release := range releases {
		repo := release.GetBase().GetRepo()
		name = repo.GetName()
		owner = repo.GetOwner().GetLogin()
		appendStr = fmt.Sprintf("\nCurrent Repo: %v/%v", owner, name)

		var err error
		release, _, err = c.ghClient.PullRequests.Get(ctx, owner, name, release.GetNumber())
		if err != nil {
			return nil, fmt.Errorf("check mergeable: %v", err.Error())
		}

		if strings.ToLower(release.GetMergeableState()) != "clean" {
			bar.Incr()
			continue
		}

		if !dryRun {
			res, _, err := c.ghClient.PullRequests.Merge(ctx, owner, name, release.GetNumber(), "release automerged by train", nil)
			if err != nil {
				return nil, fmt.Errorf("merge: %v", err.Error())
			}

			if res.GetMerged() {
				released = append(released, release.GetHTMLURL())
			}
		} else {
			released = append(released, release.GetHTMLURL())
		}

		bar.Incr()
	}

	appendStr = ""

	return released, nil
}

func (c *Client) getReleases(ctx context.Context, repos []*github.Repository) ([]*github.PullRequest, error) {
	var releases []*github.PullRequest

	count := len(repos)
	name := repos[0].GetName()
	owner := repos[0].GetOwner().GetLogin()
	appendStr := fmt.Sprintf("\nCurrent Repo: %v/%v", owner, name)

	bar := uiprogress.AddBar(count).
		AppendCompleted().
		PrependElapsed().
		PrependFunc(func(b *uiprogress.Bar) string {
			return fmt.Sprintf("Collecting Releases (%d/%d)", b.Current(), count)
		}).
		AppendFunc(func(b *uiprogress.Bar) string {
			return appendStr
		})

	for _, repo := range repos {
		owner = repo.GetOwner().GetLogin()
		name = repo.GetName()
		appendStr = fmt.Sprintf("\nCurrent Repo: %v/%v", owner, name)
		head := repo.GetDefaultBranch()

		opts := &github.PullRequestListOptions{
			Head: head,
			Base: c.base,
		}

		rs, _, err := c.ghClient.PullRequests.List(ctx, owner, name, opts)
		if err != nil {
			return nil, fmt.Errorf("pull requests: %v", err.Error())
		}

		releases = append(releases, rs...)
		bar.Incr()
	}

	appendStr = ""

	return releases, nil
}
