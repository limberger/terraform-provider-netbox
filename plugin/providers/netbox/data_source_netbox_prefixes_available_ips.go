package netbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/digitalocean/go-netbox/netbox/client/ipam"
	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/digitalocean/go-netbox/netbox/client/ipam"
	// "github.com/digitalocean/go-netbox/netbox/client"
)

func resourceNetboxPrefixesAvailableIps() *schema.Resource {
	return &schema.Resource{
		Read:   resourceNetboxPrefixesAvailableIpsRead,
		Create: resourceNetboxPrefixesAvailableIpsCreate,
		Update: resourceNetboxPrefixesAvailableIpsUpdate,
		Delete: resourceNetboxPrefixesAvailableIpsDelete,
		Schema: resourcePrefixesAvailableIpsSchema(),
	}
}

func barePrefixesAvailableIpsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"prefixes_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"address_id": &schema.Schema{
			Type: schema.TypeString,
		},
		"description": &schema.Schema{
			Type: schema.TypeString,
		},
		"family": &schema.Schema{
			Type: schema.TypeInt,
		},
		"address": &schema.Schema{
			Type: schema.TypeString,
		},
		"ip": &schema.Schema{
			Type: schema.TypeString,
		},
		"mask": &schema.Schema{
			Type: schema.TypeString,
		},
		"custom_fields": &schema.Schema{
			Type: schema.TypeMap,
		},
		"status": &schema.Schema{
			Type: schema.TypeString,
		},
		"created": &schema.Schema{
			Type: schema.TypeString,
		},
		"last_updated": &schema.Schema{
			Type: schema.TypeString,
		},
		"interface_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"interface_label": &schema.Schema{
			Type: schema.TypeString,
		},
		"role_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"role_label": &schema.Schema{
			Type: schema.TypeString,
		},
		"status_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"status_label": &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func resourcePrefixesAvailableIpsSchema() map[string]*schema.Schema {
	s := barePrefixesAvailableIpsSchema()
	for k, v := range s {
		switch k {
		case "address_id":
			v.Optional = true
			v.Computed = true
		case "prefixes_id":
			v.Optional = true
			v.Computed = true
		case "description":
			v.Optional = true
		case "family":
			v.Optional = true
			v.Computed = true
		case "address":
			v.Optional = true
			v.Computed = true
		case "ip":
			v.Optional = true
			v.Computed = true
		case "mask":
			v.Optional = true
			v.Computed = true
		case "custom_fields":
			v.Optional = true
			v.Computed = true
		case "status":
			v.Optional = true
			v.Computed = true
		case "created":
			v.Optional = true

		default:
			v.Computed = true
		}
	}
	// Add the custom_field_filter item to the schema. This is a meta-parameter
	// that allows searching for a custom field value in the data source.
	//s["custom_field_filter"] = customFieldFilterSchema([]string{"ip_address"})

	return s
}

func resourceNetboxPrefixesAvailableIpsCreate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] JP dataNetboxPrefixesAvailableIpsCreate  ** : %v\n", d)
	// c := meta.(*ProviderNetboxClient).client

	// ip_address
	if d.Get("prefixes_id") == nil {
		log.Println("[ERROR] JP - prefixes_id == nil...")
		return errors.New("prefixes_id not informed")
	}
	prefixes_id := d.Get("prefixes_id").(int)
	log.Printf("[DEBUG] JP prefixes_id %v\n", prefixes_id)
	//parm.Data.ID =
	if d.Get("description") == nil {
		log.Println("[ERROR] JP - description == nil...")
		return errors.New("description not informed")
	}
	description := d.Get("description").(string)
	log.Printf("[DEBUG] Inclusao prefixo     %v\n", prefixes_id)
	log.Printf("[DEBUG] Inclusao description %v\n", description)

	c := meta.(*ProviderNetboxClient).configuration
	log.Printf("[DEBUG] Configuration [%v]\n", c)
	log.Printf("[DEBUG] AppID [%v]\n", c.AppID)
	log.Printf("[DEBUG] Endpoint [%v]\n", c.Endpoint)

	// log.Printf("%v", *Config.cfg)
	// log.Printf("[DEBUG] Config.AppID %v\n", Config.AppID)
	// log.Printf("[DEBUG] Config.Endpoint %v\n", *Config.Endpoint)

	url := "http://" + c.Endpoint + "/api/ipam/prefixes/" + strconv.Itoa(prefixes_id) + "/available-ips/"
	jsonData := map[string]string{"description": description}
	jsonValue, _ := json.Marshal(jsonData)
	log.Printf("[DEBUG] JSON: [%v]\n", string(jsonValue))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("[ERROR] Error occurred in POST to url [%v]\n", url)
		log.Printf("[ERROR] Erro : %v\n", err)
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("authorization", "Token "+c.AppID)
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/json")
	log.Println("[DEBUG] http.Client create")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[DEBUG] Error occurred after post")
		log.Printf("[ERROR] Erro retorno http. %v \n", err)
		return err
	}
	log.Printf("[DEBUG] Http Code Response: %v\n", resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("[DEBUG] Body Read")
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		return errors.New("Return http code: " + strconv.Itoa(resp.StatusCode))
	}
	// Dinamico ...
	var i map[string]interface{}
	log.Print("[DEBUG] Will Unmarshal the body")
	jsonErr := json.Unmarshal(body, &i)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	log.Println("[DEBUG] Unmarshal Ok")
	log.Printf("[DEBUG] [%v]", string(body))
	log.Println("[DEBUG] ID setting")
	d.SetId(strconv.FormatFloat(i["id"].(float64), 'f', -1, 64))
	log.Println("[DEBUG] family")
	// log.Printf("[DEBUG] family [%v]\n", i["family"].(int))
	// d.Set("family", i["family"].(int))
	log.Println("[DEBUG] address")

	d.Set("address", i["address"].(string))
	d.Set("mask", strings.Split(i["address"].(string), "/")[1])
	d.Set("ip", strings.Split(i["address"].(string), "/")[0])
	d.Set("address_id", strconv.FormatFloat(i["id"].(float64), 'f', -1, 64))
	log.Println("[DEBUG] description")
	d.Set("description", i["description"].(string))
	d.Set("status", i["status"])
	d.Set("created", i["created"].(string))
	d.Set("last_updated", i["last_updated"].(string))
	log.Printf("Incluido id: %v\n", d.Id())
	return nil
}

