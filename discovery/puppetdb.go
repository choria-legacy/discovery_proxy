package discovery

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/choria-io/pdbproxy/models"
)

type PuppetDB struct{}

type puppetDbResult struct {
	Certname    string
	Deactivated *string
}

// Discover nodes against PuppetDB
func (p PuppetDB) Discover(request *models.DiscoveryRequest) ([]string, error) {
	query := p.parseMCollectiveQuery(request)

	log.Debugf("Query: %s", query)

	if result, err := p.queryPuppetdb(url.QueryEscape(query)); err == nil {
		if discovered, err := p.extractCertnames(result); err == nil {
			log.Debugf("Discovered %d nodes", len(discovered))

			return discovered, nil
		} else {
			return []string{}, err
		}
	} else {
		return []string{}, err
	}
}

// Create a SSL context for use when communicating with PuppetDB
func (p PuppetDB) sslTransport() *http.Transport {
	cert, err := tls.LoadX509KeyPair(config.Certificate, config.PrivateKey)

	if err != nil {
		log.Fatal("Could not load certificate ", config.Certificate, " and key ", config.PrivateKey, ": ", err)
	}

	caCert, err := ioutil.ReadFile(config.Ca)

	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12,
	}

	tlsConfig.BuildNameToCertificate()

	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return transport
}

// Performs a PQL query against PuppetDB
func (p PuppetDB) queryPuppetdb(pql string) ([]byte, error) {
	client := &http.Client{Transport: p.sslTransport()}

	log.Debugf("Querying %s:%d", config.PuppetDBHost, config.PuppetDBPort)

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s:%d/pdb/query/v4?query=%s", config.PuppetDBHost, config.PuppetDBPort, pql), nil)
	req.Header.Set("Accept-Encoding", "*")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to query PuppetDB: %s", err.Error()))
		return []byte(""), err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return []byte(body), nil
	} else {
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			log.Print(string(body))
		}

		return []byte("{}"), errors.New(resp.Status)
	}
}

// Capitalize a Puppet resource like apache::vhost into the format Apache::Vhost as expected by PuppetDB
func (p PuppetDB) capitalizeResource(resource string) string {
	var parts = strings.Split(resource, "::")
	var caps = make([]string, len(parts))

	for i, class := range parts {
		caps[i] = strings.Title(class)
	}

	return strings.Join(caps, "::")
}

// Create a PQL case insensitive regex match
func (p PuppetDB) stringRegexi(needle string) string {
	var derived string
	var buffer bytes.Buffer

	if regexp.MustCompile(`^/(.+)/$`).MatchString(needle) {
		matches := regexp.MustCompile(`^/(.+)/$`).FindStringSubmatch(needle)
		derived = matches[1]
	} else {
		derived = needle
	}

	for idx := range derived {
		if regexp.MustCompile(`[[:alpha:]]`).MatchString(string(derived[idx])) {
			buffer.WriteString(fmt.Sprintf(`[%s%s]`, strings.ToUpper(string(derived[idx])), strings.ToLower(string(derived[idx]))))
		} else {
			buffer.WriteString(string(derived[idx]))
		}

	}

	return buffer.String()
}

// Create a PQL string that combines other query strings into one and extracts certname and deactivated
func (p PuppetDB) nodeQueryString(queries []string) string {
	var parts = make([]string, 0)

	for _, query := range queries {
		if query != "" {
			parts = append(parts, fmt.Sprintf(`(%s)`, query))
		}
	}

	return fmt.Sprintf(`nodes[certname, deactivated] { %s }`, strings.Join(parts, ` and `))
}

// Create a PQL query string to find nodes in a certain sub collective
func (p PuppetDB) discoverCollective(collective models.CollectiveFilter) string {
	return fmt.Sprintf(`certname in inventory[certname] { facts.mcollective.server.collectives.match("\d+") = "%s" }`, string(collective.Word))
}

// Create a PQL query string to find nodes that has a certain class or classes
func (p PuppetDB) discoverClasses(classes models.ClassesFilter) string {
	var queries = make([]string, 0)
	re := regexp.MustCompile(`^/(.+)/$`)

	for _, c := range classes {
		class := string(c)

		if re.MatchString(class) {
			matches := re.FindStringSubmatch(class)
			queries = append(queries, fmt.Sprintf(`resources {type = "Class" and title ~ "%s"}`, p.stringRegexi(matches[1])))
		} else {
			queries = append(queries, fmt.Sprintf(`resources {type = "Class" and title = "%s"}`, p.capitalizeResource(class)))
		}
	}

	if len(queries) > 0 {
		return strings.Join(queries, " and ")
	}

	return ""
}

