package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/gomicro/crawl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(releaseCmd)
}

var releaseCmd = &cobra.Command{
	Use:               "release [org_name|user_name]",
	Short:             "Release PRs for an org or user's repos that can be merged",
	PersistentPreRun:  setupClient,
	RunE:              releaseFunc,
	ValidArgsFunction: releaseCmdValidArgsFunc,
}

func releaseFunc(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	progress := crawl.New(ctx, os.Stdout)
	progress.Start()

	if dryRun {
		fmt.Println()
		fmt.Println("===============")
		fmt.Println("Doing a dry run")
		fmt.Println("===============")
	}

	fmt.Println()

	repos, err := clt.GetRepos(ctx, progress, args[0])
	if err != nil {
		cmd.SilenceUsage = true
		return fmt.Errorf("release: %w", err)
	}

	urls, err := clt.ReleaseRepos(ctx, progress, repos, dryRun)
	if err != nil {
		cmd.SilenceUsage = true
		return fmt.Errorf("release: %w", err)
	}

	progress.Stop()

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
	}

	return nil
}

func releaseCmdValidArgsFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	setupClient(cmd, args)

	valid, err := clt.GetLogins(context.Background())
	if err != nil {
		valid = []string{"error fetching"}
	}

	return valid, cobra.ShellCompDirectiveNoFileComp
}
