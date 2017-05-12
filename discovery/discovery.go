package discovery

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type Config struct {
	Listen       string
	Port         int
	Logfile      string
	Debug        bool
	PuppetDBHost string
	PuppetDBPort int
	Certificate  string
	PrivateKey   string
	Ca           string
}

type factFilter struct {
	Fact     string `json: "fact"`
	Operator string `json: "operator"`
	Value    string `json: "value"`
}

type mcoFilter struct {
	Facts      []factFilter `json:"facts"`
	Classes    []string     `json:"classes"`
	Agents     []string     `json:"agents"`
	Identities []string     `json:"identities"`
	Collective string       `json: "collective"`
	Query      string       `json: "query"`
	NodeSet    string       `json:"node_set"`
}

var config Config
var logger log.Entry

func SetConfig(c Config) {
	config = c
}

func MCollectiveDiscover(response http.ResponseWriter, request *http.Request) {
	logger := log.WithField("remote", request.RemoteAddr)

	logger.Infof("serving request")

	req, err := newRequest(request.Body)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Could not parse incoming request data: " + err.Error()))
		return
	}

	provider := PuppetDB{Log: logger}

	if discovered, err := provider.Discover(req); err == nil {
		if data, err := json.Marshal(discovered); err == nil {
			response.Write(data)
		} else {
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte("failed to json encode results: " + err.Error()))
		}
	} else {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Discovery failed: " + err.Error()))
		return
	}
}

func newRequest(query io.Reader) (mcoFilter, error) {
	req := mcoFilter{}
	req.Facts = []factFilter{}
	req.Classes = []string{}
	req.Agents = []string{}
	req.Identities = []string{}
	req.Collective = ""
	req.Query = ""
	req.NodeSet = ""

	if err := json.NewDecoder(query).Decode(&req); err != nil {
		return req, errors.New("Could not decode JSON request: " + err.Error())
	}

	return req, nil
}
