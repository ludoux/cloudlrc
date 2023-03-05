package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var neteaseLoginGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gen called")
	},
}

func init() {
	neteaseLoginCmd.AddCommand(neteaseLoginGenCmd)
}
