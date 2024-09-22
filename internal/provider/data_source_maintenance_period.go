package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the desired interfaces.
func NewMaintenancePeriodDataSource() datasource.DataSource {
	return &MaintenancePeriodDataSource{}
}

var _ datasource.DataSource = &MaintenancePeriodDataSource{}

type MaintenancePeriodDataSource struct {
	client *EndPointMonitorClient
}

func (d *MaintenancePeriodDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_maintenance_period"
}

// Schema defines the schema for the data source.
func (d *MaintenancePeriodDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Search for an individual Scheduled Maintenance Period. This will only allow a single result to be returned.",
		Attributes: map[string]schema.Attribute{
			"search": schema.StringAttribute{
				Required: true,
			},
			"id": schema.Int32Attribute{
				Computed: true,
			},
		},
	}
}

func (d *MaintenancePeriodDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GenericSingleDataSource

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	ids, err := d.client.SearchMaintenancePeriods(data.Search.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error searching maintenance periods",
			"Could not search maintenance periods, unexpected error: "+err.Error(),
		)
		return
	}

	if len(ids) != 1 {
		resp.Diagnostics.AddError(
			"None or more than one matching maintenance period found",
			"None or more than one matching maintenance period found when searching for single id",
		)
		return
	}

	data.Id = ids[0]

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MaintenancePeriodDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*EndPointMonitorClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *EndPointMonitorClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}
