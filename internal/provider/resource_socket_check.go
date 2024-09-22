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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &SocketCheckResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewSocketCheckResource() resource.Resource {
	return &SocketCheckResource{}
}

// orderResource is the resource implementation.
type SocketCheckResource struct {
	client *EndPointMonitorClient
}

// Metadata returns the resource type name.
func (r *SocketCheckResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_socket_check"
}

// Schema defines the schema for the resource.
func (r *SocketCheckResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage socket checks which test to ensure a hostname is listening on a pre-defined port.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "A name to describe in the check, used throughout EndPoint Monitor to describe this check, including in notifications.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(3),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A space to provide a longer description of the check if needed. Will default to the name if not set.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Allows the enabling/disabling of the check from executing.",
			},
			"check_frequency": schema.Int32Attribute{
				Required:    true,
				Description: "The frequency the check will be run in seconds.",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"maintenance_override": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "If set true then notifications and alerts will be suppressed for the check.",
			},
			"trigger_count": schema.Int32Attribute{
				Required:    true,
				Description: "The sequential number of failures that need to occur for a check to trigger an alert or notification.",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"result_retention": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int32default.StaticInt32(366),
				Description: "The number of days to store historic results of the check.",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"check_host_id": schema.Int32Attribute{
				Optional:    true,
				Description: "The id of the Check Host to run the check on.",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"check_host_group_id": schema.Int32Attribute{
				Optional:    true,
				Description: "The id of the Check Host Group to run the check on.",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"check_group_id": schema.Int32Attribute{
				Required:    true,
				Description: "The id of the Check Group the check belongs to. This also determines check frequency.",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"proxy_host_id": schema.Int32Attribute{
				Optional:    true,
				Description: "The id of the Proxy Host the check should use for a HTTP proxy if needed.",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"hostname": schema.StringAttribute{
				Required:    true,
				Description: "The hostname to check.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"port": schema.Int32Attribute{
				Required:    true,
				Description: "The TCP port to check is listening.",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *SocketCheckResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SocketCheckModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	check, error := r.client.CreateSocketCheck(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating check",
			"Could not create check, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	plan.Id = check.Id

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *SocketCheckResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SocketCheckModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed check from EPM
	check, err := r.client.GetSocketCheck(state.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Fetching Check",
			"Could not read check by id "+strconv.Itoa(int(state.Id.ValueInt64()))+": "+err.Error(),
		)
		return
	}

	if check == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set state from returned data from EPM.
	state = *check

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *SocketCheckResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SocketCheckModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	check, error := r.client.UpdateSocketCheck(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating check",
			"Could not create check, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	plan.Id = check.Id

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *SocketCheckResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan SocketCheckModel
	req.State.Get(ctx, &plan)
	err := r.client.DeleteCheck(plan.Id.ValueInt64())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing check",
			"Request to EPM to delete check returned an error: "+err.Error(),
		)
		return
	}
}

func (r *SocketCheckResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SocketCheckResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, _ := strconv.Atoi(req.ID)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
