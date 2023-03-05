package cmd

import (
	"fmt"
	"os"

	"github.com/ludoux/cloudlrc/netease"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// neteaseCmd represents the netease command
var lrcMusicCmd = &cobra.Command{
	Use:   "music [id]",
	Short: "下载单曲的歌词",
	Long:  `下载单曲的歌词`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 || len(args) > 100 {
			fmt.Println("Error: 单曲数目应在 1~100 之间，输入数目:", len(args))
			os.Exit(1)
		}
		for _, v := range args {
			id, err := cast.ToInt64E(v)
			if err != nil {
				fmt.Println("Error: 输入的ID无法转换为数组:", v)
				os.Exit(1)
			}
			netease.DownloadSingleMusicLrc(id)
		}
	},
}

func init() {
	lrcCmd.AddCommand(lrcMusicCmd)
}
