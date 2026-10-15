---
layout: "vultr"
page_title: "Vultr: vultr_vpn"
sidebar_current: "docs-vultr-resource-vpn"
description: |-
  Provides a Vultr VPN resource. This can be used to manage a VPN on your Vultr account.
---

# vultr_vpn

Provides a Vultr VPN resource. This can be used to manage a VPN on your Vultr account.

## Example Usage

A VPC with a NAT gateway is required for this resource.  So to create a new VPN:

```hcl
resource "vultr_vpc" "my-vpc" {
  region      = "ewr"
  description = "my-vpc"
}

resource "vultr_nat_gateway" "my-nat" {
  vpc_id = resource.vultr_vpc.my-vpc.id
  label = "my-nat-gw"
}

resource "vultr_vpn" "my-vpn" {
  vpc_id = resource.vultr_vpc.my-vpc.id
  nat_gateway_id = resource.vultr_nat_gateway.my-nat.id
  label = "my-vpn"
  protocol = "openvpn"
  subnet = "10.2.0.0"
  subnet_prefix = 16
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) The ID of an existing VPC with a NAT gateway attached to it.
* `nat_gateway_id` - (Required) The ID of the NAT gateway.
* `region` - (Required) The ID of the region in which to create the VPN.
* `protocol` - (Required) The protocol used to initialize the VPN. Must be one of: `wireguard`; `openvpn`; `ipsec`.
* `subnet` - (Optional) The IPv4 subnet to be used when creating the VPN.
* `subnet_prefix` - (Optional) The number of bits for the subnet mask. Example: 24

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

## Import

VPNs can be imported using the VPN `ID`, e.g.

```
terraform import vultr_vpn.my_vpn 00e5401d-8673-46a7-86bf-c7fb5fa3b96b
```
