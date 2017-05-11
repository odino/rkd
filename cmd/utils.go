package cmd

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	utils "../utils"
)

func getAciPath(env string) string {
	return filepath.Join(utils.HomeDir(), ".rkd", utils.Hash(env+".rkd")+".aci")
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
