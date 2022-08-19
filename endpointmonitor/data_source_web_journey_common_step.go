package endpointmonitor

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceWebJourneyCommonStep() *schema.Resource {
	return &schema.Resource{
		Description: "Search for an individual Common Web Journey Step. This will only allow a single result to be returned.",
		ReadContext: dataSourceWebJourneyCommonStepRead,
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

func dataSourceWebJourneyCommonSteps() *schema.Resource {
	return &schema.Resource{
		Description: "Search for multiple Common Web Journey Steps. A list of ids will be returned for all matches found.",
		ReadContext: dataSourceWebJourneyCommonStepsRead,
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

func dataSourceWebJourneyCommonStepRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	steps, err := c.SearchWebJourneyCommonSteps(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(*steps) > 1 {
		err := errors.New("more than one Web Journey Common Step from given search value")
		return diag.FromErr(err)
	}

	if len(*steps) < 1 {
		err := errors.New("no Web Journey Common Step from given search value")
		return diag.FromErr(err)
	}

	for _, step := range *steps {
		d.SetId(strconv.Itoa(step.Id))
		d.Set("id", step.Id)
	}

	return diags
}

func dataSourceWebJourneyCommonStepsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	steps, err := c.SearchWebJourneyCommonSteps(d.Get("search").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(*steps) < 1 {
		err := errors.New("no Web Journey Common Steps from given search value")
		return diag.FromErr(err)
	}

	var ids []int

	for _, step := range *steps {
		ids = append(ids, step.Id)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("ids", ids)

	return diags
}
