provider "netbox" {
	app_id = "0123456789abcdef0123456789abcdef01234567"
	endpoint = "0.0.0.0:32768"
}

# Busca pelo VLAN ID
data "netbox_prefixes" "search_by_vid" {
  vlan = {
    vid = 16
  }
}
# Mostra o prefixo
output "search_by_vid_output" {
  value = "${data.netbox_prefixes.search_by_vid.prefix}"
}
data "netbox_prefixes" "prefixes" {
	prefixes_id = 1
}

output "out_prefixes" {
	value = "${data.netbox_prefixes.prefixes.prefix}"
}
output "out_prefixes_created" {
  value = "${data.netbox_prefixes.prefixes.created}"
}
output "out_prefixes_description" {
  value = "${data.netbox_prefixes.prefixes.description}"
}

data "netbox_first_free_address" "next_address" {
	prefixes_id = "${data.netbox_prefixes.prefixes.prefixes_id}"
}
