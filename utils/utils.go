package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"

	streams "github.com/odino/rkd/io"
)

// Hash the contents of a file at the
// given path
func HashFile(path string) string {
	h := md5.New()
	reader := Open(path)
	defer reader.Close()

	if _, err := io.Copy(h, reader); err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(h.Sum(nil))
}

// Hashes a string
func Hash(s string) string {
	h := md5.New()
	io.WriteString(h, s)

	return hex.EncodeToString(h.Sum(nil))
}

// Open and return a file
func Open(filePath string) *os.File {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Unable to open '" + filePath + "', aborting")
		panic(err)
	}

	return file
}

// Execute a command
func Execute(args []string, io streams.IO) {
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

// Returns the path to the home
// directory for the user currently
// running the program
func HomeDir() string {
	u, _ := user.Current()

	return u.HomeDir
}
