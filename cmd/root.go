package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cloudlrc",
	Short: "网易云音乐歌词下载工具",
	Long: `CloudLRC 是一个用于从网易云音乐批量下载歌词的命令行工具。
支持单曲、歌单、专辑的批量下载，并提供歌词翻译时间调整和自定义分隔符功能。

使用示例:
  cloudlrc lrc music 123456
  cloudlrc lrc playlist 123456789 --trans-delay 500
  cloudlrc lrc album 123456789 --split " | " --trans_first`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.test.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
