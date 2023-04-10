package clienttest

import (
	"context"
	"fmt"

	"github.com/gomicro/crawl"
	"github.com/google/go-github/github"
)

type ClientTest struct {
	cfg *Config
}

type Config struct {
	BaseBranchName    string
	Logins            []string
	LoginsError       error
	Repos             []*github.Repository
	ReposError        error
	ProcessReposError error
}

func New(cfg *Config) *ClientTest {
	return &ClientTest{
		cfg: cfg,
	}
}

func (ct *ClientTest) GetBaseBranchName() string {
	if ct.cfg != nil {
		return ct.cfg.BaseBranchName
	}

	return ""
}

func (ct *ClientTest) GetLogins(context.Context) ([]string, error) {
	if ct.cfg.LoginsError != nil {
		return nil, ct.cfg.LoginsError
	}

	return ct.cfg.Logins, nil
}

func (ct *ClientTest) GetRepos(ctx context.Context, progress *crawl.Progress, name string) ([]*github.Repository, error) {
	if ct.cfg.ReposError != nil {
		return nil, ct.cfg.ReposError
	}

	return ct.cfg.Repos, nil
}

func (ct *ClientTest) ProcessRepos(ctx context.Context, progress *crawl.Progress, repos []*github.Repository, dryRun bool) ([]string, error) {
	if ct.cfg.ProcessReposError != nil {
		return nil, ct.cfg.ProcessReposError
	}

	urls := make([]string, len(ct.cfg.Repos))

	for i, r := range repos {
		name := r.GetName()
		owner := r.GetOwner().GetLogin()
		head := r.GetDefaultBranch()

		if !dryRun {
			urls = append(urls, fmt.Sprintf("https://github.com/%s/%s/pull/%d", owner, name, i))
			continue
		}

		urls = append(urls, fmt.Sprintf("https://github.com/%s/%s/compare/%s...%s", owner, name, ct.cfg.BaseBranchName, head))
	}

	return urls, nil
}

func (ct *ClientTest) ReleaseRepos(ctx context.Context, progress *crawl.Progress, repos []*github.Repository, dryRun bool) ([]string, error) {
	return nil, nil
}
