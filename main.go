package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Run an acbuild command.
//
// Accepts a list or arguments to
// be appended to the acbuild command
// ie. acbuild arg1 arg2
func acbuild(args []string) {
	cmd := exec.Command("acbuild", args...)
	fmt.Println("acbuild", strings.Join(args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()

	if err != nil {
		panic(err)
	}
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
		log.Fatal(err)
	}
}

// Run the dev.aci.
//
// This function runs dev.aci
// through rkt, with some
// default configurations (ie. --interactive)
// that make it easy for dev
// environments.
func run() {
	args := make([]string, 0)
	args = append(args, "--insecure-options=image")
	args = append(args, "run")
	args = append(args, "--interactive")
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
			args = append(args, "--volume")
			args = append(args, parts[2]+",kind=host,source="+cwd+"/"+parts[3])
		}
	}

	args = append(args, "dev.aci")
	cmd := exec.Command("rkt", args...)
	fmt.Println("rkt", strings.Join(args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()

	if err != nil {
		panic(err)
	}
}

// Start the fun!
func main() {
	buildAci("prod")
	buildAci("dev")
	run()
}
