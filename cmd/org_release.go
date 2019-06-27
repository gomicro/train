package cmd

import (
	"os"

	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
)

func init() {
	orgCmd.AddCommand(orgReleaseCmd)
}

var orgReleaseCmd = &cobra.Command{
	Use:   "release [org_name]",
	Short: "Release PRs for an org's repos that can be merged",
	Run:   orgReleaseFunc,
}

func orgReleaseFunc(cmd *cobra.Command, args []string) {
	uiprogress.Start()

	repos, err := getOrgRepos(args[0])
	if err != nil {
		printf("org repos: %v", err.Error())
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
