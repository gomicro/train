package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	defaultShell = "zsh"
)

var (
	shell          string
	ErrUknownShell = fmt.Errorf("unrecognized shell")
)

func init() {
	rootCmd.AddCommand(NewCompletionCmd(os.Stdout))
}

func NewCompletionCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate completion files for the train cli",
		RunE:  completionRun(out),
	}

	cmd.Flags().StringVar(&shell, "shell", defaultShell, "desired shell to generate completions for")

	cmd.SetOut(out)

	return cmd
}

func completionRun(out io.Writer) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var err error

		switch strings.ToLower(shell) {
		case "bash":
			err = rootCmd.GenBashCompletion(out)
		case "fish":
			err = rootCmd.GenFishCompletion(out, false)
		case "ps", "powershell", "power_shell":
			err = rootCmd.GenPowerShellCompletion(out)
		case "zsh":
			err = rootCmd.GenZshCompletion(out)
		default:
			err = ErrUknownShell
		}

		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("completion: %w", err)
		}

		return nil
	}
}
