package cmd

import (
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/choria-io/pdbproxy/discovery"
	"github.com/choria-io/pdbproxy/restapi"
	"github.com/choria-io/pdbproxy/restapi/operations"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	flags "github.com/jessevdk/go-flags"
)

type serverCommand struct {
	command

	config *discovery.Config
}

func (s *serverCommand) Setup() error {
	s.config = &discovery.Config{}

	s.Cmd = cli.app.Command("server", "Runs a Proxy Server")
	s.Cmd.Flag("listen", "address to bind to for client requests").Short('l').Default("0.0.0.0").OverrideDefaultFromEnvar("LISTEN").StringVar(&s.config.Listen)
	s.Cmd.Flag("port", "HTTP port to bind to for client requests").Short('p').Default("0").OverrideDefaultFromEnvar("PORT").IntVar(&s.config.Port)
	s.Cmd.Flag("tlsport", "HTTPS port to bind to for client requests").Default("8085").OverrideDefaultFromEnvar("TLSPORT").IntVar(&s.config.TLSPort)
	s.Cmd.Flag("puppetdb-host", "PuppetDB host").Short('H').Default("puppet").OverrideDefaultFromEnvar("PUPPETDB_HOST").StringVar(&s.config.PuppetDBHost)
	s.Cmd.Flag("puppetdb-port", "PuppetDB port").Short('P').Default("8081").OverrideDefaultFromEnvar("PUPPETDB_PORT").IntVar(&s.config.PuppetDBPort)
	s.Cmd.Flag("ca", "Certificate Authority file").OverrideDefaultFromEnvar("CA").Required().ExistingFileVar(&s.config.Ca)
	s.Cmd.Flag("cert", "Public certificate file").OverrideDefaultFromEnvar("CERT").Required().ExistingFileVar(&s.config.Certificate)
	s.Cmd.Flag("key", "Public key file").OverrideDefaultFromEnvar("KEY").Required().ExistingFileVar(&s.config.PrivateKey)
	s.Cmd.Flag("logfile", "File to log to, STDOUT when not set").OverrideDefaultFromEnvar("LOGFILE").StringVar(&s.config.Logfile)
	s.Cmd.Flag("db", "Path to the database file to write").OverrideDefaultFromEnvar("DB").Required().StringVar(&s.config.DBFile)

	return nil
}

func (s *serverCommand) Run() error {
	if s.config.Logfile != "" {
		file, err := os.OpenFile(s.config.Logfile, os.O_CREATE|os.O_WRONLY, 0666)

		if err == nil {
			log.SetOutput(file)
		} else {
			log.Fatalf("Cannot log to file %s: %s", s.config.Logfile, err.Error())
			return err
		}
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if err := discovery.SetConfig(*s.config); err != nil {
		return err
	}

	server, err := s.server()
	if err != nil {
		return err
	}

	defer server.Shutdown()

	if err := server.Serve(); err != nil {
		log.Fatalf("Could not start webserver: %s", err.Error())
		return err
	}

	return nil
}

func (s *serverCommand) server() (*restapi.Server, error) {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")

	if err != nil {
		log.Fatalf("Could not initialize Swagger libraries: %s", err.Error())
		return nil, err
	}

	api := operations.NewPdbproxyAPI(swaggerSpec)
	api.Logger = log.Printf

	server := restapi.NewServer(api)

	if s.config.Port != 0 {
		server.EnabledListeners = []string{"https", "http"}
	} else {
		server.EnabledListeners = []string{"https"}
	}

	server.Port = s.config.Port
	server.TLSPort = s.config.TLSPort
	server.TLSHost = s.config.Listen
	server.TLSCACertificate = flags.Filename(s.config.Ca)
	server.TLSCertificate = flags.Filename(s.config.Certificate)
	server.TLSCertificateKey = flags.Filename(s.config.PrivateKey)

	setupHandlers(api)

	return server, nil
}

func setupHandlers(api *operations.PdbproxyAPI) {
	api.GetDiscoverHandler = operations.GetDiscoverHandlerFunc(
		func(params operations.GetDiscoverParams) middleware.Responder {
			return discovery.Discover(params)
		})

	api.PostSetHandler = operations.PostSetHandlerFunc(
		func(params operations.PostSetParams) middleware.Responder {
			return discovery.CreateSet(params)
		})

	api.PutSetSetHandler = operations.PutSetSetHandlerFunc(
		func(params operations.PutSetSetParams) middleware.Responder {
			return discovery.UpdateSet(params)
		})

	api.DeleteSetSetHandler = operations.DeleteSetSetHandlerFunc(
		func(params operations.DeleteSetSetParams) middleware.Responder {
			return discovery.DeleteSet(params)
		})

	api.GetSetSetHandler = operations.GetSetSetHandlerFunc(
		func(params operations.GetSetSetParams) middleware.Responder {
			return discovery.GetSet(params)
		})

	api.GetSetsHandler = operations.GetSetsHandlerFunc(
		func(params operations.GetSetsParams) middleware.Responder {
			return discovery.ListSets(params)
		})

	api.GetBackupHandler = operations.GetBackupHandlerFunc(
		func(params operations.GetBackupParams) middleware.Responder {
			return discovery.BackupSets(params)
		})
}
