package endpointmonitor

import (
	"context"
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func checkHost() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage the hosts that checks are to be run on.",
		CreateContext: resourceCheckHostCreate,
		ReadContext:   resourceCheckHostRead,
		UpdateContext: resourceCheckHostUpdate,
		DeleteContext: resourceCheckHostDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "A friendly name to describe the host. This is what the host will be refered to as on all screens and alerts.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A place to provide more detail about the host if required.",
				Required:    true,
			},
			"hostname": {
				Type:        schema.TypeString,
				Description: "The hostname of the host. Must match what the host believes it's own hostname is.",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "If disabled checks set to run against this host will be paused.",
				Required:    true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Must be either CONTROLLER or AGENT. CONTROLLER is used for hosts that expose the Web GUI and required database access. AGENT is used for hosts that purely just run checks.",
				ValidateFunc: validateCheckHostType(),
			},
			"max_checks": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The maximum number of concurrent Web Journey checks the host can run. Default is 1.",
				Default:     1,
			},
			"send_check_files": {
				Type:        schema.TypeBool,
				Description: "For agents only. Indicates if it is to send check files such as screenshots back to the controller through the controller API. Should be enabled if there isn't a common file share between agent and controllers.",
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCheckHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkHostId := d.Id()

	checkHost, err := c.GetCheckHost(checkHostId)
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.IsNewResource() && checkHost == nil {
		d.SetId("")
		return nil
	}

	mapCheckHostSchema(*checkHost, d)

	return diags
}

func resourceCheckHostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkHost := mapCheckHost(d)

	if *checkHost.Type != "CONTROLLER" && *checkHost.Type != "AGENT" {
		err := errors.New("Check Host type must be CONTROLLER or AGENT")
		return diag.FromErr(err)
	}

	o, err := c.CreateCheckHost(checkHost)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(o.Id)))

	resourceCheckHostRead(ctx, d, m)

	return diags
}

func resourceCheckHostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	_, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChangesExcept() {
		checkHost := mapCheckHost(d)

		_, err := c.UpdateCheckHost(checkHost)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceCheckHostRead(ctx, d, m)
}

func resourceCheckHostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteCheckHost(checkId)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func mapCheckHost(d *schema.ResourceData) CheckHost {
	checkId, err := strconv.Atoi(d.Id())
	if err != nil {
		checkId = 0
	}

	hostType := d.Get("type").(string)

	return CheckHost{
		Id:                  checkId,
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Hostname:            d.Get("hostname").(string),
		Type:                &hostType,
		Enabled:             d.Get("enabled").(bool),
		MaxWebJourneyChecks: d.Get("max_checks").(int),
		SendCheckFiles:      d.Get("send_check_files").(bool),
	}
}

func mapCheckHostSchema(checkHost CheckHost, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(checkHost.Id))
	d.Set("name", checkHost.Name)
	d.Set("description", checkHost.Description)
	d.Set("hostname", checkHost.Hostname)
	d.Set("type", checkHost.Type)
	d.Set("enabled", checkHost.Enabled)
	d.Set("max_checks", checkHost.MaxWebJourneyChecks)
	d.Set("send_check_files", checkHost.SendCheckFiles)
}

func validateCheckHostType() schema.SchemaValidateFunc {
	types := []string{
		"CONTROLLER",
		"AGENT",
	}
	return validation.StringInSlice(types, false)
}
