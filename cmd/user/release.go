package user

import (
	"fmt"
	"os"

	"github.com/gomicro/train/repositories"
	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
)

func init() {
	UserCmd.AddCommand(userReleaseCmd)
}

var userReleaseCmd = &cobra.Command{
	Use:   "release [username]",
	Short: "Release PRs for a user's repos that can be merged",
	Run:   userReleaseFunc,
}

func userReleaseFunc(cmd *cobra.Command, args []string) {
	uiprogress.Start()

	repos, err := getUserRepos(args[0])
	if err != nil {
		fmt.Printf("user repos: %v\n", err.Error())
		os.Exit(1)
	}

	if len(repos) < 1 {
		fmt.Println("github: no repos found")
		return
	}

	releases, err := repositories.GetReleases(clientCtx, client, repos)
	if err != nil {
		fmt.Printf("releases: %v\n", err.Error())
		os.Exit(1)
	}

	urls, err := repositories.Release(clientCtx, client, releases, dryRun)
	if err != nil {
		fmt.Printf("releasing: %v\n", err.Error())
		os.Exit(1)
	}

	uiprogress.Stop()

	if len(urls) > 0 {
		fmt.Println()
		if dryRun {
			fmt.Println("(Dryrun) Repos Released:")
		} else {
			fmt.Println("Repos Released:")
		}

		for _, url := range urls {
			fmt.Println(url)
		}

		return
	}
}
