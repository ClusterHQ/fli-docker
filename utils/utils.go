package utils

import (
	"os/exec"
	"log"
)


func CheckForTool(cliPath string) (result bool, err error) {
	path, err := exec.LookPath(cliPath)
	if err != nil {
		return false, err
	}
	log.Println("Found path: " + path + "\n")
	return true, nil
}
