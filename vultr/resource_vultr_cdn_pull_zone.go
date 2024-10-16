package vultr

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func resourceVultrCDNPullZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVultrCDNPullZoneCreate,
		ReadContext:   resourceVultrCDNPullZoneRead,
		UpdateContext: resourceVultrCDNPullZoneUpdate,
		DeleteContext: resourceVultrCDNPullZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"origin_scheme": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"origin_domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cors": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"gzip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"block_ai": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"block_bad_bots": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// computed fields
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cache_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"requests": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bytes_in": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bytes_out": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"packets_per_second": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"date_purged": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceVultrCDNPullZoneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	req := &govultr.CDNZoneReq{
		Label:        d.Get("label").(string),
		OriginScheme: d.Get("origin_scheme").(string),
		OriginDomain: d.Get("origin_domain").(string),
		CORS:         d.Get("cors").(bool),
		GZIP:         d.Get("gzip").(bool),
		BlockAI:      d.Get("block_ai").(bool),
		BlockBadBots: d.Get("block_bad_bots").(bool),
	}

	pz, _, err := client.CDN.CreatePullZone(ctx, req)
	if err != nil {
		return diag.Errorf("error creating cdn pull zone: %v", err)
	}

	d.SetId(pz.ID)

	log.Printf("[INFO] CDN Pull Zone ID: %s", d.Id())

	return resourceVultrCDNPullZoneRead(ctx, d, meta)
}

func resourceVultrCDNPullZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	pz, _, err := client.CDN.GetPullZone(ctx, d.Id())
	if err != nil {
		return diag.Errorf("error getting cdn pull zone: %v", err)
	}

	if pz == nil {
		log.Printf("[WARN] Vultr pull zone (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	if err := d.Set("label", pz.Label); err != nil {
		return diag.Errorf("unable to set resource pull zone `label` read value: %v", err)
	}

	if err := d.Set("origin_scheme", pz.OriginScheme); err != nil {
		return diag.Errorf("unable to set resource pull zone `origin_scheme` read value: %v", err)
	}

	if err := d.Set("origin_domain", pz.OriginDomain); err != nil {
		return diag.Errorf("unable to set resource pull zone `origin_domain` read value: %v", err)
	}

	if err := d.Set("cors", pz.CORS); err != nil {
		return diag.Errorf("unable to set resource pull zone `cors` read value: %v", err)
	}

	if err := d.Set("gzip", pz.GZIP); err != nil {
		return diag.Errorf("unable to set resource pull zone `gzip` read value: %v", err)
	}

	if err := d.Set("block_ai", pz.BlockAI); err != nil {
		return diag.Errorf("unable to set resource pull zone `block_ai` read value: %v", err)
	}

	if err := d.Set("block_bad_bots", pz.BlockBadBots); err != nil {
		return diag.Errorf("unable to set resource pull zone `block_bad_bots` read value: %v", err)
	}

	if err := d.Set("date_created", pz.DateCreated); err != nil {
		return diag.Errorf("unable to set resource pull zone `date_created` read value: %v", err)
	}

	if err := d.Set("status", pz.Status); err != nil {
		return diag.Errorf("unable to set resource pull zone `status` read value: %v", err)
	}

	if err := d.Set("url", pz.CDNURL); err != nil {
		return diag.Errorf("unable to set resource pull zone `url` read value: %v", err)
	}

	if err := d.Set("cache_size", pz.CacheSize); err != nil {
		return diag.Errorf("unable to set resource pull zone `cache_size` read value: %v", err)
	}

	if err := d.Set("requests", pz.Requests); err != nil {
		return diag.Errorf("unable to set resource pull zone `requests` read value: %v", err)
	}

	if err := d.Set("bytes_in", pz.BytesIn); err != nil {
		return diag.Errorf("unable to set resource pull zone `bytes_in` read value: %v", err)
	}

	if err := d.Set("bytes_out", pz.BytesOut); err != nil {
		return diag.Errorf("unable to set resource pull zone `bytes_out` read value: %v", err)
	}

	if err := d.Set("packets_per_second", pz.PacketsPerSec); err != nil {
		return diag.Errorf("unable to set resource pull zone `packets_per_sec` read value: %v", err)
	}

	if err := d.Set("date_purged", pz.DatePurged); err != nil {
		return diag.Errorf("unable to set resource pull zone `date_purged` read value: %v", err)
	}

	if err := d.Set("regions", pz.Regions); err != nil {
		return diag.Errorf("unable to set resource pull zone `regions` read value: %v", err)
	}

	return nil
}

func resourceVultrCDNPullZoneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	req := &govultr.CDNZoneReq{}

	if d.HasChange("label") {
		log.Print("[INFO] Updating pull zone `label`")
		req.Label = d.Get("label").(string)
	}

	if d.HasChange("origin_scheme") {
		log.Print("[INFO] Updating pull zone `origin_scheme`")
		req.OriginScheme = d.Get("origin_scheme").(string)
	}

	if d.HasChange("origin_domain") {
		log.Print("[INFO] Updating pull zone `origin_domain`")
		req.OriginDomain = d.Get("origin_domain").(string)
	}

	if d.HasChange("cors") {
		log.Print("[INFO] Updating pull zone `cors`")
		req.CORS = d.Get("cors").(bool)
	}

	if d.HasChange("gzip") {
		log.Print("[INFO] Updating pull zone `gzip`")
		req.GZIP = d.Get("gzip").(bool)
	}

	if d.HasChange("block_ai") {
		log.Print("[INFO] Updating pull zone `block_ai`")
		req.BlockAI = d.Get("block_ai").(bool)
	}

	if d.HasChange("block_bad_bots") {
		log.Print("[INFO] Updating pull zone `block_bad_bots`")
		req.BlockBadBots = d.Get("block_bad_bots").(bool)
	}

	pz, _, err := client.CDN.UpdatePullZone(ctx, d.Id(), req)
	if err != nil {
		return diag.Errorf("error updating cdn pull zone: %v", err)
	}

	d.SetId(pz.ID)

	log.Printf("[INFO] CDN Pull Zone ID: %s", d.Id())

	return resourceVultrCDNPullZoneRead(ctx, d, meta)
}

func resourceVultrCDNPullZoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	log.Printf("[INFO] Deleting pull zone: %s", d.Id())
	if err := client.CDN.DeletePullZone(ctx, d.Id()); err != nil {
		return diag.Errorf("error destroying pull zone (%s): %v", d.Id(), err)
	}

	return nil
}
