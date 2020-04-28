// Package netbox provides access to NetBox services.
// The only resource in this version is a custom one to allow allocating prefixes
// from a specified range of IP supernets, referred to as pools, in specified VRFs.
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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func resourcePoolPrefixesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pool": &schema.Schema{
			Type:     schema.TypeString,
			Description: "CIDR block of the pool from which the prefix will be allocated",
			Required: true,
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
			Required: true,
		},
		"environment": &schema.Schema{
			Type: schema.TypeString,
			Description: "The environment the prefix belongs to",
			Required: true,
			ForceNew: true,
		},
	}
}

var netboxMinValidPrefix = 18
var netboxMaxValidPrefix = 28

// isLengthValid checks that the requested prefix lengths falls in the specified bounds.
func isLengthValid(length int) bool {
	return length >= netboxMinValidPrefix && length <= netboxMaxValidPrefix
}

// isPoolValid checks that the requested prefix pool is one that we support.
func isPoolValid(pool string) bool {
	isValid := false
	validPools := []string{"10.0.0.0/8", "172.16.0.0/12", "100.64.0.0/10"}
	for _, p := range validPools {
		if p == pool {
			isValid = true
			break
		}
	}
	return isValid
}

// tagMapToSlice converts the input tags, represented as a map, to a slice of strings,
// which is the way that NetBox handles them. A map element like tag1: "value1" will
// be converted to the tag string "tag1=value1".
func tagMapToSlice(tags map[string]interface{}) []string {
	var tagSlice []string

	// Add tags
	for key, value := range tags {
		tagSlice = append(tagSlice, fmt.Sprintf("%s=%s", key, value))

	}
	return tagSlice
}

// tagSliceToMap converts the NetBox representation of tags, a slice of strings,
// into a map. Only tags of the form 'key=value' are represented in the map version.
func tagSliceToMap(tagSlice []string) map[string]string {
	tagMap := make(map[string]string)
	for _, tagPair := range tagSlice {
		parts := strings.Split(tagPair, "=")
		// Only process tags that are of the pattern key=value
		if len(parts) == 2 {
			tagMap[parts[0]] = parts[1]
		}
	}
	return tagMap
}

// isTagMapValid verifies that the required tags name and unique are included.
// Other tags may be added as well.
func isTagMapValid(tags map[string]interface{}) bool {
	nameFound := false
	uniqueFound := false
	if len(tags) != 0 {
		for key, _ := range tags {
			if key == "name" {
				nameFound = true
			}
			if key == "unique" {
				uniqueFound = true
			}
		}
	}
	return nameFound && uniqueFound
}

// verifyCreateInput checks that the necessary elements for allocating a prefix are
// specified.
func verifyCreateInput(d *schema.ResourceData) error {
	if d.Get("prefix_length").(int) == 0 {
		log.Println("[ERROR] prefix_length not specified.")
		return errors.New("prefix_length not specified")
	}
	if d.Get("environment").(string) == "" {
		log.Println("[ERROR] environment not specified.")
		return errors.New("environment not specified")
	}
	if d.Get("pool").(string) == "" {
		log.Println("[ERROR] pool not specified.")
		return errors.New("pool not specified")
	}
	return nil
}

// extractInputId is used when accessing an already-existing prefix, and if the ID
// matches expectations will return it.
func extractInputId(d *schema.ResourceData) (int64, error) {
	if d.Id() == "" {
		return -1, errors.New("Id must be provided")
	}
	return strconv.ParseInt(d.Id(), 10, 64)
}

