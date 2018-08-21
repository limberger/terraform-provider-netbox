package netbox

const testAccDataSourceNetboxPrefixesConfig = `
data "netbox_vlans" "vlans_by_name" {
  name="VLAN-16"
}
`
