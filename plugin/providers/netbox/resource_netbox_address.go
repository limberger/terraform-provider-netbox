package netbox

import (
	// "errors"
	// "fmt"
	// "strconv"

	"github.com/hashicorp/terraform/helper/schema"
)


// resourceNetboxAddress returns the resource structure for the netbox_address
// resource.
//
// Note that we use the data source read function here to pull down data, as
// read workflow is identical for both the resource and the data source.
func resourceNetboxAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxAddressCreate,
		Read:   dataSourceNetboxAddressRead,
		Update: resourceNetboxAddressUpdate,
		Delete: resourceNetboxAddressDelete,
//		Schema: resourceAddressSchema(),
	}
}

func resourceNetboxAddressCreate(d *schema.ResourceData, meta interface{}) error {

                
	// c := meta.(*ProviderNetboxClient).client
	// in := expandAddress(d)

	// Assert the ID field here is empty. If this is not empty the request will fail.
	// in.ID = 0

	// if _, err := c.IPAM.IPAMIPAddressesCreate(in,nil); err != nil {
	// 	return err
	// }

	// If we have custom fields, set them now. We need to get the IP address's ID
	// beforehand.
	// if customFields, ok := d.GetOk("custom_fields"); ok {
		// addrs, err := c.GetAddressesByIP(in.IPAddress)
		// if err != nil {
		// 	return fmt.Errorf("Could not read IP address after creating: %s", err)
		// }

		// if len(addrs) != 1 {
		// 	return errors.New("IP address either missing or multiple results returned by reading IP after creation")
		// }

		// d.SetId(strconv.Itoa(addrs[0].ID))

		// if _, err := c.UpdateAddressCustomFields(addrs[0].ID, customFields.(map[string]interface{})); err != nil {
		// 	return err
		// }
	// }

	// return dataSourcePHPIPAMAddressRead(d, meta)
	return nil
}

func resourceNetboxAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	// c := meta.(*ProviderNetboxClient).client
	// in := expandAddress(d)

	// IPAddress and SubnetID need to be removed for update requests.
	// in.IPAddress = ""
	// in.SubnetID = 0
	// if _, err := c.UpdateAddress(in); err != nil {
	// 	return err
	// }

	// if err := updateCustomFields(d, c); err != nil {
	// 	return err
	// }

	// return dataSourcePHPIPAMAddressRead(d, meta)
	return nil
}

func resourceNetboxAddressDelete(d *schema.ResourceData, meta interface{}) error {
	// c := meta.(*ProviderNetboxClient).client
	// in := expandAddress(d)

	// if _, err := c.DeleteAddress(in.ID, phpipam.BoolIntString(d.Get("remove_dns_on_delete").(bool))); err != nil {
	// 	return err
	// }
	// d.SetId("")
	return nil
}


