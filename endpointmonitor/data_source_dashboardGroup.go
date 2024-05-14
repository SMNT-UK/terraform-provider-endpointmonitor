package endpointmonitor

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDashboardGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Search for an individual Dashboard Group. This will only allow a single result to be returned.",
		ReadContext: dataSourceDashboardGroupRead,
		Schema: map[string]*schema.Schema{
			"search": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDashboardGroups() *schema.Resource {
	return &schema.Resource{
		Description: "Search for multiple Dashboard Groups. A list of ids will be returned for all matches found.",
		ReadContext: dataSourceDashboardGroupsRead,
		Schema: map[string]*schema.Schema{
			"search": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourceDashboardGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//filter := d.Get("filter").(*schema.Set)

	dashboardGroups, err := c.SearchDashboardGroups(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(*dashboardGroups) > 1 {
		err := errors.New("more than one dashboard group found from given search value")
		return diag.FromErr(err)
	}

	for _, dashboardGroup := range *dashboardGroups {
		d.SetId(strconv.Itoa(dashboardGroup.Id))
	}

	return diags
}

func dataSourceDashboardGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//filter := d.Get("filter").(*schema.Set)

	dashboardGroups, err := c.SearchDashboardGroups(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []int

	for _, dashboardGroup := range *dashboardGroups {
		ids = append(ids, dashboardGroup.Id)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("ids", ids)

	return diags
}
