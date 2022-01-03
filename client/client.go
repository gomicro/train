package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/gomicro/train/config"
	"github.com/gomicro/trust"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

type Client struct {
	cfg      *config.Config
	ghClient *github.Client
	rate     *rate.Limiter

	ignoreRepoMap  map[string]struct{}
	ignoreTopicMap map[string]struct{}
}

func New(cfg *config.Config) (*Client, error) {
	pool := trust.New()

	certs, err := pool.CACerts()
	if err != nil {
		return nil, fmt.Errorf("failed to create cert pool: %v\n", err.Error())
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: certs},
		},
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: cfg.Github.Token,
		},
	)

	rl := rate.NewLimiter(
		rate.Limit(cfg.Github.Limits.RequestsPerSecond),
		cfg.Github.Limits.Burst,
	)

	irMap := map[string]struct{}{}
	for i := range cfg.Github.Ignores.Repos {
		irMap[strings.ToLower(cfg.Github.Ignores.Repos[i])] = struct{}{}
	}

	itMap := map[string]struct{}{}
	for i := range cfg.Github.Ignores.Topics {
		itMap[strings.ToLower(cfg.Github.Ignores.Topics[i])] = struct{}{}
	}

	return &Client{
		cfg:      cfg,
		ghClient: github.NewClient(oauth2.NewClient(ctx, ts)),
		rate:     rl,

		ignoreRepoMap:  irMap,
		ignoreTopicMap: itMap,
	}, nil
}
