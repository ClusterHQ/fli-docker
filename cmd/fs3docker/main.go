package main

import (
	"fmt"
	"flag"
	_ "github.com/wallnerryan/fs3todocker/utils"
)

func main() {
    var user string
    var token string
    var endpoint string
    var manifest string
    var composeOpts string 


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
