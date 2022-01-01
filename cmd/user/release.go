package user

import (
	"context"
	"fmt"
	"os"

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
	ctx := context.Background()

	uiprogress.Start()

	repos, err := clt.GetUserRepos(ctx, args[0])
	if err != nil {
		fmt.Printf("user repos: %v\n", err.Error())
		os.Exit(1)
	}

	if len(repos) < 1 {
		fmt.Println("github: no repos found")
		return
	}

	urls, err := clt.ReleaseRepos(ctx, repos, base, dryRun)
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
