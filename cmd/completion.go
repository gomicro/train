package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	defaultShell = "zsh"
)

var (
	shell string
)

func init() {
	rootCmd.AddCommand(CompletionCmd)

	CompletionCmd.Flags().StringVar(&shell, "shell", defaultShell, "desired shell to generate completions for")
}

// CompletionCmd represents the command for generating completion files for the
// train cli.
var CompletionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate completion files for the train cli",
	Run:   completionFunc,
}

func completionFunc(cmd *cobra.Command, args []string) {
	var err error
	switch strings.ToLower(shell) {
	case "bash":
		err = rootCmd.GenBashCompletion(os.Stdout)
	case "fish":
		err = rootCmd.GenFishCompletion(os.Stdout, false)
	case "ps", "powershell", "power_shell":
		err = rootCmd.GenPowerShellCompletion(os.Stdout)
	case "zsh":
		err = rootCmd.GenZshCompletion(os.Stdout)
	default:
	}

	if err != nil {
		fmt.Printf("error generating completion output: %v", err.Error())
		os.Exit(1)
	}
}
