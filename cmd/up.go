package cmd

import (
	"fmt"
	"strings"

	streams "../io"
	utils "../utils"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run the container",
	Long:  `Brings up the container and start the fun`,
	Run: func(cmd *cobra.Command, args []string) {
		buildCmd.Run(cmd, args)
		runAci()
	},
}

// Run the dev.aci.
//
// This function runs dev.aci
// through rkt, with some
// default configurations (ie. --interactive)
// that make it easy for dev
// environments.
func runAci() {
	command := "rkt --insecure-options=image --net=host run --interactive " + getMountConfig() + " " + getAciPath("dev")
	fmt.Println(command)
	utils.Execute(strings.Split(command, " "), streams.NewStdIO())
}

func init() {
	RootCmd.AddCommand(upCmd)
}
