package netbox

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"testing"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/cmgreivel/go-netbox/netbox/client/ipam"
	"github.com/cmgreivel/go-netbox/netbox/models"
)

const testPrefixName = "netbox_pool_prefixes.test_prefix"

const testAccResourceNetboxPoolPrefixesConfig = `
resource "netbox_pool_prefixes" "test_prefix" {
  environment   = "test"
  resource_type = "core"
  prefix_length = 28
  tags = {
	name   = "some_name"
	unique = "some_unique_string"
  }
}
`

const testAccResourceNetboxPoolPrefixesEditTags = `
resource "netbox_pool_prefixes" "test_prefix" {
  environment   = "test"
  resource_type = "core"
  prefix_length = 28
  tags = {
	name   = "new_name"
	unique = "new_unique_string"
  }
}
`

const testAccResourceNetboxPoolPrefixesEditEnvironment = `
resource "netbox_pool_prefixes" "test_prefix" {
  environment   = "dev"
  resource_type = "core"
  prefix_length = 28
  tags = {
	name   = "new_name"
	unique = "new_unique_string"
  }
}
`

const testAccResourceNetboxPoolPrefixesEditResourceType = `
resource "netbox_pool_prefixes" "test_prefix" {
  environment   = "dev"
  resource_type = "vpn"
  prefix_length = 28
  tags = {
	name   = "new_name"
	unique = "new_unique_string"
  }
}
`

const testAccResourceNetboxPoolPrefixesEditLength = `
resource "netbox_pool_prefixes" "test_prefix" {
  environment   = "dev"
  resource_type = "vpn"
  prefix_length = 26
  tags = {
	name   = "new_name"
	unique = "new_unique_string"
  }
}
`

// This config has no order significance
const testAccResourceNetboxPoolPrefixesOtherPool = `
resource "netbox_pool_prefixes" "other_pool" {
  environment   = "test"
  resource_type = "vpn"
  prefix_length = 28
  tags = {
	name   = "some_name"
	unique = "some_unique_string"
  }
}
`

func testAccPrefixDestroy(s *terraform.State) error {
	// retrieve the client information
	client := testAccProvider.Meta().(*ProviderNetboxClient).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_pool_prefixes" {
			continue
		}
		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return err
		}
		parm := ipam.NewIpamPrefixesReadParams().WithID(id)
		result, err := client.Ipam.IpamPrefixesRead(parm, nil)
		if err == nil {
			if result.Payload != nil && result.Payload.ID == id {
				return fmt.Errorf("Prefix with ID %s still exists.", rs.Primary.ID)
			}
			return nil
		}
	}
	return nil
}

func testAccCheckPrefixExists(resourceName string, prefix *models.Prefix,
	vrf string, resourceType string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Prefix ID is not set")
		}

		// retrieve the client information
		client := testAccProvider.Meta().(*ProviderNetboxClient).client
		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return err
		}
		parm := ipam.NewIpamPrefixesReadParams().WithID(id)
		result, err := client.Ipam.IpamPrefixesRead(parm, nil)
		if err != nil {
			return err
		}

		if result.Payload == nil {
			return fmt.Errorf("Empty payload for prefix with ID %s.", rs.Primary.ID)
		}
		if result.Payload.IsPool {
			// These allocated prefixes should never be pools
			return fmt.Errorf("%s [%s] allocated as pool", *result.Payload.Prefix, rs.Primary.ID)
		}

		if len(vrf) > 0 {
			// If we specify a VRF to cehck, we need to confirm that we've allocated the prefix
			// using that VRF and that it falls in an expected pool.

			vrfParm := ipam.NewIpamVrfsListParams().WithName(&vrf)
			vrfResult, err := client.Ipam.IpamVrfsList(vrfParm, nil)
			if err != nil {
				return err
			}
			if *vrfResult.Payload.Count != 1 {
				return fmt.Errorf("Found %d vrfs %s\n",	*vrfResult.Payload.Count, vrf)
			}
			if vrfResult.Payload.Results[0].ID != result.Payload.Vrf.ID {
				return fmt.Errorf("Prefix has vrf ID %d. Expected vrf %s has ID %d\n",
					result.Payload.Vrf.ID, vrf, vrfResult.Payload.Results[0].ID)
			}
		}

		if len(resourceType) > 0 {
			inSupernet := false

			// This comes from the netbox provider code, so might be bit of a chicken-and-egg
			// issue here.
			supernets := resourceTypeToSupernets(resourceType)
			if len(supernets) <= 0 {
				return fmt.Errorf("No supernets found for resource %s", resourceType)
			}
			_, pnet, err := net.ParseCIDR(*result.Payload.Prefix)
			if err != nil {
				return err
			}
			for _, super := range supernets {
				_, rnet, err := net.ParseCIDR(super)
				if err != nil {
					return err
				}
				err = cidr.VerifyNoOverlap([]*net.IPNet{pnet}, rnet)
				if err == nil {
					// pnet falls in rnet
					inSupernet = true
					break
				}
			}
			if ! inSupernet {
				return fmt.Errorf("%s not in supernets %v", *result.Payload.Prefix, supernets)
			}
		}
		// store the resulting prefix
		*prefix = *result.Payload
		return nil
	}
}

