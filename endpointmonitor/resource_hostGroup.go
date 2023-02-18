package endpointmonitor

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func hostGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage groups of Check Hosts that can assigned to checks to execute on.",
		CreateContext: resourceHostGroupCreate,
		ReadContext:   resourceHostGroupRead,
		UpdateContext: resourceHostGroupUpdate,
		DeleteContext: resourceHostGroupDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "A name to identify this group of hosts by. This will show in searches and alerts.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Space for a longer description to define this group of hosts by. Not required.",
				Optional:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Enable or disable checks assigned to this host group.",
				Required:    true,
			},
			"check_host_ids": {
				Type:        schema.TypeList,
				Description: "List of Check Host id's that are to be part of this Host Group.",
				Required:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeInt,
					Description: "The id of a Check Host to be part of this Host Group.",
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceHostGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	hostGroupId := d.Id()

	hostGroup, err := c.GetHostGroup(hostGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.IsNewResource() && hostGroup == nil {
		d.SetId("")
		return nil
	}

	mapHostGroupSchema(*hostGroup, d)

	return diags
}

func resourceHostGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	hostGroup := mapHostGroup(d)
	o, err := c.CreateHostGroup(hostGroup)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(o.Id)))

	resourceHostGroupRead(ctx, d, m)

	return diags
}

func resourceHostGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	_, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChangesExcept() {
		hostGroup := mapHostGroup(d)

		_, err := c.UpdateHostGroup(hostGroup)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceHostGroupRead(ctx, d, m)
}

func resourceHostGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	hostGroupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteHostGroup(hostGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func mapHostGroup(d *schema.ResourceData) HostGroup {
	hostGroupId, err := strconv.Atoi(d.Id())
	if err != nil {
		hostGroupId = 0
	}

	checkHosts := d.Get("check_host_ids").([]interface{})
	check_host_ids := make([]CheckHost, len(checkHosts))

	for index, check_host_id := range checkHosts {
		check_host_ids[index] = CheckHost{Id: check_host_id.(int)}
	}

	return HostGroup{
		Id:          hostGroupId,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Enabled:     d.Get("enabled").(bool),
		Hosts:       check_host_ids,
	}
}

func mapHostGroupSchema(hostGroup HostGroup, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(hostGroup.Id))
	d.Set("name", hostGroup.Name)
	d.Set("description", hostGroup.Description)
	d.Set("enabled", hostGroup.Enabled)

	check_host_ids := make([]int, len(hostGroup.Hosts))

	for index, checkHost := range hostGroup.Hosts {
		check_host_ids[index] = checkHost.Id
	}

	d.Set("check_host_ids", check_host_ids)
}
