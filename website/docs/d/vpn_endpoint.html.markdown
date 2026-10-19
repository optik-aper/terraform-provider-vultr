---
layout: "vultr"
page_title: "Vultr: vultr_vpn_endpoint"
sidebar_current: "docs-vultr-datasource-vpn-endpoint"
description: |-
  Get information about a Vultr VPN endpoint.
---

# vultr_vpn_endpoint

Get information about a Vultr VPN endpoint.

## Example Usage

Query for a VPN endpoint via its label:

```hcl
data "vultr_vpn_endpoint" "my_vpn" {
  vpn_id = "e97bdee9-2781-4f31-be03-60fc75f399ae"
  filter {
    name = "label"
    values = ["my-vpn-description"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpn_id` - (Required) The ID of an existing VPN.
* `filter` - (Required) The field name and value for which you are querying. 

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the endpoint.
* `vpn_id` - The ID of the parent VPN.
* `label` - The label of the endpoint.
* `pre_shared_key` - A 32 byte, base64 encoded key.
* `ip_address` - The IP address of the endpoint.
* `remote_subnets` - A set of subnets containing:

    `subnet` - The subnet of the network
    `prefix` - The subnet mask of the network

* `remote_identifier` - The remote IP address of the endpoint.
* `public_key` - The public key for the endpoint, if available.
* `private_key` - The private key for the endpoint, if available.
* `server_public_key` - The server public key for the endpoint, if available.
* `date_created` - The date of creation of the endpoint.
* `date_updated` - The date of last modification for the endpoint.
