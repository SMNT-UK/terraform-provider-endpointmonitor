package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &MaintenancePeriodResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewMaintenancePeriodResource() resource.Resource {
	return &MaintenancePeriodResource{}
}

// orderResource is the resource implementation.
type MaintenancePeriodResource struct {
	client *EndPointMonitorClient
}

// Metadata returns the resource type name.
func (r *MaintenancePeriodResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_maintenance_period"
}

// Schema defines the schema for the resource.
func (r *MaintenancePeriodResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage scheduled maintenance periods to prevent checks from alerting during certain periods of the day or week.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Required:    true,
				Description: "Space for a description of the maintenance periods purpose.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: "Enable or disable the maintenance period from applying to attached checks.",
			},
			"day_of_week": schema.StringAttribute{
				Required:    true,
				Description: "The day of week the maintenance period applies to. Set as ALL for every day of the week. Must otherwise be SUNDAY, MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY or SATURDAY.",
				Validators: []validator.String{
					stringvalidator.OneOf("SUNDAY", "MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY", "ALL"),
				},
			},
			"start_time": schema.StringAttribute{
				Required:    true,
				Description: "The start time of the maintenance period in format 24HH:MM.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^[0-2]{1}[0-9]{1}:[0-5]{1}[0-9]{1}$"), "time must be in format 24HH:MM"),
				},
			},
			"end_time": schema.StringAttribute{
				Required:    true,
				Description: "The end of time the maintenance period in format 24HH:MM.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^[0-2]{1}[0-9]{1}:[0-5]{1}[0-9]{1}$"), "time must be in format 24HH:MM"),
				},
			},
			"check_ids": schema.ListAttribute{
				Optional:    true,
				Description: "A list of ids of Checks that are directly linked to the maintenance period.",
				ElementType: types.Int32Type,
			},
			"check_group_ids": schema.ListAttribute{
				Optional:    true,
				Description: "A list of ids of Check Groups that are directly linked to the maintenance period.",
				ElementType: types.Int32Type,
			},
			"dashboard_group_ids": schema.ListAttribute{
				Optional:    true,
				Description: "A list of ids of Dashboard Groups that are linked to this maintenance period.",
				ElementType: types.Int32Type,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *MaintenancePeriodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MaintenancePeriodModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	maintenancePeriod, error := r.client.CreateMaintenancePeriod(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating maintenance period",
			"Could not create maintenance period, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	plan.Id = maintenancePeriod.Id

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *MaintenancePeriodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MaintenancePeriodModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed data from EPM
	maintenancePeriod, err := r.client.GetMaintenancePeriod(state.Id.ValueInt32())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Fetching Maintenance Period",
			"Could not read maintenance period by id "+strconv.Itoa(int(state.Id.ValueInt32()))+": "+err.Error(),
		)
		return
	}

	if maintenancePeriod == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set state from returned data from EPM.
	state = *maintenancePeriod

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *MaintenancePeriodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MaintenancePeriodModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	maintenancePeriod, error := r.client.UpdateMaintenancePeriod(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating maintenance period",
			"Could not create maintenance period, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	plan.Id = maintenancePeriod.Id

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *MaintenancePeriodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan MaintenancePeriodModel
	req.State.Get(ctx, &plan)
	err := r.client.DeleteMaintenancePeriod(plan.Id.ValueInt32())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing maintenance period",
			"Request to EPM to delete maintenance period returned an error: "+err.Error(),
		)
		return
	}
}

func (r *MaintenancePeriodResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MaintenancePeriodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, _ := strconv.Atoi(req.ID)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
