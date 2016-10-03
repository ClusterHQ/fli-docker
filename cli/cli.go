package cli

import (
	"os/exec"
	"regexp"
	"fmt"

	"github.com/ClusterHQ/fli-docker/types"
	"github.com/ClusterHQ/fli-docker/logger"
)

/*
	Bindings to the FlockerHub CLI
*/

func SetFlockerHubEndpoint(endpoint string) {
	logger.Info.Println("Setting FlockerHub Endpoint %s", endpoint)
	out, err := exec.Command("/opt/clusterhq/bin/dpcli", "set", "volumehub", endpoint).Output()
	if err != nil {
		logger.Error.Println("Could not set endpoint, reason: ", out)
		logger.Error.Fatal(err)
	}
	logger.Info.Println(out)
}

func GetFlockerHubEndpoint() (flockerhubEndpoint string, err error) {
	logger.Info.Println("Getting FlockerHub Endpoint")
	out, err := exec.Command("/opt/clusterhq/bin/dpcli", "get", "volumehub").Output()
	if err != nil {
		logger.Error.Println("Could not get endpoint, reason: ", out)
		return "", err
	}
	logger.Info.Println(out)
	//TODO Parse, "out" to get specific volumehub string
	//TODO return "https://someurl:8084", nil
	return string(out), nil
}

func SetFlockerHubTokenFile(tokenFile string) {
	logger.Info.Println("Setting FlockerHub Tokenfile %s", tokenFile)
	out, err := exec.Command("/opt/clusterhq/bin/dpcli", "set", "tokenfile", tokenFile).Output()
	if err != nil {
		logger.Error.Println("Could not set tokenfile, reason: ", out)
		logger.Error.Fatal(err)
	}
	logger.Info.Println(out)
}

func GetFlockerHubTokenFile() (flockerHubTokenFile string, err error) {
	logger.Info.Println("Getting FlockerHub Tokenfile")
	out, err := exec.Command("/opt/clusterhq/bin/dpcli", "get", "tokenfile").Output()
	if err != nil {
		logger.Error.Println("Could not get tokenfile, reason: ", out)
		return "", err
	}
	logger.Info.Println(out)
	//TODO Parse, "out" to get specific tokenfile string
	//TODO return "/root/vhut.txt", nil
	return string(out), nil
}

// Run the command to sync a volumeset
func syncVolumeset(volumeSetId string) {
	logger.Info.Println("Syncing Volumeset %s", volumeSetId)
	out, err := exec.Command("/opt/clusterhq/bin/dpcli", "sync",  "volumeset", volumeSetId).Output()
	if err != nil {
		logger.Error.Println("Could not sync dataset, reason: ", out)
		logger.Error.Fatal(err)
	}
	logger.Info.Println(out)
}

// Run the command to pull a specific snapshot
func pullSnapshot(snapshotId string){
	logger.Info.Println("Pulling Snapshot %s", snapshotId)
	out, err := exec.Command("/opt/clusterhq/bin/dpcli", "pull", "snapshot", snapshotId).Output()
	if err != nil {
		logger.Error.Println("Could not pull dataset, reason: ", out)
		logger.Error.Fatal(err)
	}
	logger.Info.Println(out)
}

// Wrapper for sync and pull which takes
// a List of type Volume above to pull.
func PullSnapshots(volumes []types.Volume) {
	for _, volume := range volumes {
		syncVolumeset(volume.VolumeSet)
		// maybe worth traking if we already sync'd a volumset
		// and skipping another sync during the same PullSnapshots call.
		pullSnapshot(volume.Snapshot)
	}
}

// Created a volume and returns it.
func createVolumeFromSnapshot(volumeName string, snapshotId string) (vol types.NewVolume, err error){
	logger.Info.Println("Creating Volume from Snapshot: %s", snapshotId)
	var attrString = fmt.Sprintf("created_by=fli-docker,from_snap=%s", snapshotId)
	cmd := exec.Command("/opt/clusterhq/bin/dpcli", "create", "volume", "--snapshot", snapshotId, "-a", attrString)
	createOut, err := cmd.Output()
	if err != nil {
		logger.Error.Fatal(err)
	}
	logger.Info.Println(string(createOut))
	r, _ := regexp.Compile("/chq/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
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