package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gomicro/train/config"
	"github.com/gomicro/trust"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

type Client struct {
	base     string
	ghClient *github.Client
	rate     *rate.Limiter
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

	return &Client{
		base:     cfg.ReleaseBranch,
		ghClient: github.NewClient(oauth2.NewClient(ctx, ts)),
		rate:     rl,
	}, nil
}
