package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/cmgreivel/terraform-provider-netbox/plugin/providers/netbox"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: netbox.Provider,
	})
}
