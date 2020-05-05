# Terraform Provider Plugin for Netbox

This repository is a fork of a generic [Terraform provider for NetBox](https://github.com/limberger/terraform-provider-netbox).

It has changes for a specific use case where prefixes are allocated from select IP ranges.
Please see the code for details.

Access some plugin documentation [here](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk)
and [here](https://www.terraform.io/docs/plugins/basics.html)

## About Netbox

[Netbox](https://github.com/netbox-community/netbox) is an open source IP address management
system written in Python. Through our Go integration provided by 
[Go-Netbox](https://github.com/netbox-community/go-netbox), we will integrate it into 
[Terraform](https://www.terraform.io/), allowing for the management of prefix pools.


## Installing

See the [Plugin Basics](https://www.terraform.io/docs/plugins/basics.html) page of the Terraform
docs to see how to drop this into your config.

## Building

The following makefile targets are available. They have only been tested on MacOS.

### test

`make test` will run unit tests that do not require a running NetBox.

### testacc

`make testacc` will run acceptance tests that do require a running NetBox. Three conditions
must be met to run these successfully:
1. Netbox configured correctly (see [reconfigure](#reconfigure), below)
2. `NETBOX_APP_ID` environment variable must be set to indicate a valid Netbox API access token.
3. `NETBOX_ENDPOINT_ADDR` environment variable must be set to indicate the `<hostname>:<port_number>` for
Netbox.

### reconfigure
`make reconfigure` will make and run the `setup` utility in the `utils` directory, specifying the `-destroyAll`
option. Be careful with this one! It will remove all VRFs and Prefixes from the specified Netbox.
**Note:** The setup app should be updated to allow specifying the Netbox endpoint and token like
`testacc` does.

### reconfigure_testacc
`make reconfigure_testacc` is a combination of `reconfigure` and `testacc`. It will configure Netbox
for tests and run the tests.

### provider
`make provider` will build the `terraform-provider-netbox` in the `build` directory.

### clean
`make clean` will remove the `build` directory and the `utils/setup` utility.

## Provider Usage

### Provider Input
Credentials can be supplied via configuration variables to the `netbox`
provider instance, or via environment variables.

The options for the plugin are as follows:

 * `app_id` - The API application ID, configured in the NETBOX  panel. This
   application ID should have read/write access if you are planning to use the
   resources, but read-only access should be sufficient if you are only using
   the data sources. Can also be supplied by the `NETBOX_APP_ID` environment
   variable.
 * `endpoint` - The server, protocol and port to access the NETBOX API, such as
   `https://netbox.example.com/api`. Can also be supplied by the
   `NETBOX_ENDPOINT_ADDR` environment variable.

```
provider "netbox" {
  app_id = "0123456789abcdef0123456789abcdef01234567"
  endpoint = "0.0.0.0:32768"
}
```

### Resource netbox_pool_prefixes 
After installation, to use the plugin, simply use the resource `netbox_pool_prefixes` in
a Terraform configuration. 
#### Input
The input parameters to the resource are:
| Name          | Type        | Description                                                                                | Example                                            |
|---------------|-------------|--------------------------------------------------------------------------------------------|----------------------------------------------------|
| environment   | string      | The deployment environment (dev, test, staging, production)                                | "dev"                                              |
| resource_type | string      | One of "core" (hubs and peered concerns), "depot", "edge", "vpn", or "servicedelivery"     | "core"                                             |
| prefix_length | number      | Length of the prefix to allocate. Must be between 18 and 28, inclusive.                    | 28                                                 |
| tags          | map(string) | Map with strings as values. Required keys are "name" and "unique". Others may be included. | `{ name = "prod_name", unique = "unique_string" }` |

```
resource netbox_pool_prefixes new {
  environment   = "dev"
  resource_type = "core"
  prefix_length = 28
  tags = {
    name   = "Some_alias"
    unique = "new_unique_string"
  }
}
```

#### Output
The output from the resource is a single object called prefix with the following attributes:
| Name          | Type        | Description                                                                       | Example                                            |
|---------------|-------------|-----------------------------------------------------------------------------------|----------------------------------------------------|
| environment   | string      | The deployment environment specified at input.                                    | "dev"                                              |
| resource_type | string      | The resource_type specified at input.                                             | "core"                                             |
| id            | string      | The Terraform ID of the prefix resource. Identical to prefix_id, but as a string. | "343"                                              |
| prefix        | string      | The prefix allocated from the specified pool.                                     | "172.16.0.16/28"                                   |
| prefix_id     | number      | The ID of the prefix allocated from the specified pool.                           | 343                                                |
| prefix_length | number      | Length of the prefix to allocate. Must be between 18 and 28, inclusive.           | 28                                                 |
| tags          | map(string) | The tags specified at input.                                                      | `{ name = "prod_name", unique = "unique_string" }` |


```
output prefix { value = netbox_pool_prefixes.new }
```

```
prefix = {
  "environment" = "dev"
  "id" = "343"
  "resource_type" = "core"
  "prefix" = "172.16.0.16/28"
  "prefix_id" = 343
  "prefix_length" = 28
  "tags" = {
    "name" = "prod_name"
    "unique" = "unique_string"
  }
}
```

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
