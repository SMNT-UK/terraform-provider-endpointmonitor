package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &DashboardGroupResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewDashboardGroupResource() resource.Resource {
	return &DashboardGroupResource{}
}

// orderResource is the resource implementation.
type DashboardGroupResource struct {
	client *EndPointMonitorClient
}

// Metadata returns the resource type name.
func (r *DashboardGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dashboard_group"
}

// Schema defines the schema for the resource.
func (r *DashboardGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Dashboard Groups, the top-level organisational groups of checks running in EndPoint Monitor.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Dashboard Group. This will be used in alerts and notifications.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Required:    true,
				Description: "Space to provide a longer description of this Dashboard Group.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *DashboardGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DashboardGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dashboardGroup, error := r.client.CreateDashboardGroup(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating dashboard group",
			"Could not create dashboard group, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	plan.Id = dashboardGroup.Id

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *DashboardGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DashboardGroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed data from EPM
	dashboardGroup, err := r.client.GetDashboardGroup(state.Id.ValueInt32())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Fetching Dashboard Group",
			"Could not read dashboard group by id "+strconv.Itoa(int(state.Id.ValueInt32()))+": "+err.Error(),
		)
		return
	}

	if dashboardGroup == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set state from returned data from EPM.
	state = *dashboardGroup

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *DashboardGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DashboardGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dashboardGroup, error := r.client.UpdateDashboardGroup(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating dashboard group",
			"Could not create dashboard group, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	plan.Id = dashboardGroup.Id

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *DashboardGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan DashboardGroupModel
	req.State.Get(ctx, &plan)
	err := r.client.DeleteDashboardGroup(plan.Id.ValueInt32())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing dashboard group",
			"Request to EPM to delete dashboard group returned an error: "+err.Error(),
		)
		return
	}
}

func (r *DashboardGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DashboardGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, _ := strconv.Atoi(req.ID)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
