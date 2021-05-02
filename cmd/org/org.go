package org

import (
	"github.com/spf13/cobra"
)

var OrgCmd = &cobra.Command{
	Use:   "org [flags]",
	Short: "Org specific release train commands",
	Long:  `Org specific release train commands`,
}
