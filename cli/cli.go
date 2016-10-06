package cli

import (
	"os/exec"
	"regexp"
	"fmt"
	"strings"

	"github.com/ClusterHQ/fli-docker/types"
	"github.com/ClusterHQ/fli-docker/logger"
	"github.com/ClusterHQ/fli-docker/utils"
)

/*
	Bindings to the FlockerHub CLI
*/

func SetFlockerHubEndpoint(endpoint string) {
	logger.Info.Println("Setting FlockerHub Endpoint: ", endpoint)
	out, err := exec.Command("/opt/clusterhq/bin/dpcli", "set", "volumehub", endpoint).Output()
	if err != nil {
		logger.Error.Println("Could not set endpoint, reason: ", string(out))
		logger.Error.Fatal(err)
	}
	logger.Info.Println(string(out))
}

func GetFlockerHubEndpoint() (flockerhubEndpoint string, err error) {
	logger.Info.Println("Getting FlockerHub Endpoint")
	out, err := exec.Command("/opt/clusterhq/bin/dpcli", "get", "volumehub").Output()
	if err != nil {
		logger.Error.Println("Could not get endpoint, reason: ", string(out))
		return "", err
	}
	logger.Info.Println(string(out))
	//TODO Parse, output to get specific volumehub string
	return string(out), nil
}

func SetFlockerHubTokenFile(tokenFile string) {
	logger.Info.Println("Setting FlockerHub Tokenfile: ", tokenFile)
	out, err := exec.Command("/opt/clusterhq/bin/dpcli", "set", "tokenfile", tokenFile).Output()
	if err != nil {
		logger.Error.Println("Could not set tokenfile, reason: ", string(out))
		logger.Error.Fatal(err)
	}
	logger.Info.Println(string(out))
}

func GetFlockerHubTokenFile() (flockerHubTokenFile string, err error) {
	logger.Info.Println("Getting FlockerHub Tokenfile")
	out, err := exec.Command("/opt/clusterhq/bin/dpcli", "get", "tokenfile").Output()
	if err != nil {
		logger.Error.Println("Could not get tokenfile, reason: ", string(out))
		return "", err
	}
	logger.Info.Println(string(out))
	//TODO Parse, output to get specific tokenfile string
	return string(out), nil
}

// Run the command to sync a volumeset
func syncVolumeset(volumeSetId string) {
	logger.Info.Println("Syncing Volumeset: ", volumeSetId)
	out, err := exec.Command("/opt/clusterhq/bin/dpcli", "sync",  "volumeset", volumeSetId).CombinedOutput()
	if err != nil {
		logger.Error.Println("Could not sync dataset, reason: ", string(out))
		logger.Error.Fatal(err)
	// sometimes errors dont get sent to STDERR?
	// update from abhishek 10/4/16 that this will be fixed, so this will
	// not be needed after we update fli-docker to use fli grammer/later cli.
	}else if strings.Contains(strings.ToLower(string(out)), "error"){
		logger.Error.Println("Could not sync dataset, reason: ", string(out))
		logger.Error.Fatal(err)
	}
	logger.Info.Println(string(out))
}

// Run the command to pull a specific snapshot
func pullSnapshot(snapshotId string){
	logger.Info.Println("Pulling Snapshot: ", snapshotId)
	out, err := exec.Command("/opt/clusterhq/bin/dpcli", "pull", "snapshot", snapshotId).CombinedOutput()
	if err != nil {
		logger.Error.Println("Could not pull dataset, reason: ", string(out))
		logger.Error.Fatal(err)
	// sometimes errors dont get sent to STDERR?
	// update from abhishek 10/4/16 that this will be fixed, so this will
	// not be needed after we update fli-docker to use fli grammer/later cli.
	}else if strings.Contains(strings.ToLower(string(out)), "error"){
		logger.Error.Println("Could not pull dataset, reason: ", string(out))
		logger.Error.Fatal(err)
	}
	logger.Info.Println(string(out))
}

// Wrapper for sync and pull which takes
// a list of type Volume
func PullSnapshots(volumes []types.Volume) {
	for _, volume := range volumes {
		syncVolumeset(volume.VolumeSet)
		pullSnapshot(volume.Snapshot)
	}
}

// Created a volume and returns it.
func createVolumeFromSnapshot(volumeName string, snapshotId string) (vol types.NewVolume, err error){
	logger.Info.Println("Creating Volume from Snapshot: ", snapshotId)
	var attrString = fmt.Sprintf("created_by=fli-docker,from_snap=%s", snapshotId)
	uuid, err := utils.GenUUID()
	if err != nil {
		logger.Error.Fatal(err)
	}

	var volName = fmt.Sprintf("fli-%s", uuid)
	cmd := exec.Command("/opt/clusterhq/bin/dpcli", "create", "volume", 
		"--snapshot", snapshotId, "-a", attrString, volName)
	createOut, err := cmd.Output()
	if err != nil {
		logger.Error.Fatal(err)
	}

	logger.Info.Println(string(createOut))
	// This is where dpcli volume path <volume-name> would be handy.
	r, _ := regexp.Compile("/chq/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
	path := r.FindString(string(createOut))
	if path == "" {
			logger.Error.Fatal("Could not find volume path")
	 }
	return types.NewVolume{Name: volumeName, VolumePath: path}, nil
}

func CreateVolumesFromSnapshots(volumes []types.Volume) (newVols []types.NewVolume, err error) {
	vols := []types.NewVolume{}
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