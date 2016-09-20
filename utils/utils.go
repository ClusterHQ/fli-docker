package utils

import (
	"os/exec"
	"os"
	"log"
	"fmt"
	"regexp"

	"gopkg.in/yaml.v2"

    "github.com/docker/libcompose/docker"
    "github.com/docker/libcompose/docker/ctx"
    "github.com/docker/libcompose/project"
)

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

// Represents {compose_volume_name : "/chq/<vol_path>"}
// for volume names in compose to their new path
// after dpcli creates them.
type NewVolume struct {
	Name string
	VolumePath string
}

// Parse a raw yaml file.
// TODO this should really return the manifest back
// to the CLI so it can use it.
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
	syncCmdStr := fmt.Sprintf("/opt/clusterhq/bin/dpcli sync volumeset %s", volumeSetId)
	syncCmd := exec.Command(syncCmdStr)
	err := syncCmd.Run()
	if err != nil {
        log.Fatal(err)
    }
}

// Run the command to pull a specific snapshot
func pullSnapshot(snapshotId string){
	log.Printf("Pulling Snapshot %s", snapshotId)
	pullCmdStr := fmt.Sprintf("/opt/clusterhq/bin/dpcli pull snapshot %s", snapshotId)
	pullCmd := exec.Command(pullCmdStr)
	err := pullCmd.Run()
	if err != nil {
        log.Fatal(err)
    }
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

//func authenticateWithFlockerHub(user string, token string, endpoint string) {}

// Created a volume and returns it.
func createVolumeFromSnapshot(volumeName string, snapshotId string) (vol NewVolume, err error){
	log.Printf("Creating Volume from %s", snapshotId)
	createCmdStr := fmt.Sprintf("/opt/clusterhq/bin/dpcli create volume --snapshot %s", snapshotId)
	out, err := exec.Command(createCmdStr).Output()
	if err != nil {
        log.Fatal(err)
    }
    log.Print(out)
    output := fmt.Sprintf("%s", out)
    r, _ := regexp.Compile("/chq/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
    path := r.FindString(output)
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
// 		This is crappy and platform specific, we could use a
// 		native yaml reader/writer to do this more properly.
func MapVolumeToCompose(volume string, path string, composeFile string) {
    sedCmdStr := fmt.Sprintf("/usr/bin/sed -i 's@%s:@%s:@' %s", volume, path, composeFile)
	sedCmd := exec.Command(sedCmdStr)
	err := sedCmd.Run()
	if err != nil {
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
