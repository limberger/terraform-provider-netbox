package netbox

import (
	"log"
	"os"
	"strings"

	api "github.com/digitalocean/go-netbox/netbox"
	"github.com/digitalocean/go-netbox/netbox/client"
)

// The default PHPIPAM API endpoint.
const defaultAPIAddress = "http://localhost/api"

// Config provides the configuration for the NETBOX providerr.
type Config struct {
	// The application ID required for API requests. This needs to be created
	// in the NETBOX console. It can also be supplied via the NETBOX_APP_ID
	// environment variable.
	AppID string

	// The API endpoint. This defaults to http://localhost/api, and can also be
	// supplied via the NETBOX_ENDPOINT_ADDR environment variable.
	Endpoint string

	// Optional timeout to API calls
	Timeout string
}

type ProviderNetboxClient struct {
	client *client.NetBox
}

// ProviderNetboxClient is a structure that contains the client connections
// necessary to interface with the Go-Netbox API
//type ProviderNetboxClient struct {
//		client *Client
//}

// DefaultConfigProvider supplies a default configuration:
//  * AppID defaults to PHPIPAM_APP_ID, if set, otherwise empty
//  * Endpoint defaults to PHPIPAM_ENDPOINT_ADDR, otherwise http://localhost/api
//  * Password defaults to PHPIPAM_PASSWORD, if set, otherwise empty
//  * Username defaults to PHPIPAM_USER_NAME, if set, otherwise empty
//
// This essentially loads an initial config state for any given
// API service.
func DefaultConfigProvider() Config {
	env := os.Environ()
	cfg := Config{
		Endpoint: defaultAPIAddress,
	}

	for _, v := range env {
		d := strings.Split(v, "=")
		switch d[0] {
		case "NETBOX_APP_ID":
			cfg.AppID = d[1]
		case "NETBOX_ENDPOINT_ADDR":
			cfg.Endpoint = d[1]
		case "NETBOX_TIMEOUT":
			cfg.Timeout = d[1]
		}
	}
	return cfg
}

func (c *Config) Client() (interface{}, error) {
	log.Printf("[DEBUG] config.go Client() AppID: %s", c.AppID)

	cfg := Config{
		AppID:    c.AppID,
		Endpoint: c.Endpoint,
		Timeout:  c.Timeout,
	}

	// Se o endpoint for vazio ou o AppId for vazio,
	// Busca configuração Default (ENV)
	if c.Endpoint == "" || c.AppID == "" {
		cfg = DefaultConfigProvider()
	}

	log.Printf("[DEBUG] Initializing Netbox controllers")
	// sess := session.NewSession(cfg)
	// Create the Client
	cli := api.NewNetboxWithAPIKey(cfg.Endpoint, cfg.AppID)

	// Validate that our connection is okay
	if err := c.ValidateConnection(cli); err != nil {
		log.Printf("[DEBUG] config.go Client() Erro")
		return nil, err
	}
	cs := ProviderNetboxClient{
		client: cli,
	}
	return &cs, nil
}

// ValidateConnection ensures that we can connect to Netbox early, so that we
// do not fail in the middle of a TF run if it can be prevented.
func (c *Config) ValidateConnection(sc *client.NetBox) error {
	log.Printf("[DEBUG] config.go ValidateConnection() validando ")
	rs, err := sc.Dcim.DcimRacksList(nil, nil)
	log.Println(rs)
	return err
}
