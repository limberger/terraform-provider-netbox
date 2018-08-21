package netbox

import (
	// "errors"
	// "fmt"
	// "strconv"
	"log"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/digitalocean/go-netbox/netbox/client/ipam"
)


// resourceNetboxAddress returns the resource structure for the netbox_address
// resource.
//
// Note that we use the data source read function here to pull down data, as
// read workflow is identical for both the resource and the data source.
func resourceNetboxPrefixes() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxPrefixesCreate,
		Read:   dataSourceNetboxPrefixesRead,
		Update: resourceNetboxPrefixesUpdate,
		Delete: resourceNetboxPrefixesDelete,
		Schema: resourcePrefixesSchema(),
	}
}

func resourceNetboxPrefixesCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] JP resourceNetboxPrefixesCreate: %v\n",d)
	c := meta.(*ProviderNetboxClient).client
	var parm = ipam.NewIPAMPrefixesCreateParams()
	log.Println("Criei o parm\n")
	parm.Data.ID = int64(d.Get("prefixes_id").(int))
	//parm.Set("ID", int64(d.Get("prefixes_id").(int)))
	//parm.SetCreated(d.Get("prefixes_created"))
	log.Println("Setei o parm")
	
	//parms = ipam.NewIPAMPrefixesListParams()
	out , err := c.IPAM.IPAMPrefixesCreate(parm,nil)
	log.Printf("- Executado...\n")
	print("out %v\n", out)
	print("err %v\n", err)

	return nil
}

func resourceNetboxPrefixesUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] JP resourceNetboxPrefixesUpdate: %v\n",d)

	return nil
}

func resourceNetboxPrefixesDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] JP resourceNetboxPrefixesDelete: %v\n",d)
	return nil
}
