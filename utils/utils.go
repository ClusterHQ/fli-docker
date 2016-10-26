/*
 *Copyright ClusterHQ Inc.  See LICENSE file for details.
 *
 */

package utils

import (
	"os/exec"
	"os"
	"io/ioutil"
	"bytes"
	"strings"

	"gopkg.in/yaml.v2"
	"golang.org/x/net/context"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"

	"github.com/ClusterHQ/fli-docker/types"
	"github.com/ClusterHQ/fli-docker/logger"
)

var ComposeHelpMessage = `
-----------------------------------------------------------------------
docker-compose is not installed, it is needed to use fli-docker\n")
docker-compose is available at https://docs.docker.com/compose/install/
-----------------------------------------------------------------------`

var FliHelpMessage = `
-------------------------------------------------------
fli is not installed, it is needed to use fli-docker
fli is available at https://clusterhq.com
Using the fli contianer? Make sure it used docker tag 'clusterhq/fli'
-------------------------------------------------------`

var FliDockerVersion = `Version: v0.0.1-dev`

var FliDockerHelp = `
Usage:
  fli-docker version  [options]  (Get current tool version)
  fli-docker run      [options]  (Run with a manifest to pull and use snapshots for the compose app)
  fli-docker snapshot [options]  (Snapshot existing FlockerHub volumes used by the compose app)
  fli-docker stop     [options]  (Just like running a docker-compose stop)
  fli-docker destroy  [options]  (Just like running a docker-compose rm -f)
  fli-docker --help   (Get this help message)

  For help on a specific command, use: $ fli-docker <subcommand> --help`

var FliDockerCmd = "docker run --rm --privileged -v /chq/:/chq/:shared -v /root:/root -v /lib/modules:/lib/modules clusterhq/fli "
var FliBinaryCmd = "fli "

func CheckForPath(path string) (result bool, err error) {
	isPath, errPath := exec.LookPath(path)
	// LookPath searches for an executable binary 
	// named file in the directories named by the PATH environment 
	if errPath != nil {
		return false, errPath
	}
	logger.Info.Println("Found path: " + isPath)
	return true, nil
}

func CheckForFile(file string) (result bool, err error) {
	_, errFile := os.Stat(file)
	if errFile != nil {
		logger.Info.Println(errFile)
		return false, errFile
	}
	logger.Info.Println("Found file: " + file)
	return true, nil
}

func CheckForCmd(cmd string) (result bool, err error) {
	_, errCmd := exec.Command("sh", "-c", cmd).Output()
	if errCmd != nil {
		logger.Info.Println(errCmd)
		return false, errCmd
	}
	logger.Info.Println("Found Command: " + cmd)
	return true, nil
}

// A function to copy a file and 
// label it as fli did it.
func MakeCopy(composeFile string) {
	srcFolder := composeFile
	destFolder := composeFile + "-fli.copy"
	exists, err := CheckForFile(destFolder)
	if err != nil {
		logger.Info.Println("No existing compose file copy.")
		logger.Info.Println(err)
	}
	if exists {
		logger.Info.Println("Copy already exists, not copying")
	}else {
		cpCmd := exec.Command("cp", "-rf", srcFolder, destFolder)
		err := cpCmd.Run()
		if err != nil {
			logger.Error.Fatal(err)
		}
	}
}

func CheckForCopy(composeFile string) {
	// If we already copied the original, we
	// want to make sure we copy back the original
	// before modifying it agian otherwise
	// correct volume names may not exist.
	srcFolder := composeFile
	destFolder := composeFile + "-fli.copy"
	exists, err := CheckForFile(destFolder)
	if err != nil {
		logger.Info.Println("No existing compose file copy.")
		logger.Info.Println(err)
	}
	if exists {
		logger.Info.Println("Refreshing compose app from copy")
		cpCopyCmd := exec.Command("cp", "-rf", destFolder, srcFolder)
		err := cpCopyCmd.Run()
		if err != nil {
			logger.Error.Fatal(err)
		}
	}
}

// Parse a raw yaml file.
func ParseManifest(yamlFile []byte) (*types.Manifest){
	var manifest types.Manifest
	err := yaml.Unmarshal(yamlFile, &manifest)
	if err != nil {
		logger.Error.Fatal(err)
	}
	return &manifest
}

// Replace volume names with associated volume paths
// Ultimately we should be able to support multiple types of volumes
// https://docs.docker.com/compose/compose-file/#/volumes-volume-driver
// where we can detect if it has a "named" volume, a path, or no "<inside>:"
// and we should modify the file accordingly, for now we only support
// "named volumes" in the form of `-[space]<volume_name>:`
func MapVolumeToCompose(volume string, path string, composeFile string) {
	input, err := ioutil.ReadFile(composeFile)
		if err != nil {
			logger.Error.Print("Trouble reading docker-compose file.")
			logger.Error.Fatal(err)
		}
	prefixQuote := "- '"
	prefixNoQuote := "- "
	postfix := ":"

	//replace the "- named_volume:" name with the Flucker Hub path. (without single quote)
	output := bytes.Replace(input, []byte(prefixNoQuote + volume + postfix),
		[]byte(prefixNoQuote + path + postfix), -1)

	//replace the "- 'named_volume:" name with the Flucker Hub path. (with single quote)
	finalOutput := bytes.Replace(output, []byte(prefixQuote + volume + postfix),
		[]byte(prefixQuote + path + postfix), -1)

	//re-write
	if err = ioutil.WriteFile(composeFile, finalOutput, 0644); err != nil {
			logger.Error.Print("Error writing docker-compose file.")
			logger.Error.Fatal(err)
		 }
}

// Parse the compose file. 
// This will validate the compose file and print it.
func ParseCompose(composeFile string) {
	project, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{composeFile},
			ProjectName:  "fli-compose",
		},
	}, nil)

	if err != nil {
		logger.Error.Fatal(err)
	}

	conf, err := project.Config()
	logger.Info.Print(conf)
}

// Run the compose file with options
func RunCompose(composeFile string, projectName string) {
	project, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{composeFile},
			ProjectName:  projectName,
		},
	}, nil)

	if err != nil {
		logger.Error.Fatal(err)
	}

	err = project.Up(context.Background(), options.Up{})

    if err != nil {
        logger.Error.Fatal(err)
    }
}

// Run the compose file with options
func StopCompose(composeFile string, projectName string) {
	project, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{composeFile},
			ProjectName:  projectName,
		},
	}, nil)

	if err != nil {
		logger.Error.Fatal(err)
	}

	err = project.Stop(context.Background(), 60)

    if err != nil {
        logger.Error.Fatal(err)
    }
}

// Run the compose file with options
func DestroyCompose(composeFile string, projectName string) {
	project, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{composeFile},
			ProjectName:  projectName,
		},
	}, nil)

	if err != nil {
		logger.Error.Fatal(err)
	}

	err = project.Delete(context.Background(), options.Delete{true,true})

    if err != nil {
        logger.Error.Fatal(err)
    }
}

//Generate UUIDs
func GenUUID() (uuid string, err error){
	out, err := exec.Command("uuidgen").Output()
    if err != nil {
        logger.Error.Println(err)
        return " ", err
    }
    return strings.ToLower(string(out)), nil
}
