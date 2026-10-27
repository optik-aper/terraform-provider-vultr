package vultr

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vultr/govultr/v3"
)

const vpnEndpointMissingError string = "Not found"

func resourceVultrVPNEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrVPNEndpointCreate,
		ReadContext:   resourceVultrVPNEndpointRead,
		UpdateContext: resourceVultrVPNEndpointUpdate,
		DeleteContext: resourceVultrVPNEndpointDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"vpn_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remote_subnets": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"prefix": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"pre_shared_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsBase64,
			},
			// computed -------------------
			"config_file": {
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

func resourceVultrVPNEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	tflog.Info(ctx, fmt.Sprintf("creating endpoint on vpn (%s)", d.Get("vpn_id").(string)))

	endpointReq := &govultr.VPNEndpointReq{
		Label:        d.Get("label").(string),
		IPAddress:    d.Get("ip_address").(string),
		PreSharedKey: d.Get("pre_shared_key").(string),
	}

	if len(d.Get("remote_subnets").(*schema.Set).List()) != 0 {
		var subnets []govultr.VPNEndpointRemoteSubnet
		nets := d.Get("remote_subnets").(*schema.Set).List()
		for i := range nets {
			net := nets[i].(map[string]interface{})
			subnets = append(subnets, govultr.VPNEndpointRemoteSubnet{
				Subnet: net["subnet"].(string),
				Prefix: net["prefix"].(int),
			})
		}

		endpointReq.RemoteSubnets = subnets
	}

	endpointRes, _, err := client.VPN.CreateEndpoint(ctx, d.Get("vpn_id").(string), endpointReq)
	if err != nil {
		return diag.Errorf("error creating vpn endpoint : %v", err)
	}

	d.SetId(endpointRes.ID)

	tflog.Info(ctx, fmt.Sprintf("created endpoint (%s) on vpn (%s)", d.Id(), d.Get("vpn_id").(string)))

	return resourceVultrVPNEndpointRead(ctx, d, meta)
}

func resourceVultrVPNEndpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	endpoint, resp, err := client.VPN.GetEndpoint(ctx, d.Get("vpn_id").(string), d.Id())
	if err != nil {
		if missing := checkIsMissing(resp, err, vpnEndpointMissingError); missing {
			tflog.Warn(
				ctx,
				fmt.Sprintf("removing endpoint (%s) from vpn (%s) because it is gone", d.Id(), d.Get("vpn_id").(string)),
			)
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting endpoint (%s) on vpn (%s) : %v", d.Id(), d.Get("vpn_id").(string), err)
	}

	if err := d.Set("label", endpoint.Label); err != nil {
		return diag.Errorf("unable to set resource vpn endpoint `label` read value: %v", err)
	}

	if err := d.Set("ip_address", endpoint.Config.IPAddress); err != nil {
		return diag.Errorf("unable to set resource vpn endpoint `ip_address` read value: %v", err)
	}

	if err := d.Set("pre_shared_key", endpoint.Config.PreSharedKey); err != nil {
		return diag.Errorf("unable to set resource vpn endpoint `pre_shared_key` read value: %v", err)
	}

	var subnets []map[string]interface{}
	for i := range endpoint.RemoteSubnets {
		subnets = append(subnets, map[string]interface{}{
			"subnet": endpoint.RemoteSubnets[i].Subnet,
			"prefix": endpoint.RemoteSubnets[i].Prefix,
		})
	}

	if err := d.Set("remote_subnets", subnets); err != nil {
		return diag.Errorf("unable to set resource vpn endpoint `remote_subnets` read value: %v", err)
	}

	if err := d.Set("remote_identifier", endpoint.Config.RemoteIdentifier); err != nil {
		return diag.Errorf("unable to set resource vpn endpoint `remote_identifier` read value: %v", err)
	}

	if err := d.Set("public_key", endpoint.Config.PublicKey); err != nil {
		return diag.Errorf("unable to set resource vpn endpoint `public_key` read value: %v", err)
	}

	if err := d.Set("private_key", endpoint.Config.PrivateKey); err != nil {
		return diag.Errorf("unable to set resource vpn endpoint `private_key` read value: %v", err)
	}

	if err := d.Set("server_public_key", endpoint.Config.ServerPublicKey); err != nil {
		return diag.Errorf("unable to set resource vpn endpoint `server_public_key` read value: %v", err)
	}

	if err := d.Set("date_created", endpoint.DateCreated); err != nil {
		return diag.Errorf("unable to set resource vpn endpoint `date_created` read value: %v", err)
	}

	if err := d.Set("date_updated", endpoint.DateCreated); err != nil {
		return diag.Errorf("unable to set resource vpn endpoint `date_updated` read value: %v", err)
	}

	configFile, _, err := client.VPN.GetEndpointConfigFile(ctx, d.Get("vpn_id").(string), d.Id())
	if err != nil {
		return diag.Errorf("error getting config file for endpoint (%s) on vpn (%s) :", d.Id(), d.Get("vpn_id").(string))
	}

	if err := d.Set("config_file", configFile.Config); err != nil {
		return diag.Errorf("unable to set resource vpn endpoint `config_file` read value: %v", err)
	}

	return nil
}

func resourceVultrVPNEndpointUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	tflog.Info(ctx, fmt.Sprintf("updating endpoint (%s) on vpn (%s)", d.Id(), d.Get("vpn_id").(string)))

	var updateReq govultr.VPNEndpointReq

	if d.HasChange("label") {
		updateReq.Label = d.Get("label").(string)
	}

	if d.HasChange("ip_address") {
		updateReq.IPAddress = d.Get("ip_address").(string)
	}

	if d.HasChange("pre_shared_key") {
		updateReq.PreSharedKey = d.Get("pre_shared_key").(string)
	}

	if d.HasChange("remote_subnets") {
		_, subsNew := d.GetChange("remote_subnets")

		subsList := subsNew.(*schema.Set).List()

		if len(subsList) != 0 {
			for i := range subsList {
				sub := subsList[i].(map[string]interface{})
				updateReq.RemoteSubnets = append(updateReq.RemoteSubnets, govultr.VPNEndpointRemoteSubnet{
					Subnet: sub["subnet"].(string),
					Prefix: sub["prefix"].(int),
				})
			}
		}
	}

	if _, _, err := client.VPN.UpdateEndpoint(ctx, d.Get("vpn_id").(string), d.Id(), &updateReq); err != nil {
		return diag.Errorf("error updating vpn endpoint : %v", err)
	}

	return resourceVultrVPNEndpointRead(ctx, d, meta)
}

func resourceVultrVPNEndpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	tflog.Info(ctx, fmt.Sprintf("destroying endpoint (%s) on vpn (%s)", d.Id(), d.Get("vpn_id").(string)))

	if err := client.VPN.DeleteEndpoint(ctx, d.Get("vpn_id").(string), d.Id()); err != nil {
		return diag.Errorf("error destroying endpoint (%s) on vpn (%s)", d.Id(), d.Get("vpn_id").(string))
	}

	return nil
}
