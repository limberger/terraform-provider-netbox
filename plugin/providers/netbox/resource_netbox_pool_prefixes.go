package netbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxPoolPrefixes() *schema.Resource {
	return &schema.Resource{
		Read:   resourceNetboxPoolPrefixesRead,
		Create: resourceNetboxPoolPrefixesCreate,
		Update: resourceNetboxPoolPrefixesUpdate,
		Delete: resourceNetboxPoolPrefixesDelete,
		Schema: resourcePoolPrefixesSchema(),
	}
}

// CMG: May need to add vrf
func resourcePoolPrefixesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pool_id": &schema.Schema{
			Type:     schema.TypeInt,
			Description: "ID of the pool from which the prefix will be allocated",
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		"pool": &schema.Schema{
			Type:     schema.TypeString,
			Description: "CIDR block of the pool from which the prefix will be allocated",
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		"prefix_id": &schema.Schema{
			Type:     schema.TypeInt,
			Description: "ID of the prefix allocated with the resource",
			Optional: true,
			Computed: true,
		},
		"prefix": &schema.Schema{
			Type:     schema.TypeString,
			Description: "CIDR block of the prefix allocated with the resource",
			Optional: true,
			Computed: true,
		},
		"prefix_length": &schema.Schema{
			Type:     schema.TypeInt,
			Description: "Length of the prefix in bits",
			Required: true,
			ForceNew: true,
		},
		"tags": &schema.Schema{
			Type: schema.TypeMap,
			Description: "Tags applied to the prefix",
			Optional: true,
			Computed: true,
		},
		"ring": &schema.Schema{
			Type: schema.TypeString,
			Description: "The ring the prefix belongs to",
			Required: true,
			ForceNew: true,
		},
	}
}

func tagMapToList(d *schema.ResourceData) []string {
	var tags map[string]interface{}
	var tagList []string

	if d.Get("tags").(map[string]interface{}) != nil {
		tags, _ = d.Get("tags").(map[string]interface{})
		log.Printf("[DEBUG] Have input tags %v\n", tags)
		if len(tags) != 0 {
			// Add tags
			for key, value := range tags {
				tagList = append(tagList, fmt.Sprintf("%s=%s", key, value))

			}
		}
	}
	return tagList
}

func resourceNetboxPoolPrefixesCreate(d *schema.ResourceData, meta interface{}) error {

	poolId := -1
	pool := ""
	log.Println("[DEBUG] dataNetboxPoolPrefixesCreate")
	client := meta.(*ProviderNetboxClient).client

	if d.Get("prefix_length").(int) == 0 {
		log.Println("[ERROR] prefix_length not specified.")
		return errors.New("prefix_length not specified")
	}
	prefixLength := d.Get("prefix_length").(int)

	if d.Get("ring").(string) == "" {
		log.Println("[ERROR] ring not specified.")
		return errors.New("ring not specified")
	}
	ring := d.Get("ring").(string)

	switch {
	case d.Get("pool_id").(int) != 0:
		poolId, _ = d.Get("pool_id").(int)
		log.Printf("[DEBUG] Have poolId %d \n", poolId)

	case d.Get("pool").(string) != "":
		pool, _ = d.Get("pool").(string)
		log.Printf("[DEBUG] Have pool %s\n", pool)

	default:
		log.Println("[ERROR] Pool not specified. Must set pool_id or pool.")
		return errors.New("pool not specified")
	}

	taglist := tagMapToList(d)

	// Find the ID of the VRF. Having trouble just using the vrf name for some reason.
	// XXX: ^^^
	ringParm := ipam.NewIPAMVrfsListParams().WithName(&ring)
	ringResult, err := client.IPAM.IPAMVrfsList(ringParm, nil)
	if err != nil {
		return err
	}
	if *ringResult.Payload.Count != 1 {
		log.Printf("Found %d vrfs for ring %s\n", *ringResult.Payload.Count, ring)
		return errors.New("Too many vrfs for ring")
	}
	ringId := strconv.FormatInt(ringResult.Payload.Results[0].ID, 10)

	if poolId == -1 {
		// We need to find the poolId
		ispool := "true"
		listParm := ipam.NewIPAMPrefixesListParams().WithPrefix(&pool).WithIsPool(&ispool).WithVrfID(&ringId)
		found, err := client.IPAM.IPAMPrefixesList(listParm, nil)
		if err != nil {
			return err
		}
		if *found.Payload.Count != 1 {
			errString := fmt.Sprintf("Found %d pools for prefix %s\n", *found.Payload.Count, pool)
			return errors.New(errString)
		}
		poolId = int(found.Payload.Results[0].ID)
		d.Set("pool_id", found.Payload.Results[0].ID)
	}

	// We cannot use go-netbox (https://github.com/netbox-community/go-netbox) to allocate
	// the prefix, because go-netbox is generated from the Netbox OpenApi documentation.
	// The NetBox OpenApi documentation is incorrect for the Available Prefixes APIs
	// (https://github.com/netbox-community/netbox/issues/2769), and
	// go-netbox is generated from the OpenApi decorators. Someone started a patch
	// at https://github.com/hellerve/netbox/commit/97e35a3194b21b71d461862a8a9bc0e174c387f0,
	// but it has not been completed.

	config := meta.(*ProviderNetboxClient).configuration
	log.Printf("[DEBUG] Configuration [%v]\n", config)
	url := "http://" + config.Endpoint + "/api/ipam/prefixes/" + strconv.Itoa(poolId) + "/available-prefixes/"

	jsonData := map[string] interface{}{"prefix_length": prefixLength}
	if len(taglist) != 0 {
		jsonData["tags"] = taglist
	}
	jsonValue, _ := json.Marshal(jsonData)
	log.Printf("[DEBUG] JSON: [%v]\n", string(jsonValue))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("[ERROR] Error occurred creating POST request to url [%v]\n", url)
		log.Printf("[ERROR] Error: %v\n", err)
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("authorization", "Token " + config.AppID)
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/json")

	log.Println("[DEBUG] http.Client submitting request")
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("[DEBUG] Error occurred after post")
		log.Printf("[ERROR] Error retorned from http. %v \n", err)
		return err
	}
	log.Printf("[DEBUG] Http Code Response: %v\n", resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("[DEBUG] Body Read")
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return errors.New("Return http code: " + strconv.Itoa(resp.StatusCode))
	}

	var respData map[string]interface{}
	log.Print("[DEBUG] Will Unmarshal the body")
	jsonErr := json.Unmarshal(body, &respData)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return errors.New("json.Unmarshal failed")
	}
	log.Println("[DEBUG] Unmarshal Ok")
	log.Printf("[DEBUG] Body: %v", string(body))

	// Numbers in json.Unmarshal() are always treated as float64.
	d.SetId(strconv.FormatFloat(respData["id"].(float64), 'f', -1, 64))
	d.Set("prefix", respData["prefix"].(string))
	d.Set("prefix_id", int(respData["id"].(float64)))

	return nil
}

func resourceNetboxPoolPrefixesRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceNetboxPoolPrefixesRead ............ ")
	client := meta.(*ProviderNetboxClient).client

	if d.Id() == "" {
		return errors.New("Id must be provided")
	}
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return errors.New("Could not convert ID to int")
	}
	var parm = ipam.NewIPAMPrefixesReadParams().WithID(id)
	result, err := client.IPAM.IPAMPrefixesRead(parm, nil)
	if err != nil {
		d.SetId("")
	} else {
		prefix := result.Payload
		d.Set("prefix", prefix.Prefix)
		d.Set("prefix_id", prefix.ID)
		d.Set("ring", prefix.Vrf.Name)
		tagMap := make(map[string]string)
		for _, tagPair := range prefix.Tags {
			parts := strings.Split(tagPair, "=")
			// Only process tags that are of the pattern key=value
			if len(parts) == 2 {
				tagMap[parts[0]] = parts[1]
			}
		}
		d.Set("tags", tagMap)
	}
	return nil
}

// https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
func resourceNetboxPoolPrefixesUpdate(d *schema.ResourceData, meta interface{}) error {

	// We can only change the set of tags without having to recreate the entire
	// perfix. All other parameter changes require the prefix to be recreated.
	if d.HasChange("tags") {
		tagList := tagMapToList(d)
		if len(tagList) != 0 {
			if d.Id() == "" {
				return errors.New("Id must be provided")
			}
			if d.Get("prefix").(string) == "" {
				return errors.New("Prefix must be provided")
			}
			prefixString := d.Get("prefix").(string)
			id, err := strconv.ParseInt(d.Id(), 10, 64)
			if err != nil {
				return errors.New("Could not convert ID to int")
			}
			client := meta.(*ProviderNetboxClient).client
			// Set the new tags in the data to send with the update
			data := models.WritablePrefix{ Prefix: &prefixString, Tags: tagList }
			updateParams := ipam.NewIPAMPrefixesUpdateParams().WithID(id).WithData(&data)
			_, err = client.IPAM.IPAMPrefixesUpdate(updateParams, nil)
			if err != nil {
				log.Printf("[ERROR] failed to update prefix tags\n")
				return err
			}
		}
	}
	return resourceNetboxPoolPrefixesRead(d, meta)
}

func resourceNetboxPoolPrefixesDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceNetboxPoolPrefixesDelete ............ ")
	log.Printf("[DEBUG] d -> [%v]\n", d)

	if d.Id() != "" {
		id, err := strconv.ParseInt(d.Id(), 10, 64)
		if err != nil {
			return errors.New("Could not convert ID to int")
		}
		var parm = ipam.NewIPAMPrefixesDeleteParams().WithID(id)
		log.Printf("[DEBUG] Deleting prefix with ID %d\n", id)

		client := meta.(*ProviderNetboxClient).client
		_, err = client.IPAM.IPAMPrefixesDelete(parm, nil)
		log.Println("[DEBUG] Executing Delete.")
		if err == nil {
			log.Printf("[DEBUG] Prefix with ID %d deleted\n", id)
		} else {
			log.Printf("error calling IPAMPrefixesDelete: %v\n", err)
			return err
		}
	} else {
		return errors.New("Id must be provided")
	}
	return nil
}
