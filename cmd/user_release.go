package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
)

func init() {
	userCmd.AddCommand(userReleaseCmd)
}

var userReleaseCmd = &cobra.Command{
	Use: "release [username]",
	Run: userReleaseFunc,
}

func userReleaseFunc(cmd *cobra.Command, args []string) {
	uiprogress.Start()

	repos, err := getUserRepos(args[0])
	if err != nil {
		printf("user repos: %v", err.Error())
		os.Exit(1)
	}

	if len(repos) < 1 {
		printf("github: no repos found")
		return
	}

	releases, err := getReleases(repos)
	if err != nil {
		printf("releases: %v", err.Error())
		os.Exit(1)
	}

	urls, err := release(releases)
	if err != nil {
		printf("releasing: %v", err.Error())
		os.Exit(1)
	}

	uiprogress.Stop()

	if len(urls) > 0 {
		printf("")
		printf("Repos Released:")
		for _, url := range urls {
			printf(url)
		}

		return
	}
}

func getReleases(repos []*github.Repository) ([]*github.PullRequest, error) {
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

		rs, _, err := client.PullRequests.List(clientCtx, owner, name, opts)
		if err != nil {
			return nil, fmt.Errorf("pull requests: %v", err.Error())
		}

		releases = append(releases, rs...)
		bar.Incr()
	}

	appendStr = ""

	return releases, nil
}

func release(releases []*github.PullRequest) ([]string, error) {
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
		release, _, err = client.PullRequests.Get(clientCtx, owner, name, release.GetNumber())
		if err != nil {
			return nil, fmt.Errorf("check mergeable: %v", err.Error())
		}

		if strings.ToLower(release.GetMergeableState()) != "clean" {
			bar.Incr()
			continue
		}

		res, _, err := client.PullRequests.Merge(clientCtx, owner, name, release.GetNumber(), "release automerged by train", nil)
		if err != nil {
			return nil, fmt.Errorf("merge: %v", err.Error())
		}

		if res.GetMerged() {
			released = append(released, release.GetHTMLURL())
		}

		bar.Incr()
	}

	appendStr = ""

	return released, nil
}
