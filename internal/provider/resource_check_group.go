package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
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
	_ resource.Resource = &CheckGroupResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewCheckGroupResource() resource.Resource {
	return &CheckGroupResource{}
}

// orderResource is the resource implementation.
type CheckGroupResource struct {
	client *EndPointMonitorClient
}

// Metadata returns the resource type name.
func (r *CheckGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_group"
}

// Schema defines the schema for the resource.
func (r *CheckGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage Check Groups, the initial grouping of checks running on EndPoint Monitor.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "A meaningful name of what this group contains. This will be used in alerts and notifications.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Required:    true,
				Description: "A space to provide a longer description of this group.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"dashboard_group_id": schema.Int32Attribute{
				Required:    true,
				Description: "The id of the Dashboard Group this belongs to.",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *CheckGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CheckGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkGroup, error := r.client.CreateCheckGroup(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating check group",
			"Could not create check group, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	plan.Id = checkGroup.Id

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *CheckGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CheckGroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed check from EPM
	checkGroup, err := r.client.GetCheckGroup(state.Id.ValueInt32())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Fetching Check Group",
			"Could not read check group by id "+strconv.Itoa(int(state.Id.ValueInt32()))+": "+err.Error(),
		)
		return
	}

	if checkGroup == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set state from returned data from EPM.
	state = *checkGroup

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *CheckGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CheckGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkGroup, error := r.client.UpdateCheckGroup(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating check group",
			"Could not create check group, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	plan.Id = checkGroup.Id

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *CheckGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan CheckGroupModel
	req.State.Get(ctx, &plan)
	err := r.client.DeleteCheckGroup(plan.Id.ValueInt32())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing check group",
			"Request to EPM to delete check group returned an error: "+err.Error(),
		)
		return
	}
}

func (r *CheckGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CheckGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, _ := strconv.Atoi(req.ID)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
