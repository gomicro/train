package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/gomicro/crawl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewCreateCmd(os.Stdout))
}

func NewCreateCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create [org_name|user_name]",
		Short:             "Create release PRs for an org or user's repos",
		Args:              cobra.ExactArgs(1),
		PersistentPreRun:  setupClient,
		RunE:              createRun(out),
		ValidArgsFunction: createCmdValidArgsFunc,
	}

	return cmd
}

func createRun(out io.Writer) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		progress := crawl.New(ctx, out)
		progress.SetOut(out)
		progress.Start()

		entity := args[0]

		fmt.Fprintf(out, "Entity: %s\n", entity)
		fmt.Fprintf(out, "Base: %s\n", clt.GetBaseBranchName())

		if dryRun {
			fmt.Fprintln(out)
			fmt.Fprintln(out, "===============")
			fmt.Fprintln(out, "Doing a dry run")
			fmt.Fprintln(out, "===============")
		}

		fmt.Fprintln(out)

		repos, err := clt.GetRepos(ctx, progress, entity)
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("create: %w", err)
		}

		urls, err := clt.ProcessRepos(ctx, progress, repos, dryRun)
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("create: %w", err)
		}

		progress.Stop()

		if len(urls) > 0 {
			fmt.Fprintln(out)
			if dryRun {
				fmt.Fprintln(out, "(Dryrun) Release PRs Created:")
			} else {
				fmt.Fprintln(out, "Release PRs Created:")
			}

			for _, url := range urls {
				fmt.Fprintln(out, url)
			}
		}

		return nil
	}
}

func createCmdValidArgsFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	setupClient(cmd, args)

	valid, err := clt.GetLogins(context.Background())
	if err != nil {
		valid = []string{"error fetching"}
	}

	return valid, cobra.ShellCompDirectiveNoFileComp
}
