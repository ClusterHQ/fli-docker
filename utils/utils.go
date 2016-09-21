package utils

import (
	"os/exec"
	"os"
	"log"
	"regexp"
	"io/ioutil"
	"bytes"

	"gopkg.in/yaml.v2"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
)

var ComposeHelpMessage = `
		-----------------------------------------------------------------------
		docker-compose is not installed, it is needed to use fli-docker\n")
		docker-compose is available at https://docs.docker.com/compose/install/
		-----------------------------------------------------------------------
`

var FliHelpMessage = `
		-------------------------------------------------------
		fli is not installed, it is needed to use fli-docker
		fli is available at https://clusterhq.com
		-------------------------------------------------------
`

func CheckForPath(path string) (result bool, err error) {
	isPath, errPath := exec.LookPath(path)
	// LookPath searches for an executable binary 
	// named file in the directories named by the PATH environment 
	if errPath != nil {
		return false, errPath
	}
	log.Print("Found path: " + isPath + "\n")
	return true, nil
}

func CheckForFile(file string) (result bool, err error) {
	_, errFile := os.Stat(file)
	if errFile != nil {
		return false, errFile
	}
	log.Print("Found file: " + file + "\n")
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

// Represents {compose_volume_name : "/chq/<vol_path>"}
// for volume names in compose to their new path
// after dpcli creates them.
type NewVolume struct {
	Name string
	VolumePath string
}

// Parse a raw yaml file.
func ParseManifest(yamlFile []byte) (*Manifest){
	var manifest Manifest
	err := yaml.Unmarshal(yamlFile, &manifest)
	if err != nil {
		panic(err)
	}
	//log.Print("Manifest: %#v\n", manifest)
	return &manifest
}

// Run the command to sync a volumeset
func syncVolumeset(volumeSetId string) {
	log.Printf("Syncing Volumeset %s", volumeSetId)
	log.Printf("Running sync on volumeset %s", volumeSetId)
	out, err := exec.Command("/opt/clusterhq/bin/dpcli", "sync",  "volumeset", volumeSetId).Output()
	if err != nil {
		log.Print("Could not sync dataset, reason: ", out)
		log.Fatal(err)
	}
	log.Print(out)
}

// Run the command to pull a specific snapshot
func pullSnapshot(snapshotId string){
	log.Printf("Pulling Snapshot %s", snapshotId)
	log.Printf("Running pull for snapshot: %s", snapshotId)
	out, err := exec.Command("/opt/clusterhq/bin/dpcli", "pull", "snapshot", snapshotId).Output()
	if err != nil {
		log.Print("Could not pull dataset, reason: ", out)
		log.Fatal(err)
	}
	log.Print(out)
}

// Wrapper for sync and pull which takes
// a List of type Volume above to pull.
func PullSnapshots(volumes []Volume) {
	for _, volume := range volumes {
		syncVolumeset(volume.VolumeSet)
		// maybe worth traking if we already sync'd a volumset
		// and skipping another sync during the same PullSnapshots call.
		pullSnapshot(volume.Snapshot)
	}
}

// ************************** TODO ****************************
//func authenticateWithFlockerHub(user string, token string, endpoint string) {}
// ************************** TODO ****************************

// Created a volume and returns it.
func createVolumeFromSnapshot(volumeName string, snapshotId string) (vol NewVolume, err error){
	log.Printf("Creating Volume from Snapshot: %s", snapshotId)
	cmd := exec.Command("/opt/clusterhq/bin/dpcli", "create", "volume", "--snapshot", snapshotId)
	combinedOut, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
		log.Print(string(combinedOut))
		r, _ := regexp.Compile("/chq/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
		path := r.FindString(string(combinedOut))
	if path == "" {
			log.Fatal("Could not find volume path")
	 }
		return NewVolume{Name: volumeName, VolumePath: path}, nil
}

func CreateVolumesFromSnapshots(volumes []Volume) (newVols []NewVolume, err error) {
	vols := []NewVolume{}
	for _, volume := range volumes {
		vol, err := createVolumeFromSnapshot(volume.Name, volume.Snapshot)
		if err != nil {
			return nil, err
		}else {
			vols = append(vols, vol)
		}
	}
	return vols, nil
}

// Replace volume names with associated volume paths
// 	Ultimately we should be able to support multiple types of volumes
// https://docs.docker.com/compose/compose-file/#/volumes-volume-driver
// where we can detect if it has a "named" volume, a path, or no "<inside>:"
// and we should modify the file accordingly, for now we only support
// "named volumes" in the form of `-[space]<volume_name>:`
func MapVolumeToCompose(volume string, path string, composeFile string) {
	input, err := ioutil.ReadFile(composeFile)
		if err != nil {
			log.Print("Trouble reading docker-compose file.")
			log.Fatal(err)
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
				log.Print("Error writing docker-compose file.")
				log.Fatal(err)
		 }

}

// Parse the compose file, this will validate
// the compose file and print it.
func ParseCompose(composeFile string) {
	project, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{composeFile},
			ProjectName:  "my-compose", // configurable?
		},
	}, nil)

	if err != nil {
		log.Fatal(err)
	}

	conf, err := project.Config()
	log.Print(conf)
}

// A function to copy a file and 
// label it as fli did it.
func MakeCopy(composeFile string) {
	srcFolder := composeFile
	destFolder := composeFile + "-fli-copy"
	cpCmd := exec.Command("cp", "-rf", srcFolder, destFolder)
	err := cpCmd.Run()
	if err !=nil {
		log.Fatal(err)
	}
}

//
