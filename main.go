package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	loads "github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	flags "github.com/jessevdk/go-flags"

	"github.com/choria-io/pdbproxy/discovery"
	"github.com/choria-io/pdbproxy/restapi"
	"github.com/choria-io/pdbproxy/restapi/operations"
	"gopkg.in/alecthomas/kingpin.v2"
)

const version = "0.0.1"

var config = discovery.Config{}

func main() {
	if err := configure(); err != nil {
		log.Fatalf("Could not configure pdbproxy: %s", err.Error())
		os.Exit(1)
	}

	server := newServer()
	defer server.Shutdown()

	if err := server.Serve(); err != nil {
		log.Fatalf("Could not start webserver: %s", err.Error())
	}
}

func setupHandlers(api *operations.PdbproxyAPI) {
	api.GetDiscoverHandler = operations.GetDiscoverHandlerFunc(
		func(params operations.GetDiscoverParams) middleware.Responder {
			return discovery.SwaggerDiscovery(params)
		})
}

func newServer() *restapi.Server {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")

	if err != nil {
		log.Fatalf("Could not initialize Swagger libraries: %s", err.Error())
		os.Exit(1)
	}

	api := operations.NewPdbproxyAPI(swaggerSpec)
	api.Logger = log.Printf

	server := restapi.NewServer(api)

	if config.Port != 0 {
		server.EnabledListeners = []string{"https", "http"}
	} else {
		server.EnabledListeners = []string{"https"}
	}

	server.Port = config.Port
	server.TLSPort = config.TLSPort
	server.TLSHost = config.Listen
	server.TLSCACertificate = flags.Filename(config.Ca)
	server.TLSCertificate = flags.Filename(config.Certificate)
	server.TLSCertificateKey = flags.Filename(config.PrivateKey)

	setupHandlers(api)

	return server
}

func configure() error {
	kingpin.UsageTemplate(kingpin.SeparateOptionalFlagsUsageTemplate)
	kingpin.Version(version)
	kingpin.CommandLine.Author("R.I. Pienaar <rip@devco.net>")
	kingpin.CommandLine.Help = "Service performing Choria discovery requests securely against PuppetDB"

	kingpin.Flag("debug", "enable debug logging").Short('d').BoolVar(&config.Debug)
	kingpin.Flag("listen", "address to bind to for client requests").Short('l').Default("0.0.0.0").OverrideDefaultFromEnvar("LISTEN").StringVar(&config.Listen)
	kingpin.Flag("port", "port to bind to for client requests").Short('p').Default("0").OverrideDefaultFromEnvar("PORT").IntVar(&config.Port)
	kingpin.Flag("tlsport", "port to bind to for client HTTPS requests").Default("8081").OverrideDefaultFromEnvar("TLSPORT").IntVar(&config.TLSPort)
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
