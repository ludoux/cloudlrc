package cmd

import (
	"fmt"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// neteaseCmd represents the netease command
var lrcAlbumCmd = &cobra.Command{
	Use:   "album",
	Short: "下载专辑的歌词",
	Long:  `下载专辑的歌词`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cast.ToString(trans_first))
	},
}

func init() {
	lrcCmd.AddCommand(lrcAlbumCmd)
}
