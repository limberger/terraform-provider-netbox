package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/go-openapi/strfmt"
	runtimeclient "github.com/go-openapi/runtime/client"
	"github.com/cmgreivel/go-netbox/netbox/client"
	"github.com/cmgreivel/go-netbox/netbox/client/ipam"
	"github.com/cmgreivel/go-netbox/netbox/models"
)

// Users: (not sure I see APIs for this)
// typhon-dev
//
// Can add, change, delete, view prefixes
// typhon-ippm-dev
// typhon-ippm-stage
// typhon-ippm-production
// typhon-ippm-test

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

type setupArgs struct {
	destroyAll *bool
	endpoint string
}

const authHeaderName = "Authorization"
const authHeaderFormat = "Token %v"

func NewNetboxWithAPIKey(host string, apiToken string) *client.NetBox {
	t := runtimeclient.New(host, client.DefaultBasePath, client.DefaultSchemes)
	t.DefaultAuthentication = runtimeclient.APIKeyAuth(authHeaderName, "header", fmt.Sprintf(authHeaderFormat, apiToken))
	return client.New(t, strfmt.Default)
}

func (c *Config) Client() (*client.NetBox, error) {
	log.Printf("[DEBUG] config.go Client() AppID: %s", c.AppID)
	log.Printf("[DEBUG] config.go Client() Endpoint: %s", c.Endpoint)
	cfg := Config{
		AppID:    c.AppID,
		Endpoint: c.Endpoint,
	}
	log.Printf("[DEBUG] Initializing Netbox controllers")
	// sess := session.NewSession(cfg)
	// Create the Client
	nb := NewNetboxWithAPIKey(cfg.Endpoint, cfg.AppID)

	// Validate that our connection is okay
	if err := c.ValidateConnection(nb); err != nil {
		log.Printf("[DEBUG] config.go Client() Error [%v]", err)
		return nil, err
	}
	return nb, nil
}

// ValidateConnection ensures that we can connect to Netbox early, so that we
// do not fail in the middle of a TF run if it can be prevented.
func (c *Config) ValidateConnection(nb *client.NetBox) error {
	log.Printf("[DEBUG] config.go ValidateConnection() with call to list DCIM Racks ")
	rs, err := nb.Dcim.DcimRacksList(nil, nil)
	log.Println(rs)
	return err
}

// Could make these methods of a setup class
func addPrefixes(nb *client.NetBox) error {
	log.Println("[DEBUG] in addVrfs()")
	pools_to_add := []string{
		"10.224.0.0/16",
		"10.225.0.0/16",
		"10.226.0.0/16",
		"10.227.0.0/16",
		"100.64.0.0/10",
		"172.16.0.0/12",
	}
	sdPool := "10.228.0.0/16"
	vpnClientPrefix := "172.16.6.0/24"

	standard_vrfs := []string{"pre-dev", "dev", "prod", "stage", "test"}

	addParm := ipam.NewIpamPrefixesCreateParams()
	// Initialize the common parameters for Prefix creation
	var data models.WritablePrefix
	data.IsPool = true
	data.Tags = []string{}
	// Just get all the VRFs
	foundVrfs, err := nb.Ipam.IpamVrfsList(nil, nil)
	if err != nil {
		return err
	}
	for _, v := range standard_vrfs {
		var vrfIdPtr *int64 = nil
		for _, f := range foundVrfs.Payload.Results {
			if v == *f.Name {
				vrfIdPtr = &f.ID
			}
		}
		if vrfIdPtr == nil {
			return errors.New(fmt.Sprintf("Could not find VRF %s", v))
		}
		log.Printf("[DEBUG] Creating prefixes for VRF %s", v)
		data.Vrf = vrfIdPtr
		for _, p := range pools_to_add {
			data.Prefix = &p
			addParm.SetData(&data)
			created, err := nb.Ipam.IpamPrefixesCreate(addParm, nil)
			if err != nil {
				return err
			}
			log.Printf("[DEBUG] Created Prefix %d:%s", created.Payload.ID, *created.Payload.Prefix)

			// If this is the 172, add the vpn_client_prefix now
			if p == "172.16.0.0/12" {
				var vcpData models.WritablePrefix
				vcpData.IsPool = false
				vcpData.Vrf = vrfIdPtr
				vcpData.Prefix = &vpnClientPrefix
				vcpData.Tags = []string{}
				addParm.SetData(&vcpData)
				vpnCreated, err := nb.Ipam.IpamPrefixesCreate(addParm, nil)
				if err != nil {
					return err
				}
				log.Printf("[DEBUG] Created VPN Client Prefix %d:%s",
					vpnCreated.Payload.ID, *vpnCreated.Payload.Prefix)
			}
		}

		if v == "pre-dev" {
			// We add a servicedelivery pool to the pre-dev VRF to avoid
			// exhausting the production servicedelivery prefixes
			data.Prefix = &sdPool
			addParm.SetData(&data)
			created, err := nb.Ipam.IpamPrefixesCreate(addParm, nil)
			if err != nil {
				return err
			}
			log.Printf("[DEBUG] Created pre-dev servicedelivery prefix %d:%s",
				created.Payload.ID, *created.Payload.Prefix)
		}
	}
	// Add the servicedelivery prefix in its own VRF
	var vrfIdPtr *int64 = nil
	for _, f := range foundVrfs.Payload.Results {
		if *f.Name == "servicedelivery" {
			vrfIdPtr = &f.ID
			break
		}
	}
	if vrfIdPtr == nil {
		return errors.New(fmt.Sprintf("Could not find servicedelivery VRF"))
	}
	log.Printf("[DEBUG] Creating prefixes for VRF servicedelivery")
	data.Vrf = vrfIdPtr
	data.Prefix = &sdPool
	addParm.SetData(&data)
	created, err := nb.Ipam.IpamPrefixesCreate(addParm, nil)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Created Prefix %d:%s", created.Payload.ID, *created.Payload.Prefix)

	return nil
}

