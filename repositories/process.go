package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
)

// Process takes a github client, context for the client, and a repo to process.
// It returns the resulting pull request URL or an error if any are encountered.
func Process(ctx context.Context, client *github.Client, repo *github.Repository, prBodyTemplate string) (string, error) {
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

		changes, _ := changeLog(ctx, client, owner, name, base, head)
		body := prBody(prBodyTemplate, changes)

		pr.Body = github.String(body)

		pr, _, _ = client.PullRequests.Edit(ctx, owner, name, pr.GetNumber(), pr)
		return pr.GetHTMLURL(), nil
	}

	changes, err := changeLog(ctx, client, owner, name, base, head)
	if err != nil {
		return "", err
	}

	body := prBody(prBodyTemplate, changes)

	newPR := &github.NewPullRequest{
		Title:               github.String("Release"),
		Head:                &head,
		Base:                &base,
		Body:                github.String(body),
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := client.PullRequests.Create(ctx, owner, name, newPR)
	if err != nil {
		return "", fmt.Errorf("create pr: %v", err.Error())
	}

	return pr.GetHTMLURL(), nil
}

func changeLog(ctx context.Context, client *github.Client, owner, name, base, head string) (map[string][]string, error) {
	changes := map[string][]string{
		"added":      []string{},
		"changed":    []string{},
		"deprecated": []string{},
		"removed":    []string{},
		"fixed":      []string{},
		"security":   []string{},
	}

	comp, _, err := client.Repositories.CompareCommits(ctx, owner, name, base, head)
	if err != nil {
		return nil, fmt.Errorf("compare commits: %v", err.Error())
	}

	if len(comp.Commits) == 0 {
		return nil, fmt.Errorf("no commits")
	}

	for _, commit := range comp.Commits {
		c := strings.Split(strings.ToLower(*commit.Commit.Message), "\n")[0]

		for synonim, label := range changeMapping {
			if strings.HasPrefix(c, synonim) {
				changes[label] = append(changes[label], strings.TrimPrefix(*commit.Commit.Message, synonim))
				break
			}

			if strings.Contains(c, synonim) {
				changes[label] = append(changes[label], strings.TrimPrefix(*commit.Commit.Message, synonim))
				break
			}
		}
	}

	return changes, nil
}

func prBody(prBodyTemplate string, changes map[string][]string) string {
	body := ""

	for _, label := range changeOrder {
		for _, log := range changes[label] {
			body = body + fmt.Sprintf("* `%v` %v\n", strings.ToTitle(label), log)
		}
	}

	if body == "" {
		body = "no change log detected &mdash; try favoring words like `added`, `changed`, or `removed`\n"
	}

	return body + prBodyTemplate
}

var changeOrder = []string{
	"added",
	"changed",
	"deprecated",
	"removed",
	"fixed",
	"security",
}

var changeMapping = map[string]string{
	"added":    "added",
	"adding":   "added",
	"created":  "added",
	"creating": "added",

	"altering":  "changed",
	"changed":   "changed",
	"changing":  "changed",
	"replace":   "changed",
	"replaced":  "changed",
	"replacing": "changed",

	"deprecated":  "deprecated",
	"deprecating": "deprecated",

	"detached":  "removed",
	"detaching": "removed",
	"removed":   "removed",
	"removing":  "removed",

	"corrected":  "fixed",
	"correcting": "fixed",
	"fixed":      "fixed",
	"fixing":     "fixed",
	"updating":  "fixed",
	"resolved":   "fixed",
	"resolving":  "fixed",

	"security": "security",
	"securing": "security",
}
