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
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/cmgreivel/go-netbox/netbox/client/ipam"
	"github.com/cmgreivel/go-netbox/netbox/models"
)

// ********** NOTE **********
// There is a very special case for servicedelivery. Generally a non-standard
// environment is translated to the pre-dev VRF, but servicedelivery
// is its own VRF, regardless of the specified environment. This can be problematic
// if we want a pre-dev allocation for servicedelivery concerns. We don't
// want to consume 'real' servicedelivery IP ranges.
//
// If resource_type is servicedelivery:
//   - If environment is dev, test, stage, production, then VRF is servicedelivery
//   - If environment is servicedelivery, then VRF is servicedelivery
//   - If environment is anything else, then VRF is pre-dev
//
// If resource_type is any valid type other than servicedelivery:
//   - If environment is dev, test, stage, production, then VRF is the specified environment
//   - If environment is servicedelivery, then an error is returned. This is an invalid combination.
//   - If environment is anything else, then VRF is pre-dev


type basicPrefix struct {
	Prefix string
	Id float64
}

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
		"resource_type": &schema.Schema{
			Type:        schema.TypeString,
			Description: "Type of resource for which a prefix is being allocated, determines supernet",
			Required:    true,
			ForceNew:    true,
		},
		"prefix_id": &schema.Schema{
			Type:        schema.TypeInt,
			Description: "ID of the prefix allocated with the resource (output from NetBox)",
			Optional:    true,
			Computed:    true,
		},
		"prefix": &schema.Schema{
			Type:        schema.TypeString,
			Description: "CIDR block of the prefix allocated with the resource",
			Optional:    true,
			Computed:    true,
		},
		"prefix_length": &schema.Schema{
			Type:        schema.TypeInt,
			Description: "Length of the prefix in bits",
			Required:    true,
			ForceNew:    true,
		},
		"tags": &schema.Schema{
			Type:        schema.TypeMap,
			Description: "Tags applied to the prefix",
			Required:    true,
		},
		"environment": &schema.Schema{
			Type:        schema.TypeString,
			Description: "The environment the prefix belongs to",
			Required:    true,
			ForceNew:    true,
		},
	}
}

const netboxMinValidPrefix = 18
const netboxMaxValidPrefix = 28

var envVrfs = [4]string{"dev", "test", "stage", "prod"}

// Resource types that rely on env to determine VRF
var envResourceTypes = [3]string{"core", "depot", "edge"}

// Usually gets sdVrf for VRF, except for pre-dev environment
const sdResourceType = "servicedelivery"
const sdVrf = "servicedelivery"

// Map of resource types to supernet ranges.
var resourceSupernetMap = map[string][]string {
	"core":            {"100.64.0.0/10"},
	"depot":           {"10.224.0.0/16", "10.225.0.0/16"},
	"edge":            {"10.226.0.0/16", "10.227.0.0/16"},
	"servicedelivery": {"10.228.0.0/16"},
	"vpn":             {"172.16.0.0/12"},
}

// isLengthValid checks that the requested prefix lengths falls in the specified bounds.
func isLengthValid(length int) bool {
	return length >= netboxMinValidPrefix && length <= netboxMaxValidPrefix
}

// Returns the correct slice of supernets for a given resource type.
// Returns empty slice if an unexpected resource type is specified.
func resourceTypeToSupernets(resourceType string) [](string) {
	lcType := strings.ToLower(resourceType)
	return resourceSupernetMap[lcType]
}

