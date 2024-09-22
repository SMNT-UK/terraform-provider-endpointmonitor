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
	_ resource.Resource = &ProxyHostResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewProxyHostResource() resource.Resource {
	return &ProxyHostResource{}
}

// orderResource is the resource implementation.
type ProxyHostResource struct {
	client *EndPointMonitorClient
}

// Metadata returns the resource type name.
func (r *ProxyHostResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_proxy_host"
}

// Schema defines the schema for the resource.
func (r *ProxyHostResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage HTTP proxies that can be used for URL and Web Journey checks if a proxy is required to access the target.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "A name to reference the Proxy Host in other areas of application.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Required:    true,
				Description: "Space for a longer description for the Proxy Host if needed.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"hostname": schema.StringAttribute{
				Required:    true,
				Description: "The FQDN of the Proxy Host.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"port": schema.Int32Attribute{
				Required:    true,
				Description: "The port the proxy is listening on.",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *ProxyHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ProxyHostModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	proxyHost, error := r.client.CreateProxyHost(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating proxy host",
			"Could not create proxy host, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	plan.Id = proxyHost.Id

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *ProxyHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ProxyHostModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed data from EPM
	proxyHost, err := r.client.GetProxyHost(state.Id.ValueInt32())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Fetching Proxy Host",
			"Could not read proxy host by id "+strconv.Itoa(int(state.Id.ValueInt32()))+": "+err.Error(),
		)
		return
	}

	if proxyHost == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set state from returned data from EPM.
	state = *proxyHost

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ProxyHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ProxyHostModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	proxyHost, error := r.client.UpdateProxyHost(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating proxy host",
			"Could not create proxy host, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	plan.Id = proxyHost.Id

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ProxyHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan ProxyHostModel
	req.State.Get(ctx, &plan)
	err := r.client.DeleteProxyHost(plan.Id.ValueInt32())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing proxy host",
			"Request to EPM to delete proxy host returned an error: "+err.Error(),
		)
		return
	}
}

func (r *ProxyHostResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProxyHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, _ := strconv.Atoi(req.ID)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
