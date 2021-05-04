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
	var err error
	client, err = config.GetClient()
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		os.Exit(1)
	}

	dryRun = viper.GetBool("dryRun")
}
