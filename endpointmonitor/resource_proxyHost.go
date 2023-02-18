package endpointmonitor

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func proxyHost() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage HTTP proxies that can be used for URL and Web Journey checks if a proxy is required to access the target.",
		CreateContext: resourceProxyHostCreate,
		ReadContext:   resourceProxyHostRead,
		UpdateContext: resourceProxyHostUpdate,
		DeleteContext: resourceProxyHostDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "A name to reference the Proxy Host in other areas of application.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Space for a longer description for the Proxy Host if needed.",
				Required:    true,
			},
			"hostname": {
				Type:        schema.TypeString,
				Description: "The FQDN of the Proxy Host.",
				Required:    true,
			},
			"port": {
				Type:        schema.TypeInt,
				Description: "The port the proxy is listening on.",
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceProxyHostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	proxyHost := ProxyHost{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Hostname:    d.Get("hostname").(string),
		Port:        d.Get("port").(int),
	}

	o, err := c.CreateProxyHost(proxyHost)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(o.Id)))

	resourceProxyHostRead(ctx, d, m)

	return diags
}

func resourceProxyHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	proxyId := d.Id()

	proxyHost, err := c.GetProxyHost(proxyId)
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.IsNewResource() && proxyHost == nil {
		d.SetId("")
		return nil
	}

	d.SetId(strconv.Itoa(proxyHost.Id))
	d.Set("name", proxyHost.Name)
	d.Set("description", proxyHost.Description)
	d.Set("hostname", proxyHost.Hostname)
	d.Set("port", proxyHost.Port)

	return diags
}

func resourceProxyHostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	proxyId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("name", "description", "hostname", "port") {
		proxyHost := ProxyHost{
			Id:          proxyId,
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Hostname:    d.Get("hostname").(string),
			Port:        d.Get("port").(int),
		}

		_, err := c.UpdateProxyHost(proxyHost)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceProxyHostRead(ctx, d, m)
}

func resourceProxyHostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	proxyId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteProxyHost(proxyId)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
