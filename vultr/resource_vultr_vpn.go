package vultr

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

const vpnMissingError string = "VPN not found"

func resourceVultrVPN() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrVPNCreate,
		ReadContext:   resourceVultrVPNRead,
		UpdateContext: resourceVultrVPNUpdate,
		DeleteContext: resourceVultrVPNDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"subnet_prefix": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			// computed -------------------
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

func resourceVultrVPNCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	vpnReq := &govultr.VPNCreateReq{
		VPCID:        d.Get("vpc_id").(string),
		Label:        d.Get("label").(string),
		Protocol:     d.Get("protocol").(string),
		Subnet:       d.Get("subnet").(string),
		SubnetPrefix: d.Get("subnet_prefix").(int),
	}

	vpnRes, _, err := client.VPN.Create(ctx, vpnReq)
	if err != nil {
		return diag.Errorf("error creating vpn : %v", err)
	}

	d.SetId(vpnRes.ID)

	tflog.Info(ctx, "waiting for vpn to become active")

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"active"},
		Refresh: func() (interface{}, string, error) {
			vpnState, _, stateErr := client.VPN.Get(ctx, d.Id())
			if stateErr != nil {
				return nil, "", fmt.Errorf("error retrieving vpn (%s) : %v", d.Id(), stateErr)
			}

			tflog.Info(ctx, fmt.Sprintf("the vpn (%s) status is %q", d.Id(), vpnState.Status))

			return vpnState, vpnState.Status, nil
		},
		Timeout:        d.Timeout(schema.TimeoutCreate),
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 10,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error while waiting for vpn (%s) to become active : %v", d.Id(), err)
	}

	tflog.Info(ctx, fmt.Sprintf("created vpn (%s)", d.Id()))

	return resourceVultrVPNRead(ctx, d, meta)
}

func resourceVultrVPNRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	vpn, resp, err := client.VPN.Get(ctx, d.Id())
	if err != nil {
		if missing := checkIsMissing(resp, err, vpnMissingError); missing {
			tflog.Warn(ctx, fmt.Sprintf("removing vpn (%s) because it is gone", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting vpn : %v", err)
	}

	if err := d.Set("vpc_id", vpn.VPCID); err != nil {
		return diag.Errorf("unable to set resource vpn `vpc_id` read value: %v", err)
	}

	if err := d.Set("label", vpn.Label); err != nil {
		return diag.Errorf("unable to set resource vpn `label` read value: %v", err)
	}

	if err := d.Set("protocol", vpn.Protocol); err != nil {
		return diag.Errorf("unable to set resource vpn `protocol` read value: %v", err)
	}

	if err := d.Set("subnet", vpn.Subnet); err != nil {
		return diag.Errorf("unable to set resource vpn `subnet` read value: %v", err)
	}

	if err := d.Set("subnet_prefix", vpn.SubnetPrefix); err != nil {
		return diag.Errorf("unable to set resource vpn `subnet_prefix` read value: %v", err)
	}

	if err := d.Set("status", vpn.Status); err != nil {
		return diag.Errorf("unable to set resource vpn `status` read value: %v", err)
	}

	if err := d.Set("nat_gateway_id", vpn.NATGatewayID); err != nil {
		return diag.Errorf("unable to set resource vpn `nat_gateway_id` read value: %v", err)
	}

	if err := d.Set("public_key", vpn.PublicKey); err != nil {
		return diag.Errorf("unable to set resource vpn `nat_gateway_id` read value: %v", err)
	}

	if err := d.Set("charges_pending", vpn.Billing.Charges); err != nil {
		return diag.Errorf("unable to set resource vpn `charges_pending` read value: %v", err)
	}

	if err := d.Set("charges_monthly", vpn.Billing.Monthly); err != nil {
		return diag.Errorf("unable to set resource vpn `charges_monthly` read value: %v", err)
	}

	return nil
}

func resourceVultrVPNUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	if d.HasChange("label") {
		if _, _, err := client.VPN.Update(ctx, d.Id(), &govultr.VPNUpdateReq{Label: d.Get("label").(string)}); err != nil {
			return diag.Errorf("error updating vpn : %v", err)
		}
	}

	return resourceVultrVPNRead(ctx, d, meta)
}

func resourceVultrVPNDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	tflog.Info(ctx, fmt.Sprintf("destroying vpn (%s)", d.Id()))

	retryErr := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		err := client.VPN.Delete(ctx, d.Id())
		if err == nil {
			return nil
		}

		return retry.NonRetryableError(err)
	})

	if retryErr != nil {
		return diag.Errorf("error destroying vpn (%s) : %v", d.Id(), retryErr)
	}

	return nil
}
