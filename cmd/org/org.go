package org

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/gomicro/trust"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var (
	dryRun    bool
	client    *github.Client
	clientCtx context.Context
)

// OrgCmd represents the root of the org command
var OrgCmd = &cobra.Command{
	Use:              "org [flags]",
	Short:            "Org specific release train commands",
	PersistentPreRun: configClient,
}

func configClient(cmd *cobra.Command, args []string) {
	pool := trust.New()

	certs, err := pool.CACerts()
	if err != nil {
		fmt.Printf("failed to create cert pool: %v\n", err.Error())
		os.Exit(1)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: certs},
		},
	}

	token := os.Getenv("TRAIN_GHTOKEN")

	if token == "" {
		fmt.Printf("warning: TRAIN_GHTOKEN missing\n")
	}

	clientCtx = context.Background()
	clientCtx = context.WithValue(clientCtx, oauth2.HTTPClient, httpClient)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)

	client = github.NewClient(oauth2.NewClient(clientCtx, ts))

	dryRun = viper.GetBool("dryRun")
}
