package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gomicro/train/repositories"

	"github.com/google/go-github/github"
	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
)

func init() {
	userCmd.AddCommand(userCreateCmd)
}

var userCreateCmd = &cobra.Command{
	Use:   "create [username]",
	Short: "Create release PRs for a user's repos",
	Long:  `Create release PRs for a user's repos`,
	Args:  cobra.ExactArgs(1),
	Run:   userCreateFunc,
}

func userCreateFunc(cmd *cobra.Command, args []string) {
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

	count := len(repos)
	name := repos[0].GetName()
	owner := repos[0].GetOwner().GetLogin()

	bar := uiprogress.AddBar(count).
		AppendCompleted().
		PrependElapsed().
		PrependFunc(func(b *uiprogress.Bar) string {
			return fmt.Sprintf("Processing (%d/%d)", b.Current(), count)
		}).
		AppendFunc(func(b *uiprogress.Bar) string {
			return fmt.Sprintf("\nCurrent Repo: %v/%v", owner, name)
		})

	urls := []string{}
	for _, repo := range repos {
		name = repo.GetName()
		owner = repo.GetOwner().GetLogin()

		url, err := repositories.Process(clientCtx, client, repo)
		if err != nil {
			if strings.HasPrefix(err.Error(), "get branch: ") || strings.HasPrefix(err.Error(), "no commits") {
				bar.Incr()
				continue
			}

			printf("process repo: %v", err.Error())
			os.Exit(1)
		}

		urls = append(urls, url)
		bar.Incr()
	}

	uiprogress.Stop()

	if len(urls) > 0 {
		printf("")
		printf("Release PRs Created:")
		for _, url := range urls {
			printf(url)
		}

		return
	}
}

func getUserRepos(username string) ([]*github.Repository, error) {
	u, _, err := client.Users.Get(clientCtx, username)
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
		rs, resp, err := client.Repositories.List(clientCtx, username, nil)
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
