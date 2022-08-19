package endpointmonitor

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func maintenancePeriod() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage scheduled maintenance periods to prevent checks from alerting during certain periods of the day or week.",
		CreateContext: resourceMaintenancePeriodCreate,
		ReadContext:   resourceMaintenancePeriodRead,
		UpdateContext: resourceMaintenancePeriodUpdate,
		DeleteContext: resourceMaintenancePeriodDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Space for a description of the maintenance periods purpose.",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Enable or disable the maintenance period from applying to attached checks.",
				Optional:    true,
				Default:     true,
			},
			"day_of_week": {
				Type:        schema.TypeString,
				Description: "The day of week the maintenance period applies to. Set as ALL for every day of the week. Must otherwise be SUNDAY, MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY or SATURDAY.",
				Required:    true,
			},
			"start_time": {
				Type:        schema.TypeString,
				Description: "The start time of the maintenance period in format 24HH:MM.",
				Required:    true,
			},
			"end_time": {
				Type:        schema.TypeString,
				Description: "The end of time the maintenance period in format 24HH:MM.",
				Required:    true,
			},
			"check_ids": {
				Type:        schema.TypeList,
				Description: "A list of ids of Checks that are directly linked to the maintenance period.",
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validatePositiveInt(),
				},
			},
			"check_group_ids": {
				Type:        schema.TypeList,
				Description: "A list of ids of Check Groups that are directly linked to the maintenance period.",
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validatePositiveInt(),
				},
			},
			"app_group_ids": {
				Type:        schema.TypeList,
				Description: "A list of ids of App Groups that are linked to this maintenance period.",
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validatePositiveInt(),
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceMaintenancePeriodRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	maintenancePeriodId := d.Id()

	maintenancePeriod, err := c.GetMaintenancePeriod(maintenancePeriodId)
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.IsNewResource() && maintenancePeriod == nil {
		d.SetId("")
		return nil
	}

	mapMaintenancePeriodSchema(*maintenancePeriod, d)

	return diags
}

func resourceMaintenancePeriodCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	maintenancePeriod := mapMaintenancePeriod(d)

	o, err := c.CreateMaintenancePeriod(maintenancePeriod)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(o.Id)))

	resourceMaintenancePeriodRead(ctx, d, m)

	return diags
}

func resourceMaintenancePeriodUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	_, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChangesExcept() {
		maintenancePeriod := mapMaintenancePeriod(d)

		_, err := c.UpdateMaintenancePeriod(maintenancePeriod)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceMaintenancePeriodRead(ctx, d, m)
}

func resourceMaintenancePeriodDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteMaintenancePeriod(id)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func mapMaintenancePeriod(d *schema.ResourceData) MaintenancePeriod {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		id = 0
	}

	return MaintenancePeriod{
		Id:          id,
		Description: d.Get("description").(string),
		Enabled:     d.Get("enabled").(bool),
		DayOfWeek:   d.Get("day_of_week").(string),
		StartTime:   d.Get("start_time").(string),
		EndTime:     d.Get("end_time").(string),
		Checks:      d.Get("check_ids").([]int),
		CheckGroups: d.Get("check_group_ids").([]int),
		AppGroups:   d.Get("app_group_ids").([]int),
	}
}

func mapMaintenancePeriodSchema(maintenancePeriod MaintenancePeriod, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(maintenancePeriod.Id))
	d.Set("description", maintenancePeriod.Description)
	d.Set("enabled", maintenancePeriod.Enabled)
	d.Set("day_of_week", maintenancePeriod.DayOfWeek)
	d.Set("start_time", maintenancePeriod.StartTime)
	d.Set("end_time", maintenancePeriod.EndTime)
	d.Set("check_ids", maintenancePeriod.Checks)
	d.Set("check_group_ids", maintenancePeriod.CheckGroups)
	d.Set("app_group_ids", maintenancePeriod.AppGroups)
}
