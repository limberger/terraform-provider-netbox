package netbox

import (
	"log"
	"reflect"

	// "errors"

	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/digitalocean/go-netbox/netbox/client/ipam"
	// "github.com/digitalocean/go-netbox/netbox/client"
)

func bareFirstFreeAddressSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"prefixes_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"ip_address": &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func dataSourceFirstFreeAddressSchema() map[string]*schema.Schema {
	s := bareFirstFreeAddressSchema()
	for k, v := range s {
		switch k {
		case "ip_address":
			log.Printf("[DEBUG] JP dataSourceFirstFreeAddressSchema()\n")
			v.Optional = true
			// v.Computed = true
			//v.ConflictsWith = []string{"ip_address", "subnet_id", "description", "hostname", "custom_field_filter"}
		case "prefixes_id":
			log.Printf("[DEBUG] JP dataSourceFirstFreeAddressSchema - prefixes_id\n")
			v.Optional = true
		default:
			v.Computed = true
		}
	}
	// Add the custom_field_filter item to the schema. This is a meta-parameter
	// that allows searching for a custom field value in the data source.
	s["custom_field_filter"] = customFieldFilterSchema([]string{"ip_address"})

	log.Printf("*FIM* data_source_netbox_first_free_address dataSourceFirstFreeAddresschema() -> %v", s)

	return s
}

func dataSourceNetboxFirstFreeAddress() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNetboxFirstFreeAddressRead,
		Schema: dataSourceFirstFreeAddressSchema(),
	}
}

func dataSourceNetboxFirstFreeAddressRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] JP : dataSourceNetboxFirstFreeAddressRead %v\n", d)
	log.Printf("Tipo do d : %s\n", reflect.TypeOf(d))

	// switch {
	// 	case d.Get("prefixes_id").(int) != 0:
	// 		log.Println("Tem prefixo ...")
	// 		log.Printf("data_source_netbox_prefixes.go dataSourceNetboxPrefixesRead - Prefixo: %i\n", d.Get("prefixes_id").(int))

	// 		c := meta.(*ProviderNetboxClient).client
	// 		log.Printf("Tipo do c: %s\n", reflect.TypeOf(c))
	// 		log.Printf("Obtive o client\n")
	// 		//parms = ipam.NewIPAMPrefixesListParams()
	// 		// out , err := c.IPAM.IPAMPrefixesRead(parm,nil)

	// }

	return nil
}
