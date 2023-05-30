package cmd

import (
	"fmt"

	"github.com/ludoux/cloudlrc/netease"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var qrfile bool
var neteaseCmd = &cobra.Command{
	Use:   "netease",
	Short: "接口: 网易云音乐",
	Long:  `接口: 网易云音乐`,
	Run: func(cmd *cobra.Command, args []string) {
		//
	},
}
var neteaseLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("login called")
	},
}
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
var neteaseLoginGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		unikey := netease.LoginGen(qrfile)
		fmt.Println("Unikey:", unikey, ".\nQR 码链接: https://music.163.com/login?codekey="+unikey, "\n网易云扫码成功后，后续请使用 cloudlrc netease login check <unikey> 方式以持久化。")
	},
}

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
	rootCmd.AddCommand(neteaseCmd)
	neteaseCmd.AddCommand(neteaseLoginCmd)
	neteaseLoginCmd.AddCommand(neteaseLoginStatusCmd)
	neteaseLoginCmd.AddCommand(neteaseLoginGenCmd)
	neteaseLoginGenCmd.PersistentFlags().BoolVar(&qrfile, "qrfile", false, "生成 qr.png")
	neteaseLoginCmd.AddCommand(neteaseLoginCheckCmd)
}
