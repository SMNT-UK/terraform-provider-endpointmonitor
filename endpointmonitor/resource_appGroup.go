package endpointmonitor

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func appGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage App Groups, the top-level organisational groups of checks running in EndPoint Monitor.",
		CreateContext: resourceAppGroupCreate,
		ReadContext:   resourceAppGroupRead,
		UpdateContext: resourceAppGroupUpdate,
		DeleteContext: resourceAppGroupDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the App Group. This will be used in alerts and notifications.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Space to provide a longer description of this App Group.",
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceAppGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	appGroupId := d.Id()

	appGroup, err := c.GetAppGroup(appGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.IsNewResource() && appGroup == nil {
		d.SetId("")
		return nil
	}

	mapAppGroupSchema(*appGroup, d)

	return diags
}

func resourceAppGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	appGroup := mapAppGroup(d)

	o, err := c.CreateAppGroup(appGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(o.Id)))

	resourceAppGroupRead(ctx, d, m)

	return diags
}

func resourceAppGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	_, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChangesExcept() {
		appGroup := mapAppGroup(d)

		_, err := c.UpdateAppGroup(appGroup)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceAppGroupRead(ctx, d, m)
}

func resourceAppGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteAppGroup(checkId)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func mapAppGroup(d *schema.ResourceData) AppGroup {
	appGroupId, err := strconv.Atoi(d.Id())
	if err != nil {
		appGroupId = 0
	}

	return AppGroup{
		Id:          appGroupId,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
}

func mapAppGroupSchema(appGroup AppGroup, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(appGroup.Id))
	d.Set("name", appGroup.Name)
	d.Set("description", appGroup.Description)
}
