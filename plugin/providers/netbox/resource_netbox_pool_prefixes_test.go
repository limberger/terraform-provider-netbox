package netbox

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
)

const testAccResourceNetboxPoolPrefixesConfig = `
resource "netbox_pool_prefixes" "test_prefix" {
  ring          = "test"
  pool          = "10.0.0.0/8"
  prefix_length = 28
  tags = {
    name   = "some_name"
    unique = "some_unique_string"
  }
}
`

const testAccResourceNetboxPoolPrefixesEditTags = `
resource "netbox_pool_prefixes" "test_prefix" {
  ring          = "test"
  pool          = "10.0.0.0/8"
  prefix_length = 28
  tags = {
    name   = "new_name"
    unique = "new_unique_string"
  }
}
`

const testAccResourceNetboxPoolPrefixesEditRing = `
resource "netbox_pool_prefixes" "test_prefix" {
  ring          = "dev"
  pool          = "10.0.0.0/8"
  prefix_length = 28
  tags = {
    name   = "new_name"
    unique = "new_unique_string"
  }
}
`

const testAccResourceNetboxPoolPrefixesEditPool = `
resource "netbox_pool_prefixes" "test_prefix" {
  ring          = "dev"
  pool          = "172.16.0.0/12"
  prefix_length = 28
  tags = {
    name   = "new_name"
    unique = "new_unique_string"
  }
}
`

const testAccResourceNetboxPoolPrefixesEditLength = `
resource "netbox_pool_prefixes" "test_prefix" {
  ring          = "dev"
  pool          = "172.16.0.0/12"
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
  ring          = "test"
  pool          = "172.16.0.0/12"
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
		parm := ipam.NewIPAMPrefixesReadParams().WithID(id)
		result, err := client.IPAM.IPAMPrefixesRead(parm, nil)
		if err == nil {
			if result.Payload != nil && result.Payload.ID == id {
				return fmt.Errorf("Prefix with ID %s still exists.", rs.Primary.ID)
			}
			return nil
		}
		// CMG: Need to discriminate between prefix not there and other errors (like connection).
	}
	return nil
}

func testAccCheckPrefixExists(resourceName string, prefix *models.Prefix) resource.TestCheckFunc {
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
		parm := ipam.NewIPAMPrefixesReadParams().WithID(id)
		result, err := client.IPAM.IPAMPrefixesRead(parm, nil)
        if err != nil {
            return err
        }

		if result.Payload == nil {
			return fmt.Errorf("Empty payload for prefix with ID %s.", rs.Primary.ID)
		}

		// CMG: Can add additional checks later

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
		// CMG: Can add additional checks later
        return nil
    }
}

func testAccCheckPrefixRecreated(orig *models.Prefix, new *models.Prefix) resource.TestCheckFunc {
    return func(s *terraform.State) error {
		if orig.ID == new.ID {
			return fmt.Errorf("Prefix ID did not change from %d to %d.", orig.ID, new.ID)
		}
		// CMG: Can add additional checks later
        return nil
    }
}

func TestAccResourceNetboxPoolPrefixes_allocate(t *testing.T) {
	var prefix1 models.Prefix
	resource.Test(t, resource.TestCase {
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccPrefixDestroy,
		Steps: []resource.TestStep {
			{
                Config: testAccResourceNetboxPoolPrefixesConfig,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckPrefixExists("netbox_pool_prefixes.test_prefix", &prefix1),
                ),
            },
			{
                Config: testAccResourceNetboxPoolPrefixesOtherPool,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckPrefixExists("netbox_pool_prefixes.other_pool", &prefix1),
                ),
            },
		},
	})
}

func TestAccResourceNetboxPoolPrefixes_edit(t *testing.T) {
	var prefix1, prefix2 models.Prefix
	resource.Test(t, resource.TestCase {
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccPrefixDestroy,
		Steps: []resource.TestStep {
			{
                Config: testAccResourceNetboxPoolPrefixesConfig,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckPrefixExists("netbox_pool_prefixes.test_prefix", &prefix1),
                ),
            },
            {
                Config: testAccResourceNetboxPoolPrefixesEditTags,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckPrefixExists("netbox_pool_prefixes.test_prefix", &prefix2),
					testAccCheckPrefixModified(&prefix1, &prefix2),
                ),
            },
            {
                Config: testAccResourceNetboxPoolPrefixesEditRing,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckPrefixExists("netbox_pool_prefixes.test_prefix", &prefix1),
					testAccCheckPrefixRecreated(&prefix2, &prefix1),
                ),
            },
            {
                Config: testAccResourceNetboxPoolPrefixesEditPool,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckPrefixExists("netbox_pool_prefixes.test_prefix", &prefix2),
					testAccCheckPrefixRecreated(&prefix1, &prefix2),
                ),
            },
            {
                Config: testAccResourceNetboxPoolPrefixesEditLength,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckPrefixExists("netbox_pool_prefixes.test_prefix", &prefix1),
					testAccCheckPrefixRecreated(&prefix2, &prefix1),
                ),
            },
		},
	})
}

const testAccResourceNetboxPoolPrefixesBadRing = `
resource "netbox_pool_prefixes" "test_prefix" {
  ring          = "no-such-ring"
  pool          = "10.0.0.0/8"
  prefix_length = 28
  tags = {
    name   = "some_name"
    unique = "some_unique_string"
  }
}
`

func TestAccResourceNetboxPoolPrefixes_ringErrors(t *testing.T) {
	ringError, _ := regexp.Compile("ring")

	resource.Test(t, resource.TestCase {
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep {
			{
                Config: testAccResourceNetboxPoolPrefixesBadRing,
				ExpectError: ringError,
            },
		},
	})
}

const testAccResourceNetboxPoolPrefixesBadPool = `
resource "netbox_pool_prefixes" "test_prefix" {
  ring          = "dev"
  pool          = "10.2.0.0/8"
  prefix_length = 28
  tags = {
    name   = "some_name"
    unique = "some_unique_string"
  }
}
`

func TestAccResourceNetboxPoolPrefixes_poolErrors(t *testing.T) {
	poolError, _ := regexp.Compile("pool")

	resource.Test(t, resource.TestCase {
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep {
			{
                Config: testAccResourceNetboxPoolPrefixesBadPool,
				ExpectError: poolError,
            },
		},
	})
}

const testAccResourceNetboxPoolPrefixesLengthTooBig = `
resource "netbox_pool_prefixes" "test_prefix" {
  ring          = "dev"
  pool          = "10.0.0.0/8"
  prefix_length = 16
  tags = {
    name   = "some_name"
    unique = "some_unique_string"
  }
}
`

const testAccResourceNetboxPoolPrefixesLengthTooSmall = `
resource "netbox_pool_prefixes" "test_prefix" {
  ring          = "dev"
  pool          = "10.0.0.0/8"
  prefix_length = -7
  tags = {
    name   = "some_name"
    unique = "some_unique_string"
  }
}
`

func TestAccResourceNetboxPoolPrefixes_lengthErrors(t *testing.T) {
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
  ring          = "dev"
  pool          = "10.0.0.0/8"
  prefix_length = 26
  tags = {
    unique = "some_unique_string"
  }
}
`

const testAccResourceNetboxPoolPrefixesNoUniqueTag = `
resource "netbox_pool_prefixes" "test_prefix" {
  ring          = "dev"
  pool          = "10.0.0.0/8"
  prefix_length = 26
  tags = {
    name = "I have a name"
  }
}
`

func TestAccResourceNetboxPoolPrefixes_tagErrors(t *testing.T) {
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
