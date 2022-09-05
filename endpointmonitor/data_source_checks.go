package endpointmonitor

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCheck() *schema.Resource {
	return &schema.Resource{
		Description: "Search for an individual Check. This will only allow a single result to be returned.",
		ReadContext: dataSourceCheckRead,
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

func dataSourceChecks() *schema.Resource {
	return &schema.Resource{
		Description: "Search for multiple Checks. A list of ids will be returned for all matches found.",
		ReadContext: dataSourceChecksRead,
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

func dataSourceCheckRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//filter := d.Get("filter").(*schema.Set)

	checks, err := c.SearchChecks(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(*checks) > 1 {
		err := errors.New("more than one check found from given search value")
		return diag.FromErr(err)
	}

	for _, check := range *checks {
		d.SetId(strconv.Itoa(check.Id))
	}

	return diags
}

func dataSourceChecksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//filter := d.Get("filter").(*schema.Set)

	checks, err := c.SearchChecks(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []int

	for _, check := range *checks {
		ids = append(ids, check.Id)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("ids", ids)

	return diags
}
