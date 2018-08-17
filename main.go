package main

import (
	"log"

	"github.com/hashicorp/terraform/plugin"
	"github.com/limberger/terraform-provider-netbox/plugin/providers/netbox"
)

func main() {
	log.Println("[DEBUG] jp main.go main()")
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: netbox.Provider,
	})
}
