package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/certifi/gocertifi"
	"github.com/gomicro/train/cmd/org"
	"github.com/gomicro/train/cmd/user"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

const (
	defaultMaxIdleConns        = 25
	defaultMaxIdleConnsPerHost = 25
)

var (
	verbose   bool
	dryRun    bool
	client    *github.Client
	clientCtx context.Context
)

func init() {
	cobra.OnInitialize(initEnvs)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show more verbose output")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dryRun", "d", false, "attempt the specified command without actually making live changes")

	pool, err := gocertifi.CACerts()
	if err != nil {
		fmt.Printf("failed to create cert pool: %v\n", err.Error())
		os.Exit(1)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        defaultMaxIdleConns,
			MaxIdleConnsPerHost: defaultMaxIdleConnsPerHost,
			TLSClientConfig: &tls.Config{
				RootCAs: pool,
			},
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

	rootCmd.AddCommand(user.UserCmd)
	rootCmd.AddCommand(org.OrgCmd)
}

func initEnvs() {
}

var rootCmd = &cobra.Command{
	Use:   "train [flags]",
	Short: "Lightweight for managing release PRs",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Failed to execute: %v\n", err.Error())
		os.Exit(1)
	}
}
