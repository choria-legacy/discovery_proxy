package discovery

import (
	log "github.com/Sirupsen/logrus"

	"github.com/boltdb/bolt"
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
	DBFile       string
}

var config Config
var db *bolt.DB

func SetConfig(c Config) error {
	config = c

	if err := openDB(); err != nil {
		log.Fatalf("Could not open db %s: %s", config.DBFile, err.Error())
		return err
	}

	return nil
}

func Discover(params operations.GetDiscoverParams) middleware.Responder {
	provider := PuppetDB{}
	discovered, err := provider.Discover(params.Request)

	if err == nil {
		return operations.NewGetDiscoverOK().WithPayload(&models.DiscoverySuccessModel{Status: 200, Nodes: discovered})
	}

	return operations.NewGetDiscoverBadRequest().WithPayload(&models.ErrorModel{Status: 400, Message: "Could not discover nodes: " + err.Error()})
}

func UpdateSet(params operations.PostSetParams) middleware.Responder {
	set := Sets{DB: db}

	err := set.Update(params.Set)

	if err == nil {
		return operations.NewPostSetOK().WithPayload(&models.SuccessModel{Status: 200, Detail: "Updated node set"})
	}

	return operations.NewPostSetBadRequest().WithPayload(&models.ErrorModel{Status: 400, Message: "Could not update set: " + err.Error()})

}

func DeleteSet(params operations.DeleteSetSetParams) middleware.Responder {
	set := Sets{DB: db}

	err := set.Delete(params.Set)

	if err == nil {
		return operations.NewPostSetOK().WithPayload(&models.SuccessModel{Status: 200, Detail: "Deleted node set"})
	}

	return operations.NewPostSetBadRequest().WithPayload(&models.ErrorModel{Status: 400, Message: "Could not update set: " + err.Error()})

}

func GetSet(params operations.GetSetSetParams) middleware.Responder {
	set := Sets{DB: db}

	answer, err := set.GetSet(params.Set)

	if err != nil {
		log.Infof("Searching for set %s failed: %s", params.Set, err.Error())
		return operations.NewGetSetSetNotFound()
	}

	if answer.Set == "" {
		return operations.NewGetSetSetNotFound()
	}

	if *params.Discover {
		req := models.DiscoveryRequest{Query: *answer.Query}
		provider := PuppetDB{}
		discovered, err := provider.Discover(&req)

		if err == nil {
			answer.Nodes = discovered
		} else {
			return operations.NewGetSetSetBadRequest().WithPayload(&models.ErrorModel{Status: 400, Message: "Could not discover set nodes: " + err.Error()})
		}
	}

	return operations.NewGetSetSetOK().WithPayload(answer)
}

func openDB() error {
	bdb, err := bolt.Open(config.DBFile, 0600, nil)

	if err != nil {
		log.Fatalf("Could not open db %s: %s", config.DBFile, err.Error())
		return err
	}

	db = bdb

	return nil
}
