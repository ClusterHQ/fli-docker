/*
 *Copyright ClusterHQ Inc.  See LICENSE file for details.
 *
 */

package main

import (
	"os"
	"flag"
	"fmt"
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
	var zpool string
	var manifest string
	var compose bool
	var verbose bool
	var project string
	var push bool
	var clean bool

	// FlagSets for SubCommands
	runSet := flag.NewFlagSet("fli-docker run", flag.ExitOnError)
	snapSet := flag.NewFlagSet("fli-docker snapshot", flag.ExitOnError)
	stopSet := flag.NewFlagSet("fli-docker stop", flag.ExitOnError)
	destroySet := flag.NewFlagSet("fli-docker destroy", flag.ExitOnError)

	// runSet
	runSet.StringVar(&tokenfile, "t", "/root/token.txt", "[OPTIONAL if using binary] FlockerHub Auth Token")
	runSet.StringVar(&flockerhub, "e", "", "[OPTIONAL] Flocker Hub endpoint, optionally set it in the manifest YAML")
	runSet.StringVar(&manifest, "f", "manifest.yml", "[OPTIONAL] Stateful application manifest file")
	runSet.BoolVar(&compose, "c", false, "[OPTIONAL] if flag is present, fli-docker will start the compose services")
	runSet.BoolVar(&verbose, "verbose", false, "[OPTIONAL] verbose logging")
	runSet.StringVar(&project, "p", "fli-compose", "[OPTIONAL] project name for compose if using -c")
	runSet.StringVar(&zpool, "z", "chq", "[OPTIONAL if using binary] Name of the zpool, must proivide if using fli docker image & zpool other than 'chq'")


	// snapSet
	snapSet.StringVar(&tokenfile, "t", "", "[OPTIONAL if using binary] FlockerHub Auth Token")
	snapSet.StringVar(&flockerhub, "e", "", "[OPTIONAL] Flocker Hub endpoint, optionally set it in the manifest YAML")
	snapSet.BoolVar(&push, "push", false, "[OPTIONAL] if flag is present, fli-docker will push new snapshots back to FlockerHub")
	snapSet.BoolVar(&verbose, "verbose", false, "[OPTIONAL] verbose logging")
	snapSet.StringVar(&zpool, "z", "chq", "[OPTIONAL if using binary] Name of the zpool, must proivide if using fli docker image & zpool other than 'chq'")

	// stopSet
	stopSet.BoolVar(&verbose, "verbose", false, "[OPTIONAL] verbose logging")
	stopSet.StringVar(&manifest, "f", "manifest.yml", "[OPTIONAL] Stateful application manifest file")
	stopSet.StringVar(&project, "p", "fli-compose", "[OPTIONAL] project name for compose if using -c")

	// destroySet
	destroySet.BoolVar(&verbose, "verbose", false, "[OPTIONAL] verbose logging")
	destroySet.StringVar(&manifest, "f", "manifest.yml", "[OPTIONAL] Stateful application manifest file")
	destroySet.StringVar(&project, "p", "fli-compose", "[OPTIONAL] project name for compose if using -c")
	destroySet.BoolVar(&clean, "clean", false, "[OPTIONAL] places docker-compose file back to original state")

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
    			destroySet.Parse(os.Args[2:])
    			if verbose {
					logger.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
				}else{
					logger.Init(os.Stdout, ioutil.Discard, ioutil.Discard, os.Stderr)
				}
    		case "stop":
    			stopSet.Parse(os.Args[2:])
    			if verbose {
					logger.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
				}else{
					logger.Init(os.Stdout, ioutil.Discard, ioutil.Discard, os.Stderr)
				}
    		case "help":
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

	var dockerCmd string
	dockerCmd = "docker version"

	var fliCmd1 string
	var fliCmd2 string
	fliCmd1 = "fli version"
	fliCmd2 = "docker inspect clusterhq/fli"

	// check if needed dependencies are available
	isDockerAvail, err := utils.CheckForCmd(dockerCmd)
	if (!isDockerAvail){
		logger.Info.Println(err)
		logger.Info.Println(utils.DockerHelpMessage)
		logger.Error.Fatal("Could not find docker, please install docker. ")
	}else{
		logger.Info.Println("Docker Ready!")
	}

	isFliAvail1, _ := utils.CheckForCmd(fliCmd1)
	isFliAvail2, _ := utils.CheckForCmd(fliCmd2)
	var binary bool
	var docker bool
	var fliCmd string
	binary = true
	docker = false
	if (!isFliAvail1 && !isFliAvail2){
		logger.Info.Println(utils.FliHelpMessage)
		logger.Error.Fatal("Fli not detected, please install / configure Fli.")
	}else{
		if (!isFliAvail1) {
			binary = false
			docker = true
			fliCmd = "docker run --rm --privileged -v /chq/:/chq/:shared -v /var/log/:/var/log/ -v /root:/root -v /lib/modules:/lib/modules clusterhq/fli "
		}else{
			fliCmd = "fli "
		}
		logger.Info.Println("using fli container: ", docker)
		logger.Info.Println("using fli binary: ", binary)
	}

	if docker {
		// Figure out what the appropriate docker command is is using docker
		// docker needs certain paths to work.
		if (os.Args[1] == "run" || os.Args[1] == "snapshot") {

			resp, err := utils.GetZPool(zpool)
			if err != nil {
				logger.Info.Println(resp)
	    		logger.Message.Fatal("ZPool doesnt not exist: ", zpool)
	    	}

			if os.Args[1] == "run" {
				if tokenfile == "" {
					logger.Info.Println("token not specified with -t")
					logger.Message.Fatal("Need to set '-tokenfile'")
				}

				if zpool != "chq" {
					if strings.Contains(tokenfile, "/root/"){
						cmd := fmt.Sprintf("docker run --rm --privileged -v /%s/:/%s/:shared -v /var/log/:/var/log/ -v /root:/root -v /lib/modules:/lib/modules clusterhq/fli ", zpool, zpool)
						fliCmd = cmd
					}else {
						path, err := utils.GetBasePath(tokenfile)
						if err != nil {
	    					logger.Message.Fatal("Could not get path of auth token")
	    				}
						cmd := fmt.Sprintf("docker run --rm --privileged -v /%s/:/%s/:shared -v /var/log/:/var/log/ -v /root:/root -v %s:%s -v /lib/modules:/lib/modules clusterhq/fli ", zpool, zpool, path, path)
						fliCmd = cmd
					}
				}else {
					if !(strings.Contains(tokenfile, "/root/")){
						path, err := utils.GetBasePath(tokenfile)
						if err != nil {
	    					logger.Message.Fatal("Could not get path of auth token")
	    				}
						cmd := fmt.Sprintf("docker run --rm --privileged -v /chq/:/chq/:shared -v /var/log/:/var/log/ -v /root:/root -v %s:%s -v /lib/modules:/lib/modules clusterhq/fli ", path, path)
						fliCmd = cmd
					}
				}
				logger.Info.Println("Using docker command: ", fliCmd)
				cli.SetFlockerHubTokenFile(tokenfile, fliCmd)
			}

			if os.Args[1] == "snapshot" {
				if push {
					if tokenfile == "" {
						logger.Info.Println("token not specified with -t")
						logger.Message.Fatal("Need to set '-tokenfile'")
					}

					if zpool != "chq" {
						if strings.Contains(tokenfile, "/root/"){
							cmd := fmt.Sprintf("docker run --rm --privileged -v /%s/:/%s/:shared -v /var/log/:/var/log/ -v /root:/root -v /lib/modules:/lib/modules clusterhq/fli ", zpool, zpool)
							logger.Info.Println("Using docker command: ", cmd)
							fliCmd = cmd
						}else {
							path, err := utils.GetBasePath(tokenfile)
							if err != nil {
	    						logger.Message.Fatal("Could not get path of auth token")
	    					}
							cmd := fmt.Sprintf("docker run --rm --privileged -v /%s/:/%s/:shared -v /var/log/:/var/log/ -v /root:/root -v %s:%s -v /lib/modules:/lib/modules clusterhq/fli ", zpool, zpool, path, path)
							logger.Info.Println("Using docker command: ", cmd)
							fliCmd = cmd
						}
					}else {
						if !(strings.Contains(tokenfile, "/root/")){
							path, err := utils.GetBasePath(tokenfile)
							if err != nil {
	    						logger.Message.Fatal("Could not get path of auth token")
	    					}
							cmd := fmt.Sprintf("docker run --rm --privileged -v /chq/:/chq/:shared -v /var/log/:/var/log/ -v /root:/root -v %s:%s -v /lib/modules:/lib/modules clusterhq/fli ", path, path)
							fliCmd = cmd
						}
					}
				}else {
					//only need to check zpool as token doesnt matter if not pushing.
					if zpool != "chq" {
						cmd := fmt.Sprintf("docker run --rm --privileged -v /%s/:/%s/:shared -v /var/log/:/var/log/ -v /root:/root -v /lib/modules:/lib/modules clusterhq/fli ", zpool, zpool)
						fliCmd = cmd
					}else {
						if !(strings.Contains(tokenfile, "/root/")){
							path, err := utils.GetBasePath(tokenfile)
							if err != nil {
	    						logger.Message.Fatal("Could not get path of auth token")
	    					}
							cmd := fmt.Sprintf("docker run --rm --privileged -v /chq/:/chq/:shared -v /var/log/:/var/log/ -v /root:/root -v %s:%s -v /lib/modules:/lib/modules clusterhq/fli ", path, path)
							fliCmd = cmd
						}
					}
				}
				logger.Info.Println("Using docker command: ", fliCmd)
				cli.SetFlockerHubTokenFile(tokenfile, fliCmd)
			}
		}
	}

	if binary {
		// Check that tokenfile is in use if using binary, or need to set one.
		if (os.Args[1] == "run" || os.Args[1] == "snapshot") {
			tf, err := cli.GetFlockerHubTokenFile(fliCmd)
			if err != nil {
				logger.Info.Println("Could not get tokenfile config")
				if tokenfile == "" {
					logger.Info.Println("token not specified with -t")
					logger.Message.Fatal("Need to set '-tokenfile'")
				}else{
					cli.SetFlockerHubTokenFile(tokenfile, fliCmd)
				}
			}else{
				logger.Info.Println("Using existing token: ", tf)
			}
		}
	}

	if os.Args[1] == "run" {
		logger.Info.Println("Running: `fli-docker run`")

		if manifest == "manifest.yml" {
			logger.Warning.Println("Using default 'manifest.yml`, otherwise specify differently with -f")
		}

		// verify that the manifest exists
		isManifestAvail, err := utils.CheckForFile(manifest)
		if (!isManifestAvail){
			logger.Info.Println(err.Error())
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
					logger.Message.Fatal("Must set FlockerHub Endpoint")
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

		// verify that the compose file exists.
		isComposeFileAvail, err := utils.CheckForFile(m.DockerApp)
		if (!isComposeFileAvail){
			logger.Info.Println(err.Error())
			logger.Message.Fatal("Docker Compose file doesnt exist.")
		}

		// try and pull snapshots
		logger.Message.Println("Pulling FlockerHub volumes...")
		cli.PullSnapshots(m.Volumes, fliCmd)

		// create volumes from snapshots and map them to 
		// `newVolPaths = {compose_volume_name : "/chq/<vol_path>"...}`
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
			// Dont need token path for snapshotting locally
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
			logger.Info.Println(err.Error())
			logger.Message.Fatal("Missing manifest, either name it 'manifest.yml' or pass in file with '-f'.")
		}

		// get the yaml file passed in the args.
		filename, _ := filepath.Abs(manifest)
		// read the file.
		yamlFile, err := ioutil.ReadFile(filename)
		if err != nil {
			logger.Info.Println(err.Error())
			logger.Message.Fatal("Docker Compose file doesnt exist.")
		}

		// pass the file to the ParseManifest
		logger.Message.Println("Parsing the fli manifest...")
		m := utils.ParseManifest(yamlFile)

		logger.Info.Println("Destroying compose application")
		utils.DestroyCompose(m.DockerApp, project)

		if clean {
			logger.Message.Println("Cleaning files...")
			utils.CleanEnv(m.DockerApp)
		}

	} else if os.Args[1] == "stop" {
		logger.Info.Println("Running: `fli-docker stop`")

		if manifest == "manifest.yml" {
			logger.Warning.Println("Using default 'manifest.yml`, otherwise specify differently with -f")
		}

		// verify that the manifest exists
		isManifestAvail, err := utils.CheckForFile(manifest)
		if (!isManifestAvail){
			logger.Info.Println(err.Error())
			logger.Message.Fatal("Missing manifest, either name it 'manifest.yml' or pass in file with '-f'.")
		}

		// get the yaml file passed in the args.
		filename, _ := filepath.Abs(manifest)
		// read the file.
		yamlFile, err := ioutil.ReadFile(filename)
		if err != nil {
			logger.Info.Println(err.Error())
			logger.Message.Fatal("Docker Compose file doesnt exist.")
		}

		// pass the file to the ParseManifest
		logger.Message.Println("Parsing the fli manifest...")
		m := utils.ParseManifest(yamlFile)

		logger.Info.Println("Stopping compose application")
		utils.StopCompose(m.DockerApp, project)
	}
}
