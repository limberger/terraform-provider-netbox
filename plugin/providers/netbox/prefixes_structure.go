
package netbox
import (
	// "errors"
	// "strconv"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func barePrefixesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"prefixes_id": &schema.Schema{
			Type: schema.TypeInt,
		},
	}
}



// dataSourceAddressSchema returns the schema for the phpipam_address data
// source. It sets the searchable fields and sets up the attribute conflicts
// between IP address and address ID. It also ensures that all fields are
// computed as well.
func dataSourcePrefixesSchema() map[string]*schema.Schema {
	s := barePrefixesSchema()
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

	return s
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
	s["remove_dns_on_delete"] = &schema.Schema {
		Type: schema.TypeBool,
		Optional: true,
		Default: true,
	}
	return s

}
