package endpointmonitor

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func checkGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage Check Groups, the initial grouping of checks running on EndPoint Monitor.",
		CreateContext: resourceCheckGroupCreate,
		ReadContext:   resourceCheckGroupRead,
		UpdateContext: resourceCheckGroupUpdate,
		DeleteContext: resourceCheckGroupDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "A meaningful name of what this group contains. This will be used in alerts and notifications.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A space to provide a longer description of this group.",
				Required:    true,
			},
			"check_frequency": {
				Type:        schema.TypeInt,
				Description: "The frequency in seconds that checks within this group should be run.",
				Required:    true,
			},
			"app_group_id": {
				Type:        schema.TypeInt,
				Description: "The id of the App Group this belongs to.",
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCheckGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkGroupId := d.Id()

	checkGroup, err := c.GetCheckGroup(checkGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.IsNewResource() && checkGroup == nil {
		d.SetId("")
		return nil
	}

	mapCheckGroupSchema(*checkGroup, d)

	return diags
}

func resourceCheckGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkGroup := mapCheckGroup(d)

	o, err := c.CreateCheckGroup(checkGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(o.Id)))

	resourceCheckGroupRead(ctx, d, m)

	return diags
}

func resourceCheckGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	_, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChangesExcept() {
		checkGroup := mapCheckGroup(d)

		_, err := c.UpdateCheckGroup(checkGroup)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceCheckGroupRead(ctx, d, m)
}

func resourceCheckGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkGroupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteCheckGroup(checkGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func mapCheckGroup(d *schema.ResourceData) CheckGroup {
	checkGroupId, err := strconv.Atoi(d.Id())
	if err != nil {
		checkGroupId = 0
	}

	return CheckGroup{
		Id:             checkGroupId,
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		CheckFrequency: d.Get("check_frequency").(int),
		AppGroup: AppGroup{
			Id: d.Get("app_group_id").(int),
		},
	}
}

func mapCheckGroupSchema(checkGroup CheckGroup, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(checkGroup.Id))
	d.Set("name", checkGroup.Name)
	d.Set("description", checkGroup.Description)
	d.Set("check_frequency", checkGroup.CheckFrequency)
	d.Set("app_group_id", checkGroup.AppGroup.Id)
}
