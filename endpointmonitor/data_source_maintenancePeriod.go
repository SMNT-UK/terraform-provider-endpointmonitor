package endpointmonitor

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMaintenancePeriod() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMaintenancePeriodRead,
		Description: "Search for an individual Scheduled Maintenance Period. This will only allow a single result to be returned.",
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

func dataSourceMaintenancePeriods() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMaintenancePeriodsRead,
		Description: "Search for multiple Scheduled Maintenance Periods. A list of ids will be returned for all matches found.",
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

func dataSourceMaintenancePeriodRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	maintenancePeriods, err := c.SearchMaintenancePeriods(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(*maintenancePeriods) > 1 {
		err := errors.New("more than one maintenance period found from given search value")
		return diag.FromErr(err)
	}

	for _, maintenancePeriod := range *maintenancePeriods {
		d.SetId(strconv.Itoa(maintenancePeriod.Id))
	}

	return diags
}

func dataSourceMaintenancePeriodsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	maintenancePeriods, err := c.SearchMaintenancePeriods(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []int

	for _, maintenancePeriod := range *maintenancePeriods {
		ids = append(ids, maintenancePeriod.Id)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("ids", ids)

	return diags
}
