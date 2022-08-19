package endpointmonitor

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dnsCheck() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage DNS checks which check that a hostname reolves to an known set of addresses.",
		CreateContext: resourceDNSCheckCreate,
		ReadContext:   resourceDNSCheckRead,
		UpdateContext: resourceDNSCheckUpdate,
		DeleteContext: resourceDNSCheckDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "A name to describe in the check, used throughout EndPoint Monitor to describe this check, including in notifications.",
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A space to provide a longer description of the check if needed. Will default to the name if not set.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Allows the enabling/disabling of the check from executing.",
				Optional:    true,
				Default:     true,
			},
			"maintenance_override": {
				Type:        schema.TypeBool,
				Description: "If set true then notifications and alerts will be suppressed for the check.",
				Optional:    true,
				Default:     true,
			},
			"hostname": {
				Type:         schema.TypeString,
				Description:  "The hostname to check.",
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"expected_addresses": {
				Type:        schema.TypeList,
				Description: "The list of addresses expected to be returned for the given hostname. Addresses returned outside of this list will result in the check reporting a failure.",
				Required:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			"trigger_count": {
				Type:         schema.TypeInt,
				Description:  "The sequential number of failures that need to occur for a check to trigger an alert or notification.",
				Required:     true,
				ValidateFunc: validatePositiveInt(),
			},
			"result_retention": {
				Type:         schema.TypeInt,
				Description:  "The number of days to store historic results of the check.",
				Optional:     true,
				Default:      366,
				ValidateFunc: validatePositiveInt(),
			},
			"check_host_id": {
				Type:         schema.TypeInt,
				Description:  "The id of the Check Host to run the check on.",
				Required:     true,
				ValidateFunc: validatePositiveInt(),
			},
			"check_group_id": {
				Type:         schema.TypeInt,
				Description:  "The id of the Check Group the check belongs to. This also determines check frequency.",
				Required:     true,
				ValidateFunc: validatePositiveInt(),
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceDNSCheckRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkId := d.Id()

	check, err := c.GetDNSCheck(checkId)
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.IsNewResource() && check == nil {
		d.SetId("")
		return nil
	}

	mapDNSCheckSchema(*check, d)

	return diags
}

func resourceDNSCheckCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	check := mapDNSCheck(d)

	o, err := c.CreateDNSCheck(check)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(o.Id)))

	resourceDNSCheckRead(ctx, d, m)

	return diags
}

func resourceDNSCheckUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	_, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChangesExcept() {
		check := mapDNSCheck(d)

		_, err := c.UpdateDNSCheck(check)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDNSCheckRead(ctx, d, m)
}

func resourceDNSCheckDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteCheck(checkId)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func mapDNSCheck(d *schema.ResourceData) DNSCheck {
	checkId, err := strconv.Atoi(d.Id())
	if err != nil {
		checkId = 0
	}

	expected_addresses := d.Get("expected_addresses").([]interface{})
	addresses := make([]string, len(expected_addresses))

	for index, address := range expected_addresses {
		addresses[index] = address.(string)
	}

	return DNSCheck{
		Id:                  checkId,
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Enabled:             d.Get("enabled").(bool),
		CheckType:           "DNS",
		MaintenanceOverride: d.Get("maintenance_override").(bool),
		Hostname:            d.Get("hostname").(string),
		ExpectedAddresses:   addresses,
		TriggerCount:        d.Get("trigger_count").(int),
		ResultRetentionDays: d.Get("result_retention").(int),
		CheckHost: CheckHost{
			Id: d.Get("check_host_id").(int),
		},
		CheckGroup: CheckGroup{
			Id: d.Get("check_group_id").(int),
		},
	}
}

func mapDNSCheckSchema(check DNSCheck, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(check.Id))
	d.Set("name", check.Name)
	d.Set("description", check.Description)
	d.Set("enabled", check.Enabled)
	d.Set("maintenance_override", check.MaintenanceOverride)
	d.Set("hostname", check.Hostname)
	d.Set("expected_addresses", check.ExpectedAddresses)
	d.Set("trigger_count", check.TriggerCount)
	d.Set("result_retention", check.ResultRetentionDays)
	d.Set("check_host_id", check.CheckHost.Id)
	d.Set("check_group_id", check.CheckGroup.Id)
}
