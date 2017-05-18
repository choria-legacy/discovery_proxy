package discovery

import (
	"os"
	"testing"

	"github.com/choria-io/pdbproxy/models"
	"github.com/stretchr/testify/assert"
)

var provider = PuppetDB{}

func TestStringRegexi(t *testing.T) {
	assert.Equal(t, provider.stringRegexi("/test123/"), "[Tt][Ee][Ss][Tt]123", "should be equal")
	assert.Equal(t, provider.stringRegexi("test123"), "[Tt][Ee][Ss][Tt]123", "should be equal")
}

func TestNodeQueryString(t *testing.T) {
	var queries = make([]string, 0)

	queries = append(queries, provider.discoverClasses([]models.WordOrRegex{models.WordOrRegex("klass")}))
	queries = append(queries, provider.discoverAgents([]models.WordOrRegex{models.WordOrRegex("puppet")}))

	pql := provider.nodeQueryString(queries)
	expected := `nodes[certname, deactivated] { (resources {type = "Class" and title = "Klass"}) and (resources {type = "File" and tag = "mcollective_agent_puppet_server"}) }`

	assert.Equal(t, expected, pql, "should be equal")
}

func TestDiscoverCollective(t *testing.T) {
	expected := `certname in inventory[certname] { facts.mcollective.server.collectives.match("\d+") = "mcollective" }`
	assert.Equal(t, provider.discoverCollective(models.CollectiveFilter{models.Word("mcollective")}), expected, "should be equal")
}

func TestDiscoverClasses(t *testing.T) {
	pql := provider.discoverClasses([]models.WordOrRegex{models.WordOrRegex("klass"), models.WordOrRegex("/regex/")})
	expected := `resources {type = "Class" and title = "Klass"} and resources {type = "Class" and title ~ "[Rr][Ee][Gg][Ee][Xx]"}`

	assert.Equal(t, expected, pql, "should be equal")
}

func TestDiscoverAgents(t *testing.T) {
	pql := provider.discoverAgents([]models.WordOrRegex{models.WordOrRegex("puppet"), models.WordOrRegex("rpcutil"), models.WordOrRegex("/pup/")})
	expected := `resources {type = "File" and tag = "mcollective_agent_puppet_server"} and resources {type = "Class" and title = "Mcollective::Service"} and resources {type = "File" and tag ~ "mcollective_agent_.*?[Pp][Uu][Pp].*?_server"}`

	assert.Equal(t, expected, pql, "should be equal")
}

func TestDiscoverIdentities(t *testing.T) {
	pql := provider.discoverIdentities([]models.WordOrRegex{models.WordOrRegex("string"), models.WordOrRegex("/regex/"), models.WordOrRegex("pql: nodes[certname] { }")})
	expected := `certname = "string" or certname ~ "[Rr][Ee][Gg][Ee][Xx]" or certname in nodes[certname] { }`

	assert.Equal(t, expected, pql, "should be equal")
}

func TestDiscoverFacts(t *testing.T) {
	ops := []string{"=~", "==", "!=", ">=", "<="}
	facts := models.FactsFilter{}
	queries := []models.FactFilter{
		models.FactFilter{Fact: models.Word("re"), Operator: &ops[0], Value: models.WordOrRegex("revalue")},
		models.FactFilter{Fact: models.Word("eq"), Operator: &ops[1], Value: models.WordOrRegex("eqvalue")},
		models.FactFilter{Fact: models.Word("ne"), Operator: &ops[2], Value: models.WordOrRegex("nevalue")},
		models.FactFilter{Fact: models.Word("ge"), Operator: &ops[3], Value: models.WordOrRegex("1")},
		models.FactFilter{Fact: models.Word("lt"), Operator: &ops[4], Value: models.WordOrRegex("ltvalue")},
	}

	for i := range queries {
		facts = append(facts, &queries[i])
	}

	pql := provider.discoverFacts(facts)
	expected := `facts {name = "re" and value ~ "[Rr][Ee][Vv][Aa][Ll][Uu][Ee]"} and facts {name = "eq" and value = "eqvalue"} and facts {name = "ne" and !(value = "nevalue")} and facts {name = "ge" and value >= 1} and facts {name = "lt" and value <= "ltvalue"}`

	assert.Equal(t, expected, pql, "should be equal")
}

func TestExtractCertNames(t *testing.T) {
	fixture, err := readFixture("testdata/ok_pdb_results.json")

	if assert.Nil(t, err, "failed to load fixture") {
		nodes, err := provider.extractCertnames(fixture)

		if assert.Nil(t, err, "extracting names failed") {
			expected := []string{"dev1.devco.net", "dev4.devco.net", "dev3.devco.net", "dev2.devco.net", "dev5.devco.net"}
			assert.Equal(t, expected, nodes, "should be equal")
		}
	}
}

func readFixture(file string) ([]byte, error) {
	if f, err := os.Open(file); err == nil {
		stat, _ := f.Stat()
		fixture := make([]byte, stat.Size())

		if _, err = f.Read(fixture); err == nil {
			return fixture, nil
		} else {
			return []byte{}, err
		}
	} else {
		return []byte{}, err
	}
}
