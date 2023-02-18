package endpointmonitor

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHostGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Search for an individual Check Host Group. This will only allow a single result to be returned.",
		ReadContext: dataSourceHostGroupRead,
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

func dataSourceHostGroups() *schema.Resource {
	return &schema.Resource{
		Description: "Search for multiple Check Host Groups. A list of ids will be returned for all matches found.",
		ReadContext: dataSourceHostGroupsRead,
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

func dataSourceHostGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//filter := d.Get("filter").(*schema.Set)

	hostGroups, err := c.SearchHostGroups(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(*hostGroups) > 1 {
		err := errors.New("more than one host group found from given search value")
		return diag.FromErr(err)
	}

	for _, host := range *hostGroups {
		d.SetId(strconv.Itoa(host.Id))
	}

	return diags
}

func dataSourceHostGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//filter := d.Get("filter").(*schema.Set)

	hostGroups, err := c.SearchHostGroups(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []int

	for _, hostGroup := range *hostGroups {
		ids = append(ids, hostGroup.Id)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("ids", ids)

	return diags
}
