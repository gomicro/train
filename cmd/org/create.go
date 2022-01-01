package org

import (
	"context"
	"fmt"
	"os"

	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
)

func init() {
	OrgCmd.AddCommand(orgCreateCmd)
}

var orgCreateCmd = &cobra.Command{
	Use:   "create [org_name]",
	Short: "Create release PRs for an org's repos",
	Args:  cobra.ExactArgs(1),
	Run:   orgCreateFunc,
}

func orgCreateFunc(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	uiprogress.Start()

	fmt.Printf("Org: %v\n", args[0])
	fmt.Printf("Base: %v\n", base)

	if dryRun {
		fmt.Println()
		fmt.Println("===============")
		fmt.Println("Doing a dry run")
		fmt.Println("===============")
	}

	fmt.Println()

	repos, err := clt.GetOrgRepos(ctx, args[0])
	if err != nil {
		fmt.Printf("org repos: %v\n", err.Error())
		os.Exit(1)
	}

	if len(repos) < 1 {
		fmt.Printf("github: no repos found\n")
		return
	}

	urls, err := clt.ProcessRepos(ctx, repos, base, dryRun)
	if err != nil {
		fmt.Printf("org process repos: %v\n", err.Error())
		os.Exit(1)
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

		return
	}
}
