package netbox

import (
	"testing"
	//  "log"

	"github.com/hashicorp/terraform/helper/resource"
)

const testAccDataSourceNetboxVlansConfig = `
data "netbox_vlans" "vlans_by_name" {
  name = "VLAN 17"
}

data "netbox_vlans" "vlans_by_id" {
  vid = "${data.netbox_vlans.vlans_by_name.vid}"
}

data "netbox_vlans" "vlans_by_name16" {
  name = "VLAN 16"
}

data "netbox_vlans" "vlans_by_id16" {
  vid = "${data.netbox_vlans.vlans_by_name16.vid}"
}
`

//		* data.netbox_vlans.vlans_by_id: Resource 'data.netbox_vlans.vlans_by_name' not found for variable 'data.netbox_vlans.vlans_by_name.vid'

// const testAccDataSourceNetboxVlanCustomFieldConfig = `
// resource "netbox_vlans" "vlan" {
//   vid   = 3
//   name = "VLAN-16"
//
//   custom_fields = {
//     CustomTestVlan  = "vlan-test"
//     CustomTestVlan2 = "vlan2-test"
//   }
// }
//
// data "netbox_vlans" "custom_search" {
//   subnet_id = "${netbox_vlans.vlans_by_id.vid}"
//
//   custom_field_filter = {
//     CustomTestAddresses  = ".*terraform.*"
//     CustomTestAddresses2 = ".*terraform2.*"
//   }
// }
// `

func TestAccDataSourceNetboxVlans(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceNetboxVlansConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.netbox_vlans.vlans_by_name", "vid", "data.netbox_vlans.vlans_by_id", "vid"),
					// resource.TestCheckResourceAttrPair("data.netbox_vlans.address_by_address", "ip_address", "data.netbox_address.address_by_id", "ip_address"),
					// resource.TestCheckResourceAttr("data.netbox_vlans.address_by_address", "description", "Gateway"),
					// resource.TestCheckResourceAttr("data.netbox_vlans", "ip_address", "10.10.1.3"),
					// resource.TestCheckResourceAttr("data.netbox_vlans.address_by_description", "ip_address", "10.10.1.4"),
				),
			},
		},
	})
}

// func TestAccDataSourceNetboxVlans_CustomField(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:  func() { testAccPreCheck(t) },
// 		Providers: testAccProviders,
// 		Steps: []resource.TestStep{
// 			resource.TestStep{
// 				Config: testAccDataSourceNetboxVlanCustomFieldConfig,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("data.netbox_vlans.custom_search", "subnet_id", "3"),
// 					resource.TestCheckResourceAttr("data.netbox_vlans.custom_search", "ip_address", "10.10.1.10"),
// 					resource.TestCheckResourceAttr("data.netbox_vlans.custom_search", "description", "Terraform test address (custom fields)"),
// 					resource.TestCheckResourceAttr("data.netbox_vlans.custom_search", "hostname", "tf-test.cust1.local"),
// 					resource.TestCheckResourceAttr("data.netbox_vlans.custom_search", "custom_fields.CustomTestAddresses", "terraform-test"),
// 					resource.TestCheckResourceAttr("data.netbox_vlans.custom_search", "custom_fields.CustomTestAddresses2", "terraform2-test"),
// 				),
// 			},
// 		},
// 	})
// }
