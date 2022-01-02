package client

import (
	"context"
	"fmt"
	"strings"
)

var prBodyTemplate = `
----
Release PR created with ` + "`train`"

func (c *Client) createChangeLog(ctx context.Context, owner, name, base, head string) (map[string][]string, error) {
	changes := map[string][]string{
		"added":      []string{},
		"changed":    []string{},
		"deprecated": []string{},
		"removed":    []string{},
		"fixed":      []string{},
		"security":   []string{},
	}

	comp, _, err := c.ghClient.Repositories.CompareCommits(ctx, owner, name, base, head)
	if err != nil {
		return nil, fmt.Errorf("compare commits: %v", err.Error())
	}

	if len(comp.Commits) == 0 {
		return nil, fmt.Errorf("no commits")
	}

	for _, commit := range comp.Commits {
		c := strings.Split(strings.ToLower(*commit.Commit.Message), "\n")[0]

		if strings.Contains(c, "merge pull request") {
			continue
		}

		for synonim, label := range changeMapping {
			if strings.HasPrefix(c, synonim) {
				changes[label] = append(changes[label], strings.Join(strings.Split(c, " ")[1:], " "))
				break
			}

			if strings.Contains(c, synonim) {
				changes[label] = append(changes[label], c)
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
	"add":      "added",
	"added":    "added",
	"adding":   "added",
	"adds":     "added",
	"created":  "added",
	"creating": "added",

	"altering":   "changed",
	"change":     "changed",
	"changed":    "changed",
	"changes":    "changed",
	"changing":   "changed",
	"convert":    "changed",
	"converted":  "changed",
	"converting": "changed",
	"replace":    "changed",
	"replaced":   "changed",
	"replacing":  "changed",
	"update":     "changed",
	"updating":   "changed",

	"deprecate":   "deprecated",
	"deprecated":  "deprecated",
	"deprecating": "deprecated",

	"detach":    "removed",
	"detached":  "removed",
	"detaching": "removed",
	"remove":    "removed",
	"removed":   "removed",
	"removes":   "removed",
	"removing":  "removed",

	"correct":    "fixed",
	"corrected":  "fixed",
	"correcting": "fixed",
	"fixed":      "fixed",
	"fixing":     "fixed",
	"resolved":   "fixed",
	"resolving":  "fixed",

	"security": "security",
	"securing": "security",
}
