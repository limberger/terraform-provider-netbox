
provider "netbox" {
	app_id = "2fe35cabcfe231ebc8734a798f1cac63439a7a2b"
	endpoint = "172.17.133.10:80"
}

# Busca pelo VLAN ID
data "netbox_prefixes" "search_by_vid" {
  vlan_vid = 16
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


# Testing VLAN
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


# Criando o IP

data "netbox_prefixes_available_ips" "next_address" {
	prefixes_id = "${data.netbox_prefixes.prefixes.prefixes_id}"
	description = "IP requisitado via Terraform Realmente Novo"
}
