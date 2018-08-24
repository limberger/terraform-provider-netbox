# Terraform Provider Plugin for Netbox

This repository holds a external plugin for a [Terraform][1] provider to manage
resources within [Netbox][2], an open source IP address management system. Using
the go API client for DigitalOcean's NetBox IPAM and DCIM service [Go-Netbox][3].

Good example: [Example][4]

[1]: https://www.terraform.io/
[2]: https://github.com/digitalocean/netbox
[3]: https://github.com/digitalocean/go-netbox
[4]: http://techblog.d2-si.eu/2018/02/23/my-first-terraform-provider.html

## About Netbox

[Netbox][2] is an open source IP address management system written in Python. Through our Go integration provided by [Go-Netbox[3]], we will integrate it into
[Terraform][1], allowing for the management and lookup of sections, VLANs, subnets and IP addresses, entirely withing Terraform.


## Installing

See the [Plugin Basics][5] page of the Terraform docs to see how to plunk this
into your config. Check the [releases page][6] of this repo to get releases for
Linux, OS X, and Windows.

[5]: https://www.terraform.io/docs/plugins/basics.html
[6]: https://github.com/limberger/terraform-provider-netbox/releases

## Usage

After installation, to use the plugin, simply use any of its resources or data
sources (such as `netbox_subnet` or `netbox_address` in a Terraform
configuration.

Credentials can be supplied via configuration variables to the `netbox`
provider instance, or via environment variables. These are documented in the
next section.

You can see the following example below for a simple usage example that reserves
the first available IP address in a subnet. This address could then be passed
along to the configuration for a VM, say, for example, a
[`vsphere_virtual_machine`][7] resource.

[7]: https://www.terraform.io/docs/providers/vsphere/r/virtual_machine.html

```
provider "netbox" {
  app_id = "2fe35cabcfe231ebc8734a798f1cac63439a7a2b"
  endpoint = "172.17.133.10"
}

data "netbox_prefixes" "prefixes" {
  prefixes_id = 1
}

data "netbox_prefixes_available_ips" "next_address" {
  prefixes_id = "${data.netbox_prefixes.prefixes.prefixes_id}"
}
```

### Plugin Options

The options for the plugin are as follows:

 * `app_id` - The API application ID, configured in the PHPIPAM API panel. This
   application ID should have read/write access if you are planning to use the
   resources, but read-only access should be sufficient if you are only using
   the data sources. Can also be supplied by the `NETBOX_APP_ID` environment
   variable.
 * `endpoint` - The full URL to the PHPIPAM API endpoint, such as
   `https://phpipam.example.com/api`. Can also be supplied by the
   `NETBOX_ENDPOINT_ADDR` environment variable.
 * `timeout` - Timeout to API calls.

### Data Sources

The following data sources are supplied by this plugin:

#### The `netbox_prefixes` Data Source

The `netbox_prefixes` cadastred on netbox

**Example:**

```
data "phpipam_address" "address" {
  ip_address = "10.10.1.1"
}

output "address_description" {
  value = "${data.phpipam_address.address.description}"
}
```

**Example With `description`:**

```
data "netbox_prefixes" "prefixes" {
  prefixes_id = 1
}

output "prefix_description" {
  value = "${data.netbox_prefixes.prefixes.description}"
}
```

**Example With `vlan id - vid`:**

```
data "netbox_prefix" "search_by_vid" {
  vlan = {
    vid = 16
  }
}

output "vlan_description" {
  value = "${data.netbox_prefixes.search_by_vid.name}"
}
```

##### Argument Reference

The data source takes the following parameters:

 * `address_id` - The ID of the IP address in the PHPIPAM database.
 * `ip_address` - The actual IP address in PHPIPAM.
 * `subnet_id` - The ID of the subnet that the address resides in. This is
   required to search on the `description` or `hostname` fields.
 * `description` - The description of the IP address. `subnet_id` is required
   when using this field.
 * `hostname` - The host name of the IP address. `subnet_id` is required when
   using this field.
 * `custom_field_filter` - A map of custom fields to search for. The filter
   values are regular expressions that follow the RE2 syntax for which you can
   find documentation [here](https://github.com/google/re2/wiki/Syntax). All
   fields need to match for the match to succeed.

⚠️  **NOTE:** `description`, `hostname`, and `custom_field_filter` fields return
the first match found without any warnings. If you are looking to return
multiple addresses, combine this data source with the `phpipam_addresses` data
source.

⚠️  **NOTE:** An empty or unspecified `custom_field_filter` value is the
equivalent to a regular expression that matches everything, and hence will
return the first address it sees in the subnet.

Arguments are processed in the following order of precedence:

 * `address_id`
 * `ip_address`
 * `subnet_id`, and either one of `description`, `hostname`, or
   `custom_field_filter`

##### Attribute Reference

The following attributes are exported:

 * `address_id` - The ID of the IP address in the PHPIPAM database.
 * `ip_address` - the IP address.
 * `subnet_id` - The database ID of the subnet this IP address belongs to.
 * `is_gateway` - `true` if this IP address has been designated as a gateway.
 * `description` - The description provided to this IP address.
 * `hostname` - The hostname supplied to this IP address.
 * `owner` - The owner name provided to this IP address.
 * `mac_address` - The MAC address provided to this IP address.
 * `state_tag_id` - The tag ID in the database for the IP address' specific
   state. **NOTE:** This is currently represented as an integer but may change
   to the specific string representation at a later time.
 * `skip_ptr_record` - `true` if PTR records are not being created for this IP
   address.
 * `ptr_record_id` - The ID of the associated PTR record in the PHPIPAM
   database.
 * `device_id` - The ID of the associated device in the PHPIPAM database.
 * `switch_port_label` - A string port label that is associated with this
   address.
 * `note` - The note supplied to this IP address.
 * `last_seen` - The last time this IP address answered ping probes.
 * `exclude_ping` - `true` if this address is excluded from ping probes.
 * `edit_date` - The last time this resource was modified.
 * `custom_fields` - A key/value map of custom fields for this address.


#### The `netbox_vlans` Data Source

The `netbox_vlans` data source allows you to search for vlans_by_id

**Example:**

```
data "netbox_vlans" "search_by_vid" {
  vid = 16
}

output "vlans_description" {
  value = "${data.netbox_vlans.search_by_vid.description}"
}
```

```
data "netbox_vlans" "search_by_name" {
  name = "Vlan16"
}

output "vlans_description" {
  value = "${data.netbox_vlans.search_by_name.description}"
}
```



#### End



#### The `phpipam_vlan` Resource

The `phpipam_vlan` resource can be used to manage a VLAN on PHPIPAM. Use it to
set up a VLAN through Terraform, or update details such as its name or
description. If you are just looking for information on a VLAN, use the
`phpipam_vlan` data source instead.

**Example:**

```
resource "phpipam_vlan" "vlan" {
  name        = "tf-test"
  number      = 1000
  description = "Managed by Terraform"

  custom_fields = {
    CustomTestVLANs = "terraform-test"
  }
}
```

##### Argument Reference

The resource takes the following parameters:

 * `name` (Required) - The name/label of the VLAN.
 * `number` (Required) - The number of the VLAN (the actual VLAN ID on your switch).
 * `l2_domain_id` (Optional) - The layer 2 domain ID in the PHPIPAM database.
 * `description` (Optional) - The description supplied to the VLAN.
 * `edit_date` (Optional) - The date this resource was last updated.
 * `custom_fields` (Optional) -  A key/value map of custom fields for this
   VLAN.

⚠️  **NOTE on custom fields:** PHPIPAM installations with custom fields must have
all fields set to optional when using this plugin. For more info see
[here](https://github.com/phpipam/phpipam/issues/1073). Further to this, either
ensure that your fields also do not have default values, or ensure the default
is set in your TF configuration. Diff loops may happen otherwise!

##### Attribute Reference

The following attributes are exported:

 * `vlan_id` - The ID of the VLAN to look up. **NOTE:** this is the database ID,
   not the VLAN number - if you need this, use the `number` parameter.
 * `edit_date` - The date this resource was last updated.

## LICENSE

```
Copyright 2018 BB, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
