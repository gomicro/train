package cmd

import (
	"fmt"
	"os"
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
	Run:       configFunc,
	ValidArgs: configValidArgs,
}

func configFunc(cmd *cobra.Command, args []string) {
	field := args[0]
	value := args[1]

	confFile, err := config.ParseFromFile()
	if err != nil {
		fmt.Printf("error: %v", err.Error())
	}

	switch strings.ToLower(field) {
	case "release_branch":
		confFile.ReleaseBranch = value
	default:
		fmt.Printf("unreconized config field: %v\n", field)
		os.Exit(1)
	}

	err = confFile.WriteFile()
	if err != nil {
		fmt.Printf("error: %v", err.Error())
		os.Exit(1)
	}

	fmt.Println("Config file updated")
}