// resourceNetboxPoolPrefixesCreate is invoked when a new prefix will be allocated.
func resourceNetboxPoolPrefixesCreate(d *schema.ResourceData, meta interface{}) error {

	log.Println("[DEBUG] dataNetboxPoolPrefixesCreate")
	client := meta.(*ProviderNetboxClient).client

	err := verifyCreateInput(d)
	if err != nil {
		return err
	}
	prefixLength := d.Get("prefix_length").(int)
	if !isLengthValid(prefixLength) {
		log.Println("[ERROR] invalid prefix length specified.")
		return errors.New("prefix_length must be between 18 & 28, inclusive")
	}
	environment := d.Get("environment").(string)
	pool := d.Get("pool").(string)

	if !isPoolValid(pool) {
		log.Println("[ERROR] invalid pool specified.")
		return errors.New("invalid pool specified")
	}

	if d.Get("tags").(map[string]interface{}) == nil {
		log.Println("[ERROR] tags not specified.")
		return errors.New("tags not specified")

	}
	tags, _ := d.Get("tags").(map[string]interface{})
	if !isTagMapValid(tags) {
		log.Println("[ERROR] required tags (name, unique) missing.")
		return errors.New("required tags (name, unique) missing")
	}

	tagSlice := tagMapToSlice(tags)

	// Find the ID of the VRF. Having trouble just using the vrf name for some reason.
	// XXX: ^^^
	environmentParm := ipam.NewIPAMVrfsListParams().WithName(&environment)
	environmentResult, err := client.IPAM.IPAMVrfsList(environmentParm, nil)
	if err != nil {
		return err
	}
	if *environmentResult.Payload.Count != 1 {
		return errors.New(fmt.Sprintf("Found %d vrfs for environment %s\n", *environmentResult.Payload.Count, environment))
	}
	environmentId := strconv.FormatInt(environmentResult.Payload.Results[0].ID, 10)

	// We need to find the poolId
	ispool := "true"
	listParm := ipam.NewIPAMPrefixesListParams().WithPrefix(&pool).WithIsPool(&ispool).WithVrfID(&environmentId)
	found, err := client.IPAM.IPAMPrefixesList(listParm, nil)
	if err != nil {
		return err
	}
	if *found.Payload.Count != 1 {
		return errors.New(fmt.Sprintf("Found %d pools for prefix %s\n", *found.Payload.Count, pool))
	}
	poolId := int(found.Payload.Results[0].ID)

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
	if len(tagSlice) != 0 {
		jsonData["tags"] = tagSlice
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

// resourceNetboxPoolPrefixesRead is invoked when an existing prefix is being read.
func resourceNetboxPoolPrefixesRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceNetboxPoolPrefixesRead ............ ")
	client := meta.(*ProviderNetboxClient).client

	id, err := extractInputId(d)
	if err != nil {
		return err
	}
	parm := ipam.NewIPAMPrefixesReadParams().WithID(id)
	result, err := client.IPAM.IPAMPrefixesRead(parm, nil)
	if err != nil {
		d.SetId("")
	} else {
		prefix := result.Payload
		d.Set("prefix", prefix.Prefix)
		d.Set("prefix_id", prefix.ID)
		d.Set("environment", prefix.Vrf.Name)
		d.Set("tags", tagSliceToMap(prefix.Tags))
	}
	return nil
}

// resourceNetboxPoolPrefixesUpdate is used to modify an existing prefix. Only a change in tags
// will enable an update. Other input changes, e.g., prefix length or pool, require any existing prefix
// to be deleted and then recreated.
// https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
func resourceNetboxPoolPrefixesUpdate(d *schema.ResourceData, meta interface{}) error {

	// We can only change the set of tags without having to recreate the entire
	// perfix. All other parameter changes require the prefix to be recreated.
	if d.HasChange("tags") {
		tags, _ := d.Get("tags").(map[string]interface{})
		if !isTagMapValid(tags) {
			log.Println("[ERROR] required tags missing. name and unique are required.")
			return errors.New("required tags missing. name and unique are required.")
		}

		tagSlice := tagMapToSlice(tags)
		if len(tagSlice) != 0 {
			id, err := extractInputId(d)
			if err != nil {
				return err
			}
			if d.Get("prefix").(string) == "" {
				return errors.New("Prefix must be provided")
			}
			prefixString := d.Get("prefix").(string)
			client := meta.(*ProviderNetboxClient).client
			// Set the new tags in the data to send with the update
			data := models.WritablePrefix{ Prefix: &prefixString, Tags: tagSlice }
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

// resourceNetboxPoolPrefixesDelete will remove an existing prefix.
func resourceNetboxPoolPrefixesDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceNetboxPoolPrefixesDelete ............ ")
	log.Printf("[DEBUG] d -> [%v]\n", d)

	id, err := extractInputId(d)
	if err != nil {
		return err
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
	return nil
}
