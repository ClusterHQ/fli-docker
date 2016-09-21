package main

import (
	"fmt"
	"log"
	"flag"
	"path/filepath"
	"io/ioutil"
	"github.com/wallnerryan/fli-docker/utils"
)

func main() {

	// should this be a struct?
	var user string
	var token string
	var endpoint string
	var manifest string
	var compose bool
	var verbose bool

	var composeCmd string
	composeCmd = "docker-compose version"

	var fliCmd string
	fliCmd = "/opt/clusterhq/bin/dpcli" //this will need `fli version` or somthing

	// Check if needed dependencies are available
	isComposeAvail, err := utils.CheckForCmd(composeCmd)
	if (!isComposeAvail){
		fmt.Printf(utils.ComposeHelpMessage)
		log.Fatal("Could not find `docker-compose` ", err)
	}else{
		log.Println("docker-compose Ready!\n")
	}

	isFliAvail, err := utils.CheckForPath(fliCmd)
	if (!isFliAvail){
		fmt.Printf(utils.FliHelpMessage)
		log.Fatal("Could not find `fli` ", err)
	}else{
		log.Println("fli Ready!\n")
	}

	flag.StringVar(&user, "u", "", "Flocker Hub username")
	flag.StringVar(&token, "t", "", "Flocker Hub user token")
	// Should we replace or add the above with the option to point to vhub.txt?
	flag.StringVar(&endpoint, "e", "", "Flocker Hub endpoint")
	flag.StringVar(&manifest, "f", "manifest.yml", "Stateful application manifest file")
	flag.BoolVar(&compose, "c", false, "if flag is present, fli-docker will start the compose file with 'up -d'")
	flag.BoolVar(&verbose, "v", false, "verbose logging")


	// Parse all the flags from user input
	flag.Parse()

	if manifest == "manifest.yml" {
		if verbose {
			log.Println("Using default 'manifest.yml`, otherwise specify differently with -f")
		}
	}

	// Verify that the manifest exists
	isManifestAvail, err := utils.CheckForFile(manifest)
	if (!isManifestAvail){
		log.Fatal(err.Error())
	}

	// Get the yaml file passed in the args.
	filename, _ := filepath.Abs(manifest)
	// Read the file.
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// Pass the file to the ParseManifest
	m := utils.ParseManifest(yamlFile)

	if endpoint == "" {
		if verbose {
			log.Println("endpoint not specifed with -e, checking manifest")
		}
		// TODO - then we can check m.Hub.Endpoint
	}

	if user == "" {
		if verbose {
			log.Println("user not specifed with -u, checking manifest")
		}
		// TODO - then we can check m.Hub.User
	}

	if token == "" {
		if verbose {
			log.Println("token not specifed with -t, checking manifest")
		}
		// TODO - then we can check m.Hub.AuthToken
	}

	// Verify that the compose file exists.
	isComposeFileAvail, err := utils.CheckForFile(m.DockerApp)
	if (!isComposeFileAvail){
		log.Fatal(err.Error())
	}

	// Try and pull snapshots
	// TODO need to return err and check for it?
	utils.PullSnapshots(m.Volumes)

	// Create volumes from snapshots and map them to 
	//    newVolPaths = {compose_volume_name : "/chq/<vol_path>"}
	newVolPaths, err := utils.CreateVolumesFromSnapshots(m.Volumes)

	// TODO need to return err and check for it?
	utils.MakeCopy(m.DockerApp)

	// replace volume_name with volume_name's associated "/chq/<vol_path/"
	// write file back to compose file
	for _, newVol := range newVolPaths {
		utils.MapVolumeToCompose(newVol.Name, newVol.VolumePath, m.DockerApp)
	}

	if verbose {
		// Verify Compose File
		utils.ParseCompose(m.DockerApp)
	}

	// (IF) -c is there for compose args, run compose, if not, done.
	if compose {
		if verbose {
			log.Println("compose option set, running docker-compose")
		}
		utils.RunCompose(m.DockerApp)
	}

}
