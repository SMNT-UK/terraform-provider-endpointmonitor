package endpointmonitor

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCheckGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Search for an individual Check Group. This will only allow a single result to be returned.",
		ReadContext: dataSourceCheckGroupRead,
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

func dataSourceCheckGroups() *schema.Resource {
	return &schema.Resource{
		Description: "Search for multiple Check Groups. A list of ids will be returned for all matches found.",
		ReadContext: dataSourceCheckGroupsRead,
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

func dataSourceCheckGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//filter := d.Get("filter").(*schema.Set)

	checkGroups, err := c.SearchCheckGroups(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(*checkGroups) > 1 {
		err := errors.New("more than one check group found from given search value")
		return diag.FromErr(err)
	}

	for _, checkGroup := range *checkGroups {
		d.SetId(strconv.Itoa(checkGroup.Id))
	}

	return diags
}

func dataSourceCheckGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//filter := d.Get("filter").(*schema.Set)

	checkGroups, err := c.SearchCheckGroups(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []int

	for _, checkGroup := range *checkGroups {
		ids = append(ids, checkGroup.Id)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("ids", ids)

	return diags
}