// Verify that an environment of 'servicedelivery' is not specified
// with a resourceType that is not 'servicedelivery' to avoid possible
// confusion and bad allocation of prefixes.
func areResourceAndEnvValid(resourceType string, environment string) bool {
	if strings.ToLower(environment) != "servicedelivery" {
		return true
	}
	if strings.ToLower(resourceType) == "servicedelivery" {
		// Both are servicedelivery, so no confusion here.
		return true
	}
	// environment is servicedelivery, but resourceType is not.
	return false
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

// CMG: Add sanity check on input environment and resource_type for servicedelivery
// My test had servicedelivery for environment and core for resource_type, and it
// caused oddities--allocated for pre-dev core!

// envToVrf converts the environment passed in to the VRF that manages it in NetBox.
// If it doesn't match any of our recognized environment patterns, default
// to the pre-dev (super temp) environment.
func envToVrf(env string) string {
	// If the passed-in env exactly matches a vrf, use that.
	lowerEnv := strings.ToLower(env)
	for _, vrf := range envVrfs {
		if lowerEnv == vrf {
			return vrf
		}
	}
	// Some special checks to allow flexibility in environment spec.
	prodRe := regexp.MustCompile(`^prod\d`)
	matchedProd := prodRe.MatchString(lowerEnv)
	if matchedProd || lowerEnv == "production" {
		// "production," "prod0", "prod1", etc.
		return "prod"
	}
	if lowerEnv == "staging" {
		return "stage"
	}
	return "pre-dev"
}

// getVrf converts the environment passed in to the VRF
// that manages it in NetBox. The typical mapping is based on the environment,
// but service delivery resources use the same VRF for every environment.
// If it doesn't match any of our recognized environment patters, default
// to the pre-dev (super temp) environment.
func getVrf(env string, resourceType string) string {
	envVrf := envToVrf(env)
	// If this isn't one of our standard environments, skip the servicedelivery check
	if envVrf != "pre-dev" && strings.ToLower(resourceType) == sdResourceType {
		return sdVrf
	}
	return envVrf
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
	if d.Get("resource_type").(string) == "" {
		log.Println("[ERROR] resource_type not specified.")
		return errors.New("resource_type not specified")
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

// allocatePrefix does the actual call
func allocatePrefix(provider *ProviderNetboxClient, pool *string, vrfId *string,
	prefixLength int, tags []string) (basicPrefix, error) {

	var emptyPrefix basicPrefix

	ispool := "true"
	listParm := ipam.NewIpamPrefixesListParams().WithPrefix(pool).WithIsPool(&ispool).WithVrfID(vrfId)
	found, err := provider.client.Ipam.IpamPrefixesList(listParm, nil)
	if err != nil {
		return emptyPrefix, err
	}
	if *found.Payload.Count != 1 {
		return emptyPrefix, errors.New(fmt.Sprintf("Found %d pools for prefix %s\n", *found.Payload.Count, *pool))
	}
	poolId := int(found.Payload.Results[0].ID)

	// XXX: Should consider fixing this in https://github.com/cmgreivel/go-netbox.
	// We cannot use go-netbox (https://github.com/netbox-community/go-netbox) to allocate
	// the prefix, because go-netbox is generated from the Netbox OpenApi documentation.
	// The NetBox OpenApi documentation is incorrect for the Available Prefixes APIs
	// (https://github.com/netbox-community/netbox/issues/2769), and
	// go-netbox is generated from the OpenApi decorators. Someone started a patch
	// at https://github.com/hellerve/netbox/commit/97e35a3194b21b71d461862a8a9bc0e174c387f0,
	// but it has not been completed.

	config := provider.configuration
	log.Printf("[DEBUG] Configuration [%v]\n", config)
	url := "http://" + config.Endpoint + "/api/ipam/prefixes/" + strconv.Itoa(poolId) + "/available-prefixes/"

	jsonData := map[string]interface{}{"prefix_length": prefixLength}
	if len(tags) != 0 {
		jsonData["tags"] = tags
	}
	jsonValue, _ := json.Marshal(jsonData)
	log.Printf("[DEBUG] JSON: [%v]\n", string(jsonValue))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("[ERROR] Error occurred creating POST request to url [%v]\n", url)
		log.Printf("[ERROR] Error: %v\n", err)
		return emptyPrefix, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("authorization", "Token "+config.AppID)
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/json")

	log.Println("[DEBUG] http.Client submitting request")
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("[DEBUG] Error occurred after post")
		log.Printf("[ERROR] Error retorned from http. %v \n", err)
		return emptyPrefix, err
	}
	log.Printf("[DEBUG] Http Code Response: %v\n", resp.StatusCode)

	// CMG: What will be the result of a full prefix here.
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("[DEBUG] Body Read")
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return emptyPrefix, errors.New("Return http code: " + strconv.Itoa(resp.StatusCode))
	}

	var respData map[string]interface{}
	log.Print("[DEBUG] Will Unmarshal the body")
	jsonErr := json.Unmarshal(body, &respData)
	if jsonErr != nil {
		log.Print(jsonErr)
		return emptyPrefix, errors.New("json.Unmarshal failed")
	}
	log.Println("[DEBUG] Unmarshal Ok")
	log.Printf("[DEBUG] Body: %v", string(body))

	// Numbers in json.Unmarshal() are always treated as float64.
	return basicPrefix{respData["prefix"].(string), respData["id"].(float64)}, nil
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
		return errors.New(fmt.Sprintf("prefix_length must be between %d & %d, inclusive",
			netboxMinValidPrefix, netboxMaxValidPrefix))
	}
	environment := d.Get("environment").(string)
	resourceType := d.Get("resource_type").(string)

	if !areResourceAndEnvValid(resourceType, environment) {
		e := fmt.Sprintf("resource %s cannot be used with environment %s", resourceType, environment)
		log.Printf("[ERROR] %s", e)
		return errors.New(e)
	}
	// We rely on getting the supernets to determine resourceType validity.
	pools := resourceTypeToSupernets(resourceType)
	if len(pools) == 0 {
		e := fmt.Sprintf("no pools found for resource %s.", resourceType)
		log.Printf("[ERROR] %s", e)
		return errors.New(e)
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

	vrf := getVrf(environment, resourceType)
	// Find the ID of the VRF. Having trouble just using the vrf name for some reason.
	// XXX: ^^^
	vrfParm := ipam.NewIpamVrfsListParams().WithName(&vrf)
	vrfResult, err := client.Ipam.IpamVrfsList(vrfParm, nil)
	if err != nil {
		return err
	}
	if *vrfResult.Payload.Count != 1 {
		return errors.New(fmt.Sprintf("Found %d vrfs for environment %s (vrf %s)\n",
			*vrfResult.Payload.Count, environment, vrf))
	}
	vrfId := strconv.FormatInt(vrfResult.Payload.Results[0].ID, 10)

	var prefixInfo basicPrefix
	for _, pool := range pools {
		prefixInfo, err = allocatePrefix(meta.(*ProviderNetboxClient), &pool, &vrfId, prefixLength, tagSlice)
		if err != nil {
			return err
		}
		if len(prefixInfo.Prefix) > 0 {
			break
		}
		log.Printf("Did not allocate prefix for %s in %s from pool %s",
			resourceType, environment, pool)
	}
	if len(prefixInfo.Prefix) == 0 {
		return errors.New(fmt.Sprintf("Did not allocate prefix for %s in %s", resourceType, environment))
	}
	// Set the appropriate values in the Terraform-managed structure
	d.SetId(strconv.FormatFloat(prefixInfo.Id, 'f', -1, 64))
	d.Set("prefix", prefixInfo.Prefix)
	d.Set("prefix_id", prefixInfo.Id)
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
	parm := ipam.NewIpamPrefixesReadParams().WithID(id)
	result, err := client.Ipam.IpamPrefixesRead(parm, nil)
	if err != nil {
		d.SetId("")
	} else {
		prefix := result.Payload
		// We do not set the environment here and let TF report it
		d.Set("prefix", prefix.Prefix)
		d.Set("prefix_id", prefix.ID)
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
			data := models.WritablePrefix{Prefix: &prefixString, Tags: tagSlice}
			updateParams := ipam.NewIpamPrefixesUpdateParams().WithID(id).WithData(&data)
			_, err = client.Ipam.IpamPrefixesUpdate(updateParams, nil)
			if err != nil {
				log.Printf("[ERROR] failed to update prefix tags: %v\n", err)
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
	var parm = ipam.NewIpamPrefixesDeleteParams().WithID(id)
	log.Printf("[DEBUG] Deleting prefix with ID %d\n", id)

	client := meta.(*ProviderNetboxClient).client
	_, err = client.Ipam.IpamPrefixesDelete(parm, nil)
	log.Println("[DEBUG] Executing Delete.")
	if err == nil {
		log.Printf("[DEBUG] Prefix with ID %d deleted\n", id)
	} else {
		log.Printf("[ERROR] error calling IpamPrefixesDelete: %v\n", err)
		return err
	}
	return nil
}
