package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gomicro/trust"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func GetClient() (*github.Client, error) {
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

	config, err := ParseFromFile()
	if err != nil {
		return nil, err
	}

	clientCtx := context.Background()
	clientCtx = context.WithValue(clientCtx, oauth2.HTTPClient, httpClient)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: config.Github.Token,
		},
	)

	return github.NewClient(oauth2.NewClient(clientCtx, ts)), nil
}
