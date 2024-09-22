package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the desired interfaces.
func NewDashboardGroupsDataSource() datasource.DataSource {
	return &DashboardGroupsDataSource{}
}

var _ datasource.DataSource = &DashboardGroupsDataSource{}

type DashboardGroupsDataSource struct {
	client *EndPointMonitorClient
}

func (d *DashboardGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dashboard_groups"
}

// Schema defines the schema for the data source.
func (d *DashboardGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Search for multiple Dashboard Groups. A list of ids will be returned for all matches found.",
		Attributes: map[string]schema.Attribute{
			"search": schema.StringAttribute{
				Required: true,
			},
			"ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.Int32Type,
			},
		},
	}
}

func (d *DashboardGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GenericMultipleDataSource

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	ids, err := d.client.SearchDashboardGroups(data.Search.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error searching dsahboard groups",
			"Could not search dashboard groups, unexpected error: "+err.Error(),
		)
		return
	}

	if len(ids) < 1 {
		resp.Diagnostics.AddError(
			"No matching dashboard groups found",
			"No matching dashboard groups found in EPM",
		)
		return
	}

	data.Ids = ids

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DashboardGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
