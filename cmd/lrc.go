package cmd

import (
	"fmt"

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

func init() {
	rootCmd.AddCommand(lrcCmd)
	lrcCmd.PersistentFlags().BoolVar(&trans_first, "trans_first", false, "翻译在原文前")
}
