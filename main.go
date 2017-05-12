package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/choria-io/pdbproxy/discovery"
	"gopkg.in/alecthomas/kingpin.v2"
)

const Version = "0.0.1"

var config = discovery.Config{}

func main() {
	if err := configure(); err != nil {
		log.Fatalf("Could not configure pdbproxy: %s", err.Error())
		os.Exit(1)
	}

	if err := serve(); err != nil {
		log.Fatalf("Could not start webserver: %s", err.Error())
	}
}

func serve() error {
	http.HandleFunc("/v1/discover", discovery.MCollectiveDiscover)

	log.Infof("Starting pdbproxy version %s listener on %s:%d", Version, config.Listen, config.Port)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.Listen, config.Port), nil); err != nil {
		return err
	}

	return nil
}

func configure() error {
	kingpin.UsageTemplate(kingpin.SeparateOptionalFlagsUsageTemplate)
	kingpin.Version(Version)
	kingpin.CommandLine.Author("R.I. Pienaar <rip@devco.net>")
	kingpin.CommandLine.Help = "Service performing Choria discovery requests securely against PuppetDB"

	kingpin.Flag("debug", "enable debug logging").Short('d').BoolVar(&config.Debug)
	kingpin.Flag("listen", "address to bind to for client requests").Short('l').Default("0.0.0.0").OverrideDefaultFromEnvar("LISTEN").StringVar(&config.Listen)
	kingpin.Flag("port", "port to bind to for client requests").Short('p').Default("8080").OverrideDefaultFromEnvar("PORT").IntVar(&config.Port)
	kingpin.Flag("puppetdb-host", "PuppetDB host").Short('H').Default("puppet").OverrideDefaultFromEnvar("PUPPETDB_HOST").StringVar(&config.PuppetDBHost)
	kingpin.Flag("puppetdb-port", "PuppetDB port").Short('P').Default("8081").OverrideDefaultFromEnvar("PUPPETDB_PORT").IntVar(&config.PuppetDBPort)
	kingpin.Flag("ca", "Certificate Authority file").OverrideDefaultFromEnvar("CA").Required().ExistingFileVar(&config.Ca)
	kingpin.Flag("cert", "Public certificate file").OverrideDefaultFromEnvar("CERT").Required().ExistingFileVar(&config.Certificate)
	kingpin.Flag("key", "Public key file").OverrideDefaultFromEnvar("KEY").Required().ExistingFileVar(&config.PrivateKey)
	kingpin.Flag("logfile", "File to log to, STDOUT when not set").OverrideDefaultFromEnvar("LOGFILE").StringVar(&config.Logfile)

	kingpin.Parse()

	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})

	if config.Logfile != "" {
		file, err := os.OpenFile(config.Logfile, os.O_CREATE|os.O_WRONLY, 0666)

		if err == nil {
			log.SetOutput(file)
		} else {
			log.Fatalf("Cannot log to file %s: %s", config.Logfile, err.Error())
			os.Exit(1)
		}
	}

	if config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	discovery.SetConfig(config)

	return nil
}
