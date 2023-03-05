package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// neteaseCmd represents the netease command
var neteaseCmd = &cobra.Command{
	Use:   "netease",
	Short: "接口: 网易云音乐",
	Long:  `接口: 网易云音乐`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("netease called")
	},
}

func init() {
	rootCmd.AddCommand(neteaseCmd)
}
