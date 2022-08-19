package endpointmonitor

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCheckHost() *schema.Resource {
	return &schema.Resource{
		Description: "Search for an individual Check Host. This will only allow a single result to be returned.",
		ReadContext: dataSourceCheckHostRead,
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

func dataSourceCheckHosts() *schema.Resource {
	return &schema.Resource{
		Description: "Search for multiple Check Hosts. A list of ids will be returned for all matches found.",
		ReadContext: dataSourceCheckHostsRead,
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

func dataSourceCheckHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//filter := d.Get("filter").(*schema.Set)

	hosts, err := c.SearchCheckHosts(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(*hosts) > 1 {
		err := errors.New("more than one check host found from given search value")
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	for _, host := range *hosts {
		d.Set("id", host.Id)
	}

	return diags
}

func dataSourceCheckHostsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//filter := d.Get("filter").(*schema.Set)

	hosts, err := c.SearchCheckHosts(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []int

	for _, host := range *hosts {
		ids = append(ids, host.Id)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("ids", ids)

	return diags
}
