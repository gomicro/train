package repositories

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

// Process takes a github client, context for the client, and a repo to process.
// It returns the resulting pull request URL or an error if any are encountered.
func Process(ctx context.Context, client *github.Client, repo *github.Repository) (string, error) {
	name := repo.GetName()
	owner := repo.GetOwner().GetLogin()
	head := repo.GetDefaultBranch()
	base := "release"

	_, _, err := client.Repositories.GetBranch(ctx, owner, name, base)
	if err != nil {
		return "", fmt.Errorf("get branch: %v", err.Error())
	}

	opts := &github.PullRequestListOptions{
		Head: head,
		Base: base,
	}

	prs, _, err := client.PullRequests.List(ctx, owner, name, opts)
	if err != nil {
		return "", fmt.Errorf("list prs: %v", err.Error())
	}

	if len(prs) > 0 {
		pr := prs[0]

		pr.Title = github.String("Release")
		pr.Body = github.String("Release created with `train`")

		pr, _, _ = client.PullRequests.Edit(ctx, owner, name, pr.GetNumber(), pr)
		return pr.GetHTMLURL(), nil
	}

	comp, _, err := client.Repositories.CompareCommits(ctx, owner, name, base, head)
	if err != nil {
		return "", fmt.Errorf("compare commits: %v", err.Error())
	}

	if len(comp.Commits) == 0 {
		return "", fmt.Errorf("no commits")
	}

	newPR := &github.NewPullRequest{
		Title:               github.String("Release"),
		Head:                &head,
		Base:                &base,
		Body:                github.String("Release created with `train`"),
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := client.PullRequests.Create(ctx, owner, name, newPR)
	if err != nil {
		return "", fmt.Errorf("create pr: %v", err.Error())
	}

	return pr.GetHTMLURL(), nil
}
