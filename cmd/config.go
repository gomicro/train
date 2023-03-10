package cmd

import (
	"fmt"
	"strings"

	"github.com/gomicro/train/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configValidArgs = []string{
	"release_branch\tthe head branch name to use for creating the release PRs",
}

var configCmd = &cobra.Command{
	Use:       "config [config_field] [value]",
	Short:     "config train",
	Long:      `configure train`,
	Args:      cobra.ExactArgs(2),
	RunE:      configFunc,
	ValidArgs: configValidArgs,
}

func configFunc(cmd *cobra.Command, args []string) error {
	field := args[0]
	value := args[1]

	confFile, err := config.ParseFromFile()
	if err != nil {
		cmd.SilenceUsage = true
		return fmt.Errorf("config: %w", err)
	}

	switch strings.ToLower(field) {
	case "release_branch":
		confFile.ReleaseBranch = value
	default:
		cmd.SilenceUsage = true
		return fmt.Errorf("config: unreconized config field: %s", field)
	}

	err = confFile.WriteFile()
	if err != nil {
		cmd.SilenceUsage = true
		return fmt.Errorf("config: %w", err)
	}

	// TODO: change to verbose output
	fmt.Println("Config file updated")

	return nil
}
