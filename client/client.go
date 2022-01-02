package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gomicro/trust"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Client struct {
	base     string
	ghClient *github.Client
}

func New(ghToken, base string) (*Client, error) {
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
			AccessToken: ghToken,
		},
	)

	return &Client{
		base:     base,
		ghClient: github.NewClient(oauth2.NewClient(ctx, ts)),
	}, nil
}
