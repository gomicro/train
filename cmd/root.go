package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/certifi/gocertifi"
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
		printf("failed to create cert pool: %v", err.Error())
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
		printf("warning: TRAIN_GHTOKEN missing")
	}

	clientCtx = context.Background()
	clientCtx = context.WithValue(clientCtx, oauth2.HTTPClient, httpClient)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)

	client = github.NewClient(oauth2.NewClient(clientCtx, ts))
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
		printf("Failed to execute: %v", err.Error())
		os.Exit(1)
	}
}

func printf(f string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(f, args...))
}

func verbosef(f string, args ...interface{}) {
	if verbose {
		fmt.Println(fmt.Sprintf(f, args...))
	}
}
