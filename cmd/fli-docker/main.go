package main

import (
	"fmt"
	"log"
	"flag"
	"path/filepath"
	"io/ioutil"
	"github.com/ClusterHQ/fli-docker/utils"
)

func main() {

	// should this be a struct?
	var tokenfile string
	var flockerhub string
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

	flag.StringVar(&tokenfile, "t", "", "Flocker Hub user token")
	// Should we replace or add the above with the option to point to vhub.txt?
	flag.StringVar(&flockerhub, "e", "", "Flocker Hub endpoint")
	flag.StringVar(&manifest, "f", "manifest.yml", "Stateful application manifest file")
	flag.BoolVar(&compose, "c", false, "if flag is present, fli-docker will start the compose services")
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

	if flockerhub == "" {
		if verbose {
			log.Println("FlockerHub endpoint not specifed with -e, checking if set, or setting from manifest")
		}
		// TODO check from `dpcli get volumehub` if set.
		// utils.GetFlockerHubEndpoint()
		// IF ITS NOT SET
			// check for endpoint in m.Hub.Endpoint from manifest.
		log.Println("Found FlockerHub Endpoint " + m.Hub.Endpoint + "in manifest")
			// TODO check if blank, exit if blank "must set FlockerHub endpoint"
		flockerhub = m.Hub.Endpoint
			// TODO if not, set the endpoint
			// utils.SetFlockerHubEndpoint(flockerhub)
	}

	if tokenfile == "" {
		if verbose {
			log.Println("token not specifed with -t, checking if set, or setting from manifest")
		}
		// TODO  check from `dpcli get tokenfile` if set
		// utils.GetFlockerHubTokenFile()
		// IF ITS NOT SET
			// check for tokenfile in m.Hub.AuthToken from manifest.
		log.Println("Found tokenfile " + m.Hub.AuthToken + "in manifest")
			// TODO check if blank, exit if blank "must set FlockerHub tokenfile"
		flockerhub = m.Hub.AuthToken
			// TODO set dpcli tokenfile
			// utils.SetFlockerHubTokenFile(tokenfile)
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
