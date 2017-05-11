package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of rkd",
	Long:  `All software has versions. This is rkd's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vTooAlphaToBother")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
