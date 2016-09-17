package utils

import (
	"os/exec"
	"log"
)

func CheckForPath(path string) (result bool, err error) {
	isPath, err := exec.LookPath(path)
	if err != nil {
		return false, err
	}
	log.Println("Found path: " + isPath + "\n")
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
