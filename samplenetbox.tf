provider "netbox" {
	app_id = "2fe35cabcfe231ebc8734a798f1cac63439a7a2b"
	endpoint = "172.17.133.10"
}

data "netbox_prefixes" "prefixes" {
	prefixes_id = 1
}

data "netbox_first_free_address" "next_address" {
	prefixes_id = "${data.netbox_prefixes.prefixes.prefixes_id}"
}



#data "netbox_subnet" "subnet" {
#  subnet_address = "10.10.2.0"
#  subnet_mask    = 24
#}


#data "netbox_first_free_address" "next_address" {
#  subnet_id = "${data.netbox_subnet.subnet.subnet_id}"
#}

#resource "netbox_address" {
#	subnet_id = "${data.netbox_subnet.subnet.subnet_id}"
#	ip_address = "${data.netbox_first_free_address.next_address.ip_address}"
#	hostame = "tf-test-host.io.bb.com.br"
#	description = "Managed by Terraform"
#
#	lifecycle {
#		ignore_changes = [
#			"subnet_id",
#			"ip_address",
#		]
#	}
#}
