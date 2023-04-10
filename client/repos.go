package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gomicro/crawl"
	"github.com/gomicro/crawl/bar"
	"github.com/google/go-github/github"
)

var ErrGetBranch = errors.New("get branch")

func (c *Client) GetRepos(ctx context.Context, progress *crawl.Progress, name string) ([]*github.Repository, error) {
	count := 0
	orgFound := true

	c.rate.Wait(ctx) //nolint: errcheck
	org, resp, err := c.ghClient.Organizations.Get(ctx, name)
	if resp == nil && err != nil {

		if _, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("github: hit rate limit")
		}

		return nil, fmt.Errorf("get org: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		orgFound = false

		c.rate.Wait(ctx) //nolint: errcheck
		user, _, err := c.ghClient.Users.Get(ctx, name)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return nil, fmt.Errorf("github: hit rate limit")
			}

			return nil, fmt.Errorf("get user: %v", err.Error())
		}

		count = user.GetPublicRepos() + user.GetTotalPrivateRepos()
	} else {
		count = org.GetPublicRepos() + org.GetTotalPrivateRepos()
	}

	if count < 1 {
		return nil, fmt.Errorf("no repos found")
	}

	theme := bar.NewThemeFromTheme(bar.DefaultTheme)
	theme.Append(func(b *bar.Bar) string {
		return fmt.Sprintf(" %0.2f", b.CompletedPercent())
	})
	theme.Prepend(func(b *bar.Bar) string {
		return fmt.Sprintf("Fetching (%d/%d) %s", b.Current(), b.Total(), b.Elapsed())
	})

	repoBar := bar.New(theme, count)
	progress.AddBar(repoBar)

	orgOpts := &github.RepositoryListByOrgOptions{
		Type: "all",
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: 100,
		},
	}

	userOpts := &github.RepositoryListOptions{
		Type: "all",
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: 100,
		},
	}

	var repos []*github.Repository
	for {
		var rs []*github.Repository
		c.rate.Wait(ctx) //nolint: errcheck
		if orgFound {
			rs, resp, err = c.ghClient.Repositories.ListByOrg(ctx, name, orgOpts)
		} else {
			rs, resp, err = c.ghClient.Repositories.List(ctx, name, userOpts)
		}

		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return nil, fmt.Errorf("github: hit rate limit")
			}

			return nil, fmt.Errorf("list repos: %v", err.Error())
		}

		for i := range rs {
			repoBar.Incr()

			if rs[i].GetArchived() {
				continue
			}

			name := strings.ToLower(rs[i].GetName())
			_, looseMatch := c.ignoreRepoMap[name]

			fullName := fmt.Sprintf("%v/%v", strings.ToLower(rs[i].GetOwner().GetLogin()), name)
			_, exactMatch := c.ignoreRepoMap[fullName]

			if looseMatch || exactMatch {
				continue
			}

			topics := rs[i].Topics

			topicSkip := false
			for _, t := range topics {
				_, topicMatch := c.ignoreTopicMap[strings.ToLower(t)]
				if topicMatch {
					topicSkip = true
					break
				}
			}

			if topicSkip {
				continue
			}

			repos = append(repos, rs[i])
		}

		if resp.NextPage == 0 {
			break
		}

		if orgFound {
			orgOpts.Page = resp.NextPage
		} else {
			userOpts.Page = resp.NextPage
		}
	}

	return repos, nil
}

func (c *Client) ProcessRepos(ctx context.Context, progress *crawl.Progress, repos []*github.Repository, dryRun bool) ([]string, error) {
	count := len(repos)
	name := repos[0].GetName()
	owner := repos[0].GetOwner().GetLogin()
	appendStr := fmt.Sprintf("\nCurrent Repo: %v/%v", owner, name)

	theme := bar.NewThemeFromTheme(bar.DefaultTheme)
	theme.Append(func(b *bar.Bar) string {
		return fmt.Sprintf(" %0.2f", b.CompletedPercent())
	})
	theme.Append(func(b *bar.Bar) string {
		return appendStr
	})
	theme.Prepend(func(b *bar.Bar) string {
		return fmt.Sprintf("Processing (%d/%d) %s", b.Current(), b.Total(), b.Elapsed())
	})

	repoBar := bar.New(theme, count)
	progress.AddBar(repoBar)

	urls := []string{}
	for _, repo := range repos {
		name = repo.GetName()
		owner = repo.GetOwner().GetLogin()
		appendStr = fmt.Sprintf("\nCurrent Repo: %v/%v", owner, name)

		url, err := c.processRepo(ctx, repo, dryRun)
		if err != nil {
			if errors.Is(err, ErrGetBranch) || errors.Is(err, ErrNoCommits) {
				repoBar.Incr()
				continue
			}

			return nil, fmt.Errorf("process repo: %w", err)
		}

		urls = append(urls, url)
		repoBar.Incr()
	}

	appendStr = ""

	sort.Strings(urls)

	return urls, nil
}

