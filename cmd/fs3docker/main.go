package main

import (
	"fmt"
	"log"
	"flag"
	"github.com/wallnerryan/fs3todocker/utils"
)

func main() {
    var user string
    var token string
    var endpoint string
    var manifest string
    var composeOpts string 

    // `docker-compose` path
    // TODO configurable in the future
    var composePath string
    composePath = "/usr/local/bin/docker-compose"

    // `dpcli` Path
    // TODO maybe configurable in the future
    var dpcliPath string
    dpcliPath = "/opt/clusterhq/bin/dpcli"

    // Check if needed dependencies are available
    isComposeAvail, err := utils.CheckForTool(composePath)
    if (!isComposeAvail){
    	fmt.Printf("-----------------------------------------------------------------------\n")
    	fmt.Printf("docker-compose is not installed, it is needed to use fs3todocker\n")
		fmt.Printf("docker-compose is available at https://docs.docker.com/compose/install/\n")
		fmt.Printf("-----------------------------------------------------------------------\n")
		log.Fatal(err.Error())
    }else{
		log.Println("docker-compose Ready!\n")
    }

    isDpcliAvail, err := utils.CheckForTool(dpcliPath)
    if (!isDpcliAvail){
    	fmt.Printf("-------------------------------------------------------\n")
    	fmt.Printf("dpcli is not installed, it is needed to use fs3todocker\n")
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

    flag.Parse()

    fmt.Printf("user = %s\n", user)
    fmt.Printf("token = %s\n", token)
    fmt.Printf("endpoint = %s\n", endpoint)
    fmt.Printf("manifest = %s\n", manifest)
    fmt.Printf("composeOpts = %s\n", composeOpts)
}
