package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	VERSION    = "0.0.0"
	GIT_HASH   = ""
	BUILD_TIME = ""
	GO_VER     = ""
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  `显示版本信息`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cloudlrc Version:", VERSION)
		fmt.Println("GIT_HASH:", GIT_HASH)
		fmt.Println("BUILD_TIME:", BUILD_TIME)
		fmt.Println("GO_VER:", GO_VER)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
