package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrVPN() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrVPNRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_prefix": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charges_pending": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"charges_monthly": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrVPNRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var vpnList []govultr.VPN
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	for {
		vpns, meta, _, err := client.VPN.List(ctx, options)
		if err != nil {
			return diag.Errorf("error getting list of vpns : %v", err)
		}

		for i := range vpns {
			sm, err := structToMap(vpns[i])
			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				vpnList = append(vpnList, vpns[i])
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(vpnList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(vpnList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(vpnList[0].ID)

	if err := d.Set("vpc_id", vpnList[0].VPCID); err != nil {
		return diag.Errorf("unable to set vpn `vpc_id` read value: %v", err)
	}

	if err := d.Set("label", vpnList[0].Label); err != nil {
		return diag.Errorf("unable to set vpn `label` read value: %v", err)
	}

	if err := d.Set("protocol", vpnList[0].Protocol); err != nil {
		return diag.Errorf("unable to set vpn `protocol` read value: %v", err)
	}

	if err := d.Set("subnet", vpnList[0].Subnet); err != nil {
		return diag.Errorf("unable to set vpn `subnet` read value: %v", err)
	}

	if err := d.Set("subnet_prefix", vpnList[0].SubnetPrefix); err != nil {
		return diag.Errorf("unable to set vpn `subnet_prefix` read value: %v", err)
	}

	if err := d.Set("status", vpnList[0].Status); err != nil {
		return diag.Errorf("unable to set vpn `status` read value: %v", err)
	}

	if err := d.Set("nat_gateway_id", vpnList[0].NATGatewayID); err != nil {
		return diag.Errorf("unable to set vpn `nat_gateway_id` read value: %v", err)
	}

	if err := d.Set("public_key", vpnList[0].PublicKey); err != nil {
		return diag.Errorf("unable to set vpn `nat_gateway_id` read value: %v", err)
	}

	if err := d.Set("charges_pending", vpnList[0].Billing.Charges); err != nil {
		return diag.Errorf("unable to set vpn `charges_pending` read value: %v", err)
	}

	if err := d.Set("charges_monthly", vpnList[0].Billing.Monthly); err != nil {
		return diag.Errorf("unable to set vpn `charges_monthly` read value: %v", err)
	}

	return nil
}
