package user

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gomicro/train/repositories"

	"github.com/google/go-github/github"
	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
)

func init() {
	UserCmd.AddCommand(userCreateCmd)
}

var userCreateCmd = &cobra.Command{
	Use:   "create [username]",
	Short: "Create release PRs for a user's repos",
	Args:  cobra.ExactArgs(1),
	Run:   userCreateFunc,
}

func userCreateFunc(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	uiprogress.Start()

	repos, err := getUserRepos(ctx, args[0])
	if err != nil {
		fmt.Printf("user repos: %v\n", err.Error())
		os.Exit(1)
	}

	if len(repos) < 1 {
		fmt.Println("github: no repos found")
		return
	}

	count := len(repos)
	name := repos[0].GetName()
	owner := repos[0].GetOwner().GetLogin()
	appendStr := fmt.Sprintf("\nCurrent Repo: %v/%v", owner, name)

	bar := uiprogress.AddBar(count).
		AppendCompleted().
		PrependElapsed().
		PrependFunc(func(b *uiprogress.Bar) string {
			return fmt.Sprintf("Processing (%d/%d)", b.Current(), count)
		}).
		AppendFunc(func(b *uiprogress.Bar) string {
			return appendStr
		})

	urls := []string{}
	for _, repo := range repos {
		name = repo.GetName()
		owner = repo.GetOwner().GetLogin()
		appendStr = fmt.Sprintf("\nCurrent Repo: %v/%v", owner, name)

		url, err := repositories.Process(ctx, client, repo, base, dryRun)
		if err != nil {
			if strings.HasPrefix(err.Error(), "get branch: ") || strings.HasPrefix(err.Error(), "no commits") {
				bar.Incr()
				continue
			}

			fmt.Printf("process repo: %v\n", err.Error())
			os.Exit(1)
		}

		urls = append(urls, url)
		bar.Incr()
	}

	appendStr = ""

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

func getUserRepos(ctx context.Context, username string) ([]*github.Repository, error) {
	u, _, err := client.Users.Get(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get user: %v", err.Error())
	}

	count := u.GetPublicRepos() + u.GetTotalPrivateRepos()

	bar := uiprogress.AddBar(count).
		AppendCompleted().
		PrependElapsed().
		PrependFunc(func(b *uiprogress.Bar) string {
			return fmt.Sprintf("Fetching (%d/%d)", b.Current(), count)
		})

	var repos []*github.Repository

	opts := &github.RepositoryListByOrgOptions{
		Type: "all",
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: 100,
		},
	}

	for {
		rs, resp, err := client.Repositories.List(ctx, username, nil)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return nil, fmt.Errorf("github: hit rate limit")
			}

			return nil, fmt.Errorf("list repos: %v", err.Error())
		}

		for range rs {
			bar.Incr()
		}

		repos = append(repos, rs...)

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return repos, nil
}
