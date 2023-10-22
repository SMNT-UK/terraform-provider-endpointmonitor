package endpointmonitor

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func socketCheck() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage socket checks which test to ensure a hostname is listening on a pre-defined port.",
		CreateContext: resourceSocketCheckCreate,
		ReadContext:   resourceSocketCheckRead,
		UpdateContext: resourceSocketCheckUpdate,
		DeleteContext: resourceSocketCheckDelete,
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
			"hostname": {
				Type:         schema.TypeString,
				Description:  "The hostname to check the socket is open on.",
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"port": {
				Type:         schema.TypeInt,
				Description:  "The TCP port to check is listening.",
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

func resourceSocketCheckRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkId := d.Id()

	check, err := c.GetSocketCheck(checkId)
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.IsNewResource() && check == nil {
		d.SetId("")
		return nil
	}

	mapSocketCheckSchema(*check, d)

	return diags
}

func resourceSocketCheckCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	check := mapSocketCheck(d)

	o, err := c.CreateSocketCheck(check)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(o.Id)))

	resourceSocketCheckRead(ctx, d, m)

	return diags
}

func resourceSocketCheckUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	_, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChangesExcept() {
		check := mapSocketCheck(d)

		_, err := c.UpdateSocketCheck(check)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSocketCheckRead(ctx, d, m)
}

func resourceSocketCheckDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func mapSocketCheck(d *schema.ResourceData) SocketCheck {
	checkId, err := strconv.Atoi(d.Id())
	if err != nil {
		checkId = 0
	}

	check := SocketCheck{
		Id:                  checkId,
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Enabled:             d.Get("enabled").(bool),
		CheckFrequency:      d.Get("check_frequency").(int),
		CheckType:           "SOCKET",
		MaintenanceOverride: d.Get("maintenance_override").(bool),
		Hostname:            d.Get("hostname").(string),
		Port:                d.Get("port").(int),
		TriggerCount:        d.Get("trigger_count").(int),
		ResultRetentionDays: d.Get("result_retention").(int),
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

func mapSocketCheckSchema(check SocketCheck, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(check.Id))
	d.Set("name", check.Name)
	d.Set("description", check.Description)
	d.Set("enabled", check.Enabled)
	d.Set("check_frequency", check.CheckFrequency)
	d.Set("maintenance_override", check.MaintenanceOverride)
	d.Set("hostname", check.Hostname)
	d.Set("port", check.Port)
	d.Set("trigger_count", check.TriggerCount)
	d.Set("result_retention", check.ResultRetentionDays)
	d.Set("check_group_id", check.CheckGroup.Id)

	if check.CheckHost != nil {
		d.Set("check_host_id", check.CheckHost.Id)
	}

	if check.HostGroup != nil {
		d.Set("check_host_group_id", check.HostGroup.Id)
	}
}
