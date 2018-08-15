package netbox

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"app_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["app_id"],
			},
			"endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["endpoint"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"netbox_address": resourcePHPIPAMAddress(),
			"netbox_section": resourcePHPIPAMSection(),
			"netbox_subnet":  resourcePHPIPAMSubnet(),
			"netbox_vlan":    resourcePHPIPAMVLAN(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"netbox_address":            dataSourcePHPIPAMAddress(),
			"netbox_addresses":          dataSourcePHPIPAMAddresses(),
			"netbox_first_free_address": dataSourcePHPIPAMFirstFreeAddress(),
			"netbox_section":            dataSourcePHPIPAMSection(),
			"netbox_subnet":             dataSourcePHPIPAMSubnet(),
			"netbox_subnets":            dataSourcePHPIPAMSubnets(),
			"netbox_vlan":               dataSourcePHPIPAMVLAN(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"app_id":   "The application ID required for API requests",
		"endpoint": "The full URL (plus path) to the API endpoint",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AppID:    d.Get("app_id").(string),
		Endpoint: d.Get("endpoint").(string),
	}
	return config.Client()
}
