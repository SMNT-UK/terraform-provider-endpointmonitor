package endpointmonitor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Description: "The API base path of your EndPoint Monitor installation. This is usually the path you would normally access EndPoint Monitor through with /api appended. This can also be passed in through the environment variable EPM_URL.",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("EPM_URL", nil),
			},
			"key": {
				Type:        schema.TypeString,
				Description: "An API key issued from your EndPoint Monitor installation under the API Keys section. Make sure the API key used has write access. This should be passed in using an environment variable with name EPM_API_KEY. Do not store this key in any configuration.",
				Sensitive:   true,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("EPM_API_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"endpointmonitor_proxy_host":         proxyHost(),
			"endpointmonitor_check_host":         checkHost(),
			"endpointmonitor_app_group":          appGroup(),
			"endpointmonitor_check_group":        checkGroup(),
			"endpointmonitor_url_check":          urlCheck(),
			"endpointmonitor_socket_check":       socketCheck(),
			"endpointmonitor_dns_check":          dnsCheck(),
			"endpointmonitor_ping_check":         pingCheck(),
			"endpointmonitor_web_journey_check":  webJourneyCheck(),
			"endpointmonitor_maintenance_period": maintenancePeriod(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"endpointmonitor_proxy_host":               dataSourceProxyHost(),
			"endpointmonitor_proxy_hosts":              dataSourceProxyHosts(),
			"endpointmonitor_app_group":                dataSourceAppGroup(),
			"endpointmonitor_app_groups":               dataSourceAppGroups(),
			"endpointmonitor_check_group":              dataSourceCheckGroup(),
			"endpointmonitor_check_groups":             dataSourceCheckGroups(),
			"endpointmonitor_check_host":               dataSourceCheckHost(),
			"endpointmonitor_check_hosts":              dataSourceCheckHosts(),
			"endpointmonitor_web_journey_common_step":  dataSourceWebJourneyCommonStep(),
			"endpointmonitor_web_journey_common_steps": dataSourceWebJourneyCommonSteps(),
			"endpointmonitor_maintenance_period":       dataSourceMaintenancePeriod(),
			"endpointmonitor_maintenance_periods":      dataSourceMaintenancePeriods(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("key").(string)

	var url *string

	hVal, ok := d.GetOk("url")
	if ok {
		tempUrl := hVal.(string)
		url = &tempUrl
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if apiKey != "" {
		c, err := NewClient(*url, &apiKey)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create EndPoint Monitor client",
				Detail:   "Unable to authenticate user for authenticated EndPoint Monitor client",
			})

			return nil, diags
		}

		return c, diags
	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Unable to create EndPoint Monitor client",
		Detail:   "Unable to create anonymous EndPoint Monitor client",
	})

	return nil, diags
}
