package cmd

import (
	"fmt"
	"os"

	"github.com/ludoux/cloudlrc/netease"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var trans_first bool

var lrcCmd = &cobra.Command{
	Use:   "lrc",
	Short: "下载歌词",
	Long:  `下载歌词`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("lrc called")
		fmt.Println(cast.ToString(trans_first))
	},
}

var lrcMusicCmd = &cobra.Command{
	Use:   "music <id...>",
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
var lrcAlbumCmd = &cobra.Command{
	Use:   "album",
	Short: "下载专辑的歌词",
	Long:  `下载专辑的歌词`,
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
		netease.DownloadAlbumLrc(id)
	},
}

func init() {
	rootCmd.AddCommand(lrcCmd)
	lrcCmd.PersistentFlags().BoolVar(&trans_first, "trans_first", false, "翻译在原文前")

	lrcCmd.AddCommand(lrcMusicCmd)
	lrcCmd.AddCommand(lrcPlaylistCmd)
	lrcCmd.AddCommand(lrcAlbumCmd)
}
