package cmd

import (
	"fmt"
	"os"

	"github.com/ludoux/cloudlrc/netease"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var (
	trans_first bool
	trans_delay int64
	split       string
)

var lrcCmd = &cobra.Command{
	Use:   "lrc",
	Short: "下载歌词",
	Long:  `下载歌词，支持单曲、歌单、专辑批量下载`,
	Run: func(cmd *cobra.Command, args []string) {
		//
	},
}

var lrcMusicCmd = &cobra.Command{
	Use:   "music <id...>",
	Short: "下载单曲的歌词",
	Long:  `下载单曲的歌词，支持批量下载多个单曲`,
	Example: `  cloudlrc lrc music 123456 789012
  cloudlrc lrc music 123456 --trans-delay -500
  cloudlrc lrc music 123456 --split " | " --trans_first`,
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
			netease.DownloadSingleMusicLrcWithConfig(id, trans_first, trans_delay, split)
		}
	},
}
var lrcPlaylistCmd = &cobra.Command{
	Use:   "playlist <id>",
	Short: "下载歌单的歌词",
	Long:  `下载歌单的歌词，批量下载歌单内所有歌曲的歌词`,
	Example: `  cloudlrc lrc playlist 123456789
  cloudlrc lrc playlist 123456789 --trans-delay 500
  cloudlrc lrc playlist 123456789 --split " / " --trans_first`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Error: 歌单ID数目应为 1，输入数目:", len(args))
			os.Exit(1)
		}
		id, err := cast.ToInt64E(args[0])
		if err != nil {
			fmt.Println("Error: 输入的ID无法转换为数组:", args[0])
			os.Exit(1)
		}
		netease.DownloadPlaylistLrcWithConfig(id, trans_first, trans_delay, split)
	},
}
var lrcAlbumCmd = &cobra.Command{
	Use:   "album <id>",
	Short: "下载专辑的歌词",
	Long:  `下载专辑的歌词，批量下载专辑内所有歌曲的歌词`,
	Example: `  cloudlrc lrc album 123456789
  cloudlrc lrc album 123456789 --trans-delay -300
  cloudlrc lrc album 123456789 --split "\n" --trans_first`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Error: 专辑ID数目应为 1，输入数目:", len(args))
			os.Exit(1)
		}
		id, err := cast.ToInt64E(args[0])
		if err != nil {
			fmt.Println("Error: 输入的ID无法转换为数组:", args[0])
			os.Exit(1)
		}
		netease.DownloadAlbumLrcWithConfig(id, trans_first, trans_delay, split)
	},
}

func init() {
	rootCmd.AddCommand(lrcCmd)
	lrcCmd.PersistentFlags().BoolVar(&trans_first, "trans_first", false, "翻译在原文前（优先级互换）")
	lrcCmd.PersistentFlags().Int64Var(&trans_delay, "trans-delay", 0, "翻译时间偏移量（毫秒），正数为延后，负数为提前")
	lrcCmd.PersistentFlags().StringVar(&split, "split", "", "自定义分隔符，用于合并原文和翻译（覆盖配置文件设置）")

	lrcCmd.AddCommand(lrcMusicCmd)
	lrcCmd.AddCommand(lrcPlaylistCmd)
	lrcCmd.AddCommand(lrcAlbumCmd)
}
