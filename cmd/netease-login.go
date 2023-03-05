package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// neteaseLoginCmd represents the login command
var neteaseLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("login called")
	},
}

func init() {
	neteaseCmd.AddCommand(neteaseLoginCmd)
}
