package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(orgCmd)
}

var orgCmd = &cobra.Command{
	Use:   "org [flags]",
	Short: "Org specific release train commands",
	Long:  `Org specific release train commands`,
}
