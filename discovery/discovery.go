package discovery

import (
	log "github.com/Sirupsen/logrus"

	"github.com/boltdb/bolt"
	"github.com/choria-io/pdbproxy/models"
	"github.com/choria-io/pdbproxy/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

// Config for the discovery service
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

// SetConfig allows callers to configure the service
func SetConfig(c Config) error {
	config = c

	if err := openDB(); err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

// BackupSets produces a BoltDB backup in a file called .bak
func BackupSets(params operations.GetBackupParams) middleware.Responder {
	set := Sets{DB: db}

	target := config.DBFile + ".bak"

	err := set.Backup(&target)

	if err == nil {
		return operations.NewGetBackupOK().WithPayload(&models.SuccessModel{Code: 200, Message: "Created backup file " + target})
	}

	return operations.NewGetBackupInternalServerError().WithPayload(&models.ErrorModel{Code: 500, Message: "Could not create backup: " + err.Error()})
}

// Discover does basic MCollective like discovery against PuppetDB
func Discover(params operations.GetDiscoverParams) middleware.Responder {
	provider := PuppetDB{}
	discovered, err := provider.Discover(params.Request)

	if err == nil {
		return operations.NewGetDiscoverOK().WithPayload(&models.DiscoverySuccessModel{Code: 200, Nodes: discovered})
	}

	return operations.NewGetDiscoverBadRequest().WithPayload(&models.ErrorModel{Code: 400, Message: "Could not discover nodes: " + err.Error()})
}

// CreateSet will make new sets
func CreateSet(params operations.PostSetParams) middleware.Responder {
	set := Sets{DB: db}

	if set.Exists(string(params.Set.Set)) {
		return operations.NewPostSetBadRequest().WithPayload(&models.ErrorModel{Code: 400, Message: "A set called " + string(params.Set.Set) + " already exist"})
	}

	err := set.Update(params.Set)

	if err == nil {
		return operations.NewPostSetOK().WithPayload(&models.SuccessModel{Code: 200, Message: "Created node set"})
	}

	return operations.NewPostSetInternalServerError().WithPayload(&models.ErrorModel{Code: 500, Message: "Could not create set: " + err.Error()})
}

// UpdateSet will update an existing set with new parameters
func UpdateSet(params operations.PutSetSetParams) middleware.Responder {
	set := Sets{DB: db}

	if !set.Exists(string(params.Set)) {
		return operations.NewPutSetSetNotFound()
	}

	err := set.Update(params.NewSet)

	if err == nil {
		return operations.NewPutSetSetOK().WithPayload(&models.SuccessModel{Code: 200, Message: "Updated node set"})
	}

	return operations.NewPutSetSetInternalServerError().WithPayload(&models.ErrorModel{Code: 500, Message: "Could not update set: " + err.Error()})
}

// DeleteSet will remove an already existing set
func DeleteSet(params operations.DeleteSetSetParams) middleware.Responder {
	set := Sets{DB: db}

	if !set.Exists(string(params.Set)) {
		return operations.NewDeleteSetSetNotFound()
	}

	err := set.Delete(params.Set)

	if err == nil {
		return operations.NewDeleteSetSetOK().WithPayload(&models.SuccessModel{Code: 200, Message: "Deleted node set"})
	}

	return operations.NewDeleteSetSetInternalServerError().WithPayload(&models.ErrorModel{Code: 500, Message: "Could not update set: " + err.Error()})
}

// GetSet retrieves the parameters of a set and optionally does a discovery
func GetSet(params operations.GetSetSetParams) middleware.Responder {
	set := Sets{DB: db}

	if !set.Exists(string(params.Set)) {
		return operations.NewGetSetSetNotFound()
	}

	answer, err := set.Get(params.Set)

	if err != nil {
		log.Errorf("Retrieving set %s failed: %s", params.Set, err.Error())
		return operations.NewGetSetSetInternalServerError().WithPayload(&models.ErrorModel{Code: 500, Message: "Retrieving set failed: " + err.Error()})
	}

	if *params.Discover {
		req := models.DiscoveryRequest{Query: *answer.Query}
		provider := PuppetDB{}
		discovered, err := provider.Discover(&req)

		if err == nil {
			answer.Nodes = discovered
		} else {
			return operations.NewGetSetSetInternalServerError().WithPayload(&models.ErrorModel{Code: 500, Message: "Could not discover set nodes: " + err.Error()})
		}
	}

	return operations.NewGetSetSetOK().WithPayload(answer)
}

// ListSets finds all known sets
func ListSets(params operations.GetSetsParams) middleware.Responder {
	set := Sets{DB: db}

	sets := set.Sets()

	return operations.NewGetSetsOK().WithPayload(&models.Sets{Code: 200, Sets: sets})
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
