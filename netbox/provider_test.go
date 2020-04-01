package netbox

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]terraform.ResourceProvider

const envErrMsg = `NETBOX_APP_ID, NETBOX_ENDPOINT_ADDR must be set for acceptance tests`

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"netbox": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	switch {
	case os.Getenv("NETBOX_APP_ID") == "":
		t.Fatal(envErrMsg)
	case os.Getenv("NETBOX_ENDPOINT_ADDR") == "":
		t.Fatal(envErrMsg)
	}
}