// Create a PQL query string to find nodes with certain agents
func (p PuppetDB) discoverAgents(agents models.AgentsFilter) string {
	var queries = make([]string, 0)

	re := regexp.MustCompile(`^/(.+)/$`)

	for _, a := range agents {
		agent := string(a)

		if agent == "rpcutil" {
			var k = make(models.ClassesFilter, 0)
			k = append(k, models.WordOrRegex("mcollective::service"))

			queries = append(queries, p.discoverClasses(k))
		} else if re.MatchString(agent) {
			matches := re.FindStringSubmatch(string(agent))
			queries = append(queries, fmt.Sprintf(`resources {type = "File" and tag ~ "mcollective_agent_.*?%s.*?_server"}`, p.stringRegexi(matches[1])))
		} else {
			queries = append(queries, fmt.Sprintf(`resources {type = "File" and tag = "mcollective_agent_%s_server"}`, agent))
		}
	}

	if len(queries) > 0 {
		return strings.Join(queries, " and ")
	}

	return ""
}

// Create a PQL query string to find nodes with certain certnames
//
// A special identity like pql:<PQL QUERY> is accepted which will
// perform custom PQL queries, these queries must return just certnames
func (p PuppetDB) discoverIdentities(identities models.IdentitiesFilter) string {
	if len(identities) == 0 {
		return ""
	}

	var queries = make([]string, 0)
	pqlRe := regexp.MustCompile(`^pql:\s*(.+)$`)
	regexIdentRe := regexp.MustCompile(`^\/(.+)\/$`)

	for _, i := range identities {
		identity := string(i)

		if pqlRe.MatchString(identity) {
			matches := pqlRe.FindStringSubmatch(identity)
			queries = append(queries, fmt.Sprintf(`certname in %s`, matches[1]))
		} else if regexIdentRe.MatchString(identity) {
			matches := regexIdentRe.FindStringSubmatch(identity)
			queries = append(queries, fmt.Sprintf(`certname ~ "%s"`, p.stringRegexi(matches[1])))
		} else {
			queries = append(queries, fmt.Sprintf(`certname = "%s"`, identity))
		}
	}

	if len(queries) > 0 {
		return strings.Join(queries, " or ")
	}

	return ""
}

// Creates a PQL querty string to find facts with MCollective operators supported and mapped
func (p PuppetDB) discoverFacts(facts models.FactsFilter) string {
	var queries []string

	for _, fact := range facts {
		value := string(fact.Value)

		switch *fact.Operator {
		case "=~":
			queries = append(queries, fmt.Sprintf(`facts {name = "%s" and value ~ "%s"}`, fact.Fact, p.stringRegexi(value)))
		case "==":
			queries = append(queries, fmt.Sprintf(`facts {name = "%s" and value = "%s"}`, fact.Fact, value))
		case "!=":
			queries = append(queries, fmt.Sprintf(`facts {name = "%s" and !(value = "%s")}`, fact.Fact, value))
		case ">=", ">", "<=", "<":
			if regexp.MustCompile(`^\d+$`).MatchString(value) {
				queries = append(queries, fmt.Sprintf(`facts {name = "%s" and value %s %s}`, fact.Fact, *fact.Operator, value))
			} else {
				queries = append(queries, fmt.Sprintf(`facts {name = "%s" and value %s "%s"}`, fact.Fact, *fact.Operator, value))
			}
		}
	}

	if len(queries) > 0 {
		return strings.Join(queries, " and ")
	}

	return ""
}

// // Extract all active certnames from quety results
func (p PuppetDB) extractCertnames(discovered []byte) ([]string, error) {
	var results []puppetDbResult
	var nodes = make([]string, 0)

	if err := json.Unmarshal([]byte(discovered), &results); err != nil {
		log.Errorf("Could not parse PuppetDB result: %s", err.Error())
		return nodes, err
	}

	for _, node := range results {
		if node.Deactivated == nil {
			nodes = append(nodes, node.Certname)
		}
	}

	return nodes, nil
}

// Parse the incoming MCollective discovery request and turns it into a PQL query
func (p PuppetDB) parseMCollectiveQuery(query *models.DiscoveryRequest) string {
	if query.Query != "" {
		return (query.Query)
	}

	var queries []string

	if string(query.Collective.Word) != "" {
		queries = append(queries, p.discoverCollective(query.Collective))
	}

	if len(query.Classes) > 0 {
		queries = append(queries, p.discoverClasses(query.Classes))
	}

	if len(query.Agents) > 0 {
		queries = append(queries, p.discoverAgents(query.Agents))
	}

	if len(query.Facts) > 0 {
		queries = append(queries, p.discoverFacts(query.Facts))
	}

	if len(query.Identities) > 0 {
		queries = append(queries, p.discoverIdentities(query.Identities))
	}

	return p.nodeQueryString(queries)
}
