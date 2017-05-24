package choria

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/go-openapi/strfmt"

	apiclient "github.com/choria-io/pdbproxy/client"
	httptransport "github.com/go-openapi/runtime/client"
)

// Choria is a utilty encompasing mcollective and choria config and various utilities
type Choria struct {
	Config *MCollectiveConfig
	Sets *Sets
}

// Server is a representation of a network server host and port
type Server struct {
	Host string
	Port int
}

// New sets up a Choria with all its config loaded and so forth
func New(path string) (*Choria, error) {
	// TODO check SSL sanity

	c := Choria{}

	config, err := NewConfig(path)
	if err != nil {
		return &c, err
	}

	c.Config = config

	return &c, nil
}

// DiscoveryServer is the server configured as a discovery proxy
// TODO srv lookup
func (c *Choria) DiscoveryServer() (Server, error) {
	s := Server{
		Host: c.Config.Choria.DiscoveryHost,
		Port: c.Config.Choria.DiscoveryPort,
	}

	if !c.ProxiedDiscovery() {
		return s, errors.New("Proxy discovery is not enabled")
	}
	return s, nil
}

// ProxiedDiscovery determines if a client is configured for proxied discover
func (c *Choria) ProxiedDiscovery() bool {
	if c.Config.HasOption("plugin.choria.discovery_host") || c.Config.HasOption("plugin.choria.discovery_port") {
		return true
	}

	return c.Config.Choria.DiscoveryProxy
}

// Certname determines the choria certname
func (c *Choria) Certname() string {
	certname := c.Config.Identity

	currentUser, _ := user.Current()

	if currentUser.Uid != "0" {
		if u, ok := os.LookupEnv("USER"); ok {
			certname = fmt.Sprintf("%s.mcollective", u)
		}
	}

	if u, ok := os.LookupEnv("MCOLLECTIVE_CERTNAME"); ok {
		certname = u
	}

	return certname
}

// CAPath determines the path to the CA file
func (c *Choria) CAPath() (string, error) {
	ssl, err := c.SSLDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(ssl, "certs", "ca.pem"), nil
}

// ClientPrivateKey determines the location to the client cert
func (c *Choria) ClientPrivateKey() (string, error) {
	ssl, err := c.SSLDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(ssl, "private_keys", fmt.Sprintf("%s.pem", c.Certname())), nil
}

// ClientPublicCert determines the location to the client cert
func (c *Choria) ClientPublicCert() (string, error) {
	ssl, err := c.SSLDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(ssl, "certs", fmt.Sprintf("%s.pem", c.Certname())), nil
}

// SSLDir determines the AIO SSL directory
func (c *Choria) SSLDir() (string, error) {
	if c.Config.Choria.SSLDir != "" {
		return c.Config.Choria.SSLDir, nil
	}

	u, _ := user.Current()
	if u.Uid == "0" {
		path, err := c.PuppetSetting("ssldir")
		if err != nil {
			return "", err
		}

		return path, nil
	}

	return filepath.Join(u.HomeDir, ".puppetlabs", "etc", "puppet", "ssl"), nil
}

// PuppetSetting retrieves a config setting by shelling out to puppet apply --configprint
func (c *Choria) PuppetSetting(setting string) (string, error) {
	args := []string{"apply", "--configprint", setting}

	out, err := exec.Command("puppet", args...).Output()
	if err != nil {
		return "", err
	}

	return strings.Replace(string(out), "\n", "", -1), nil
}

// DiscoveryProxyClient is a client for the discovery REST service
func (c *Choria) DiscoveryProxyClient() (*apiclient.Pdbproxy, error) {
	server, err := c.DiscoveryServer()
	if err != nil {
		return nil, err
	}

	host := fmt.Sprintf("%s:%d", server.Host, server.Port)

	context, err := c.SSLContext()
	if err != nil {
		return nil, err
	}

	http := &http.Client{Transport: context}
	transport := httptransport.NewWithClient(host, apiclient.DefaultBasePath, apiclient.DefaultSchemes, http)

	return apiclient.New(transport, strfmt.NewFormats()), nil
}

// SSLContext creates a SSL context loaded with our certs and ca
func (c *Choria) SSLContext() (*http.Transport, error) {
	pub, _ := c.ClientPublicCert()
	pri, _ := c.ClientPrivateKey()
	ca, _ := c.CAPath()

	cert, err := tls.LoadX509KeyPair(pub, pri)

	if err != nil {
		return &http.Transport{}, errors.New("Could not load certificate " + pub + " and key " + pri + ": " + err.Error())
	}

	caCert, err := ioutil.ReadFile(ca)

	if err != nil {
		return &http.Transport{}, err
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

	return transport, nil
}
