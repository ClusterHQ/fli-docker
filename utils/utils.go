package utils

import (
	"os/exec"
	"log"
	"fmt"
	"gopkg.in/yaml.v2"
)

func CheckForPath(path string) (result bool, err error) {
	isPath, errPath := exec.LookPath(path)
	if errPath != nil {
		return false, errPath
	}
	log.Print("Found path: " + isPath + "\n")
	return true, nil
}

func CheckForCmd(cmd string) (result bool, err error) {
	_, errCmd := exec.Command("sh", "-c", cmd).Output()
	if errCmd != nil {
		return false, errCmd
	}
	log.Print("Found Command: " + cmd + "\n")
	return true, nil
}

/* Place holders and descriptions for what needs to be
   added in order for the CLI to process the `fli`
   manifest and docker-compose file.
*/

// Need to create a struct for the Compose YAML file
// Use `libcompose` - https://github.com/docker/libcompose

// Need to create a struct for the Manifest file
// Need to startt from scratch - https://github.com/go-yaml/yaml
type Manifest struct {
	DockerApp string      `yaml:"docker_app"`
	Hub FlockerHub        `yaml:"flocker_hub"`
	Volumes []Volume      `yaml:"volumes"`
}

type FlockerHub struct { 
		Endpoint string   `yaml:"endpoint"`
		AuthToken string  `yaml:"auth_token"`
}

type Volume struct {
	Name string      `yaml:"name"`
	Snapshot string  `yaml:"snapshot"`
	VolumeSet string `yaml:"volumeset"`
}

// Parse a raw yaml file.
// TODO this should really return the manifest back
// to the CLI so it can use it.
func ParseManifest(yamlFile []byte) {
	var manifest Manifest
	err := yaml.Unmarshal(yamlFile, &manifest)
	if err != nil {
        panic(err)
    }
	fmt.Printf("Manifest: %#v\n", manifest)
}


//func authenticateWithFlockerHub(user string, token string, endpoint string) {}

//func syncVolumeset() {volumeSetId string}

//func pullSnapshot() {snapshotId string}

//func createVolumeFromSnapshot(snapshotId string) (path string) {}

//func mapVolumeToCompose(composeFile file) {}
/*  read the file
	map the YAML to a struct
	(in manifest)
	parse the volumeset-id and snapshotid
	parse the volume names
	(in compose)
	replace volume names with associated volume paths
*/

//
