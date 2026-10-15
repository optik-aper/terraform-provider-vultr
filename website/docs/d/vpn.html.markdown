---
layout: "vultr"
page_title: "Vultr: vultr_vpn"
sidebar_current: "docs-vultr-datasource-vpn"
description: |-
  Get information about a Vultr VPN.
---

# vultr_vpn

Get information about a Vultr VPN.

## Example Usage

Query for a VPN via its label:

```hcl
data "vultr_vpn" "my_vpn" {
  filter {
    name = "label"
    values = ["my-vpn-description"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) The field name and value for which you are querying. 

The `filter` block supports the following:

* `name` - Attribute name to filter with.
* `values` - One or more values filter with.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the VPN.
* `vpc_id` - The ID of an existing VPC with a NAT gateway attached to it.
* `nat_gateway_id` - The ID of the NAT gateway.
* `region` - The ID of the region in which to create the VPN.
* `protocol` - The protocol used to initialize the VPN. Must be one of: `wireguard`; `openvpn`; `ipsec`.
* `subnet` - The IPv4 subnet to be used when creating the VPN.
* `subnet_prefix` - The number of bits for the subnet mask. Example: 24
* `status` - The current status of the VPN.
* `public_key` - The public key used by the VPN.
* `charges_pending` - The pending charges for the VPN.
* `charges_monthly` - The monthly cost of the VPN.
* `date_created` - The date of creation of the VPN.