func addVrfs(nb *client.NetBox) error {
	log.Println("[DEBUG] in addVrfs()")
	to_add := []string{"pre-dev", "dev", "prod", "stage", "test", "servicedelivery"}


	// Sanity check to make sure these don't exist already.
	// Netbox does not prevent VRFs with duplicate names, and this
	// will cause headaches for us.
	foundVrfs, err := nb.Ipam.IpamVrfsList(nil, nil)
	if err != nil {
		return err
	}
	for _, f := range foundVrfs.Payload.Results {
		for _, v := range to_add {
			if v == *f.Name {
				return errors.New(fmt.Sprintf("VRF %s already in NetBox. Exiting!", v))
			}
		}
	}

	addParm := ipam.NewIpamVrfsCreateParams()
	// Initialize the common parameters for VRF creation
	var data models.WritableVRF
	data.Tags = []string{}
	data.EnforceUnique = true
	for _, v := range to_add {
		data.Name = &v
		addParm.SetData(&data)
		created, err := nb.Ipam.IpamVrfsCreate(addParm, nil)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Created VRF %d:%s", created.Payload.ID, *created.Payload.Name)
	}
	return nil
}

func destroyAll(nb *client.NetBox) error {
	log.Println("[DEBUG] in destroyAll()")

	// Find all the prefixes
	foundPrefixes, err := nb.Ipam.IpamPrefixesList(nil, nil)
	if err != nil {
		return err
	}
	// found.Payload.Count
	// found.Payload.Results (list of *models.Prefix)
	//   ID int64
	//   Prefix *string
	//   Tags [] string

	log.Printf("Deleting %d Prefixes", *foundPrefixes.Payload.Count)
	delParm := ipam.NewIpamPrefixesDeleteParams()
	for _, p := range foundPrefixes.Payload.Results {
		log.Printf("Deleting %d: %s", p.ID, *p.Prefix)
		delParm.SetID(p.ID)
		_, err = nb.Ipam.IpamPrefixesDelete(delParm, nil)
		if err != nil {
			return err
		}
	}

	// Same with VRFs
	foundVrfs, err := nb.Ipam.IpamVrfsList(nil, nil)
	if err != nil {
		return err
	}
	log.Printf("Deleting %d Vrfs", *foundVrfs.Payload.Count)
	vrfsDelParm := ipam.NewIpamVrfsDeleteParams()
	for _, v := range foundVrfs.Payload.Results {
		log.Printf("Deleting %d: %s", v.ID, *v.Name)
		vrfsDelParm.SetID(v.ID)
		_, err = nb.Ipam.IpamVrfsDelete(vrfsDelParm, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseArgs() setupArgs {
	var args setupArgs

	args.destroyAll = flag.Bool("destroyAll", false, "If set will destroy all VRFs and prefixes in the system")
	flag.StringVar(&args.endpoint, "endpoint", "localhost:32777", "The NetBox endpoint.")

	flag.Parse()
	return args
}

func main() {
	fmt.Println("Hello")

	args := parseArgs()

	log.Printf("[DEBUG] destroyAll %t", *args.destroyAll)
	log.Printf("[DEBUG] endpoint %s", args.endpoint)
	config := Config{
		AppID:    "0123456789abcdef0123456789abcdef01234567",
		Endpoint: args.endpoint,
	}
	nb, err := config.Client()

	if err != nil {
		log.Fatalf("[FATAL] error creating client %v", err)
	}
	if *args.destroyAll {
		err  = destroyAll(nb)
		if err != nil {
			log.Fatalf("[FATAL] error destroying all %v", err)
		}
	}
	err = addVrfs(nb)
	if err != nil {
		log.Fatalf("[FATAL] error adding VRFs %v", err)
	}

	err = addPrefixes(nb)
	if err != nil {
		log.Fatalf("[FATAL] error adding VRFs %v", err)
	}
}