func testAccCheckPrefixModified(orig *models.Prefix, new *models.Prefix) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if orig.ID != new.ID {
			return fmt.Errorf("Prefix ID unexpectedly changed from %d to %d.", orig.ID, new.ID)
		}
		// XXX: Can add additional checks later
		return nil
	}
}

func testAccCheckPrefixRecreated(orig *models.Prefix, new *models.Prefix) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if orig.ID == new.ID {
			return fmt.Errorf("Prefix ID did not change from %d to %d.", orig.ID, new.ID)
		}
		// XXX: Can add additional checks later
		return nil
	}
}

func TestAccResourceNetboxPoolPrefixesAllocate(t *testing.T) {
	var prefix1 models.Prefix
	resource.Test(t, resource.TestCase {
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccPrefixDestroy,
		Steps: []resource.TestStep {
			{
				Config: testAccResourceNetboxPoolPrefixesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrefixExists(testPrefixName, &prefix1, "test", "core"),
				),
			},
			{
				Config: testAccResourceNetboxPoolPrefixesOtherPool,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrefixExists("netbox_pool_prefixes.other_pool", &prefix1, "test", "vpn"),
				),
			},
		},
	})
}

func TestAccResourceNetboxPoolPrefixesEdit(t *testing.T) {
	var prefix1, prefix2 models.Prefix
	resource.Test(t, resource.TestCase {
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccPrefixDestroy,
		Steps: []resource.TestStep {
			{
				Config: testAccResourceNetboxPoolPrefixesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrefixExists(testPrefixName, &prefix1, "test", "core"),
				),
			},
			{
				Config: testAccResourceNetboxPoolPrefixesEditTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrefixExists(testPrefixName, &prefix2, "test", "core"),
					testAccCheckPrefixModified(&prefix1, &prefix2),
				),
			},
			{
				Config: testAccResourceNetboxPoolPrefixesEditEnvironment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrefixExists(testPrefixName, &prefix1, "dev", "core"),
					testAccCheckPrefixRecreated(&prefix2, &prefix1),
				),
			},
			{
				Config: testAccResourceNetboxPoolPrefixesEditResourceType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrefixExists(testPrefixName, &prefix2, "dev", "vpn"),
					testAccCheckPrefixRecreated(&prefix1, &prefix2),
				),
			},
			{
				Config: testAccResourceNetboxPoolPrefixesEditLength,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrefixExists(testPrefixName, &prefix1, "dev", "vpn"),
					testAccCheckPrefixRecreated(&prefix2, &prefix1),
				),
			},
		},
	})
}

const testAccResourceNetboxPoolPrefixesCustomEnvironment = `
resource "netbox_pool_prefixes" "test_prefix" {
  environment   = "no-such-environment"
  resource_type = "core"
  prefix_length = 28
  tags = {
	name   = "some_name"
	unique = "some_unique_string"
  }
}
`

func TestAccResourceNetboxPoolPrefixesCustomEnvironment(t *testing.T) {
	var prefix1 models.Prefix

	resource.Test(t, resource.TestCase {
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep {
			{
				Config: testAccResourceNetboxPoolPrefixesCustomEnvironment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrefixExists(testPrefixName, &prefix1, "pre-dev", "core"),
				),
			},
		},
	})
}

const testAccResourceNetboxPoolPrefixesBadType = `
resource "netbox_pool_prefixes" "test_prefix" {
  environment   = "dev"
  resource_type = "unsupported"
  prefix_length = 28
  tags = {
	name   = "some_name"
	unique = "some_unique_string"
  }
}
`

func TestAccResourceNetboxPoolPrefixesTypeErrors(t *testing.T) {
	poolError, _ := regexp.Compile("no pools found")

	resource.Test(t, resource.TestCase {
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep {
			{
				Config: testAccResourceNetboxPoolPrefixesBadType,
				ExpectError: poolError,
			},
		},
	})
}

const testAccResourceNetboxPoolPrefixesLengthTooBig = `
resource "netbox_pool_prefixes" "test_prefix" {
  environment   = "dev"
  resource_type = "core"
  prefix_length = 16
  tags = {
	name   = "some_name"
	unique = "some_unique_string"
  }
}
`

