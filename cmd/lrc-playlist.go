package cmd

import (
	"fmt"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// neteaseCmd represents the netease command
var lrcPlaylistCmd = &cobra.Command{
	Use:   "playlist",
	Short: "下载歌单的歌词",
	Long:  `下载歌单的歌词`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cast.ToString(trans_first))
	},
}

func init() {
	lrcCmd.AddCommand(lrcPlaylistCmd)
}
