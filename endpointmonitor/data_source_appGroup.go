package endpointmonitor

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Search for an individual App Group. This will only allow a single result to be returned.",
		ReadContext: dataSourceAppGroupRead,
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

func dataSourceAppGroups() *schema.Resource {
	return &schema.Resource{
		Description: "Search for multiple App Groups. A list of ids will be returned for all matches found.",
		ReadContext: dataSourceAppGroupsRead,
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

func dataSourceAppGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//filter := d.Get("filter").(*schema.Set)

	appGroups, err := c.SearchAppGroups(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(*appGroups) > 1 {
		err := errors.New("more than one app group found from given search value")
		return diag.FromErr(err)
	}

	for _, appGroup := range *appGroups {
		d.SetId(strconv.Itoa(appGroup.Id))
	}

	return diags
}

func dataSourceAppGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//filter := d.Get("filter").(*schema.Set)

	appGroups, err := c.SearchAppGroups(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []int

	for _, appGroup := range *appGroups {
		ids = append(ids, appGroup.Id)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("ids", ids)

	return diags
}
