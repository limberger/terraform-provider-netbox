package netbox

import (
	"log"
	"reflect"

	// "errors"

	"github.com/digitalocean/go-netbox/netbox/client/ipam"
	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/digitalocean/go-netbox/netbox/client"
)

func dataSourceNetboxPrefixes() *schema.Resource {
	log.Println("[DEBUG] jp data_source_netbox_prefixes.go dataSourceNetboxPrefixes()")
	return &schema.Resource{
		Read:   dataSourceNetboxPrefixesRead,
		Schema: dataSourcePrefixesSchema(),
	}
}

// Read will fetch the data of a resource.
func dataSourceNetboxPrefixesRead(d *schema.ResourceData, meta interface{}) error {
	//out := ipam.NewIPAMPrefixesListParams()

	log.Printf("[DEBUG] jp data_source_netbox_prefixes.go : dataSourceNetboxPrefixesRead %v\n", d)
	switch {
	case d.Get("prefixes_id").(int) != 0:
		log.Println("É um prefixo ...")
		log.Printf("data_source_netbox_prefixes.go dataSourceNetboxPrefixesRead - Prefixo: %d\n", d.Get("prefixes_id").(int))
		var parm = ipam.NewIPAMPrefixesReadParams()
		log.Println("Criei o parm")
		parm.SetID(int64(d.Get("prefixes_id").(int)))
		log.Println("Setei o parm")
		log.Printf("Tipo do parm [meta] : %s", reflect.TypeOf(meta))
		//(&&meta).IPAM.IPAMPrefixesRead(parm,nil)

		c := meta.(*ProviderNetboxClient).client
		log.Printf("Obtive o client\n")
		//parms = ipam.NewIPAMPrefixesListParams()
		out, err := c.IPAM.IPAMPrefixesRead(parm, nil)
		log.Printf("- Executado...\n")
		if err == nil {
			log.Printf("Ok na chamada do IPAMPrefixesList\n")
			log.Printf("Out: %v\n", out)
			log.Printf("Created: %v\n", &out.Payload.Created)
			d.Set("created", out.Payload.Created)
			log.Printf("Description: %v\n", out.Payload.Description)
			d.Set("description", out.Payload.Description)
			log.Printf("Family: %v\n", out.Payload.Family)
			d.Set("family", out.Payload.Family)
			log.Printf("ID: %v\n", out.Payload.ID)
			d.Set("is_pool", out.Payload.IsPool)
			log.Printf("LastUpdated: %v\n", out.Payload.LastUpdated)
			d.Set("last_updated", out.Payload.LastUpdated)
			log.Print("\n")
		} else {
			log.Printf("erro na chamada do IPAMPrefixesList\n")
			log.Printf("Err: %v\n", err)
			log.Print("\n")
			return err
		}
	}
	log.Printf("data_source_netbox_prefixes.go dataSourceNetboxPrefixesRead %v\n", d)
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
		"family": &schema.Schema{
			Type: schema.TypeString,
		},
		"is_pool": &schema.Schema{
			Type: schema.TypeBool,
		},
		"last_updated": &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func resourcePrefixesSchema() map[string]*schema.Schema {
	s := barePrefixesSchema()

	for k, v := range s {
		switch k {
		case "prefixes_id":
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
func dataSourcePrefixesSchema() map[string]*schema.Schema {
	s := barePrefixesSchema()
	log.Printf("[DEBUG] ANTES: %v\n", s)
	for k, v := range s {
		switch k {
		case "prefixes_id":
			log.Printf("[DEBUG] JP dataSourcePrefixesSchema()\n")
			v.Optional = true
			v.Computed = true
			//v.ConflictsWith = []string{"ip_address", "subnet_id", "description", "hostname", "custom_field_filter"}
		default:
			v.Computed = true
		}
	}
	// Add the custom_field_filter item to the schema. This is a meta-parameter
	// that allows searching for a custom field value in the data source.
	s["custom_field_filter"] = customFieldFilterSchema([]string{"prefixes_id"})
	log.Printf("[DEBUG]DEPOIS: %v\n", s)

	return s
}
