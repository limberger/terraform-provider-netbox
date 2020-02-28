package netbox

import (
//	"bytes"
//	"encoding/json"
//	"errors"
//	"io/ioutil"
	"log"
//	"net/http"
//	"strconv"
//	"strings"

	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNetboxPoolPrefixes() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNetboxPoolPrefixesRead,
		Schema: dataSourcePoolPrefixesSchema(),
	}
}

func dataSourcePoolPrefixesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pool_id": &schema.Schema{
			Type: schema.TypeInt,
			// We may search by the prefix itself and not the ID
			Optional: true,
			// We will return the pool_id during create if
			// the prefix was specified.
			Computed: true,
		},
		"pool": &schema.Schema{
			Type: schema.TypeString,
			// We may search by the Id and not the prefix
			Optional: true,
			// We will return the pool during create if the
			// ID was specified.
			Computed: true,
		},
		"prefix_id": &schema.Schema{
			Type: schema.TypeInt,
			Optional: true,
			// prefix_id is generated when the prefix is allocated
			Computed: true,
		},
		"prefix": &schema.Schema{
			Type: schema.TypeString,
			Optional: true,
			// prefix is generated when allocated
			Computed: true,
		},
	}
}

func dataSourceNetboxPoolPrefixesRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("dataSourceNetboxPoolPrefixesRead ............ ")
	switch {
		// We either need the pool or the pool_id
	case d.Get("pool").(string) != "":

		ispool := "true"
		pool, _ := d.Get("pool_prefix").(string)
		var listParm = ipam.NewIPAMPrefixesListParams().WithPrefix(&pool).WithIsPool(&ispool)
		// Find the Pool Prefix that the prefix belongs to


		// CMG HERE
		log.Printf("- Listing pool prefix\n")
		c := meta.(*ProviderNetboxClient).client
		pool_out, err := c.IPAM.IPAMPrefixesList(listParm, nil)

		if err == nil {
			// https://github.com/netbox-community/go-netbox/blob/master/netbox/ipam/ip_a_m_prefixes_list_responses.go#L67
			// https://github.com/netbox-community/go-netbox/blob/master/netbox/ipam/ip_a_m_prefixes_list_responses.go#L90

			if *pool_out.Payload.Count == 1 {
				log.Printf("Found one prefix!\n")
				log.Printf(" -> %v\n", pool_out.Payload.Results[0])
				d.Set("pool_prefix_id", pool_out.Payload.Results[0].ID)
			} else {
				log.Printf("Found %d prefixes\n", pool_out.Payload.Count)
			}

/* CMG
			var parm = ipam.NewIPAMIPPrefixesReadParams()

			parm.SetID(id)
			//(&&meta).IPAM.IPAMPrefixesRead(parm,nil)

			d.Set("address_id", string(out.Payload.ID))
			d.Set("address", out.Payload.Address)

			d.Set("mask", strings.Split(*out.Payload.Address, "/")[1])
			d.Set("ip", strings.Split(*out.Payload.Address, "/")[0])

			log.Printf("Setando Address_id %v\n", out.Payload.ID)
			d.Set("created", out.Payload.Created)
			if out.Payload.CustomFields != nil {
				d.Set("custom_fields", out.Payload.CustomFields)
			}
			d.Set("description", out.Payload.Description)
			d.Set("family", out.Payload.Family)
			if out.Payload.Interface != nil {
				d.Set("interface_id", out.Payload.Interface.ID)
				d.Set("interface_name", out.Payload.Interface.Name)
			}
			if out.Payload.Role != nil {
				d.Set("role_id", out.Payload.Role.Value)
				d.Set("role_label", out.Payload.Role.Label)
			}
			if out.Payload.Status != nil {
				d.Set("status_id", out.Payload.Status.Value)
				d.Set("status_label", out.Payload.Status.Value)
			}

			d.Set("last_updated", out.Payload.LastUpdated)
			log.Print("\n")
CMG */
		} else {
			log.Printf("Error with IPAMPrefixesList()\n")
			log.Printf("Err: %v\n", err)
			log.Print("\n")
			return err
		}

	case d.Get("pool_id").(int) != 0:
		log.Printf("pool_id is SET!")
		return nil
	default:
		//return errors.New("Neither pool nor pool_id is defined.")
		log.Printf("Neither pool nor pool_id is defined.")
		//d.SetId("")
		return nil
	}
	return nil
}

