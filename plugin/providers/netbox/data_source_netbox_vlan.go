package netbox

import (
	"errors"
	"log"
	"strconv"

	// "errors"

	"github.com/digitalocean/go-netbox/netbox/client/ipam"
	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/digitalocean/go-netbox/netbox/client"
)

func dataSourceNetboxVlans() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNetboxVlansRead,
		Schema: dataSourceVlanSchema(),
	}
}

// Read will fetch the data of a resource.
func dataSourceNetboxVlansRead(d *schema.ResourceData, meta interface{}) error {
	//out := ipam.NewIPAMVlansListParams()

	var parm = ipam.NewIPAMVlansListParams()
	switch {
	case d.Get("vid").(int) != 0:
		// func (a *Client) IPAMVlansRead(params *IPAMVlansReadParams, authInfo runtime.ClientAuthInfoWriter) (*IPAMVlansReadOK, error) {
		log.Printf("Ok... localizando por vid: %v\n", d.Get("vid").(int))
		vid := float64(d.Get("vid").(int))
		parm.SetVid(&vid)
		c := meta.(*ProviderNetboxClient).client
		//parms = ipam.NewIPAMVlansListParams()
		out, err := c.IPAM.IPAMVlansList(parm, nil)
		log.Printf("- Executado...\n")
		if err == nil {

			if *out.Payload.Count == 0 {
				return errors.New("Vid not found")
			} else if *out.Payload.Count > 1 {
				return errors.New("More than one Vid found with name " + d.Get("name").(string))
			}
			result := out.Payload.Results[0]
			d.SetId(strconv.Itoa(int(result.ID)))
			d.Set("created", result.Created)
			d.Set("description", result.Description)
			d.Set("display_name", result.DisplayName)
			d.Set("group", result.Group)
			d.Set("vid", *result.Vid)
			d.Set("last_updated", result.LastUpdated)
			d.Set("name", *result.Name)
			d.Set("role", result.Role)
			d.Set("nested_site", result.Site)
			d.Set("status", result.Status)
			d.Set("nested_tenant", result.Tenant)
			d.Set("custom_fields", result.CustomFields)

		} else {
			log.Printf("erro na chamada do IPAMVlansList\n")
			log.Printf("Err: %v\n", err)
			log.Print("\n")
			return err
		}
	case d.Get("name").(string) != "":
		c := meta.(*ProviderNetboxClient).client
		parmsl := ipam.NewIPAMVlansListParams()
		name := d.Get("name").(string)
		log.Printf("Nome: %v\n", name)
		parmsl.SetName(&name)
		log.Printf("Parmsl: %v\n", parmsl)
		out, err := c.IPAM.IPAMVlansList(parmsl, nil)
		if err == nil {
			if *out.Payload.Count == 0 {
				return errors.New("Name not found - Payload = 0 - need one of vid or name")
			} else if *out.Payload.Count > 1 {
				return errors.New("More than one vlan found with name " + d.Get("name").(string))
			}
			result := out.Payload.Results[0]
			d.SetId(strconv.Itoa(int(result.ID)))
			d.Set("created", result.Created)
			d.Set("description", result.Description)
			d.Set("display_name", result.DisplayName)
			d.Set("group", result.Group)
			log.Printf("Por nome: result.Vid %v Nome: %v Buscado: %v\n", *result.Vid, *result.Name, name)
			d.Set("vid", *result.Vid)
			d.Set("last_updated", result.LastUpdated)
			d.Set("name", *result.Name)
			d.Set("role", result.Role)
			d.Set("nested_site", result.Site)
			d.Set("status", result.Status)
			d.Set("nested_tenant", result.Tenant)
			d.Set("custom_fields", result.CustomFields)
			log.Printf("Custom Fields: %v\n", result.CustomFields)
		} else {
			log.Printf("erro na chamada do IPAMVlansList\n")
			log.Printf("Err: %v\n", err)
			return err
		}

	case d.Get("vid").(int) == 0 && d.Get("name").(string) == "":
		log.Printf("Informado: vid %v\n", d.Get("vid").(int))
		log.Printf("Informado: name %v\n", d.Get("name").(string))
		log.Printf("Informado: d %v\n", d)
		return errors.New("No valid combination of parameters found - need one of vid or name ...")
	}
	return nil
}

func bareVlanSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vid": &schema.Schema{
			Type: schema.TypeInt,
		},
		"name": &schema.Schema{
			Type: schema.TypeString,
		},
		"id_in": &schema.Schema{
			Type: schema.TypeString,
		},
		"q": &schema.Schema{
			Type: schema.TypeString,
		},
		"site_id": &schema.Schema{
			Type: schema.TypeString,
		},
		"site": &schema.Schema{
			Type: schema.TypeString,
		},
		"group_id": &schema.Schema{
			Type: schema.TypeString,
		},
		"group": &schema.Schema{
			Type: schema.TypeString,
		},
		"tenant_id": &schema.Schema{
			Type: schema.TypeString,
		},
		"tenant": &schema.Schema{
			Type: schema.TypeString,
		},
		"role": &schema.Schema{
			Type: schema.TypeMap,
		},
		"status": &schema.Schema{
			Type: schema.TypeString,
		},
		"tag": &schema.Schema{
			Type: schema.TypeString,
		},
		// Number of results to return per page.
		"limit": &schema.Schema{
			Type: schema.TypeInt,
		},
		// The initial index from which to return the results.
		"offset": &schema.Schema{
			Type: schema.TypeInt,
		},
		"custom_fields": &schema.Schema{
			Type: schema.TypeMap,
		},
	}
}

func resourceVlansSchema() map[string]*schema.Schema {
	s := bareVlanSchema()

	for k, v := range s {
		switch k {
		case "vid":
			v.Computed = true
			v.Optional = true
		case "name":
			v.Optional = true
			v.Computed = true
		case "custom_fields":
			v.Optional = true
		case "status":
			v.Optional = true
		case "role":
			v.Optional = true
			//v.ConflictsWith = []string{"ip_address", "subnet_id", "description", "hostname", "custom_field_filter"}
		default:
			v.Computed = true
		}
	}
	// Add the remove_dns_on_delete item to the schema. This is a meta-parameter
	// that is not part of the API resource and exists to instruct NETBOX to
	// gracefully remove the address from its DNS integrations as well when it is
	// removed. The default on this option is true.
	s["remove_dns_on_delete"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	}
	return s
}

// dataSourceAddressSchema returns the schema for the NETBOX_VLANS data
// source. It sets the searchable fields and sets up the attribute conflicts
// between IP address and address ID. It also ensures that all fields are
// computed as well.
func dataSourceVlanSchema() map[string]*schema.Schema {
	s := bareVlanSchema()
	for k, v := range s {
		switch k {
		case "vid":
			v.Computed = true
			v.Optional = true
		case "name":
			v.Optional = true
			v.Computed = true
		case "custom_fields":
			v.Optional = true
		case "status":
			v.Optional = true

			//v.ConflictsWith = []string{"ip_address", "subnet_id", "description", "hostname", "custom_field_filter"}
		default:
			v.Computed = true
		}
	}
	// Add the custom_field_filter item to the schema. This is a meta-parameter
	// that allows searching for a custom field value in the data source.
	s["custom_field_filter"] = customFieldFilterSchema([]string{"vid"})

	return s
}
