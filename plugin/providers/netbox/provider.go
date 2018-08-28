package netbox

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var descriptions map[string]string

// Run before to initiliza vars ...
func init() {
	descriptions = map[string]string{
		"app_id":   "The application ID required for API requests",
		"endpoint": "The full URL (plus path) to the API endpoint",
		// "timeout":  "Max. wait time should wait for a successful connection to the API",
	}
}

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema:         providerSchema(),
		DataSourcesMap: providerDataSourcesMap(),
		ResourcesMap:   providerResources(),
		ConfigureFunc:  providerConfigure,
	}
}

// List of supported configuration fields for your provider.
// Here we define a linked list of all the fields that we want to
// support in our provider (api_key, endpoint, timeout & max_retries).
// More info in https://github.com/hashicorp/terraform/blob/v0.6.6/helper/schema/schema.go#L29-L142
func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"app_id": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: descriptions["key created on netbox"],
		},
		"endpoint": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: descriptions["endpoint of netbox (without http:// and / )"],
		},
		"timeout": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: descriptions["Max. wait time should wait for a successful connection to the API"],
		},
	}
}

// List of supported resources and their configuration fields.
// Here we define da linked list of all the resources that we want to
// support in our provider. As an example, if you were to write an AWS provider
// which supported resources like ec2 instances, elastic balancers and things of that sort
// then this would be the place to declare them.
// More info here https://github.com/hashicorp/terraform/blob/v0.6.6/helper/schema/resource.go#L17-L81
func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"netbox_vlans":                  resourceNetboxVlans(),
		"netbox_prefixes":               resourceNetboxPrefixes(),
		"netbox_prefixes_available_ips": resourceNetboxPrefixesAvailableIps(),
	}
}

// List of supported resources and their configuration fields.
// Here we define da linked list of all the resources that we want to
// support in our provider. As an example, if you were to write an AWS provider
// which supported resources like ec2 instances, elastic balancers and things of that sort
// then this would be the place to declare them.
// More info here https://github.com/hashicorp/terraform/blob/v0.6.6/helper/schema/resource.go#L17-L81

func providerDataSourcesMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"netbox_vlans":    dataSourceNetboxVlans(),
		"netbox_prefixes": dataSourceNetboxPrefixes(),
	}
}

// This is the function used to fetch the configuration params given
// to our provider which we will use to initialise a dummy client that
// interacts with the API.
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AppID:    d.Get("app_id").(string),
		Endpoint: d.Get("endpoint").(string),
		// Timeout:  d.Get("timeout").(string),
	}
	return config.Client()
}
