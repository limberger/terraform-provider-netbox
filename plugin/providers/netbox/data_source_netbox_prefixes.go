package netbox

import (
	"log"
	"reflect"
	// "errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/digitalocean/go-netbox/netbox/client/ipam"
	// "github.com/digitalocean/go-netbox/netbox/client"
)


func dataSourceNetboxPrefixes() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNetboxPrefixesRead,
		Schema: dataSourcePrefixesSchema(),
	}
}

func dataSourceNetboxPrefixesRead(d *schema.ResourceData, meta interface{}) error {
	//out := ipam.NewIPAMPrefixesListParams()

	log.Printf("[DEBUG] JP data_source_netbox_prefixes.go : dataSourceNetboxPrefixesRead %v\n",d)
	switch {
		case d.Get("prefixes_id").(int) != 0:
			log.Println("É um prefixo ...")
			log.Printf("data_source_netbox_prefixes.go dataSourceNetboxPrefixesRead - Prefixo: %i\n", d.Get("prefixes_id").(int))
			var parm = ipam.NewIPAMPrefixesReadParams()
			log.Println("Criei o parm\n")
			parm.SetID(int64(d.Get("prefixes_id").(int)))
			log.Println("Setei o parm")
			log.Printf("Tipo do parm [meta] : %s", reflect.TypeOf(meta))
			//(&&meta).IPAM.IPAMPrefixesRead(parm,nil)

			c := meta.(*ProviderNetboxClient).client
			log.Printf("Obtive o client\n")
			//parms = ipam.NewIPAMPrefixesListParams()
			out , err := c.IPAM.IPAMPrefixesRead(parm,nil)
			log.Printf("- Executado...\n")
			if err == nil {
				log.Printf("Ok na chamada do IPAMPrefixesList\n")
				log.Printf("Out: %v\n", out)
				log.Printf("Created: %v\n", &out.Payload.Created)
				d.Set("Created",out.Payload.Created)
				log.Printf("Description: %v\n", out.Payload.Description)
				d.Set("Description",out.Payload.Description)
				log.Printf("Family: %v\n", out.Payload.Family)
				d.Set("Family",out.Payload.Family)
				log.Printf("ID: %v\n", out.Payload.ID)
				d.Set("ID",out.Payload.ID)
				log.Printf("IsPool: %v\n", out.Payload.IsPool)
				d.Set("IsPool",out.Payload.IsPool)
				log.Printf("LastUpdated: %v\n", out.Payload.LastUpdated)
				d.Set("LastUpdated",out.Payload.LastUpdated)
				log.Print("\n")

			} else {
				log.Printf("erro na chamada do IPAMPrefixesList\n")
				log.Printf("Err: %v\n", err)
				log.Print("\n")
				return err
			}

	}
	log.Printf("data_source_netbox_prefixes.go dataSourceNetboxPrefixesRead %v\n",d)
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


 