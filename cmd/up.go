package cmd

import (
	"fmt"
	"strings"

	streams "github.com/odino/rkd/io"
	"github.com/odino/rkd/utils"
	"github.com/spf13/cobra"
)

var net string = "default"

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
	if net == "" {
		net = "host"
	}

	command := fmt.Sprintf("rkt --insecure-options=image --net=%s --interactive=true %s run %s", net, getMountConfig(), getAciPath("dev"))
	fmt.Println(command)
	utils.Execute(strings.Split(command, " "), streams.NewStdIO())
}

func init() {
	RootCmd.AddCommand(upCmd)

	upCmd.Flags().StringVarP(&net, "net", "n", "", "configure the pod's networking. Optionally, pass a list of user-configured networks to load and set arguments to pass to each network, respectively. Syntax: --net[=n[:args], ...]")
}
