package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"

	streams "./io"
)

// Run an acbuild command.
//
// Accepts a list or arguments to
// be appended to the acbuild command
// ie. acbuild([arg1, arg2]) executes
// $ acbuild arg1 arg2
func acbuild(args []string) {
	fmt.Println("acbuild", strings.Join(args, " "))
	execute(append([]string{"acbuild"}, args...), streams.NewStdIO())
}

// Hash the contents of a file at the
// given path
func hash(path string) string {
	h := md5.New()
	reader := open(path)
	defer reader.Close()

	if _, err := io.Copy(h, reader); err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(h.Sum(nil))
}

// Open and return a file
func open(filePath string) *os.File {
	manifest, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Unable to open '" + filePath + "', aborting")
		panic(err)
	}

	return manifest
}

// Returns the path to the home
// directory for the user currently
// running the program
func HomeDir() string {
	u, _ := user.Current()

	return u.HomeDir
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

	manifest := open(env + ".rkd")
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

func getAciPath(env string) string {
	return filepath.Join(HomeDir(), ".rkd", hash(env + ".rkd")+".aci")
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
	execute(strings.Split(command, " "), streams.NewStdIO())
}

// Returns the configuration of dev
// mounts formatted for a rkt command.
//
// What we do is parse the 'dev.rkd'
// lin by line and if we find a 'mount'
// instruction we format it for the
// rkt CLI.
func getMountConfig() string {
	config := ""
	cwd, _ := os.Getwd()
	file, err := os.Open("dev.rkd")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		cmd := scanner.Text()
		if cmd[0:6] == "mount " {
			parts := strings.Split(cmd, " ")
			config += "--volume " + parts[2] + ",kind=host,source=" + cwd + "/" + parts[3]
		}
	}

	return config
}

// Execute a command
func execute(args []string, io streams.IO) {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdin = io.Stdin
	cmd.Stdout = io.Stdout
	cmd.Stderr = io.Stderr

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

	// We have a directory to store ACIs
	err := os.Mkdir(filepath.Join(HomeDir(), ".rkd"), 0755)

	if err != nil && !os.IsExist(err) {
		panic(err)
	}
}

// Start the fun!
func main() {
	checkRequirements()
	buildAci("prod")
	buildAci("dev")
	runAci()
}
