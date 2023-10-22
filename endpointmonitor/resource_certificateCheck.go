package endpointmonitor

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func certificateCheck() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage TLS certificate checks that test a given URL for an expected response.",
		CreateContext: resourceCertificateCheckCreate,
		ReadContext:   resourceCertificateCheckRead,
		UpdateContext: resourceCertificateCheckUpdate,
		DeleteContext: resourceCertificateCheckDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "A name to describe in the check, used throughout EndPoint Monitor to describe this check, including in notifications.",
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A space to provide a longer description of the check if needed. Will default to the name if not set.",
				Optional:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Allows the enabling/disabling of the check from executing.",
				Optional:    true,
				Default:     true,
			},
			"check_frequency": {
				Type:         schema.TypeInt,
				Description:  "The frequency the check will be run in seconds.",
				Optional:     true,
				Default:      60,
				ValidateFunc: validatePositiveInt(),
			},
			"maintenance_override": {
				Type:        schema.TypeBool,
				Description: "If set true then notifications and alerts will be suppressed for the check.",
				Optional:    true,
				Default:     false,
			},
			"url": {
				Type:         schema.TypeString,
				Description:  "The URL to check the certificate for.",
				Required:     true,
				ValidateFunc: validateUrl(),
			},
			"trigger_count": {
				Type:         schema.TypeInt,
				Description:  "The sequential number of failures that need to occur for a check to trigger an alert or notification.",
				Required:     true,
				ValidateFunc: validatePositiveInt(),
			},
			"result_retention": {
				Type:         schema.TypeInt,
				Description:  "The number of days to store historic results of the check.",
				Optional:     true,
				Default:      366,
				ValidateFunc: validatePositiveInt(),
			},
			"warning_days_remaining": {
				Type:         schema.TypeInt,
				Description:  "The maximum number of remaining days on a certificate before an warning is triggered.",
				Required:     true,
				ValidateFunc: validatePositiveInt(),
			},
			"alert_days_remaining": {
				Type:         schema.TypeInt,
				Description:  "The maximum number of remaining days on a certificate before a failure is triggered.",
				Required:     true,
				ValidateFunc: validatePositiveInt(),
			},
			"check_date_only": {
				Type:        schema.TypeBool,
				Description: "If set to true, then only certificate validity period will be checked and nothing else.",
				Optional:    true,
				Default:     false,
			},
			"check_full_chain": {
				Type:        schema.TypeBool,
				Description: "If set to false, only the initially returned certificate from the given URL will be checked, and not the full certificate chain.",
				Optional:    true,
				Default:     true,
			},
			"check_host_id": {
				Type:         schema.TypeInt,
				Description:  "The id of the Check Host to run the check on.",
				Optional:     true,
				ValidateFunc: validatePositiveInt(),
			},
			"check_host_group_id": {
				Type:         schema.TypeInt,
				Description:  "The id of the Check Host Group to run the check on.",
				Optional:     true,
				ValidateFunc: validatePositiveInt(),
			},
			"check_group_id": {
				Type:         schema.TypeInt,
				Description:  "The id of the Check Group the check belongs to. This also determines check frequency.",
				Required:     true,
				ValidateFunc: validatePositiveInt(),
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCertificateCheckRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkId := d.Id()

	check, err := c.GetCertificateCheck(checkId)
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.IsNewResource() && check == nil {
		d.SetId("")
		return nil
	}

	mapCertificateCheckSchema(*check, d)

	return diags
}

func resourceCertificateCheckCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	check := mapCertificateCheck(d)

	o, err := c.CreateCertificateCheck(check)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(o.Id)))

	resourceCertificateCheckRead(ctx, d, m)

	return diags
}

func resourceCertificateCheckUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	_, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChangesExcept() {
		check := mapCertificateCheck(d)

		_, err := c.UpdateCertificateCheck(check)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceCertificateCheckRead(ctx, d, m)
}

func resourceCertificateCheckDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteCheck(checkId)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func mapCertificateCheck(d *schema.ResourceData) CertificateCheck {
	checkId, err := strconv.Atoi(d.Id())
	if err != nil {
		checkId = 0
	}

	check := CertificateCheck{
		Id:                   checkId,
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		Enabled:              d.Get("enabled").(bool),
		CheckFrequency:       d.Get("check_frequency").(int),
		CheckType:            "TLS_CERTIFICATE",
		MaintenanceOverride:  d.Get("maintenance_override").(bool),
		URL:                  d.Get("url").(string),
		TriggerCount:         d.Get("trigger_count").(int),
		ResultRetentionDays:  d.Get("result_retention").(int),
		WarningDaysRemaining: d.Get("warning_days_remaining").(int),
		AlertDaysRemaining:   d.Get("alert_days_remaining").(int),
		CheckDatesOnly:       d.Get("check_date_only").(bool),
		CheckFullChain:       d.Get("check_full_chain").(bool),
		CheckGroup: CheckGroup{
			Id: d.Get("check_group_id").(int),
		},
	}

	if d.Get("check_host_id") != nil {
		check.CheckHost = &CheckHost{
			Id: d.Get("check_host_id").(int),
		}
	}

	if d.Get("check_host_group_id") != nil {
		check.HostGroup = &HostGroup{
			Id: d.Get("check_host_group_id").(int),
		}
	}

	return check
}

func mapCertificateCheckSchema(check CertificateCheck, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(check.Id))
	d.Set("name", check.Name)
	d.Set("description", check.Description)
	d.Set("enabled", check.Enabled)
	d.Set("check_frequency", check.CheckFrequency)
	d.Set("maintenance_override", check.MaintenanceOverride)
	d.Set("trigger_count", check.TriggerCount)
	d.Set("result_retention", check.ResultRetentionDays)
	d.Set("url", check.URL)
	d.Set("warning_days_remaining", check.WarningDaysRemaining)
	d.Set("alert_days_remaining", check.AlertDaysRemaining)
	d.Set("check_date_only", check.CheckDatesOnly)
	d.Set("check_full_chain", check.CheckFullChain)
	d.Set("check_group_id", check.CheckGroup.Id)

	if check.CheckHost != nil {
		d.Set("check_host_id", check.CheckHost.Id)
	}

	if check.HostGroup != nil {
		d.Set("check_host_group_id", check.HostGroup.Id)
	}
}
