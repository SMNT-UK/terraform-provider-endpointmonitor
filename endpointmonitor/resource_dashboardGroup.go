package endpointmonitor

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dashboardGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage Dashboard Groups, the top-level organisational groups of checks running in EndPoint Monitor.",
		CreateContext: resourceDashboardGroupCreate,
		ReadContext:   resourceDashboardGroupRead,
		UpdateContext: resourceDashboardGroupUpdate,
		DeleteContext: resourceDashboardGroupDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Dashboard Group. This will be used in alerts and notifications.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Space to provide a longer description of this Dashboard Group.",
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDashboardGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	dashboardGroupId := d.Id()

	dashboardGroup, err := c.GetDashboardGroup(dashboardGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.IsNewResource() && dashboardGroup == nil {
		d.SetId("")
		return nil
	}

	mapDashboardGroupSchema(*dashboardGroup, d)

	return diags
}

func resourceDashboardGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	dashboardGroup := mapDashboardGroup(d)

	o, err := c.CreateDashboardGroup(dashboardGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(o.Id)))

	resourceDashboardGroupRead(ctx, d, m)

	return diags
}

func resourceDashboardGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	_, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChangesExcept() {
		dashboardGroup := mapDashboardGroup(d)

		_, err := c.UpdateDashboardGroup(dashboardGroup)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDashboardGroupRead(ctx, d, m)
}

func resourceDashboardGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	dashboardGroupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteDashboardGroup(dashboardGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func mapDashboardGroup(d *schema.ResourceData) DashboardGroup {
	dashboardGroupId, err := strconv.Atoi(d.Id())
	if err != nil {
		dashboardGroupId = 0
	}

	return DashboardGroup{
		Id:          dashboardGroupId,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
}

func mapDashboardGroupSchema(dashboardGroup DashboardGroup, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(dashboardGroup.Id))
	d.Set("name", dashboardGroup.Name)
	d.Set("description", dashboardGroup.Description)
}
