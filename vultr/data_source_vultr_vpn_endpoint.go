package vultr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func dataSourceVultrVPNEndpoint() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrVPNEndpointRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vpn_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_subnets": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefix": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"pre_shared_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"server_public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVultrVPNEndpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var endpointList []govultr.VPNEndpoint
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}

	vpnID := d.Get("vpn_id").(string)

	for {
		endpoints, meta, _, err := client.VPN.ListEndpoints(ctx, vpnID, options)
		if err != nil {
			return diag.Errorf("error getting list of vpn endpoints : %v", err)
		}

		for i := range endpoints {
			sm, err := structToMap(endpoints[i])
			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				endpointList = append(endpointList, endpoints[i])
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(endpointList) > 1 {
		return diag.Errorf("your search returned too many results. Please refine your search to be more specific")
	}

	if len(endpointList) < 1 {
		return diag.Errorf("no results were found")
	}

	d.SetId(endpointList[0].ID)

	if err := d.Set("label", endpointList[0].Label); err != nil {
		return diag.Errorf("unable to set vpn endpoint `label` read value: %v", err)
	}

	if err := d.Set("ip_address", endpointList[0].Config.IPAddress); err != nil {
		return diag.Errorf("unable to set vpn endpoint `ip_address` read value: %v", err)
	}

	if err := d.Set("pre_shared_key", endpointList[0].Config.PreSharedKey); err != nil {
		return diag.Errorf("unable to set vpn endpoint `pre_shared_key` read value: %v", err)
	}

	var subnets []map[string]interface{}
	for i := range endpointList[0].RemoteSubnets {
		subnets = append(subnets, map[string]interface{}{
			"subnet": endpointList[0].RemoteSubnets[i].Subnet,
			"prefix": endpointList[0].RemoteSubnets[i].Prefix,
		})
	}

	if err := d.Set("remote_subnets", subnets); err != nil {
		return diag.Errorf("unable to set vpn endpoint `remote_subnets` read value: %v", err)
	}

	if err := d.Set("remote_identifier", endpointList[0].Config.RemoteIdentifier); err != nil {
		return diag.Errorf("unable to set vpn endpoint `remote_identifier` read value: %v", err)
	}

	if err := d.Set("public_key", endpointList[0].Config.PublicKey); err != nil {
		return diag.Errorf("unable to set vpn endpoint `public_key` read value: %v", err)
	}

	if err := d.Set("private_key", endpointList[0].Config.PrivateKey); err != nil {
		return diag.Errorf("unable to set vpn endpoint `private_key` read value: %v", err)
	}

	if err := d.Set("server_public_key", endpointList[0].Config.ServerPublicKey); err != nil {
		return diag.Errorf("unable to set vpn endpoint `server_public_key` read value: %v", err)
	}

	if err := d.Set("date_created", endpointList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set vpn endpoint `date_created` read value: %v", err)
	}

	if err := d.Set("date_updated", endpointList[0].DateCreated); err != nil {
		return diag.Errorf("unable to set vpn endpoint `date_updated` read value: %v", err)
	}

	return nil
}
