package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(userCmd)
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User specific release train commands",
	Long:  `User specific release train commands`,
}
