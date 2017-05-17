package discovery

import (
	"bytes"
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
	Fact     string `json:"fact"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type mcoFilter struct {
	Facts      []factFilter `json:"facts"`
	Classes    []string     `json:"classes"`
	Agents     []string     `json:"agents"`
	Identities []string     `json:"identities"`
	Collective string       `json:"collective"`
	Query      string       `json:"query"`
	NodeSet    string       `json:"node_set"`
}

type errorModel struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

type successModel struct {
	Status int      `json:"status"`
	Nodes  []string `json:"nodes"`
}

var config Config
var logger log.Entry

func SetConfig(c Config) {
	config = c
}

func failRequest(status int, message string, detail string, response http.ResponseWriter) {
	err := errorModel{status, message, detail}

	response.WriteHeader(err.Status)
	data, _ := json.Marshal(err)
	response.Write([]byte(data))
}

func MCollectiveDiscover(response http.ResponseWriter, request *http.Request) {
	logger := log.WithField("remote", request.RemoteAddr)

	logger.Infof("serving request")

	req, err := newRequest(request.Body)
	if err != nil {
		failRequest(http.StatusBadRequest, "Could not parse incoming request data: "+err.Error(), err.Error(), response)
		return
	}

	provider := PuppetDB{Log: logger}

	if discovered, err := provider.Discover(req); err == nil {
		if data, err := json.Marshal(successModel{Status: 200, Nodes: discovered}); err == nil {
			response.Write(data)
		} else {
			failRequest(http.StatusBadRequest, "Failed to json encode results: "+err.Error(), err.Error(), response)
		}
	} else {
		failRequest(http.StatusBadRequest, "Discovery failed: "+err.Error(), err.Error(), response)
	}
}

func newRequest(query io.Reader) (mcoFilter, error) {
	req := mcoFilter{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(query)

	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		return req, errors.New("Could not decode JSON request: " + err.Error())
	}

	return req, nil
}
