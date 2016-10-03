package main

import (
	"os"
	"flag"
	"path/filepath"
	"io/ioutil"
	"github.com/ClusterHQ/fli-docker/utils"
	"github.com/ClusterHQ/fli-docker/cli"
	"github.com/ClusterHQ/fli-docker/logger"
)

func main() {

	// should this be a struct?
	var tokenfile string
	var flockerhub string
	var manifest string
	var compose bool
	var verbose bool
	var project string

	flag.StringVar(&tokenfile, "t", "", "Flocker Hub user token")
	flag.StringVar(&flockerhub, "e", "", "Flocker Hub endpoint")
	flag.StringVar(&manifest, "f", "manifest.yml", "Stateful application manifest file")
	flag.BoolVar(&compose, "c", false, "if flag is present, fli-docker will start the compose services")
	flag.BoolVar(&verbose, "verbose", false, "verbose logging")
	flag.StringVar(&project, "project", "fli-compose", "project name for compose if using -c")

	// parse all the flags from user input
	flag.Parse()

    if verbose {
    	logger.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
    }else{
    	logger.Init(os.Stdout, ioutil.Discard, ioutil.Discard, os.Stderr)
    }

	var composeCmd string
	composeCmd = "docker-compose version"

	var fliCmd string
	// this will need `fli version` or somthing
	fliCmd = "/opt/clusterhq/bin/dpcli"

	// check if needed dependencies are available
	isComposeAvail, err := utils.CheckForCmd(composeCmd)
	if (!isComposeAvail){
		logger.Info.Println(utils.ComposeHelpMessage)
		logger.Error.Fatal("Could not find `docker-compose` ", err)
	}else{
		logger.Info.Println("docker-compose Ready!")
	}

	isFliAvail, err := utils.CheckForPath(fliCmd)
	if (!isFliAvail){
		logger.Info.Println(utils.FliHelpMessage)
		logger.Error.Fatal("Could not find `fli` ", err)
	}else{
		logger.Info.Println("fli Ready!")
	}

	if manifest == "manifest.yml" {
		logger.Warning.Println("Using default 'manifest.yml`, otherwise specify differently with -f")
	}

	// verify that the manifest exists
	isManifestAvail, err := utils.CheckForFile(manifest)
	if (!isManifestAvail){
		logger.Error.Fatal(err.Error())
	}

	// get the yaml file passed in the args.
	filename, _ := filepath.Abs(manifest)
	// read the file.
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error.Fatal(err.Error())
	}

	// pass the file to the ParseManifest
	logger.Message.Println("Parsing the fli manifest...")
	m := utils.ParseManifest(yamlFile)

	if flockerhub == "" {
		logger.Warning.Println("FlockerHub endpoint not specified with -e, checking if set, or setting from manifest")
		// TODO check from `dpcli get volumehub` if set.
		// utils.GetFlockerHubEndpoint()
		// IF ITS NOT SET
			// check for endpoint in m.Hub.Endpoint from manifest.
		logger.Info.Println("Found FlockerHub Endpoint " + m.Hub.Endpoint + "in manifest")
			// TODO check if blank, exit if blank "must set FlockerHub endpoint"
		flockerhub = m.Hub.Endpoint
			// TODO if not, set the endpoint
			// utils.SetFlockerHubEndpoint(flockerhub)
	}

	if tokenfile == "" {
		logger.Warning.Println("token not specifed with -t, checking if set, or setting from manifest")
		// TODO  check from `dpcli get tokenfile` if set
		// utils.GetFlockerHubTokenFile()
		// IF ITS NOT SET
			// check for tokenfile in m.Hub.AuthToken from manifest.
		logger.Info.Println("Found tokenfile " + m.Hub.AuthToken + "in manifest")
			// TODO check if blank, exit if blank "must set FlockerHub tokenfile"
		flockerhub = m.Hub.AuthToken
			// TODO set dpcli tokenfile
			// utils.SetFlockerHubTokenFile(tokenfile)
	}

	// verify that the compose file exists.
	isComposeFileAvail, err := utils.CheckForFile(m.DockerApp)
	if (!isComposeFileAvail){
		logger.Error.Fatal(err.Error())
	}

	// try and pull snapshots
	logger.Message.Println("Pulling FlockerHub volumes...")
	cli.PullSnapshots(m.Volumes)

	// create volumes from snapshots and map them to 
	// `newVolPaths = {compose_volume_name : "/chq/<vol_path>"...}`
	logger.Message.Println("Creating volumes from snapshots...")
	newVolPaths, err := cli.CreateVolumesFromSnapshots(m.Volumes)

	// create a copy of the compose file before we edit it.
	// it will be `filename` + `-fli.copy`
	utils.MakeCopy(m.DockerApp)

	// replace volume_name with `volume_name`'s associated 
	// "/chq/<vol_path/" and modify the compose file
	logger.Message.Println("Mapping new volumes in compose file...")
	for _, newVol := range newVolPaths {
		utils.MapVolumeToCompose(newVol.Name, newVol.VolumePath, m.DockerApp)
	}

	// this just parses the compose file, not needed, 
	// but in verbose we can be thorough as it also prints it.
	utils.ParseCompose(m.DockerApp)

	// `-c` means "run compose".
	// if not, done, we only modified the compose file.
	if compose {
		logger.Info.Println("compose option set, running docker-compose")
		utils.RunCompose(m.DockerApp, project)
	}
}
