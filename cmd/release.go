package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(orgReleaseCmd)
}

var orgReleaseCmd = &cobra.Command{
	Use:               "release [org_name|user_name]",
	Short:             "Release PRs for an org or user's repos that can be merged",
	PersistentPreRun:  setupClient,
	Run:               orgReleaseFunc,
	ValidArgsFunction: releaseCmdValidArgsFunc,
}

func orgReleaseFunc(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	uiprogress.Start()

	if dryRun {
		fmt.Println()
		fmt.Println("===============")
		fmt.Println("Doing a dry run")
		fmt.Println("===============")
	}

	fmt.Println()

	repos, err := clt.GetRepos(ctx, args[0])
	if err != nil {
		fmt.Printf("repos: %v\n", err.Error())
		os.Exit(1)
	}

	urls, err := clt.ReleaseRepos(ctx, repos, dryRun)
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

func releaseCmdValidArgsFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	setupClient(cmd, args)

	valid, err := clt.GetLogins(context.Background())
	if err != nil {
		valid = []string{"error fetching"}
	}

	return valid, cobra.ShellCompDirectiveNoFileComp
}
