package main

import (
	"os"
	"flag"
	"strings"
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
	var push bool

	// FlagSets for SubCommands
	runSet := flag.NewFlagSet("fli-docker run", flag.ExitOnError)
	snapSet := flag.NewFlagSet("fli-docker snapshot", flag.ExitOnError)
	stopDestroySet := flag.NewFlagSet("fli-docker stop (or) fli-docker destroy", flag.ExitOnError)

	// runSet
	runSet.StringVar(&tokenfile, "t", "", "[OPTIONAL] Flocker Hub user token, optionally set it in the manifest YAML")
	runSet.StringVar(&flockerhub, "e", "", "[OPTIONAL] Flocker Hub endpoint, optionally set it in the manifest YAML")
	runSet.StringVar(&manifest, "f", "manifest.yml", "[OPTIONAL] Stateful application manifest file")
	runSet.BoolVar(&compose, "c", false, "[OPTIONAL] if flag is present, fli-docker will start the compose services")
	runSet.BoolVar(&verbose, "verbose", false, "[OPTIONAL] verbose logging")
	runSet.StringVar(&project, "p", "fli-compose", "[OPTIONAL] project name for compose if using -c")

	// snapSet
	snapSet.StringVar(&tokenfile, "t", "", "[OPTIONAL] Flocker Hub user token, optionally set it in the manifest YAML")
	snapSet.StringVar(&flockerhub, "e", "", "[OPTIONAL] Flocker Hub endpoint, optionally set it in the manifest YAML")
	snapSet.BoolVar(&push, "push", false, "[OPTIONAL] if flag is present, fli-docker will push new snapshots back to FlockerHub")
	snapSet.BoolVar(&verbose, "verbose", false, "[OPTIONAL] verbose logging")

	// stopSet
	stopDestroySet.BoolVar(&verbose, "verbose", false, "[OPTIONAL] verbose logging")
	stopDestroySet.StringVar(&manifest, "f", "manifest.yml", "[OPTIONAL] Stateful application manifest file")
    stopDestroySet.StringVar(&project, "p", "fli-compose", "[OPTIONAL] project name for compose if using -c")

	// Initialize logger before `verbose` is captured for
	// log messages before that conditional
	logger.Init(os.Stdout, ioutil.Discard, ioutil.Discard, os.Stderr)

	if (len(os.Args) > 1) {
		if (strings.Contains(os.Args[1], "help")) {os.Args[1] = "help"}
		switch os.Args[1] {
  			case "version":
    			logger.Message.Println(utils.FliDockerVersion)
    			os.Exit(0)
  			case "run":
    			runSet.Parse(os.Args[2:])
    			if verbose {
					logger.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
				}else{
					logger.Init(os.Stdout, ioutil.Discard, ioutil.Discard, os.Stderr)
				}
    		case "snapshot":
    			snapSet.Parse(os.Args[2:])
    			if verbose {
					logger.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
				}else{
					logger.Init(os.Stdout, ioutil.Discard, ioutil.Discard, os.Stderr)
				}
    		case "destroy":
    			stopDestroySet.Parse(os.Args[2:])
    			if verbose {
					logger.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
				}else{
					logger.Init(os.Stdout, ioutil.Discard, ioutil.Discard, os.Stderr)
				}
    		case "stop":
    			stopDestroySet.Parse(os.Args[2:])
    			if verbose {
					logger.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
				}else{
					logger.Init(os.Stdout, ioutil.Discard, ioutil.Discard, os.Stderr)
				}
    		case "help":
    			snapSet.Parse(os.Args[2:])
    			logger.Message.Println(utils.FliDockerHelp)
    			os.Exit(0)
    		default:
    			logger.Message.Println("Unrecognized Command. Use fli-docker --help.")
    			os.Exit(0)
		}
	} else {
		logger.Message.Println("Unrecognized Command. Use fli-docker --help.")
    	os.Exit(0)
	}

	var composeCmd string
	composeCmd = "docker-compose version"

	var fliCmd1 string
	var fliCmd2 string
	fliCmd1 = "fli version"
	fliCmd2 = utils.FliDockerCmd + "version"

	// check if needed dependencies are available
	isComposeAvail, err := utils.CheckForCmd(composeCmd)
	if (!isComposeAvail){
		logger.Info.Println(utils.ComposeHelpMessage)
		logger.Error.Fatal("Could not find `docker-compose` ", err)
	}else{
		logger.Info.Println("docker-compose Ready!")
	}

	isFliAvail1, err := utils.CheckForCmd(fliCmd1)
	isFliAvail2, err := utils.CheckForCmd(fliCmd2)
	var binary bool
	var docker bool
	var fliCmd string
	binary = true
	docker = false
	if (!isFliAvail1 && !isFliAvail2){
		logger.Info.Println(utils.FliHelpMessage)
		logger.Info.Fatal("fli not detected")
	}else{
		if (!isFliAvail1) {
			binary = false
			docker = true
			fliCmd = utils.FliDockerCmd
		}else{
			fliCmd = utils.FliBinaryCmd
		}
		logger.Info.Println("using fli container: ", docker)
		logger.Info.Println("using fli binary: ", binary)
	}

	if os.Args[1] == "run" {
		logger.Info.Println("Running: `fli-docker run`")

		if manifest == "manifest.yml" {
			logger.Warning.Println("Using default 'manifest.yml`, otherwise specify differently with -f")
		}

		// verify that the manifest exists
		isManifestAvail, err := utils.CheckForFile(manifest)
		if (!isManifestAvail){
			logger.Error.Println(err.Error())
			logger.Message.Fatal("Missing manifest, either name it 'manifest.yml' or pass in file with '-f'.")
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

		// was it passed with `-e`?
		if flockerhub == "" {
			logger.Info.Println("FlockerHub endpoint not specified with -e")
			fh, err := cli.GetFlockerHubEndpoint(fliCmd)
			if err != nil{
				logger.Error.Fatal("Could not get FlockerHub config")
			}
			logger.Info.Println("Existing FlockerHub Endpoint config: ", fh)
			// was it placed in manifest?
			flockerhubFromManifest := m.Hub.Endpoint
			logger.Info.Println("FlockerHub Endpoint " + m.Hub.Endpoint + " in manifest")
			if flockerhubFromManifest == "" {
				// Did the user have a pre-existing fli setup? 
				// Lets try and assume the volumes are there.
				if strings.Contains(fh, "FlockerHub URL:            -") {
					logger.Error.Fatal("Must set FlockerHub Endpoint")
				}else{
					logger.Info.Println("Trying existing FlockerHub configuration: ", fh)
				}
			}else{
				// set endpoint from manifest
				cli.SetFlockerHubEndpoint(flockerhubFromManifest, fliCmd)
			}
		}else{
			// set endpoint from fli-docker arg
			cli.SetFlockerHubEndpoint(flockerhub, fliCmd)
		}

		if tokenfile == "" {
			logger.Info.Println("token not specified with -t")
			tf, err := cli.GetFlockerHubTokenFile(fliCmd)
			if err != nil{
				logger.Error.Fatal("Could not get tokenfile config")
			}
			logger.Info.Println("Existing tokenfile config: ", tf)
			// Was is placed in the manifest?
			logger.Info.Println("tokenfile " + m.Hub.AuthToken + " in manifest")
			tokenfileFromManifest := m.Hub.AuthToken
			if tokenfileFromManifest == "" {
				if strings.Contains(tf, "Authentication Token File: -") {
					logger.Error.Fatal("Must set tokenfile")
				}else{
					logger.Info.Println("Trying existing tokenfile config: ", tf)
				}
			}else{
				cli.SetFlockerHubTokenFile(tokenfileFromManifest, fliCmd)
			}
		}else{
			cli.SetFlockerHubTokenFile(tokenfile, fliCmd)
		}

		// verify that the compose file exists.
		isComposeFileAvail, err := utils.CheckForFile(m.DockerApp)
		if (!isComposeFileAvail){
			logger.Error.Fatal(err.Error())
		}

		// try and pull snapshots
		logger.Message.Println("Pulling FlockerHub volumes...")
		cli.PullSnapshots(m.Volumes, fliCmd)

		// create volumes from snapshots and map them to 
		// `newVolPaths = {compose_volume_name : "/chq/<vol_path>"...}`
		logger.Message.Println("Creating volumes from snapshots...")
		newVolPaths, err := cli.CreateVolumesFromSnapshots(m.Volumes, fliCmd)

		// create a copy of the compose file before we edit it.
		// replace a fresh copy if we already copied before
		utils.CheckForCopy(m.DockerApp)
		// it will be `filename` + `-fli.copy`
		// will only copy if copy doesnt exist already
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

	} else if os.Args[1] == "snapshot" {
		logger.Info.Println("Running: `fli-docker snapshot`")

		// Does user want us to push snapshots back?
		if push {
			logger.Message.Println("Snapshotting and Pushing volumes to FlockerHub...")
			cli.SnapshotAndPushWorkingVolumes(fliCmd)
		}else{
			logger.Message.Println("Snapshotting volumes...")
			cli.SnapshotWorkingVolumes(fliCmd)
		}

	} else if os.Args[1] == "destroy" {
		logger.Info.Println("Running: `fli-docker destroy`")

		if manifest == "manifest.yml" {
			logger.Warning.Println("Using default 'manifest.yml`, otherwise specify differently with -f")
		}

		// verify that the manifest exists
		isManifestAvail, err := utils.CheckForFile(manifest)
		if (!isManifestAvail){
			logger.Error.Println(err.Error())
			logger.Message.Fatal("Missing manifest, either name it 'manifest.yml' or pass in file with '-f'.")
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

		logger.Info.Println("Destroying compose application")
		utils.DestroyCompose(m.DockerApp, project)

	} else if os.Args[1] == "stop" {
		logger.Info.Println("Running: `fli-docker stop`")

		if manifest == "manifest.yml" {
			logger.Warning.Println("Using default 'manifest.yml`, otherwise specify differently with -f")
		}

		// verify that the manifest exists
		isManifestAvail, err := utils.CheckForFile(manifest)
		if (!isManifestAvail){
			logger.Error.Println(err.Error())
			logger.Message.Fatal("Missing manifest, either name it 'manifest.yml' or pass in file with '-f'.")
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

		logger.Info.Println("Stopping compose application")
		utils.StopCompose(m.DockerApp, project)
	}
}
