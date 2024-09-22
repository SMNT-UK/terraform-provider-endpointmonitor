package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &HostGroupResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewHostGroupResource() resource.Resource {
	return &HostGroupResource{}
}

// orderResource is the resource implementation.
type HostGroupResource struct {
	client *EndPointMonitorClient
}

// Metadata returns the resource type name.
func (r *HostGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_host_group"
}

// Schema defines the schema for the resource.
func (r *HostGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage groups of Check Hosts that can assigned to checks to execute on.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "A name to identify this group of hosts by. This will show in searches and alerts.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Space for a longer description to define this group of hosts by. Not required.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Enable or disable checks assigned to this host group.",
				Default:     booldefault.StaticBool(true),
			},
			"check_host_ids": schema.SetAttribute{
				Required:    true,
				Description: "List of Check Host id's that are to be part of this Host Group.",
				ElementType: types.Int32Type,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *HostGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan HostGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hostGroup, error := r.client.CreateHostGroup(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating host group",
			"Could not create host group, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	plan.Id = hostGroup.Id

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *HostGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state HostGroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed data from EPM
	hostGroup, err := r.client.GetHostGroup(state.Id.ValueInt32())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Fetching Host Group",
			"Could not read host group by id "+strconv.Itoa(int(state.Id.ValueInt32()))+": "+err.Error(),
		)
		return
	}

	if hostGroup == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set state from returned data from EPM.
	state = *hostGroup

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *HostGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan HostGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hostGroup, error := r.client.UpdateHostGroup(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating host group",
			"Could not create host group, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	plan.Id = hostGroup.Id

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *HostGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan HostGroupModel
	req.State.Get(ctx, &plan)
	err := r.client.DeleteHostGroup(plan.Id.ValueInt32())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing host group",
			"Request to EPM to delete host group returned an error: "+err.Error(),
		)
		return
	}
}

func (r *HostGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *HostGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, _ := strconv.Atoi(req.ID)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
