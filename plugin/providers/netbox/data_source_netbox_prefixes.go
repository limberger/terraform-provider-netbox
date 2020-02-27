package netbox

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	// "errors"

	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/netbox-community/go-netbox/netbox/client"
)

func dataSourceNetboxPrefixes() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNetboxPrefixesRead,
		Schema: dataSourcePrefixesSchema(),
	}
}

// Read will fetch the data of a resource.
func dataSourceNetboxPrefixesRead(d *schema.ResourceData, meta interface{}) error {
	//out := ipam.NewIPAMPrefixesListParams()
	log.Printf("data_source_netbox_prefixes.go dataSourceNetboxPrefixesRead ............ ")
	switch {
	// Pega por prefix_id
	case d.Get("prefixes_id").(int) != 0:
		var parm = ipam.NewIPAMPrefixesReadParams()
		parm.SetID(int64(d.Get("prefixes_id").(int)))
		//(&&meta).IPAM.IPAMPrefixesRead(parm,nil)

		c := meta.(*ProviderNetboxClient).client
		log.Printf("Obtive o client\n")
		//parms = ipam.NewIPAMPrefixesListParams()
		out, err := c.IPAM.IPAMPrefixesRead(parm, nil)
		log.Printf("- Executado...\n")
		if err == nil {

			d.SetId(strconv.FormatInt(out.Payload.ID, 10)) // Sempre setar o ID
			d.Set("created", out.Payload.Created.String())
			d.Set("description", out.Payload.Description)
			d.Set("family", out.Payload.Family)
			d.Set("is_pool", out.Payload.IsPool)
			d.Set("prefix", out.Payload.Prefix)
			d.Set("last_updated", out.Payload.LastUpdated)
			d.Set("vlan_vid", *out.Payload.Vlan.Vid)
			log.Print("\n")
		} else {
			log.Printf("erro na chamada do IPAMPrefixesList\n")
			log.Printf("Err: %v\n", err)
			log.Print("\n")
			return err
		}
		// Pega por prefix.vlan.vid
	case d.Get("vlan_vid").(int) != 0:
		var parml = ipam.NewIPAMPrefixesListParams()
		vlan_vid := float64(d.Get("vlan_vid").(int))
		parml.SetVlanVid(&vlan_vid)
		c := meta.(*ProviderNetboxClient).client
		out, err := c.IPAM.IPAMPrefixesList(parml, nil)
		if err == nil {
			if *out.Payload.Count == 0 {
				return errors.New("Prefix not found")
			} else if *out.Payload.Count > 1 {
				return errors.New(fmt.Sprintf("More than one Prefix found with vid %v\n", d.Get("vlanvid").(int)))
			}
			result := out.Payload.Results[0]
			d.SetId(strconv.FormatInt(result.ID, 10)) // Sempre setar o ID
			d.Set("created", result.Created.String())
			d.Set("custom_fields", result.CustomFields)
			d.Set("description", result.Description)
			d.Set("family", result.Family)
			d.Set("is_pool", result.IsPool)
			d.Set("prefix", result.Prefix)
			d.Set("last_updated", result.LastUpdated)
			d.Set("vlan_Vid", *result.Vlan.Vid)
			log.Print("\n")
		} else {
			log.Printf("erro na chamada do IPAMPrefixesList\n")
			log.Printf("Err: %v\n", err)
			log.Print("\n")
			return err
		}
	default:
		return errors.New("No valid combination of parameters found - prefix_id or vlan_vid")
	}
	return nil
}

func barePrefixesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"prefixes_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"created": &schema.Schema{
			Type: schema.TypeString,
		},
		"description": &schema.Schema{
			Type: schema.TypeString,
		},
		"prefix": &schema.Schema{
			Type: schema.TypeString,
		},
		"family": &schema.Schema{
			Type: schema.TypeString,
		},
		"vlan": &schema.Schema{
			Type: schema.TypeMap,
		},
		"is_pool": &schema.Schema{
			Type: schema.TypeBool,
		},
		"last_updated": &schema.Schema{
			Type: schema.TypeString,
		},
		"vlan_vid": &schema.Schema{
			Type: schema.TypeInt,
		},
	}
}

func resourcePrefixesSchema() map[string]*schema.Schema {
	s := barePrefixesSchema()

	for k, v := range s {
		switch k {
		case "prefixes_id":
			v.Optional = true
			v.ConflictsWith = []string{"vlan_vid"}
		case "prefix":
			v.Optional = true
		case "created":
			v.Optional = true
		case "vlan_vid":
			v.Optional = true
			v.ConflictsWith = []string{"prefixes_id"}
		default:
			v.Computed = true
		}
	}
	// Add the remove_dns_on_delete item to the schema. This is a meta-parameter
	// that is not part of the API resource and exists to instruct NETBOX to
	// gracefully remove the address from its DNS integrations as well when it is
	// removed. The default on this option is true.
	s["remove_dns_on_delete"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	}
	return s

}

// dataSourceAddressSchema returns the schema for the dataSourceNetboxPrefixes data
// source. It sets the searchable fields and sets up the attribute conflicts
// between IP address and address ID. It also ensures that all fields are
// computed as well.
func dataSourcePrefixesSchema() map[string]*schema.Schema {
	s := barePrefixesSchema()
	for k, v := range s {
		switch k {
		case "prefixes_id":
			v.Optional = true
		case "vlan_vid":
			v.Optional = true
		case "prefix":
			v.Optional = true
		case "created":
			v.Optional = true
			//v.ConflictsWith = []string{"ip_address", "subnet_id", "description", "hostname", "custom_field_filter"}
		default:
			v.Computed = true
		}
	}
	// Add the custom_field_filter item to the schema. This is a meta-parameter
	// that allows searching for a custom field value in the data source.
	s["custom_field_filter"] = customFieldFilterSchema([]string{"prefixes_id"})

	return s
}
