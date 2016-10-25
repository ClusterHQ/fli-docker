package cli

import (
	"os/exec"
	"os"
	"fmt"
	"strings"
	"bufio"

	"github.com/ClusterHQ/fli-docker/types"
	"github.com/ClusterHQ/fli-docker/logger"
	"github.com/ClusterHQ/fli-docker/utils"
)

/*
	Bindings to the FlockerHub CLI
*/

func SetFlockerHubEndpoint(endpoint string, fli string) {
	logger.Info.Println("Setting FlockerHub Endpoint: ", endpoint)
	var cmd = fmt.Sprintf("%s config -u %s", fli, endpoint)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		logger.Error.Println("Could not set endpoint")
		logger.Error.Fatal(err)
	}
	logger.Info.Println(string(out))
}

func GetFlockerHubEndpoint(fli string) (flockerhubEndpoint string, err error) {
	logger.Info.Println("Getting FlockerHub Endpoint")
	var cmd = fmt.Sprintf("%s config | grep 'FlockerHub URL:' | awk '{print $3}'", fli)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		logger.Error.Println("Could not get endpoint")
		logger.Error.Println(err)
		return "", err
	}
	logger.Info.Println(string(out))
	return string(out), nil
}

func SetFlockerHubTokenFile(tokenFile string, fli string) {
	logger.Info.Println("Setting FlockerHub Tokenfile: ", tokenFile)
	var cmd = fmt.Sprintf("%s config -t %s", fli, tokenFile)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		logger.Error.Println("Could not set tokenfile")
		logger.Error.Fatal(err)
	}
	logger.Info.Println(string(out))
}

func GetFlockerHubTokenFile(fli string) (flockerHubTokenFile string, err error) {
	logger.Info.Println("Getting FlockerHub Tokenfile")
	var cmd = fmt.Sprintf("%s config | grep 'Authentication Token File:' | awk '{print $4}'", fli)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		logger.Error.Println("Could not get tokenfile")
		logger.Error.Println(err)
		return "", err
	}
	logger.Info.Println(string(out))
	return string(out), nil
}

// Run the command to sync a volumeset
func syncVolumeset(volumeSetId string, fli string) {
	logger.Info.Println("Syncing Volumeset: ", volumeSetId)
	var cmd = fmt.Sprintf("%s sync %s", fli, volumeSetId)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		logger.Error.Println("Could not sync dataset")
		logger.Error.Fatal(err)
	}
	logger.Info.Println(string(out))
}

// Run the command to pull a specific snapshot
func pullSnapshot(volumeSetId string, snapshotId string, fli string){
	logger.Info.Println("Pulling Snapshot: ", snapshotId)
	var cmd = fmt.Sprintf("%s pull %s:%s", fli, volumeSetId, snapshotId)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		logger.Error.Println("Could not pull dataset, reason")
		logger.Error.Fatal(err)
	}
	logger.Info.Println(string(out))
}

// Wrapper for sync and pull which takes
// a list of type Volume
func PullSnapshots(volumes []types.Volume, fli string) {
	for _, volume := range volumes {
		syncVolumeset(volume.VolumeSet, fli)
		pullSnapshot(volume.VolumeSet, volume.Snapshot, fli)
	}
}

// Created a volume and returns it.
func createVolumeFromSnapshot(volumeName string, volumeSet string, snapshotId string, fli string) (vol types.NewVolume, err error){
	logger.Info.Println("Creating Volume from Snapshot: ", snapshotId)
	var attrString = fmt.Sprintf("created_by=fli-docker,from_snap=%s", snapshotId)
	uuid, err := utils.GenUUID()
	if err != nil {
		logger.Error.Fatal(err)
	}

	var volName = fmt.Sprintf("fli-%s", uuid)
	var createCmd = fmt.Sprintf("%s clone %s:%s -a %s %s", fli, volumeSet, snapshotId, attrString, volName)
	cmd := exec.Command("sh", "-c", createCmd)
	createOut, err := cmd.Output()
	if err != nil {
		logger.Error.Fatal(err)
	}
	var path = strings.TrimSpace(string(createOut))
	logger.Info.Println(path)
	if path == "" {
			logger.Error.Fatal("Could not find volume path")
	 }
	return types.NewVolume{Name: volumeName, VolumePath: path, VolumeName: volName, VolumeSet: volumeSet}, nil
}

