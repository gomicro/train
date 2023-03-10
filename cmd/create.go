package cmd

import (
	"context"
	"fmt"

	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:               "create [org_name|user_name]",
	Short:             "Create release PRs for an org or user's repos",
	Args:              cobra.ExactArgs(1),
	PersistentPreRun:  setupClient,
	RunE:              createFunc,
	ValidArgsFunction: createCmdValidArgsFunc,
}

func createFunc(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	uiprogress.Start()

	fmt.Printf("Entity: %v\n", args[0])
	fmt.Printf("Base: %v\n", clt.GetBaseBranchName())

	if dryRun {
		fmt.Println()
		fmt.Println("===============")
		fmt.Println("Doing a dry run")
		fmt.Println("===============")
	}

	fmt.Println()

	repos, err := clt.GetRepos(ctx, args[0])
	if err != nil {
		cmd.SilenceUsage = true
		return fmt.Errorf("create: %w", err)
	}

	urls, err := clt.ProcessRepos(ctx, repos, dryRun)
	if err != nil {
		cmd.SilenceUsage = true
		return fmt.Errorf("create: %w", err)
	}

	uiprogress.Stop()

	if len(urls) > 0 {
		fmt.Println()
		if dryRun {
			fmt.Println("(Dryrun) Release PRs Created:")
		} else {
			fmt.Println("Release PRs Created:")
		}

		for _, url := range urls {
			fmt.Println(url)
		}
	}

	return nil
}

func createCmdValidArgsFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	setupClient(cmd, args)

	valid, err := clt.GetLogins(context.Background())
	if err != nil {
		valid = []string{"error fetching"}
	}

	return valid, cobra.ShellCompDirectiveNoFileComp
}
