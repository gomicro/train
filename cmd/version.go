package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewVersionCmd(os.Stdout, Version))
}

var (
	// Version is the current version of train, made available for use through
	// out the application.
	Version string
)

func NewVersionCmd(out io.Writer, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display the version",
		Long:  `Display the version of the CLI.`,
		Run:   versionRun(version),
	}

	cmd.SetOut(out)

	return cmd
}

func versionRun(version string) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if version == "" {
			cmd.Println("Train version dev-local")
			return
		}

		cmd.Printf("Train version %s\n", version)
	}
}
