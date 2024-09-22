package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &endPointMonitorProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &endPointMonitorProvider{
			version: version,
		}
	}
}

// hashicupsProviderModel maps provider schema data to a Go type.
type endPointMonitorProviderModel struct {
	Url types.String `tfsdk:"url"`
	Key types.String `tfsdk:"key"`
}

// hashicupsProvider is the provider implementation.
type endPointMonitorProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *endPointMonitorProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "endpointmonitor"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *endPointMonitorProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Optional: true,
			},
			"key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *endPointMonitorProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring EndPointMonitor client")

	// Retrieve provider data from configuration
	var config endPointMonitorProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Url.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Unknown EndPointMonitor URL",
			"The provider cannot create the EndPointMonitor client as there is an unknown configuration value for the EndPointMonitor URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the EPM_URL environment variable.",
		)
	}

	if config.Key.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("key"),
			"Unknown EndPointMonitor API Key",
			"The provider cannot create the EndPointMonitor client as there is an unknown configuration value for the EndPointMonitor API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the EPM_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	url := os.Getenv("EPM_URL")
	key := os.Getenv("EPM_API_KEY")

	if !config.Url.IsNull() {
		url = config.Url.ValueString()
	}

	if !config.Key.IsNull() {
		key = config.Key.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing EndPointMonitor URL",
			"The provider cannot create the EndPointMonitor client as there is a missing or empty value for the HashiCups API host. "+
				"Set the host value in the configuration or use the EPM_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing EndPointMonitor API Key",
			"The provider cannot create the EndPointMonitor client as there is a missing or empty value for the EndPointMonitor API key. "+
				"Set the username value in the configuration or use the EPM_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "endpointmonitor_url", url)
	ctx = tflog.SetField(ctx, "endpointmonitor_key", key)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "endpointmonitor_key")

	tflog.Debug(ctx, "Creating HashiCups client")

	// Create a new HashiCups client using the configuration values
	client, err := NewEPMClient(url, &key)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create EndPointMonitor Client",
			"An unexpected error occurred when creating the EndPointMonitor client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Client Error: "+err.Error(),
		)
		return
	}

	// Make the HashiCups client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured EndPointMonitor client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *endPointMonitorProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCheckGroupDataSource,
		NewCheckGroupsDataSource,
		NewCheckHostDataSource,
		NewCheckHostsDataSource,
		NewCheckDataSource,
		NewChecksDataSource,
		NewDashboardGroupDataSource,
		NewDashboardGroupsDataSource,
		NewHostGroupDataSource,
		NewHostGroupsDataSource,
		NewMaintenancePeriodDataSource,
		NewMaintenancePeriodsDataSource,
		NewProxyHostDataSource,
		NewProxyHostsDataSource,
		NewWebJourneyCommonStepDataSource,
		NewWebJourneyCommonStepsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *endPointMonitorProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUrlCheckResource,
		NewDnsCheckResource,
		NewCertificateCheckResource,
		NewPingCheckResource,
		NewSocketCheckResource,
		NewWebJourneyCheckResource,
		NewWebJourneyCommonStepResource,
		NewCheckGroupResource,
		NewCheckHostResource,
		NewDashboardGroupResource,
		NewHostGroupResource,
		NewProxyHostResource,
		NewMaintenancePeriodResource,
	}
}
