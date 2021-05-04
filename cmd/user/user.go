package user

import (
	"fmt"
	"os"

	"github.com/gomicro/train/config"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	base   string
	dryRun bool
	client *github.Client
)

// UserCmd represents the root of the user command
var UserCmd = &cobra.Command{
	Use:              "user",
	Short:            "User specific release train commands",
	Long:             `User specific release train commands`,
	PersistentPreRun: setupCommand,
}

func setupCommand(cmd *cobra.Command, args []string) {
	c, err := config.ParseFromFile()
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		os.Exit(1)
	}

	client, err = c.GetClient()
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		os.Exit(1)
	}

	base = "release"
	if c.ReleaseBranch != "" {
		base = c.ReleaseBranch
	}

	dryRun = viper.GetBool("dryRun")
}
