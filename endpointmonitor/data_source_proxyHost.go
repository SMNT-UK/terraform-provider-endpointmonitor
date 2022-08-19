package endpointmonitor

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProxyHost() *schema.Resource {
	return &schema.Resource{
		Description: "Search for an individual Proxy Host. This will only allow a single result to be returned.",
		ReadContext: dataSourceProxyHostRead,
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

func dataSourceProxyHosts() *schema.Resource {
	return &schema.Resource{
		Description: "Search for multiple Proxy Hosts. A list of ids will be returned for all matches found.",
		ReadContext: dataSourceProxyHostsRead,
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

func dataSourceProxyHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//filter := d.Get("filter").(*schema.Set)

	proxyHosts, err := c.SearchProxyHosts(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(*proxyHosts) > 1 {
		err := errors.New("more than one proxy host found from given search value")
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	for _, proxyHost := range *proxyHosts {
		d.Set("id", proxyHost.Id)
	}

	return diags
}

func dataSourceProxyHostsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//filter := d.Get("filter").(*schema.Set)

	proxyHosts, err := c.SearchProxyHosts(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []int

	for _, proxyHost := range *proxyHosts {
		ids = append(ids, proxyHost.Id)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("ids", ids)

	return diags
}
