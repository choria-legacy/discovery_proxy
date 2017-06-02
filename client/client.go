package client

import (
	"fmt"
	"net/http"

	"github.com/choria-io/discovery_proxy/choria"
	"github.com/go-openapi/strfmt"

	httptransport "github.com/go-openapi/runtime/client"
)

// NewDiscoveryProxyClient is a client for the discovery REST service
func NewDiscoveryProxyClient(c *choria.Choria) (*DiscoveryProxy, error) {
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
	transport := httptransport.NewWithClient(host, DefaultBasePath, DefaultSchemes, http)

	return New(transport, strfmt.NewFormats()), nil
}
