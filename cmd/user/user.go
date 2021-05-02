package user

import (
	"github.com/spf13/cobra"
)

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "User specific release train commands",
	Long:  `User specific release train commands`,
}
