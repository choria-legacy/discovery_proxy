package discovery

import (
	"github.com/choria-io/pdbproxy/models"
	"github.com/choria-io/pdbproxy/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

type Config struct {
	Listen       string
	Port         int
	TLSPort      int
	Logfile      string
	Debug        bool
	PuppetDBHost string
	PuppetDBPort int
	Certificate  string
	PrivateKey   string
	Ca           string
}

var config Config

func SetConfig(c Config) {
	config = c
}

func SwaggerDiscovery(params operations.GetDiscoverParams) middleware.Responder {
	provider := PuppetDB{}

	discovered, err := provider.Discover(params.Request)

	if err == nil {
		return operations.NewGetDiscoverOK().WithPayload(&models.DiscoverySuccessModel{Status: 200, Nodes: discovered})
	}

	return operations.NewGetDiscoverBadRequest().WithPayload(&models.ErrorModel{Status: 400, Message: "Could not discover nodes: " + err.Error()})
}
