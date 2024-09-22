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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &CheckHostResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewCheckHostResource() resource.Resource {
	return &CheckHostResource{}
}

// orderResource is the resource implementation.
type CheckHostResource struct {
	client *EndPointMonitorClient
}

// Metadata returns the resource type name.
func (r *CheckHostResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_host"
}

// Schema defines the schema for the resource.
func (r *CheckHostResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage the hosts that checks are to be run on.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"hostname": schema.StringAttribute{
				Required:    true,
				Description: "The hostname of the host. Must match what the host believes its own hostname is.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Required:    true,
				Description: "A place to provide more detail about the host if required.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Optional:    true,
				Description: "If disabled checks set to run against this host will be paused.",
				Default:     booldefault.StaticBool(true),
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Must be either CONTROLLER or AGENT. CONTROLLER is used for hosts that expose the Web GUI and required database access. AGENT is used for hosts that purely just run checks.",
				Validators: []validator.String{
					stringvalidator.OneOf("CONTROLLER", "AGENT"),
				},
			},
			"max_checks": schema.Int32Attribute{
				Computed:    true,
				Optional:    true,
				Description: "The maximum number of concurrent Web Journey checks the host can run. Default is 1.",
				Default:     int32default.StaticInt32(1),
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"send_check_files": schema.BoolAttribute{
				Computed:    true,
				Optional:    true,
				Description: "For agents only. Indicates if it is to send check files such as screenshots back to the controller through the controller API. Should be enabled if there isn't a common file share between agent and controllers.",
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *CheckHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CheckHostModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkHost, error := r.client.CreateCheckHost(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating check host",
			"Could not create check host, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	plan.Id = checkHost.Id

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *CheckHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CheckHostModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed check from EPM
	checkHost, err := r.client.GetCheckHost(state.Id.ValueInt32())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Fetching Check Host",
			"Could not read check host by id "+strconv.Itoa(int(state.Id.ValueInt32()))+": "+err.Error(),
		)
		return
	}

	if checkHost == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set state from returned data from EPM.
	state = *checkHost

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *CheckHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CheckHostModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkHost, error := r.client.UpdateCheckHost(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating check host",
			"Could not create check host, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	plan.Id = checkHost.Id

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *CheckHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan CheckHostModel
	req.State.Get(ctx, &plan)
	err := r.client.DeleteCheckHost(plan.Id.ValueInt32())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing check host",
			"Request to EPM to delete check host returned an error: "+err.Error(),
		)
		return
	}
}

func (r *CheckHostResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CheckHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, _ := strconv.Atoi(req.ID)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
