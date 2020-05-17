package cmd

import (
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
	switch strings.ToLower(shell) {
	case "bash":
		rootCmd.GenBashCompletion(os.Stdout)
	case "fish":
		rootCmd.GenFishCompletion(os.Stdout, false)
	case "ps", "powershell", "power_shell":
		rootCmd.GenPowerShellCompletion(os.Stdout)
	case "zsh":
		rootCmd.GenZshCompletion(os.Stdout)
	default:
	}
}
