package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/quangpq/docker-machine-driver-hetzner/driver"
)

// Version will be added once we start the build process by goreleaser
var version string

func main() {
	versionFlag := flag.Bool("v", false, "prints current docker-machine-driver-hetzner version")
	flag.Parse()
	if *versionFlag {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}
	plugin.RegisterDriver(driver.NewDriver(version))
}
