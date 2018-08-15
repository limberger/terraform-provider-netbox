package main

import (
	"github.com/hashicorp/terraform/plugin"
	"gothub.com/limberger/terraform-provider-netbox/plugin/providers/netbox"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: netbox.Provider,
	})
}