func (c *Client) processRepo(ctx context.Context, repo *github.Repository, dryRun bool) (string, error) {
	name := repo.GetName()
	owner := repo.GetOwner().GetLogin()
	head := repo.GetDefaultBranch()

	c.rate.Wait(ctx) //nolint: errcheck
	_, _, err := c.ghClient.Repositories.GetBranch(ctx, owner, name, c.cfg.ReleaseBranch)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrGetBranch, err)
	}

	opts := &github.PullRequestListOptions{
		Head: head,
		Base: c.cfg.ReleaseBranch,
	}

	c.rate.Wait(ctx) //nolint: errcheck
	prs, _, err := c.ghClient.PullRequests.List(ctx, owner, name, opts)
	if err != nil {
		return "", fmt.Errorf("list prs: %w", err)
	}

	if len(prs) > 0 {
		pr := prs[0]

		pr.Title = github.String("Release")

		changes, _ := c.createChangeLog(ctx, owner, name, c.cfg.ReleaseBranch, head)
		body := prBody(prBodyTemplate, changes)

		pr.Body = github.String(body)

		if !dryRun {
			c.rate.Wait(ctx) //nolint: errcheck
			pr, _, _ = c.ghClient.PullRequests.Edit(ctx, owner, name, pr.GetNumber(), pr)
		}

		return pr.GetHTMLURL(), nil
	}

	changes, err := c.createChangeLog(ctx, owner, name, c.cfg.ReleaseBranch, head)
	if err != nil {
		return "", err
	}

	body := prBody(prBodyTemplate, changes)

	newPR := &github.NewPullRequest{
		Title:               github.String("Release"),
		Head:                &head,
		Base:                &c.cfg.ReleaseBranch,
		Body:                github.String(body),
		MaintainerCanModify: github.Bool(true),
	}

	if !dryRun {
		c.rate.Wait(ctx) //nolint: errcheck
		pr, _, err := c.ghClient.PullRequests.Create(ctx, owner, name, newPR)
		if err != nil {
			return "", fmt.Errorf("create pr: %w", err)
		}

		return pr.GetHTMLURL(), nil
	}

	return fmt.Sprintf("https://github.com/%s/%s/compare/%s...%s", owner, name, c.cfg.ReleaseBranch, head), nil
}

func (c *Client) ReleaseRepos(ctx context.Context, progress *crawl.Progress, repos []*github.Repository, dryRun bool) ([]string, error) {
	releases, err := c.getReleases(ctx, progress, repos)
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

	theme := bar.NewThemeFromTheme(bar.DefaultTheme)
	theme.Append(func(b *bar.Bar) string {
		return fmt.Sprintf(" %0.2f", b.CompletedPercent())
	})
	theme.Append(func(b *bar.Bar) string {
		return appendStr
	})
	theme.Prepend(func(b *bar.Bar) string {
		return fmt.Sprintf("Processing Releases (%d/%d) %s", b.Current(), b.Total(), b.Elapsed())
	})

	repoBar := bar.New(theme, count)
	progress.AddBar(repoBar)

	var released []string
	for _, release := range releases {
		repo := release.GetBase().GetRepo()
		name = repo.GetName()
		owner = repo.GetOwner().GetLogin()
		appendStr = fmt.Sprintf("\nCurrent Repo: %v/%v", owner, name)

		var err error
		c.rate.Wait(ctx) //nolint: errcheck
		release, _, err = c.ghClient.PullRequests.Get(ctx, owner, name, release.GetNumber())
		if err != nil {
			return nil, fmt.Errorf("check mergeable: %v", err.Error())
		}

		if strings.ToLower(release.GetMergeableState()) != "clean" {
			repoBar.Incr()
			continue
		}

		if !dryRun {
			c.rate.Wait(ctx) //nolint: errcheck
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

		repoBar.Incr()
	}

	appendStr = ""

	sort.Strings(released)

	return released, nil
}

func (c *Client) getReleases(ctx context.Context, progress *crawl.Progress, repos []*github.Repository) ([]*github.PullRequest, error) {
	var releases []*github.PullRequest

	count := len(repos)
	name := repos[0].GetName()
	owner := repos[0].GetOwner().GetLogin()
	appendStr := fmt.Sprintf("\nCurrent Repo: %v/%v", owner, name)

	theme := bar.NewThemeFromTheme(bar.DefaultTheme)
	theme.Append(func(b *bar.Bar) string {
		return fmt.Sprintf(" %0.2f", b.CompletedPercent())
	})
	theme.Append(func(b *bar.Bar) string {
		return appendStr
	})
	theme.Prepend(func(b *bar.Bar) string {
		return fmt.Sprintf("Collecting Releases (%d/%d) %s", b.Current(), b.Total(), b.Elapsed())
	})

	repoBar := bar.New(theme, count)
	progress.AddBar(repoBar)

	for _, repo := range repos {
		owner = repo.GetOwner().GetLogin()
		name = repo.GetName()
		appendStr = fmt.Sprintf("\nCurrent Repo: %v/%v", owner, name)
		head := repo.GetDefaultBranch()

		opts := &github.PullRequestListOptions{
			Head: head,
			Base: c.cfg.ReleaseBranch,
		}

		c.rate.Wait(ctx) //nolint: errcheck
		rs, _, err := c.ghClient.PullRequests.List(ctx, owner, name, opts)
		if err != nil {
			return nil, fmt.Errorf("pull requests: %v", err.Error())
		}

		releases = append(releases, rs...)
		repoBar.Incr()
	}

	appendStr = ""

	return releases, nil
}