func saveCurrentWorkingVols(volumes []types.NewVolume) {
	// open files r and w
	exists, _ := utils.CheckForFile(".flidockervols")
	if exists {
		err := os.Remove(".flidockervols")
		if err != nil {
			logger.Error.Fatal("Could not delete .flidockervols")
		}
	}
	_, err := os.Create(".flidockervols")
	if err != nil {
        logger.Error.Fatal(err)
    }
    file, err := os.OpenFile(".flidockervols", os.O_APPEND|os.O_WRONLY,0600)
    if err != nil {
        logger.Error.Fatal(err)
    }
    defer file.Close()

    for _, newVol := range volumes {
    	var volRecord = fmt.Sprintf("%s,%s\n", strings.TrimSpace(newVol.VolumeName), newVol.VolumeSet)
    	if _, err = file.WriteString(volRecord); err != nil {
     		logger.Error.Fatal(err)
    	}
	}
	logger.Info.Println("Saves working vols to .flidockervols")
}

func CreateVolumesFromSnapshots(volumes []types.Volume, fli string) (newVols []types.NewVolume, err error) {
	vols := []types.NewVolume{}
	for _, volume := range volumes {
		vol, err := createVolumeFromSnapshot(volume.Name, volume.VolumeSet, volume.Snapshot, fli)
		if err != nil {
			return nil, err
		}else {
			vols = append(vols, vol)
		}
	}
	// Record current working volumes.
	saveCurrentWorkingVols(vols)
	return vols, nil
}

// Run the command to push a specific snapshot
func pushSnapshot(volumeSetId string, snapshotId string, fli string){
	logger.Info.Println("Pushing Snapshot: ", snapshotId)
	var cmd = fmt.Sprintf("%s push %s:%s", fli, volumeSetId, snapshotId)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		logger.Error.Println("Could not push snapshot, reason")
		logger.Error.Fatal(err)
	}
	logger.Info.Println(string(out))
}

// Run the command to push a specific snapshot
func createSnapshot(volumeSetId string, volumeId string, snapName string, fli string){
	logger.Info.Println("Creating Snapshot: ", snapName)
	var cmd = fmt.Sprintf("%s snapshot %s:%s %s", fli, volumeSetId, volumeId, snapName)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		logger.Error.Println("Could not push snapshot, reason")
		logger.Error.Fatal(err)
	}
	logger.Info.Println(string(out))
}

func SnapshotWorkingVolumes(fli string){
	file, err := os.Open(".flidockervols")
    if err != nil {
        logger.Error.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        result := strings.Split(scanner.Text(), ",")
        var uuid, _ = utils.GenUUID()
        var snap = fmt.Sprintf("%s-%s", result[0], uuid)
        logger.Message.Println("Snapshotting", result[0], "from Volumeset", result[1])
        createSnapshot(result[1], result[0], snap, fli)
    }

    if err := scanner.Err(); err != nil {
        logger.Error.Fatal(err)
    }
}

func SnapshotAndPushWorkingVolumes(fli string){
    file, err := os.Open(".flidockervols")
    if err != nil {
        logger.Error.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        result := strings.Split(scanner.Text(), ",")
        var uuid, _ = utils.GenUUID()
        var snap = fmt.Sprintf("%s-%s", result[0], uuid)
        logger.Message.Println("Snapshotting and Pushing", result[0], "from Volumeset", result[1])
        createSnapshot(result[1], result[0], snap, fli)
        syncVolumeset(result[1], fli)
        pushSnapshot(result[1], snap, fli)
    }

    if err := scanner.Err(); err != nil {
        logger.Error.Fatal(err)
    }
}
