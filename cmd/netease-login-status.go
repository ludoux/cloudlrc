package cmd

import (
	"fmt"

	"github.com/ludoux/cloudlrc/netease"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var neteaseLoginStatusCmd = &cobra.Command{
	Use:   "status", //app netease login status
	Short: "检测本地 cookie 网易云登录状态",
	Long:  `检测本地 cookie 网易云登录状态`,
	Run: func(cmd *cobra.Command, args []string) {
		msg, ornot := netease.LoginStatus()
		if ornot {
			fmt.Println("已登录，账户:", msg)
		} else {
			fmt.Println("未登录")
		}
	},
}

func init() {
	neteaseLoginCmd.AddCommand(neteaseLoginStatusCmd)
}
