package netbox

import (
	"log"
	"github.com/digitalocean/go-netbox/netbox"
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

// ProviderNetboxClient is a structure that contains the client connections
// necessary to interface with the Go-Netbox API 
type ProviderNetboxClient struct {
		client *addresses.Controller
}

func (c *Config) Client() (interface{}, error) {
	cfg = netbox.Config{
		AppID: c.AppID,
		Endpoint: c.Endpoint 
	}
	log.Printf("[DEBUG] Initializing Netbox controllers")
	sess := session.NewSession(cfg)
	// Create the Client
	client := netbox.NewNetboxWithAPIKey(cfg:Endpoint, cfg:AppID)

	// Validate that our connection is okay
	if err := c.ValidateConnections(client); err != nil {
		return nil, err
	}

	return &client, nil
}

// ValidateConnection ensures that we can connect to Netbox early, so that we
// do not fail in the middle of a TF run if it can be prevented.
func (c *Config) ValidateConnection(sc *netbox.NewNetboxWithAPIKey) error {
	rs, err := sc.Dcim.DcimRacksList(nil, nil)
	return err
}
