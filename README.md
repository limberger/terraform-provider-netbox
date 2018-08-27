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
sources (such as `netbox_prefixes`, `netbox_vlans` or `netbox_prefixes_available_ips` in a Terraform
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
  app_id = "0123456789abcdef0123456789abcdef01234567"
  endpoint = "0.0.0.0:32768"
}

data "netbox_prefixes" "prefixes" {
  prefixes_id = 1
}

resource "netbox_prefixes_available_ips" "next_address" {
	prefixes_id = "${data.netbox_prefixes.prefixes.prefixes_id}"
	description = "IP requisitado via Terraform 20180827"
}
```

### Plugin Options

The options for the plugin are as follows:

 * `app_id` - The API application ID, configured in the NETBOX  panel. This
   application ID should have read/write access if you are planning to use the
   resources, but read-only access should be sufficient if you are only using
   the data sources. Caan also be supplied by the `NETBOX_APP_ID` environment
   variable.
 * `endpoint` - The server, protocol and port to access the NETBOX API, such as
   `https://netbox.example.com/api`. Can also be supplied by the
   `NETBOX_ENDPOINT_ADDR` environment variable.

### Data Sources

The following data sources are supplied by this plugin:

#### The `netbox_prefixes` Data Source

The `netbox_prefixes` cadastred on netbox

**Example:**

```
data "netbox_prefixes" "prefixes" {
  prefixes_id = 1
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

 * `address_id` - The ID of the IP address in the NETBOX database.
 * `description` - The description of the IP address. `subnet_id` is required
   when using this field.

⚠️  **NOTE:** `description`, `hostname`, and `custom_field_filter` fields return
the first match found without any warnings.

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

 * `description` - The description provided to this IP address.


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
