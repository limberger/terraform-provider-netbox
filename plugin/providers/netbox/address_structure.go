package netbox
import (
	// "errors"
	// "strconv"

	"github.com/hashicorp/terraform/helper/schema"
)


// bareAddressSchema returns a map[string]*schema.Schema with the schema used
// to represent a Netbox address resource. This output should then be modified
// so that required and computed fields are set properly for both the data
// source and the resource.
func bareAddressSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"subnet_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"ip_address": &schema.Schema{
			Type: schema.TypeString,
		},
		"is_gateway": &schema.Schema{
			Type: schema.TypeBool,
		},
		"description": &schema.Schema{
			Type: schema.TypeString,
		},
		"hostname": &schema.Schema{
			Type: schema.TypeString,
		},
		"mac_address": &schema.Schema{
			Type: schema.TypeString,
		},
		"owner": &schema.Schema{
			Type: schema.TypeString,
		},
		"state_tag_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"skip_ptr_record": &schema.Schema{
			Type: schema.TypeBool,
		},
		"ptr_record_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"device_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"switch_port_label": &schema.Schema{
			Type: schema.TypeString,
		},
		"note": &schema.Schema{
			Type: schema.TypeString,
		},
		"last_seen": &schema.Schema{
			Type: schema.TypeString,
		},
		"exclude_ping": &schema.Schema{
			Type: schema.TypeBool,
		},
		"edit_date": &schema.Schema{
			Type: schema.TypeString,
		},
		"custom_fields": &schema.Schema{
			Type: schema.TypeMap,
		},
	}
}


// dataSourceAddressSchema returns the schema for the netbox_address data
// source. It sets the searchable fields and sets up the attribute conflicts
// between IP address and address ID. It also ensures that all fields are
// computed as well.
func dataSourceAddressSchema() map[string]*schema.Schema {
	s := bareAddressSchema()
	for k, v := range s {
		switch k {
		case "address_id":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"ip_address", "subnet_id", "description", "hostname", "custom_field_filter"}
		case "ip_address":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"address_id", "subnet_id", "description", "hostname", "custom_field_filter"}
		case "subnet_id":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"ip_address", "address_id"}
		case "description":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"ip_address", "address_id", "hostname", "custom_field_filter"}
		case "hostname":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"ip_address", "address_id", "description", "custom_field_filter"}
		default:
			v.Computed = true
		}
	}
	// Add the custom_field_filter item to the schema. This is a meta-parameter
	// that allows searching for a custom field value in the data source.
	s["custom_field_filter"] = customFieldFilterSchema([]string{"ip_address", "address_id", "hostname", "description"})

	return s
}

// expandAddress returns the addresses.Address structure for a
// phpiapm_address resource or data source. Depending on if we are dealing with
// the resource or data source, extra considerations may need to be taken.
func expandAddress(d *schema.ResourceData) string {
	// s := addresses.Address{
	// 	ID:          d.Get("address_id").(int),
	// 	SubnetID:    d.Get("subnet_id").(int),
	// 	IPAddress:   d.Get("ip_address").(string),
	// 	IsGateway:   phpipam.BoolIntString(d.Get("is_gateway").(bool)),
	// 	Description: d.Get("description").(string),
	// 	Hostname:    d.Get("hostname").(string),
	// 	MACAddress:  d.Get("mac_address").(string),
	// 	Owner:       d.Get("owner").(string),
	// 	Tag:         d.Get("state_tag_id").(int),
	// 	PTRIgnore:   phpipam.BoolIntString(d.Get("skip_ptr_record").(bool)),
	// 	PTRRecordID: d.Get("ptr_record_id").(int),
	// 	DeviceID:    d.Get("device_id").(int),
	// 	Port:        d.Get("switch_port_label").(string),
	// 	Note:        d.Get("note").(string),
	// 	LastSeen:    d.Get("last_seen").(string),
	// 	ExcludePing: phpipam.BoolIntString(d.Get("exclude_ping").(bool)),
	// }

	// return s
	return ""
}
