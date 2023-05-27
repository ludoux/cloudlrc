package cmd

import (
	"fmt"
	"os"

	"github.com/ludoux/cloudlrc/netease"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// neteaseCmd represents the netease command
var lrcPlaylistCmd = &cobra.Command{
	Use:   "playlist",
	Short: "下载歌单的歌词",
	Long:  `下载歌单的歌词`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Error: 专辑数目应为 1，输入数目:", len(args))
			os.Exit(1)
		}
		id, err := cast.ToInt64E(args[0])
		if err != nil {
			fmt.Println("Error: 输入的ID无法转换为数组:", args[0])
			os.Exit(1)
		}
		netease.DownloadPlaylistLrc(id)
	},
}

func init() {
	lrcCmd.AddCommand(lrcPlaylistCmd)
}
