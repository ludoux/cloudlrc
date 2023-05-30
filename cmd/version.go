package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	VERSION = "0.0.0"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  `显示版本信息`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cloudlrc %s\n", VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
