package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/choria-io/discovery_proxy/cmd"
)

func main() {
	err := cmd.ParseCLI()
	if err != nil {
		log.Fatalf("Could not configure: %s", err.Error())
		os.Exit(1)
	}

	err = cmd.Run()
	if err != nil {
		log.Fatalf("Could not run: %s", err.Error())
	}
}
