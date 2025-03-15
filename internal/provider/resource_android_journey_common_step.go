package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &AndroidJourneyCommonStepResource{}
)

func NewAndroidJourneyCommonStepResource() resource.Resource {
	return &AndroidJourneyCommonStepResource{}
}

type AndroidJourneyCommonStepResource struct {
	client *EndPointMonitorClient
}

// Metadata returns the resource type name.
func (r *AndroidJourneyCommonStepResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_android_journey_common_step"
}

// Schema defines the schema for the resource.
func (r *AndroidJourneyCommonStepResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Defines a shared complete step of an Android Journey, starting with the checks to perform on what is currently displayed, followed by the actions to take.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "A name to describe what the step is doing. This will be included in any alerts and notifications.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Optional longer description space to provide any supporting information about this step if needed.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"wait_time": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The number of milliseconds to wait for any loading / actions on the page to complete before any checks on this step are started.",
				Default:     int32default.StaticInt32(5000),
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"step_check": schema.ListNestedBlock{
				Description: "Defines the checks performed as part of an Android Journey Step to validate the currently displayed content of an app.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"description": schema.StringAttribute{
							Required:    true,
							Description: "A description to describe what the step_check is checking.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"warning_only": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "If true then if this check fails, then it will only produce a warning, not a full check failure. Default is false.",
							Default:     booldefault.StaticBool(false),
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "The type of check to complete. CHECK_FOR_TEXT - Checks for specific text being shown. CHECK_FOR_ELEMENT - Checks for a specific app component being shown.",
							Validators: []validator.String{
								stringvalidator.OneOf("CHECK_FOR_TEXT", "CHECK_FOR_ELEMENT"),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"check_for_text": schema.SingleNestedBlock{
							Description: "Defines the attributes needed for performing a Check for Text check as part of an Android Journey check.",
							Validators: []validator.Object{
								objectvalidator.AlsoRequires(
									path.MatchRelative().AtName("text_to_find"),
								),
							},
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Computed: true,
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.UseStateForUnknown(),
									},
								},
								"text_to_find": schema.StringAttribute{
									Optional:    true,
									Description: "The text to search for on the currently displayed Android window.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"state": schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: "Must be PRESENT or ABSENT. Defines if the textToFind should be found or not found ot the displayed Android window. PRESENT means the check will fail if the textToFind is not found. ABSENT means the check will fail of the textToFind is found.",
									Default:     stringdefault.StaticString("PRESENT"),
									Validators: []validator.String{
										stringvalidator.OneOf("PRESENT", "ABSENT"),
									},
								},
							},
						},
						"check_for_element": schema.SingleNestedBlock{
							Description: "Defines the attributes needed for performing a Check for Element check as part of an Android Journey check.",
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Computed: true,
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.UseStateForUnknown(),
									},
								},
								"component_id": schema.StringAttribute{
									Optional:    true,
									Description: "The component id of the element to look for (this does not need the android package prefix). Either this or xpath are required.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"component_type": schema.StringAttribute{
									Optional:    true,
									Description: "The component type to filter any matching elements for the xpath or component id by.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"xpath": schema.StringAttribute{
									Optional:    true,
									Description: "The xpath to use to search for the target element. Either this or componentId are required.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"state": schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: "Must be either PRESENT or ABSENT. PRESENT means if the element is not found, the check will fail. ABSENT means the element is found, the check will fail.",
									Default:     stringdefault.StaticString("PRESENT"),
									Validators: []validator.String{
										stringvalidator.OneOf("PRESENT", "ABSENT"),
									},
								},
								"attribute_name": schema.StringAttribute{
									Optional:    true,
									Description: "Optional check for testing if the found element has an attribute named by this value.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"attribute_value": schema.StringAttribute{
									Optional:    true,
									Description: "Optional check for testing if the found attributeName attribute has the value defined here.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
							},
						},
					},
				},
			},
			"step_interaction": schema.ListNestedBlock{
				Description: "Defines an interaction to make ar part of an Android Journey check.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"sequence": schema.Int32Attribute{
							Required:    true,
							Description: "The order in which to run each interaction, working in lowest number to highest.",
							Validators: []validator.Int32{
								int32validator.AtLeast(1),
							},
						},
						"description": schema.StringAttribute{
							Required:    true,
							Description: "A description to describe the action being taken. This is used as parts of alerts and reporting.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"always_required": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "If set to false, the action is deemed optional, so the check will attempt to perform it, but if it fails, the check will continue as normal.",
							Default:     booldefault.StaticBool(true),
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "The type of action to perform. Options are: CLICK, INPUT_TEXT, INPUT_PASSWORD, SAVE_SCREEN_SOURCE, ROTATE_DISPLAY, SCROLL_TO_ELEMENT, SELECT_SPINNER_OPTION, SWIPE, SCREENSHOT or WAIT.",
							Validators: []validator.String{
								stringvalidator.OneOf("CLICK", "INPUT_TEXT", "INPUT_PASSWORD", "SAVE_SCREEN_SOURCE", "ROTATE_DISPLAY", "SCROLL_TO_ELEMENT", "SELECT_SPINNER_OPTION", "SWIPE", "SCREENSHOT", "WAIT"),
							},
						},
						"wait_time": schema.Int32Attribute{
							Optional:    true,
							Description: "The number of milliseconds to wait for the WAIT interaction type.",
							Validators: []validator.Int32{
								int32validator.AtLeast(1),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"click": schema.SingleNestedBlock{
							Description: "The attributes required as part of performing a CLICK interaction during an Android Journey check. Only one attribute needs to be provided.",
							Attributes: map[string]schema.Attribute{
								"component_id": schema.StringAttribute{
									Optional:    true,
									Description: "The id of the component to click on. The id does not need to include the Android package name prefix.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"xpath": schema.StringAttribute{
									Optional:    true,
									Description: "Xpath defining the component to click on.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"search_text": schema.StringAttribute{
									Optional:    true,
									Description: "The text to search for and click on within the currently displayed app screen.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
							},
						},
						"text_input": schema.SingleNestedBlock{
							Description: "The attributes required as part of performing a INPUT_TEXT interaction during an Android Journey check.",
							Validators: []validator.Object{
								objectvalidator.AlsoRequires(
									path.MatchRelative().AtName("input_text"),
								),
							},
							Attributes: map[string]schema.Attribute{
								"component_id": schema.StringAttribute{
									Optional:    true,
									Description: "The id of the component to input the text into. This does not need to include the Android package name prefix.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"xpath": schema.StringAttribute{
									Optional:    true,
									Description: "The xpath of the component to input the text into. Either this or elementId should be given, but not both.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"input_text": schema.StringAttribute{
									Optional:    true,
									Description: "The text to input into the element defined by either component_id or xpath.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
							},
						},
						"password_input": schema.SingleNestedBlock{
							Description: "The attributes required as part of performing a INPUT_PASSWORD interaction during an Android Journey check.",
							Validators: []validator.Object{
								objectvalidator.AlsoRequires(
									path.MatchRelative().AtName("input_password"),
								),
							},
							Attributes: map[string]schema.Attribute{
								"component_id": schema.StringAttribute{
									Optional:    true,
									Description: "The id of the component to input the password into. This does not need to include the Android package name prefix.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"xpath": schema.StringAttribute{
									Optional:    true,
									Description: "The xpath of the component to input the password into. Either this or elementId should be given, but not both.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"input_password": schema.StringAttribute{
									Optional:    true,
									Description: "The password to input. This will not be stored in your Terraform state and ideally should be passed in to your Terraform as a environment variable rather than statically stored in your Terraform code.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
						"rotate_display": schema.SingleNestedBlock{
							Description: "The attributes required as part of performing a ROTATE_DISPLAY interaction during an Android Journey check.",
							Validators: []validator.Object{
								objectvalidator.AlsoRequires(
									path.MatchRelative().AtName("orientation"),
								),
							},
							Attributes: map[string]schema.Attribute{
								"orientation": schema.StringAttribute{
									Optional:    true,
									Description: "The orientation to rotate the screeen to, either PORTRAIT or LANDSCAPE.",
									Validators: []validator.String{
										stringvalidator.OneOf("PORTRAIT", "LANDSCAPE"),
									},
								},
							},
						},
						"select_spinner_option": schema.SingleNestedBlock{
							Description: "The attributes required as part of performing a SELECT_SPINNER_OPTION interaction during an Android Journey check. Only one attribute needs to be provided.",
							Attributes: map[string]schema.Attribute{
								"component_id": schema.StringAttribute{
									Optional:    true,
									Description: "The id of the spinner object to make the selection in. The id does not need to include the Android package name prefix.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"xpath": schema.StringAttribute{
									Optional:    true,
									Description: "Xpath defining the spinner object to make the selection in.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"search_text": schema.StringAttribute{
									Optional:    true,
									Description: "The text of the current selected spinner object value to search for to identify the spinner to interact with.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"option_list_position": schema.Int32Attribute{
									Optional:    true,
									Description: "The position from the list of options within the spinner to select, starting from 0.",
									Validators: []validator.Int32{
										int32validator.AtLeast(0),
									},
								},
								"option_list_text": schema.StringAttribute{
									Optional:    true,
									Description: "The value from the list of options within the spinner to select.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
							},
						},
						"swipe": schema.SingleNestedBlock{
							Description: "The attributes required as part of performing a SWIPE interaction during an Android Journey check.",
							Validators: []validator.Object{
								objectvalidator.AlsoRequires(
									path.MatchRelative().AtName("swipe_direction"),
									path.MatchRelative().AtName("swipe_length"),
								),
							},
							Attributes: map[string]schema.Attribute{
								"component_id": schema.StringAttribute{
									Optional:    true,
									Description: "The id of the component start the swipe action within. This does not need to include the Android package name prefix.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"xpath": schema.StringAttribute{
									Optional:    true,
									Description: "The xpath of the component to start the wipe action within. Either this or elementId should be given, but not both.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"start_swipe_coordinates": schema.StringAttribute{
									Optional:    true,
									Description: "The x,y position in pixels on the screen to star the swipe from if component_id or xpath not used.",
									Validators: []validator.String{
										stringvalidator.RegexMatches(regexp.MustCompile(`^\d+,\d+$`), "coordinates must be given in format xx,yy"),
									},
								},
								"swipe_direction": schema.StringAttribute{
									Optional:    true,
									Description: "The direction to swipe across the screen. Must be one of LEFT, RIGHT, UP or DOWN.",
									Validators: []validator.String{
										stringvalidator.OneOf("LEFT", "RIGHT", "UP", "DOWN"),
									},
								},
								"swipe_length": schema.Int32Attribute{
									Optional:    true,
									Computed:    true,
									Description: "The distance across the screen to swipe in pixels.",
									Default:     int32default.StaticInt32(200),
									Validators: []validator.Int32{
										int32validator.AtLeast(10),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *AndroidJourneyCommonStepResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AndroidJourneyCommonStepModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	step, error := r.client.CreateAndroidJourneyCommonStep(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating Common Step",
			"Could not create Common Step, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	// Because the response we get from the EPM API doesn't contain the any password_input values,
	// to be able to just copy its response to the state, we need to grab the passwords from the
	// plan and put them in the returned state. Easier to do than it is to go through and
	// manually populate all the ids.
	for x := 0; x < len(step.StepInteractions); x++ {
		action := step.StepInteractions[x]

		if action.PasswordInput != nil {
			action.PasswordInput.InputPassword = plan.StepInteractions[x].PasswordInput.InputPassword
		}
	}

	plan = *step

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *AndroidJourneyCommonStepResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AndroidJourneyCommonStepModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed check from EPM
	commonStep, err := r.client.GetCommonAndroidJourneyStep(state.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching Common Step",
			"Could not read Common Step by id "+strconv.Itoa(int(state.Id.ValueInt64()))+": "+err.Error(),
		)
		return
	}

	if commonStep == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set state from returned data from EPM.
	// Because the response we get from the EPM API doesn't contain the any password_input values,
	// to be able to just copy its response to the state, we need to grab the passwords from the
	// current state and put them in the returned state. Easier to do than it is to go through and
	// manually populate all the ids.
	passwordActions := make([]*AndroidStepInteractionModel, 0)

	// Get all passwords currently in state, so we use their ids to know what to match against the
	// returned passwords from the EPM API.
	for _, action := range state.StepInteractions {
		if action.PasswordInput != nil {
			foundAction := action // Copy as `action` pointer will change through next loop.
			passwordActions = append(passwordActions, &foundAction)
		}
	}

	for _, action := range commonStep.StepInteractions {
		if action.PasswordInput != nil {
			// Does this exist in state? If so, we'll overwrite with the state password.
			for _, passwordAction := range passwordActions {
				if passwordAction.Id.Equal(action.Id) {
					action.PasswordInput.InputPassword = passwordAction.PasswordInput.InputPassword
				}
			}
		}
	}

	state = *commonStep

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *AndroidJourneyCommonStepResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AndroidJourneyCommonStepModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	step, error := r.client.UpdateAndroidJourneyCommonStep(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating Common Step",
			"Could not create Common Step, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	// Because the response we get from the EPM API doesn't contain the any password_input values,
	// to be able to just copy its response to the state, we need to grab the passwords from the
	// plan and put them in the returned state. Easier to do than it is to go through and
	// manually populate all the ids.
	for x := 0; x < len(step.StepInteractions); x++ {
		action := step.StepInteractions[x]

		if action.PasswordInput != nil {
			action.PasswordInput.InputPassword = step.StepInteractions[x].PasswordInput.InputPassword
		}
	}

	plan = *step

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *AndroidJourneyCommonStepResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan AndroidJourneyCommonStepModel
	req.State.Get(ctx, &plan)
	err := r.client.DeleteAndroidCommonStep(plan.Id.ValueInt64())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing Common Step",
			"Request to EPM to delete Common Step returned an error: "+err.Error(),
		)
		return
	}
}

func (r *AndroidJourneyCommonStepResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AndroidJourneyCommonStepResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, _ := strconv.Atoi(req.ID)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
