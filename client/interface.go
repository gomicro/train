package client

import (
	"context"

	"github.com/gomicro/crawl"
	"github.com/google/go-github/github"
)

// interface for a train client
type Clienter interface {
	GetBaseBranchName() string
	GetLogins(context.Context) ([]string, error)
	GetRepos(context.Context, *crawl.Progress, string) ([]*github.Repository, error)
	ProcessRepos(context.Context, *crawl.Progress, []*github.Repository, bool) ([]string, error)
	ReleaseRepos(context.Context, *crawl.Progress, []*github.Repository, bool) ([]string, error)
}
