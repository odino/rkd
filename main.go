package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	streams "./io"
)

// Run an acbuild command.
//
// Accepts a list or arguments to
// be appended to the acbuild command
// ie. acbuild arg1 arg2
func acbuild(args []string) {
	fmt.Println("acbuild", strings.Join(args, " "))
	execute(append([]string{"acbuild"}, args...), streams.NewStdIO())
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
	if _, err := os.Stat("./" + env + ".aci"); os.IsNotExist(err) {
		fmt.Println("Building " + env + ".aci")
	} else {
		fmt.Println(env + ".aci already built")
		return
	}

	file, err := os.Open("./" + env + ".rkd")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if env == "prod" {
		acbuild([]string{"begin"})
	} else {
		acbuild([]string{"begin", "./prod.aci"})
	}

	for scanner.Scan() {
		acbuild(strings.Split(scanner.Text(), " "))
	}

	acbuild([]string{"write", env + ".aci"})
	acbuild([]string{"end"})

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

// Run the dev.aci.
//
// This function runs dev.aci
// through rkt, with some
// default configurations (ie. --interactive)
// that make it easy for dev
// environments.
func runAci() {
	command := "rkt --insecure-options=image run --interactive"
	cwd, _ := os.Getwd()

	file, err := os.Open("./dev.rkd")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		cmd := scanner.Text()
		if cmd[0:6] == "mount " {
			parts := strings.Split(cmd, " ")
			command += " --volume " + parts[2] + ",kind=host,source=" + cwd + "/" + parts[3]
		}
	}

	command += " dev.aci"
	fmt.Println(command)
	execute(strings.Split(command, " "), streams.NewStdIO())
}

// Execute a command
//
// If no IO if specified, it will
// default to the invoking process'
// stdout and stderr.
func execute(args []string, io streams.IO) {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = io.Stdout

	if cmd.Stdout == nil {
		cmd.Stdout = os.Stdout
	}

	cmd.Stderr = io.Stderr

	if cmd.Stderr == nil {
		cmd.Stderr = os.Stderr
	}

	err := cmd.Start()

	if err != nil {
		panic(err)
	}

	err = cmd.Wait()

	if err != nil {
		panic(err)
	}
}

// Make sure everything needed by rkd is
// available on the system.
func checkRequirements() {
	// acbuild is installed and can run
	execute([]string{"acbuild"}, streams.NewDevNullIO())
}

// Start the fun!
func main() {
	checkRequirements()
	buildAci("prod")
	buildAci("dev")
	runAci()
}
