package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"github.com/gosuri/uiprogress"
)

// GetReleases takes a context, github client, and respos to iterate and collects
// release pull requests for the repos.
func GetReleases(ctx context.Context, client *github.Client, repos []*github.Repository) ([]*github.PullRequest, error) {
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
		base := "release"

		opts := &github.PullRequestListOptions{
			Head: head,
			Base: base,
		}

		rs, _, err := client.PullRequests.List(ctx, owner, name, opts)
		if err != nil {
			return nil, fmt.Errorf("pull requests: %v", err.Error())
		}

		releases = append(releases, rs...)
		bar.Incr()
	}

	appendStr = ""

	return releases, nil
}

// Release takex context, githug blient, a list of releases, and whether or not
// to perform a dry run. If it is not a dry run and the pull request is
// mergeable, it will merge it.
func Release(ctx context.Context, client *github.Client, releases []*github.PullRequest, dryRun bool) ([]string, error) {
	var released []string

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

	for _, release := range releases {
		repo := release.GetBase().GetRepo()
		name = repo.GetName()
		owner = repo.GetOwner().GetLogin()
		appendStr = fmt.Sprintf("\nCurrent Repo: %v/%v", owner, name)

		var err error
		release, _, err = client.PullRequests.Get(ctx, owner, name, release.GetNumber())
		if err != nil {
			return nil, fmt.Errorf("check mergeable: %v", err.Error())
		}

		if strings.ToLower(release.GetMergeableState()) != "clean" {
			bar.Incr()
			continue
		}

		if !dryRun {
			res, _, err := client.PullRequests.Merge(ctx, owner, name, release.GetNumber(), "release automerged by train", nil)
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
