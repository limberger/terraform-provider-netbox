package netbox

import (
	"log"

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
//			"netbox_address": resourceNetboxAddress(),
		    "netbox_prefixes": resourceNetboxPrefixes(),
		},

		DataSourcesMap: map[string]*schema.Resource{
//			"netbox_address":            dataSourceNetboxAddress(),
			"netbox_prefixes":           dataSourceNetboxPrefixes(),
			"netbox_first_free_address": dataSourceNetboxFirstFreeAddress(),
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
	log.Printf("[DEBUG] provider.providerConfigure JP provider.providerConfigure APP_ID %s",d.Get("app_id"))
	log.Printf("[DEBUG] provider.providerConfigure JP provider.providerConfigure ENDPOINT %s",d.Get("endpoint"))
	config := Config{
		AppID:    d.Get("app_id").(string),
		Endpoint: d.Get("endpoint").(string),
	}
	return config.Client()
}