func resourceNetboxPrefixesAvailableIpsRead(d *schema.ResourceData, meta interface{}) error {
	//out := ipam.NewIPAMPrefixesListParams()
	log.Printf("resourceNetboxPrefixesAvailableIpsRead ............ ")
	switch {
	// Pega por prefix_id
	case d.Get("address_id").(string) != "":
		var parm = ipam.NewIPAMIPAddressesReadParams()

		id, _ := strconv.ParseInt(d.Get("address_id").(string), 10, 64)
		parm.SetID(id)
		//(&&meta).IPAM.IPAMPrefixesRead(parm,nil)

		c := meta.(*ProviderNetboxClient).client
		out, err := c.IPAM.IPAMIPAddressesRead(parm, nil)
		log.Printf("- Executado...\n")
		if err == nil {

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
		} else {
			log.Printf("erro na chamada do IPAMIPAddressesRead\n")
			log.Printf("Err: %v\n", err)
			log.Print("\n")
			return err
		}

	default:
		//return errors.New("No valid parameters found - address_id")
		log.Printf("Address_id not informed or not exist.")
		//d.SetId("")
		return nil
	}
	return nil
}

func resourceNetboxPrefixesAvailableIpsUpdate(d *schema.ResourceData, meta interface{}) error {
	//out := ipam.NewIPAMPrefixesListParams()
	log.Printf("resourceNetboxPrefixesAvailableIpsUpdate ............ ")
	// switch {
	// // Pega por prefix_id
	// case d.Get("address_id").(int) != 0: // Obrigatório
	// 	var parm = ipam.NewIPAMIPAddressesUpdateParams()
	// 	parm.SetID(int64(d.Get("prefixes_id").(int)))
	// 	//(&&meta).IPAM.IPAMPrefixesRead(parm,nil)
	//
	// 	var entrada = models.WritableIPAddress{}
	// 	entrada.Address = d.Get("address").(*string)
	// 	entrada.Created = d.Get("created").(strfmt.Date)
	// 	entrada.CustomFields = d.Get("custom_field").(interface{})
	// 	entrada.Description = d.Get("description").(string)
	// 	entrada.ID = d.Get("address_id").(int64)
	// 	parm.SetData(&entrada)
	// 	c := meta.(*ProviderNetboxClient).client
	// 	out, err := c.IPAM.IPAMIPAddressesUpdate(parm, nil)
	// 	log.Printf("- Executado...\n")
	// 	if err == nil {
	//
	// 		d.SetId(string(out.Payload.ID)) // Sempre setar o ID
	// 		d.Set("address", out.Payload.Address)
	// 		d.Set("created", out.Payload.Created)
	// 		if out.Payload.CustomFields != nil {
	// 			d.Set("custom_fields", out.Payload.CustomFields)
	// 		}
	// 		d.Set("description", out.Payload.Description)
	// 		d.Set("last_updated", out.Payload.LastUpdated)
	// 		log.Print("\n")
	// 	} else {
	// 		log.Printf("erro na chamada do IPAMIPAddressesRead\n")
	// 		log.Printf("Err: %v\n", err)
	// 		log.Print("\n")
	// 		return err
	// 	}
	// 	// Pega por prefix.vlan.vid
	// default:
	// 	//return errors.New("No valid parameters found - address_id")
	// 	log.Printf("Address_id not informed or not exist.")
	// 	d.SetId("")
	// 	return nil
	// }

	return nil
}
func resourceNetboxPrefixesAvailableIpsDelete(d *schema.ResourceData, meta interface{}) error {
	//out := ipam.NewIPAMPrefixesListParams()
	log.Printf("resourceNetboxPrefixesAvailableIpsDelete ............ ")
	log.Printf("[DEBUG] d -> [%v]\n", d)
	switch {
	// Pega por prefix_id
	case d.Id() != "": // Obrigatório
		var parm = ipam.NewIPAMIPAddressesDeleteParams()

		log.Printf("[DEBUG] Id [%v]\n", d.Id())
		id, _ := strconv.ParseInt(d.Id(), 10, 64)
		parm.SetID(id)
		log.Printf("[DEBUG] Deletando IP Address ID [%v]\n", id)

		c := meta.(*ProviderNetboxClient).client
		_, err := c.IPAM.IPAMIPAddressesDelete(parm, nil)
		log.Printf("- Executado Delete...\n")
		if err == nil {
			log.Printf("[DEBUG] Recurso %v deletado\n", d.Get("address_id").(string))
		} else {
			log.Printf("erro na chamada do IPAMIPAddressesDelete\n")
			log.Printf("Err: %v\n", err)
			log.Print("\n")
			return err
		}
		// Pega por prefix.vlan.vid
	default:
		return errors.New("Address_id must be informed or must exist.")
	}
	return nil
}