const testAccResourceNetboxPoolPrefixesLengthTooSmall = `
resource "netbox_pool_prefixes" "test_prefix" {
  environment   = "dev"
  resource_type = "core"
  prefix_length = -7
  tags = {
	name   = "some_name"
	unique = "some_unique_string"
  }
}
`

func TestAccResourceNetboxPoolPrefixesLengthErrors(t *testing.T) {
	lengthError, _ := regexp.Compile("length")

	resource.Test(t, resource.TestCase {
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep {
			{
				Config: testAccResourceNetboxPoolPrefixesLengthTooBig,
				ExpectError: lengthError,
			},
			{
				Config: testAccResourceNetboxPoolPrefixesLengthTooSmall,
				ExpectError: lengthError,
			},
		},
	})
}

const testAccResourceNetboxPoolPrefixesNoNameTag = `
resource "netbox_pool_prefixes" "test_prefix" {
  environment   = "dev"
  resource_type = "core"
  prefix_length = 26
  tags = {
	unique = "some_unique_string"
  }
}
`

const testAccResourceNetboxPoolPrefixesNoUniqueTag = `
resource "netbox_pool_prefixes" "test_prefix" {
  environment   = "dev"
  resource_type = "core"
  prefix_length = 26
  tags = {
	name = "I have a name"
  }
}
`

func TestAccResourceNetboxPoolPrefixesTagErrors(t *testing.T) {
	tagError, _ := regexp.Compile("tags")

	resource.Test(t, resource.TestCase {
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep {
			{
				Config: testAccResourceNetboxPoolPrefixesNoNameTag,
				ExpectError: tagError,
			},
			{
				Config: testAccResourceNetboxPoolPrefixesNoUniqueTag,
				ExpectError: tagError,
			},
		},
	})
}

// Test servicedelivery
func TestAccResourceNetboxPoolPrefixesServiceDeliveryOK(t *testing.T) {
	var prefix models.Prefix

	const baseConfig =
		"resource \"netbox_pool_prefixes\" \"%s\" {\n" +
		"  environment   = \"%s\"\n" +
		"  resource_type = \"servicedelivery\"\n" +
		"  prefix_length = 26\n" +
		"  tags = {\n" +
		"    name   = \"some_name\"\n" +
		"    unique = \"some_unique_string\"\n" +
		"  }\n" +
		"}\n"

	steps := []resource.TestStep{}
	prefixNames := []string{}
	configStrings := []string{}
	for i, env := range []string{"dev", "test", "stage", "prod"} {
		prefixNames = append(prefixNames, fmt.Sprintf("sd_%s", env))
		configStrings = append(configStrings, fmt.Sprintf(baseConfig, prefixNames[i], env))
		steps = append(steps,
			resource.TestStep {
				Config: configStrings[i],
				Check: testAccCheckPrefixExists("netbox_pool_prefixes." + prefixNames[i],
					&prefix, "servicedelivery", "servicedelivery"),
			},
		)
	}

	resource.Test(t, resource.TestCase {
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccPrefixDestroy,
		Steps: steps,
	})
}

func TestAccResourceNetboxPoolPrefixesServiceDeliveryPredev(t *testing.T) {
	var prefix models.Prefix

	resource.Test(t, resource.TestCase {
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccPrefixDestroy,
		Steps: []resource.TestStep {
			{
				Config: "resource \"netbox_pool_prefixes\" \"sd_predev\" {\n" +
					"  environment   = \"custom\"\n" +
					"  resource_type = \"servicedelivery\"\n" +
					"  prefix_length = 26\n" +
					"  tags = {\n" +
					"    name   = \"some_name\"\n" +
					"    unique = \"some_unique_string\"\n" +
					"  }\n" +
					"}\n",
				Check: testAccCheckPrefixExists("netbox_pool_prefixes.sd_predev",
					&prefix, "pre-dev", "servicedelivery"),
			},
		},
	})
}


func TestAccResourceNetboxPoolPrefixesServiceDeliveryInvalid(t *testing.T) {
	sdError, _ := regexp.Compile("environment servicedelivery")

	resource.Test(t, resource.TestCase {
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccPrefixDestroy,
		Steps: []resource.TestStep {
			{
				Config: "resource \"netbox_pool_prefixes\" \"sd_dev\" {\n" +
					"  environment   = \"servicedelivery\"\n" +
					"  resource_type = \"core\"\n" +
					"  prefix_length = 26\n" +
					"  tags = {\n" +
					"    name   = \"some_name\"\n" +
					"    unique = \"some_unique_string\"\n" +
					"  }\n" +
					"}\n",
				ExpectError: sdError,
			},
		},
	})
}
