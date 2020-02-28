package netbox

import (
	// "errors"
	// "fmt"
	// "strconv"
	"log"

	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// resourceNetboxAddress returns the resource structure for the netbox_address
// resource.
//
// Note that we use the data source read function here to pull down data, as
// read workflow is identical for both the resource and the data source.
func resourceNetboxVlans() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxVlansCreate,
		Read:   dataSourceNetboxVlansRead,
		Update: resourceNetboxVlansUpdate,
		Delete: resourceNetboxVlansDelete,
		Exists: resourceNetboxVlansExists,

		Schema: resourceVlansSchema(),
	}
}

// Exists is called before Read and obviously makes sure the resource exists.
func resourceNetboxVlansExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	return true, nil
}

// Create will simply create a new instance of your resource.
// The is also where you will have to set the ID (has to be an Int) of your resource.
// If the API you are using doesn’t provide an ID, you can always use a random Int.
func resourceNetboxVlansCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] JP resourceNetboxVlansCreate: %v\n", d)
	c := meta.(*ProviderNetboxClient).client
	var parm = ipam.NewIPAMVlansCreateParams()
	log.Println("Criei o parm")
	parm.Data.ID = int64(d.Get("prefixes_id").(int))
	//parm.Set("ID", int64(d.Get("prefixes_id").(int)))
	//parm.SetCreated(d.Get("prefixes_created"))
	log.Println("Setei o parm")

	//parms = ipam.NewIPAMPrefixesListParams()
	out, err := c.IPAM.IPAMVlansCreate(parm, nil)
	log.Printf("- Executado...\n")
	print("out %v\n", out)
	print("err %v\n", err)

	return nil
}

//Update is optional if your Resource doesn’t support update.
//For example, I’m not using update in the Terraform LDAP Provider.
//I just destroy and recreate the resource everytime there is a change.
func resourceNetboxVlansUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] JP resourceNetboxVlansUpdate: %v\n", d)

	return nil
}

func resourceNetboxVlansDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] JP resourceNetboxVlansDelete: %v\n", d)
	d.SetId("")
	return nil
}
