package cmd

import (
	"fmt"

	"github.com/ludoux/cloudlrc/netease"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var neteaseLoginCheckCmd = &cobra.Command{
	Use:   "check <unikey>",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			_, s := netease.LoginCheck(cast.ToString(args[0]))
			fmt.Println(s)
		}

	},
}

func init() {
	neteaseLoginCmd.AddCommand(neteaseLoginCheckCmd)
}
