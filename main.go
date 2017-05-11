package main

import (
	"fmt"
	"os"

	"github.com/odino/rkd/cmd"
)

// Start the fun!
func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(99)
	}
	return
}
