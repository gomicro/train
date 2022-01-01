package user

import (
	"fmt"
	"os"

	"github.com/gomicro/train/client"
	"github.com/gomicro/train/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	base   string
	dryRun bool
	clt    *client.Client
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

	base = "release"
	if c.ReleaseBranch != "" {
		base = c.ReleaseBranch
	}

	clt, err = client.New(c.Github.Token, base)
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		os.Exit(1)
	}

	dryRun = viper.GetBool("dryRun")
}
