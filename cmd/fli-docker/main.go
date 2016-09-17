package main

import (
    "fmt"
    "log"
    "flag"
    "github.com/wallnerryan/fli-docker/utils"
)

func main() {

    // should this be a struct?
    var user string
    var token string
    var endpoint string
    var manifest string
    var composeOpts string 

    // should be configurable in the future?
    var composePath string
    composePath = "/usr/local/bin/docker-compose"

    // should be configurable in the future?
    var dpcliPath string
    dpcliPath = "/opt/clusterhq/bin/dpcli"

    // Check if needed dependencies are available
    isComposeAvail, err := utils.CheckForPath(composePath)
    if (!isComposeAvail){
        fmt.Printf("-----------------------------------------------------------------------\n")
        fmt.Printf("docker-compose is not installed, it is needed to use flitodock\n")
        fmt.Printf("docker-compose is available at https://docs.docker.com/compose/install/\n")
        fmt.Printf("-----------------------------------------------------------------------\n")
    log.Fatal(err.Error())
    }else{
    log.Println("docker-compose Ready!\n")
    }

    isDpcliAvail, err := utils.CheckForPath(dpcliPath)
    if (!isDpcliAvail){
        fmt.Printf("-------------------------------------------------------\n")
        fmt.Printf("dpcli is not installed, it is needed to use flitodock\n")
        fmt.Printf("dpcli is available at https://clusterhq.com\n")
        fmt.Printf("-------------------------------------------------------\n")
    log.Fatal(err.Error())
    }else{
    log.Println("dpcli Ready!\n")
    }

    flag.StringVar(&user, "u", "", "Flocker Hub username")
    flag.StringVar(&token, "t", "", "Flocker Hub user token")
    flag.StringVar(&endpoint, "v", "", "Flocker Hub endpoint")
    flag.StringVar(&manifest, "f", "manifest.yml", "Stateful application manifest file")
    flag.StringVar(&composeOpts, "c", "up", "Options to pass to Docker Compose such as 'up -d'")
    /* 
    Im thinking this should be optional meaning if its not
    present then flidock will not also run the docker-compose command
    but rather will just edit the docker-compose.yml file in place
    and let the use run the docker-compose command. 
    This may be even a good option to start with instead of using
    '-c' at all.
    */

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

    // 1. Verify that the compose file exists.
    // 2. Verify that the manifest exists
    // 3. Process the manifest into a Struct in YAML
    // 4. and get a mapping of 
    //    compose_volume_name : {volumeset: <id>, snapshot: <id>}
    // 5. Try and pull snapshots
    // 6. Create volumes from snapshots and map them to 
    //    {compose_volume_name : "/chq/<vol_path>"}
    // 7. Parse the the compose file into struct YAML
    // 8. replace volume_name with volume_name's associated "/chq/<vol_path/"
    // 9. write file back to compose file
    // 10 (IF) -c is there for compose args, run compose, if not, done.

}
