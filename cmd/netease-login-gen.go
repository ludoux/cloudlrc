package cmd

import (
	"fmt"

	"github.com/ludoux/cloudlrc/netease"
	"github.com/spf13/cobra"
)

var qrfile bool

var neteaseLoginGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		unikey := netease.LoginGen(qrfile)
		fmt.Println("Unikey:", unikey, ".\nQR 码链接: https://music.163.com/login?codekey="+unikey, "\n网易云扫码成功后，后续请使用 cloudlrc netease login check <unikey> 方式以持久化。")
	},
}

func init() {
	neteaseLoginCmd.AddCommand(neteaseLoginGenCmd)
	neteaseLoginGenCmd.PersistentFlags().BoolVar(&qrfile, "qrfile", false, "生成 qr.png")
}
