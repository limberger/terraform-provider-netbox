package netbox

import (
	"errors"
	"log"

	// "errors"

	"github.com/digitalocean/go-netbox/netbox/client/ipam"
	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/digitalocean/go-netbox/netbox/client"
)

func dataSourceNetboxVlans() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNetboxVlanRead,
		Schema: dataSourceVlanSchema(),
	}
}

// Read will fetch the data of a resource.
func dataSourceNetboxVlanRead(d *schema.ResourceData, meta interface{}) error {
	//out := ipam.NewIPAMVlansListParams()

	var parm = ipam.NewIPAMVlansReadParams()
	switch {
	case d.Get("vid").(int) != 0:
		// func (a *Client) IPAMVlansRead(params *IPAMVlansReadParams, authInfo runtime.ClientAuthInfoWriter) (*IPAMVlansReadOK, error) {
		parm.SetID(int64(d.Get("vid").(int)))
		c := meta.(*ProviderNetboxClient).client
		//parms = ipam.NewIPAMVlansListParams()
		out, err := c.IPAM.IPAMVlansRead(parm, nil)
		log.Printf("- Executado...\n")
		if err == nil {
			d.Set("created", out.Payload.Created)
			d.Set("description", out.Payload.Description)
			d.Set("display-name", out.Payload.DisplayName)
			d.Set("group", out.Payload.Group)
			d.Set("vid", out.Payload.Vid)
			d.Set("last-updated", out.Payload.LastUpdated)
			d.Set("name", out.Payload.Name)
			d.Set("nested-role", out.Payload.Role)
			d.Set("nested-site", out.Payload.Site)
			d.Set("vlan-status", out.Payload.Status)
			d.Set("nested-tenant", out.Payload.Tenant)
			log.Print("\n")
		} else {
			log.Printf("erro na chamada do IPAMVlansList\n")
			log.Printf("Err: %v\n", err)
			log.Print("\n")
			return err
		}
	case d.Get("name").(string) != "":

		c := meta.(*ProviderNetboxClient).client
		parmsl := ipam.NewIPAMVlansListParams()
		parmsl.SetName(d.Get("name").(*string))
		out, err := c.IPAM.IPAMVlansList(parmsl, nil)
		log.Printf("- Executado...\n")
		if err == nil {
			if *out.Payload.Count == 0 {
				return errors.New("Name not found - need one of vid or name")
			} else if *out.Payload.Count > 1 {
				return errors.New("More than one vlan found with name " + d.Get("name").(string))
			}
			result := out.Payload.Results[0]
			d.Set("created", result.Created)
			d.Set("description", result.Description)
			d.Set("display-name", result.DisplayName)
			d.Set("group", result.Group)
			d.Set("vid", result.Vid)
			d.Set("last-updated", result.LastUpdated)
			d.Set("name", result.Name)
			d.Set("nested-role", result.Role)
			d.Set("nested-site", result.Site)
			d.Set("vlan-status", result.Status)
			d.Set("nested-tenant", result.Tenant)
			log.Print("\n")
		} else {
			log.Printf("erro na chamada do IPAMVlansList\n")
			log.Printf("Err: %v\n", err)
			log.Print("\n")
			return err
		}

	case d.Get("vid").(int) != 0 && d.Get("name").(string) == "":
		return errors.New("No valid combination of parameters found - need one of vid or name")
	}
	//	out := make([]addresses.Address, 1)
	//	var err error
	// We need to determine how to get the address. An ID search takes priority,
	// and after that addresss.
	// switch {
	// case d.Get("address_id").(int) != 0:
	// 	out[0], err = c.GetAddressByID(d.Get("address_id").(int))
	// 	if err != nil {
	// 		return err
	// 	}
	// case d.Get("ip_address").(string) != "":
	// 	out, err = c.GetAddressesByIP(d.Get("ip_address").(string))
	// 	if err != nil {
	// 		return err
	// 	}
	// case d.Get("subnet_id").(int) != 0 && (d.Get("description").(string) != "" || d.Get("hostname").(string) != "" || len(d.Get("custom_field_filter").(map[string]interface{})) > 0):
	// 	out, err = addressSearchInSubnet(d, meta)
	// 	if err != nil {
	// 		return err
	// 	}
	// default:
	// 	return errors.New("No valid combination of parameters found - need one of address_id, ip_address, or subnet_id and (description|hostname|custom_field_filter)")
	// }
	// if len(out) != 1 {
	// 	return errors.New("Your search returned zero or multiple results. Please correct your search and try again")
	// }
	// flattenAddress(out[0], d)
	// fields, err := c.GetAddressCustomFields(out[0].ID)
	// if err != nil {
	// 	return err
	// }
	// trimMap(fields)
	// if err := d.Set("custom_fields", fields); err != nil {
	// 	return err
	// }
	return nil
}

func bareVlanSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vid": &schema.Schema{
			Type: schema.TypeInt,
		},
		"name": &schema.Schema{
			Type: schema.TypeString,
		},
		"id_in": &schema.Schema{
			Type: schema.TypeString,
		},
		"q": &schema.Schema{
			Type: schema.TypeString,
		},
		"site_id": &schema.Schema{
			Type: schema.TypeString,
		},
		"site": &schema.Schema{
			Type: schema.TypeString,
		},
		"group_id": &schema.Schema{
			Type: schema.TypeString,
		},
		"group": &schema.Schema{
			Type: schema.TypeString,
		},
		"tenant_id": &schema.Schema{
			Type: schema.TypeString,
		},
		"tenant": &schema.Schema{
			Type: schema.TypeString,
		},
		"role_id": &schema.Schema{
			Type: schema.TypeString,
		},
		"role": &schema.Schema{
			Type: schema.TypeString,
		},
		"status": &schema.Schema{
			Type: schema.TypeString,
		},
		"tag": &schema.Schema{
			Type: schema.TypeString,
		},
		// Number of results to return per page.
		"limit": &schema.Schema{
			Type: schema.TypeInt,
		},
		// The initial index from which to return the results.
		"offset": &schema.Schema{
			Type: schema.TypeInt,
		},
	}
}

func resourceVlansSchema() map[string]*schema.Schema {
	s := bareVlanSchema()

	for k, v := range s {
		switch k {
		case "vid":
			v.Optional = true
			v.Computed = true
		case "name":
			v.Optional = true
			v.Computed = true
			//v.ConflictsWith = []string{"ip_address", "subnet_id", "description", "hostname", "custom_field_filter"}
		default:
			v.Computed = true
		}
	}
	// Add the remove_dns_on_delete item to the schema. This is a meta-parameter
	// that is not part of the API resource and exists to instruct PHPIPAM to
	// gracefully remove the address from its DNS integrations as well when it is
	// removed. The default on this option is true.
	s["remove_dns_on_delete"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	}
	return s
}

// dataSourceAddressSchema returns the schema for the phpipam_address data
// source. It sets the searchable fields and sets up the attribute conflicts
// between IP address and address ID. It also ensures that all fields are
// computed as well.
func dataSourceVlanSchema() map[string]*schema.Schema {
	s := bareVlanSchema()
	for k, v := range s {
		switch k {
		case "vid":
			v.Optional = true
			v.Computed = true
		case "name":
			v.Optional = true
			v.Computed = true

			//v.ConflictsWith = []string{"ip_address", "subnet_id", "description", "hostname", "custom_field_filter"}
		default:
			v.Computed = true
		}
	}
	// Add the custom_field_filter item to the schema. This is a meta-parameter
	// that allows searching for a custom field value in the data source.
	s["custom_field_filter"] = customFieldFilterSchema([]string{"vid"})

	return s
}
