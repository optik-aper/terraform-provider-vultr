---
layout: "vultr"
page_title: "Vultr: vultr_vpn_endpoint"
sidebar_current: "docs-vultr-resource-vpn-endpoint"
description: |-
  Provides a Vultr VPN endpoint resource. This can be used to manage endpoints on a VPN.
---

# vultr_vpn_endpoint

Provides a Vultr VPN endpoint resource. This can be used to manage endpoints on a VPN.

## Example Usage

A VPN is required for this resource, which requires a VPC and a NAT gateway:

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
  protocol = "ipsec"
  subnet = "10.2.0.0"
  subnet_prefix = 16
}

resource "vultr_vpn_endpoint" "my-endpoint" {
  vpn_id = resource.vultr_vpn.my-vpn.id
  label = "my-endpoint"
  ip_address = "10.2.123.1"
  pre_shared_key = base64encode("1234123412341234abcdabcdefghefgh")  

  remote_subnets { 
    subnet = "10.3.123.0"
    prefix = 24 
  }

  remote_subnets {
    subnet = "10.8.4.0"
    prefix = 24
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpn_id` - (Required) The ID of an existing VPN.
* `label` - (Required) The label of the endpoint.
* `pre_shared_key` - (Optional) A 32 byte, base64 encoded key.
* `ip_address` - (Required, if using an `ipsec` VPN) The IP address of the endpoint.
* `remote_subnets` - (Required, if using an `ipsec` VPN) A set of subnets containing:

    `subnet` - The subnet of the network
    `prefix` - The subnet mask of the network

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

## Import

VPN endpoints can be imported using the endpoint's `ID`, e.g.

```
terraform import vultr_vpn_endpoint.my_vpn e97bdee9-2781-4f31-be03-60fc75f399ae
```
