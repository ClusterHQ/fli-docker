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
    var composeOpts string 

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
    flag.StringVar(&endpoint, "v", "", "Flocker Hub endpoint")
    flag.StringVar(&manifest, "f", "manifest.yml", "Stateful application manifest file")
    flag.StringVar(&composeOpts, "c", "up", "Options to pass to Docker Compose such as 'up -d'") //optional

    // Parse all the flags from user input
    flag.Parse()

    /*
    # only for debug
    fmt.Printf("user = %s\n", user)
    fmt.Printf("token = %s\n", token)
    fmt.Printf("endpoint = %s\n", endpoint)
    fmt.Printf("manifest = %s\n", manifest)
    fmt.Printf("composeOpts = %s\n", composeOpts)
    */

    //TODO check for empty vars, or default ones.
    //  - only Opt being used now is 'manifest'
    //    others need to be checked and used
    //    such as the endpoint for vhub
    //    and user/token if we want or if its in the
    //    the manifest.

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

    // Verify Compose File
    utils.ParseCompose(m.DockerApp)

    // (IF) -c is there for compose args, run compose, if not, done.
    // TODO

}
