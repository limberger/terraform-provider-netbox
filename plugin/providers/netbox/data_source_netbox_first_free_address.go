package netbox

import (

	// "errors"

	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/digitalocean/go-netbox/netbox/client/ipam"
	// "github.com/digitalocean/go-netbox/netbox/client"
)

func dataSourceNetboxFirstFreeAddress() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNetboxFirstFreeAddressRead,
		Schema: dataSourceFirstFreeAddressSchema(),
	}
}

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
			v.Optional = true
		case "prefixes_id":
			v.Optional = true
		default:
			v.Computed = true
		}
	}
	// Add the custom_field_filter item to the schema. This is a meta-parameter
	// that allows searching for a custom field value in the data source.
	s["custom_field_filter"] = customFieldFilterSchema([]string{"ip_address"})

	return s
}

func resourceFirstFreeAddressSchema() map[string]*schema.Schema {
	s := bareFirstFreeAddressSchema()

	for k, v := range s {
		switch k {
		case "ip_address":
			v.Optional = true
		case "prefixes_id":
			v.Optional = true
			//v.ConflictsWith = []string{"ip_address", "subnet_id", "description", "hostname", "custom_field_filter"}
		default:
			v.Computed = true
		}
	}

	return s
}

func dataSourceNetboxFirstFreeAddressRead(d *schema.ResourceData, meta interface{}) error {

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
