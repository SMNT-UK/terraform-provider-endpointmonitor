package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the desired interfaces.
func NewWebJourneyCommonStepsDataSource() datasource.DataSource {
	return &WebJourneyCommonStepsDataSource{}
}

var _ datasource.DataSource = &WebJourneyCommonStepsDataSource{}

type WebJourneyCommonStepsDataSource struct {
	client *EndPointMonitorClient
}

func (d *WebJourneyCommonStepsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_web_journey_common_steps"
}

// Schema defines the schema for the data source.
func (d *WebJourneyCommonStepsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Search for multiple Common Web Journey Steps. A list of ids will be returned for all matches found.",
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

func (d *WebJourneyCommonStepsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GenericMultipleDataSource

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	ids, err := d.client.SearchWebJoureyCommonSteps(data.Search.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error searching web journey common steps",
			"Could not search web journey common steps, unexpected error: "+err.Error(),
		)
		return
	}

	data.Ids = ids

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WebJourneyCommonStepsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
