package netbox

import (
	"fmt"
	"log"

	"github.com/go-openapi/strfmt"
	runtimeclient "github.com/go-openapi/runtime/client"
	"github.com/cmgreivel/go-netbox/netbox/client"
)

// Config provides the configuration for the NETBOX providerr.
type Config struct {
	// The application ID required for API requests. This needs to be created
	// in the NETBOX console. It can also be supplied via the NETBOX_APP_ID
	// environment variable.
	AppID string

	// The API endpoint. This defaults to http://localhost/api, and can also be
	// supplied via the NETBOX_ENDPOINT_ADDR environment variable.
	Endpoint string
}

type ProviderNetboxClient struct {
	client        *client.NetBox
	configuration Config
}

// ProviderNetboxClient is a structure that contains the client connections
// necessary to interface with the Go-Netbox API
//type ProviderNetboxClient struct {
//		client *Client
//}

const authHeaderName = "Authorization"
const authHeaderFormat = "Token %v"

func NewNetboxWithAPIKey(host string, apiToken string) *client.NetBox {
	t := runtimeclient.New(host, client.DefaultBasePath, client.DefaultSchemes)
	t.DefaultAuthentication = runtimeclient.APIKeyAuth(authHeaderName, "header", fmt.Sprintf(authHeaderFormat, apiToken))
	return client.New(t, strfmt.Default)
}

func (c *Config) Client() (interface{}, error) {
	log.Printf("[DEBUG] config.go Client() AppID: %s", c.AppID)
	log.Printf("[DEBUG] config.go Client() Endpoint: %s", c.Endpoint)
	cfg := Config{
		AppID:    c.AppID,
		Endpoint: c.Endpoint,
	}
	log.Printf("[DEBUG] Initializing Netbox controllers")
	// sess := session.NewSession(cfg)
	// Create the Client
	cli := NewNetboxWithAPIKey(cfg.Endpoint, cfg.AppID)

	// Validate that our connection is okay
	if err := c.ValidateConnection(cli); err != nil {
		log.Printf("[DEBUG] config.go Client() Error")
		return nil, err
	}
	cs := ProviderNetboxClient{
		client:        cli,
		configuration: cfg,
	}
	log.Printf("[DEBUG] config.go returning ProviderNetboxClient")
	return &cs, nil
}

// ValidateConnection ensures that we can connect to Netbox early, so that we
// do not fail in the middle of a TF run if it can be prevented.
func (c *Config) ValidateConnection(sc *client.NetBox) error {
	log.Printf("[DEBUG] ValidateConnection() with call to list DCIM Racks ")
	rs, err := sc.Dcim.DcimRacksList(nil, nil)
	log.Println(rs)
	return err
}
