package main

import (
	"fmt"
	"os"

	"./cmd"
)

// Start the fun!
func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(99)
	}
	return
}
