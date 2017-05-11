package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	streams "github.com/odino/rkd/io"
	"github.com/odino/rkd/utils"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the containers",
	Long:  `Builds prod and dev ACIs for the current software`,
	Run: func(cmd *cobra.Command, args []string) {
		checkRequirements()
		buildAci("prod")
		buildAci("dev")
	},
}

// Run an acbuild command.
//
// Accepts a list or arguments to
// be appended to the acbuild command
// ie. acbuild([arg1, arg2]) executes
// $ acbuild arg1 arg2
func acbuild(args []string) {
	fmt.Println("acbuild", strings.Join(args, " "))
	utils.Execute(append([]string{"acbuild"}, args...), streams.NewStdIO())
}

// Builds an ACI.
//
// ACIs can be either "prod"
// (what you probably want to run
// in production) or "dev" (which
// probably have additional configurations
// such as a different exec command
// or a mount volume for your code).
func buildAci(env string) {
	// Let's make sure we're able to intercept
	// signals so that we shut the app down
	// gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			fmt.Println("Interrupted")
		}
	}()

	// Let's make sure whatever happens while
	// we're building the ACI we execute an
	// "acbuild end" so that the user can
	// re-trigger a build without getting the
	// "build already in progress" error...
	defer func() {
		if err := recover(); err != nil {
			acbuild([]string{"end"})
			os.Exit(1)
		}
	}()

	aciPath := getAciPath(env)

	if _, err := os.Stat(aciPath); os.IsNotExist(err) {
		fmt.Println("Building " + aciPath)
	} else {
		fmt.Println(aciPath + " already built")
		return
	}

	manifest := utils.Open(env + ".rkd")
	defer manifest.Close()
	scanner := bufio.NewScanner(manifest)

	if env == "prod" {
		acbuild([]string{"begin"})
	} else {
		acbuild([]string{"begin", "./prod.aci"})
	}

	for scanner.Scan() {
		acbuild(strings.Split(scanner.Text(), " "))
	}

	acbuild([]string{"write", aciPath})
	acbuild([]string{"end"})

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

// Make sure everything needed by rkd is
// available on the system.
func checkRequirements() {
	// acbuild is installed and can run
	utils.Execute([]string{"acbuild"}, streams.NewDevNullIO())

	// We have a directory to store ACIs
	err := os.Mkdir(filepath.Join(utils.HomeDir(), ".rkd"), 0755)

	if err != nil && !os.IsExist(err) {
		panic(err)
	}
}

func init() {
	RootCmd.AddCommand(buildCmd)
}
