package endpointmonitor

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func pingCheck() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage ping checks to check a hostname or address is online.",
		CreateContext: resourcePingCheckCreate,
		ReadContext:   resourcePingCheckRead,
		UpdateContext: resourcePingCheckUpdate,
		DeleteContext: resourcePingCheckDelete,
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
				Optional:    true,
				Description: "A space to provide a longer description of the check if needed. Will default to the name if not set.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Allows the enabling/disabling of the check from executing.",
				Optional:    true,
				Default:     true,
			},
			"maintenance_override": {
				Type:        schema.TypeBool,
				Description: "If set true then notifications and alerts will be suppressed for the check.",
				Optional:    true,
				Default:     true,
			},
			"hostname": {
				Type:         schema.TypeString,
				Description:  "The hostname to ping.",
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"warning_response_time": {
				Type:         schema.TypeInt,
				Description:  "The warning response time threshold in millisonds.",
				Required:     true,
				ValidateFunc: validatePositiveInt(),
			},
			"timeout_time": {
				Type:         schema.TypeInt,
				Description:  "The number of milliseconds to wait for a response before giving up.",
				Required:     true,
				ValidateFunc: validatePositiveInt(),
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
			"check_host_id": {
				Type:         schema.TypeInt,
				Description:  "The id of the Check Host to run the check on.",
				Required:     true,
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
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourcePingCheckRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkId := d.Id()

	check, err := c.GetPingCheck(checkId)
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.IsNewResource() && check == nil {
		d.SetId("")
		return nil
	}

	mapPingCheckSchema(*check, d)

	return diags
}

func resourcePingCheckCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	check := mapPingCheck(d)

	o, err := c.CreatePingCheck(check)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(o.Id)))

	resourcePingCheckRead(ctx, d, m)

	return diags
}

func resourcePingCheckUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	_, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChangesExcept() {
		check := mapPingCheck(d)

		_, err := c.UpdatePingCheck(check)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourcePingCheckRead(ctx, d, m)
}

func resourcePingCheckDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func mapPingCheck(d *schema.ResourceData) PingCheck {
	checkId, err := strconv.Atoi(d.Id())
	if err != nil {
		checkId = 0
	}

	return PingCheck{
		Id:                  checkId,
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Enabled:             d.Get("enabled").(bool),
		CheckType:           "PING",
		MaintenanceOverride: d.Get("maintenance_override").(bool),
		Hostname:            d.Get("hostname").(string),
		WarningRepsonseTime: d.Get("warning").(int),
		Timeout:             d.Get("timeout").(int),
		TriggerCount:        d.Get("trigger_count").(int),
		ResultRetentionDays: d.Get("result_retention").(int),
		CheckHost: CheckHost{
			Id: d.Get("check_host_id").(int),
		},
		CheckGroup: CheckGroup{
			Id: d.Get("check_group_id").(int),
		},
	}
}

func mapPingCheckSchema(check PingCheck, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(check.Id))
	d.Set("name", check.Name)
	d.Set("description", check.Description)
	d.Set("enabled", check.Enabled)
	d.Set("maintenance_override", check.MaintenanceOverride)
	d.Set("hostname", check.Hostname)
	d.Set("warning", check.WarningRepsonseTime)
	d.Set("timeout", check.Timeout)
	d.Set("trigger_count", check.TriggerCount)
	d.Set("result_retention", check.ResultRetentionDays)
	d.Set("check_host_id", check.CheckHost.Id)
	d.Set("check_group_id", check.CheckGroup.Id)
}
