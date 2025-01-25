package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &WebJourneyCheckResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewWebJourneyCheckResource() resource.Resource {
	return &WebJourneyCheckResource{}
}

// orderResource is the resource implementation.
type WebJourneyCheckResource struct {
	client *EndPointMonitorClient
}

// Metadata returns the resource type name.
func (r *WebJourneyCheckResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_web_journey_check"
}

// Schema defines the schema for the resource.
func (r *WebJourneyCheckResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create and manage web journey checks that can be set up to navigate through a website and perform period checks to ensure page elements, network calls and console logs are there or not as expected.",
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
			"start_url": schema.StringAttribute{
				Required:    true,
				Description: "The URL to load start the journey at.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^https?://"), "url must start with http:// or https://"),
				},
			},
			"window_height": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The height of the browser window used for the check.",
				Default:     int32default.StaticInt32(1080),
				Validators: []validator.Int32{
					int32validator.AtLeast(600),
				},
			},
			"window_width": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The width of the browser window used for the check.",
				Default:     int32default.StaticInt32(1920),
				Validators: []validator.Int32{
					int32validator.AtLeast(800),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"monitor_domain": schema.ListNestedBlock{
				Description: "Define a domain to monitor network calls from during the check. If no monitor_domain's are defined, then all calls will be monitored.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"domain": schema.StringAttribute{
							Required:    true,
							Description: "The domain to monitor.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(3),
							},
						},
						"include_sub_domains": schema.BoolAttribute{
							Required:    true,
							Description: "If true, all sub-domains of the domain will be monitored too. If false, just the given domain will be monitored.",
						},
					},
				},
			},
			"step": schema.ListNestedBlock{
				Description: "Defines a complete step of a web journey, starting with the checks to perform on the current page, followed by actions to take.",
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
							Description: "This indicates the order in which the steps will executing during the check.",
							Validators: []validator.Int32{
								int32validator.AtLeast(0),
							},
						},
						"name": schema.StringAttribute{
							Required:    true,
							Description: "A name to describe what the step is doing. This will be included in any alerts and notifications.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "Should be COMMON or CUSTOM. COMMON allows use of a pre-defined Web Journey step, common_step_id must be set when using this option. CUSTOM allows a custom one to be defined for this check.",
							Validators: []validator.String{
								stringvalidator.OneOf("COMMON", "CUSTOM"),
							},
						},
						"common_step_id": schema.Int64Attribute{
							Optional:    true,
							Description: "If type is set to COMMON, then this should be set. The id of the Common Web Journey Step to use.",
							Validators: []validator.Int64{
								int64validator.AtLeast(1),
							},
						},
						"wait_time": schema.Int32Attribute{
							Optional:    true,
							Computed:    true,
							Description: "The number of milliseconds to wait for any page load / actions on the page to complete before any checks on this step are started.",
							Default:     int32default.StaticInt32(5000),
							Validators: []validator.Int32{
								int32validator.AtLeast(1),
							},
						},
						"page_load_time_warning": schema.Int32Attribute{
							Optional:    true,
							Computed:    true,
							Description: "The maximum number of milliseconds that any discovered network call can take before a warning is created for it and the check is set to a warning status.",
							Default:     int32default.StaticInt32(2500),
							Validators: []validator.Int32{
								int32validator.AtLeast(1),
							},
						},
						"page_load_time_alert": schema.Int32Attribute{
							Optional:    true,
							Computed:    true,
							Description: "The maximum number of milliseconds that any discovered network call can take before an alert is created for it, and the check is set to a failed status.",
							Default:     int32default.StaticInt32(5000),
							Validators: []validator.Int32{
								int32validator.AtLeast(1),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"page_check": schema.ListNestedBlock{
							Description: "The set of checks to run against the currently loaded content.",
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
										Description: "A description of what this check is doing. This will be used in alerts and notifications.",
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
										Description: "The type of check to execute. Options are: CHECK_FOR_TEXT - Check for any string on or not on the current page. CHECK_FOR_ELEMENT - Check for an element and it's properties on the current page. CHECK_CURRENT_URL - Check the current url. CHECK_URL_RESPONSE - Check for specific network calls made after the last step. CHECK_CONSOLE_LOG - Check for console logs made after the last step.",
										Validators: []validator.String{
											stringvalidator.OneOf("CHECK_FOR_TEXT", "CHECK_FOR_ELEMENT", "CHECK_CURRENT_URL", "CHECK_URL_RESPONSE", "CHECK_CONSOLE_LOG"),
										},
									},
								},
								Blocks: map[string]schema.Block{
									"check_for_text": schema.SingleNestedBlock{
										Description: "Check a specific stirng is present or absent on the current page.",
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
												Description: "The string to search for for on the page.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"element_type": schema.StringAttribute{
												Optional:    true,
												Description: "Limit the search to specific elements.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"state": schema.StringAttribute{
												Optional:    true,
												Computed:    true,
												Description: "Must be either PRESENT or ABSENT. PRESENT means the text_to_find must be found on the page for the check to succeed. ABSENT mesns the text_to_find must not be on the page for the check to succeed.",
												Default:     stringdefault.StaticString("PRESENT"),
												Validators: []validator.String{
													stringvalidator.OneOf("PRESENT", "ABSENT"),
												},
											},
										},
									},
									"check_element_on_page": schema.SingleNestedBlock{
										Description: "Check for a specific element and it's attributes on the current page.",
										Attributes: map[string]schema.Attribute{
											"id": schema.Int64Attribute{
												Computed: true,
												PlanModifiers: []planmodifier.Int64{
													int64planmodifier.UseStateForUnknown(),
												},
											},
											"elemenet_id": schema.StringAttribute{
												Optional:    true,
												Description: "The id of the element to check.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"elemenet_name": schema.StringAttribute{
												Optional:    true,
												Description: "The name of the element to check.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"state": schema.StringAttribute{
												Optional:    true,
												Computed:    true,
												Description: "Must be either PRESENT or ABSENT. PRESENT means the element must be found on the page for the check to succeed. ABSENT means the element must not be on the page for the check to succeed.",
												Default:     stringdefault.StaticString("PRESENT"),
												Validators: []validator.String{
													stringvalidator.OneOf("PRESENT", "ABSENT"),
												},
											},
											"attribute_name": schema.StringAttribute{
												Optional:    true,
												Description: "Filter element matches out by those only containing a given attribute name.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"attribute_value": schema.StringAttribute{
												Optional:    true,
												Description: "Further filter element matches out by having a given attribute value too.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"element_content": schema.StringAttribute{
												Optional:    true,
												Description: "Filter element matches out by their content.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
										},
									},
									"check_current_url": schema.SingleNestedBlock{
										Description: "Check the URL of the current page.",
										Validators: []validator.Object{
											objectvalidator.AlsoRequires(
												path.MatchRelative().AtName("url"),
												path.MatchRelative().AtName("comparison"),
											),
										},
										Attributes: map[string]schema.Attribute{
											"id": schema.Int64Attribute{
												Computed: true,
												PlanModifiers: []planmodifier.Int64{
													int64planmodifier.UseStateForUnknown(),
												},
											},
											"url": schema.StringAttribute{
												Optional:    true,
												Description: "The URL to compare against the current URL of the page.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"comparison": schema.StringAttribute{
												Optional:    true,
												Description: "Must be one of EQUALS, DOES_NOT_EQUAL, STARTS_WITH, ENDS_WITH, CONTAINS or DOES_NOT_CONTAIN. The way to compare the given url against the current URL of the page.",
												Validators: []validator.String{
													stringvalidator.OneOf("EQUALS", "DOES_NOT_EQUAL", "STARTS_WITH", "ENDS_WITH", "CONTAINS", "DOES_NOT_CONTAIN"),
												},
											},
										},
									},
									"check_url_response": schema.SingleNestedBlock{
										Description: "Check a network request made after the previous step.",
										Validators: []validator.Object{
											objectvalidator.AlsoRequires(
												path.MatchRelative().AtName("url"),
												path.MatchRelative().AtName("comparison"),
												path.MatchRelative().AtName("warning_response_time"),
												path.MatchRelative().AtName("alert_response_time"),
											),
										},
										Attributes: map[string]schema.Attribute{
											"id": schema.Int64Attribute{
												Computed: true,
												PlanModifiers: []planmodifier.Int64{
													int64planmodifier.UseStateForUnknown(),
												},
											},
											"url": schema.StringAttribute{
												Optional:    true,
												Description: "The URL to search for.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"comparison": schema.StringAttribute{
												Optional:    true,
												Description: "Must be one of EQUALS, DOES_NOT_EQUAL, STARTS_WITH, ENDS_WITH, CONTAINS or DOES_NOT_CONTAIN. The way to compare the given url against the current URL of the page.",
												Validators: []validator.String{
													stringvalidator.OneOf("EQUALS", "DOES_NOT_EQUAL", "STARTS_WITH", "ENDS_WITH", "CONTAINS", "DOES_NOT_CONTAIN"),
												},
											},
											"warning_response_time": schema.Int32Attribute{
												Optional:    true,
												Description: "The response time in milliseconds that will trigger a warning.",
												Validators: []validator.Int32{
													int32validator.AtLeast(1),
												},
											},
											"alert_response_time": schema.Int32Attribute{
												Optional:    true,
												Description: "The response time in milliseconds that will trigger the check to fail.",
												Validators: []validator.Int32{
													int32validator.AtLeast(1),
												},
											},
											"response_code": schema.Int32Attribute{
												Optional:    true,
												Description: "The response code required for the check to be successful.",
												Validators: []validator.Int32{
													int32validator.AtLeast(100),
												},
											},
											"any_info_response": schema.BoolAttribute{
												Optional:    true,
												Computed:    true,
												Description: "Accept any response code from 100-199.",
												PlanModifiers: []planmodifier.Bool{
													boolplanmodifier.UseStateForUnknown(),
												},
											},
											"any_success_response": schema.BoolAttribute{
												Optional:    true,
												Computed:    true,
												Description: "Accept any response code from 200-299.",
												PlanModifiers: []planmodifier.Bool{
													boolplanmodifier.UseStateForUnknown(),
												},
											},
											"any_redirect_response": schema.BoolAttribute{
												Optional:    true,
												Computed:    true,
												Description: "Accept any response code from 300-399.",
												PlanModifiers: []planmodifier.Bool{
													boolplanmodifier.UseStateForUnknown(),
												},
											},
											"any_client_error_response": schema.BoolAttribute{
												Optional:    true,
												Computed:    true,
												Description: "Accept any response code from 400-499.",
												PlanModifiers: []planmodifier.Bool{
													boolplanmodifier.UseStateForUnknown(),
												},
											},
											"any_server_error_response": schema.BoolAttribute{
												Optional:    true,
												Computed:    true,
												Description: "Accept any response code from 500-599.",
												PlanModifiers: []planmodifier.Bool{
													boolplanmodifier.UseStateForUnknown(),
												},
											},
										},
									},
									"check_console_log": schema.SingleNestedBlock{
										Description: "Check for a log entry made after the past step.",
										Validators: []validator.Object{
											objectvalidator.AlsoRequires(
												path.MatchRelative().AtName("log_level"),
												path.MatchRelative().AtName("message"),
												path.MatchRelative().AtName("comparison"),
											),
										},
										Attributes: map[string]schema.Attribute{
											"id": schema.Int64Attribute{
												Computed: true,
												PlanModifiers: []planmodifier.Int64{
													int64planmodifier.UseStateForUnknown(),
												},
											},
											"log_level": schema.StringAttribute{
												Optional:    true,
												Description: "Must be one of: ANY, NORMAL, WARNING or ERROR. The level of the log to check for.",
												Validators: []validator.String{
													stringvalidator.OneOf("ANY", "NORMAL", "WARNING", "ERROR"),
												},
											},
											"message": schema.StringAttribute{
												Optional:    true,
												Description: "The full or partial log message to check for.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"comparison": schema.StringAttribute{
												Optional:    true,
												Description: "Must be one of EQUALS, DOES_NOT_EQUAL, STARTS_WITH, ENDS_WITH, CONTAINS or DOES_NOT_CONTAIN. The way to compare the given message against the console logs.",
												Validators: []validator.String{
													stringvalidator.OneOf("EQUALS", "DOES_NOT_EQUAL", "STARTS_WITH", "ENDS_WITH", "CONTAINS", "DOES_NOT_CONTAIN"),
												},
											},
										},
									},
								},
							},
						},
						"network_suppression": schema.ListNestedBlock{
							Description: "Suppress one or more network calls from causing any warnings or failures.",
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
										Description: "Space for a description of what this is supressing.",
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},
									"url": schema.StringAttribute{
										Required:    true,
										Description: "The full or part URL to suppress.",
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},
									"comparison": schema.StringAttribute{
										Required:    true,
										Description: "Must be EQUALS, STARTS_WITH, ENDS_WITH or CONTAINS. The way to compare the given url to the network calls made after the last step.",
										Validators: []validator.String{
											stringvalidator.OneOf("EQUALS", "STARTS_WITH", "ENDS_WITH", "CONTAINS"),
										},
									},
									"response_code": schema.Int32Attribute{
										Optional:    true,
										Description: "The response code for the given url that is to be suppressed for warnings or alerts.",
										Validators: []validator.Int32{
											int32validator.AtLeast(100),
										},
									},
									"any_client_error": schema.BoolAttribute{
										Optional:    true,
										Computed:    true,
										Description: "Suppress any 400-499 response code for the given url.",
									},
									"any_server_error": schema.BoolAttribute{
										Optional:    true,
										Computed:    true,
										Description: "Suppress any 500-599 response code for the given url.",
									},
								},
							},
						},
						"console_message_suppression": schema.ListNestedBlock{
							Description: "Suppress one or more cosole log messages from creating a warning or failure for a Web Journey Step.",
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
										Description: "Space for a description of what this is supressing.",
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},
									"log_level": schema.StringAttribute{
										Required:    true,
										Description: "The log level to suppress. Must be ANY, WARNING or ERROR.",
										Validators: []validator.String{
											stringvalidator.OneOf("ANY", "WARNING", "ERROR"),
										},
									},
									"message": schema.StringAttribute{
										Required:    true,
										Description: "The full log message or part of the log message to suppress.",
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},
									"comparison": schema.StringAttribute{
										Required:    true,
										Description: "Must be EQUALS, STARTS_WITH, ENDS_WITH or CONTAINS. The way to compare the given full or part message to the console logs made after the previous step.",
										Validators: []validator.String{
											stringvalidator.OneOf("EQUALS", "STARTS_WITH", "ENDS_WITH", "CONTAINS"),
										},
									},
								},
							},
						},
						"action": schema.ListNestedBlock{
							Description: "The set of actions to perform at the end of the step such as clicking on elements or enterting text.",
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
										Description: "This defines the order that actions will be taken, from number lowest first to highest number last.",
										Validators: []validator.Int32{
											int32validator.AtLeast(1),
										},
									},
									"description": schema.StringAttribute{
										Required:    true,
										Description: "Space for a description of what this action does.",
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},
									"always_required": schema.BoolAttribute{
										Optional:    true,
										Computed:    true,
										Description: "If true the the action given must be able to be completed against the current page, and if it can't the check will be marked as failed. If false, and the action can't complete, for example because the element is missing, the step will continue onto the next action regardless.",
										Default:     booldefault.StaticBool(true),
									},
									"type": schema.StringAttribute{
										Required:    true,
										Description: "The type of action to perform. Options are: CLICK, DOUBLE_CLICK, RIGHT_CLICK, TEXT_INPUT, PASSWORD_INPUT, CHANGE_WINDOW_BY_ORDER, CHANGE_WINDOW_BY_TITLE, NAVIGATE_URL, WAIT, REFRESH_PAGE, CLOSE_WINDOW, CHANGE_IFRAME_BY_ORDER, CHANGE_IFRAME_BY_XPATH, SCROLL_TO_ELEMENT, TAKE_SCREENSHOT, SAVE_DOM or SELECT_OPTION.",
										Validators: []validator.String{
											stringvalidator.OneOf("CLICK", "DOUBLE_CLICK", "RIGHT_CLICK", "TEXT_INPUT", "PASSWORD_INPUT", "CHANGE_WINDOW_BY_ORDER", "CHANGE_WINDOW_BY_TITLE", "NAVIGATE_URL", "WAIT", "REFRESH_PAGE", "CLOSE_WINDOW", "CHANGE_IFRAME_BY_ORDER", "CHANGE_IFRAME_BY_XPATH", "SCROLL_TO_ELEMENT", "TAKE_SCREENSHOT", "SAVE_DOM", "SELECT_OPTION"),
										},
									},
									"window_id": schema.Int32Attribute{
										Optional:    true,
										Description: "The opening order number of the window to change focus to for CHANGE_WINDOW_BY_ORDER action types.",
										Validators: []validator.Int32{
											int32validator.AtLeast(1),
										},
									},
									"window_title": schema.StringAttribute{
										Optional:    true,
										Description: "The title of the window to change focus to for CHANGE_WINDOW_BY_TITLE action types.",
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},
									"navigate_url": schema.StringAttribute{
										Optional:    true,
										Description: "The URL to navigate to for the NAVIGATE_URL action type.",
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},
									"wait_time": schema.Int32Attribute{
										Optional:    true,
										Description: "The number of milliseconds to wait for the WAIT action type.",
										Validators: []validator.Int32{
											int32validator.AtLeast(1),
										},
									},
									"iframe_id": schema.Int32Attribute{
										Optional:    true,
										Description: "The order number of the iframe to set focus to for the CHANGE_IFRAME_BY_ORDER action type. Set to 0 if you need to move focus back to the main page.",
										Validators: []validator.Int32{
											int32validator.AtLeast(1),
										},
									},
									"iframe_xpath": schema.StringAttribute{
										Optional:    true,
										Description: "The xpath of the iframe to set focus to for the CHANGE_IFRAME_BY_XPATH action type.",
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
										},
									},
								},
								Blocks: map[string]schema.Block{
									"click": schema.SingleNestedBlock{
										Description: "The additional details needed for a CLICK, DOUBLE_CLICK or RIGHT_CLICK action type.",
										Attributes: map[string]schema.Attribute{
											"xpath": schema.StringAttribute{
												Optional:    true,
												Description: "The xpath of the element to click on. If multiple matches, the first will be used. Can not be used with search_text.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"search_text": schema.StringAttribute{
												Optional:    true,
												Description: "The text on the page to click on. If this has multiple matches then the first will be used. Can not be used with xpath.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"element_type": schema.StringAttribute{
												Optional:    true,
												Description: "Only to be used alongside search_text. The element type/name to help target the given search_text.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
										},
									},
									"text_input": schema.SingleNestedBlock{
										Description: "The additional details needed for a TEXT_INPUT action type.",
										Validators: []validator.Object{
											objectvalidator.AlsoRequires(
												path.MatchRelative().AtName("input_text"),
											),
										},
										Attributes: map[string]schema.Attribute{
											"xpath": schema.StringAttribute{
												Optional:    true,
												Description: "The xpath of the element to input text into. If multiple matches, the first will be used. Not to be used with element_id or element_name.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"element_id": schema.StringAttribute{
												Optional:    true,
												Description: "The id of the element to input text into. Not to be used with xapth or element_name.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"element_name": schema.StringAttribute{
												Optional:    true,
												Description: "The name of the element to input text into. Not to be used with xapth or element_id.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"input_text": schema.StringAttribute{
												Optional:    true,
												Description: "The text to input.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
										},
									},
									"password_input": schema.SingleNestedBlock{
										Description: "The additional details needed for a PASSWORD_INPUT action type.",
										Validators: []validator.Object{
											objectvalidator.AlsoRequires(
												path.MatchRelative().AtName("input_password"),
											),
										},
										Attributes: map[string]schema.Attribute{
											"xpath": schema.StringAttribute{
												Optional:    true,
												Description: "The xpath of the element to input the password into. If multiple matches, the first will be used. Not to be used with element_id or element_name.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"element_id": schema.StringAttribute{
												Optional:    true,
												Description: "The id of the element to input the password into. Not to be used with xapth or element_name.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"element_name": schema.StringAttribute{
												Optional:    true,
												Description: "The name of the element to input the password into. Not to be used with xapth or element_id.",
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
									"scroll_to_element": schema.SingleNestedBlock{
										Description: "The additional details needed for the SCROLL_TO_ELEMENT action type.",
										Attributes: map[string]schema.Attribute{
											"xpath": schema.StringAttribute{
												Optional:    true,
												Description: "The xpath of the element to scroll to. If multiple matches, the first will be used. Can not be used with search_text.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"search_text": schema.StringAttribute{
												Optional:    true,
												Description: "The text on the page to scroll to. If this has multiple matches then the first will be used. Can not be used with xpath.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"element_type": schema.StringAttribute{
												Optional:    true,
												Description: "Only to be used alongside search_text. The element type/name to help target the given search_text.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
										},
									},
									"select_option": schema.SingleNestedBlock{
										Description: "Additional details needed for SELECT_OPTION action type, used to choose a value from a select element.",
										Attributes: map[string]schema.Attribute{
											"element_id": schema.StringAttribute{
												Optional:    true,
												Description: "The id of the select element to select a value from. Not to be used when xpath is set.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"xpath": schema.StringAttribute{
												Optional:    true,
												Description: "The xpath of the select element to select a value from. Not to be used when element_id is set.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"option_index": schema.Int32Attribute{
												Optional:    true,
												Description: "Choose the option to select by the order it is shown in the list, starting from 0.",
												Validators: []validator.Int32{
													int32validator.AtLeast(1),
												},
											},
											"option_name": schema.StringAttribute{
												Optional:    true,
												Description: "Choose the option to select by its name shown in the list.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
											"option_value": schema.StringAttribute{
												Optional:    true,
												Description: "Choose the option to select by its form value.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
											},
										},
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
func (r *WebJourneyCheckResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WebJourneyCheckModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	check, error := r.client.CreateWebJourneyCheck(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating check",
			"Could not create check, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	// Because the response we get from the EPM API doesn't contain the any password_input values,
	// to be able to just copy its response to the state, we need to grab the passwords from the
	// plan and put them in the returned state. Easier to do than it is to go through and
	// manually populate all the ids.
	for i := 0; i < len(check.Steps); i++ {
		step := check.Steps[i]

		for x := 0; x < len(step.Actions); x++ {
			action := step.Actions[x]

			if action.PasswordInput != nil {
				action.PasswordInput.InputPassword = plan.Steps[i].Actions[x].PasswordInput.InputPassword
			}
		}
	}

	plan = *check

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *WebJourneyCheckResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WebJourneyCheckModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed check from EPM
	check, err := r.client.GetWebJourneyCheck(state.Id.ValueInt64())
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
	// Because the response we get from the EPM API doesn't contain the any password_input values,
	// to be able to just copy its response to the state, we need to grab the passwords from the
	// current state and put them in the returned state. Easier to do than it is to go through and
	// manually populate all the ids.
	passwordActions := make([]*WebJourneyActionModel, 0)

	// Get all passwords currently in state, so we use their ids to know what to match against the
	// returned passwords from the EPM API.
	for _, step := range state.Steps {
		for _, action := range step.Actions {
			if action.PasswordInput != nil {
				foundAction := action // Copy as `action` pointer will change through next loop.
				passwordActions = append(passwordActions, &foundAction)
			}
		}
	}

	for _, step := range check.Steps {
		for _, action := range step.Actions {
			if action.PasswordInput != nil {
				// Does this exist in state? If so, we'll overwrite with the state password.
				for _, passwordAction := range passwordActions {
					if passwordAction.Id.Equal(action.Id) {
						action.PasswordInput.InputPassword = passwordAction.PasswordInput.InputPassword
					}
				}
			}
		}
	}

	state = *check

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *WebJourneyCheckResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WebJourneyCheckModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	check, error := r.client.UpdateWebJourneyCheck(plan, ctx)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error creating check",
			"Could not create check, unexpected error: "+error.Error(),
		)
		return
	}

	// Update state with any computed values.
	// Because the response we get from the EPM API doesn't contain the any password_input values,
	// to be able to just copy its response to the state, we need to grab the passwords from the
	// plan and put them in the returned state. Easier to do than it is to go through and
	// manually populate all the ids.
	for i := 0; i < len(check.Steps); i++ {
		step := check.Steps[i]

		for x := 0; x < len(step.Actions); x++ {
			action := step.Actions[x]

			if action.PasswordInput != nil {
				action.PasswordInput.InputPassword = plan.Steps[i].Actions[x].PasswordInput.InputPassword
			}
		}
	}

	plan = *check

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *WebJourneyCheckResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan WebJourneyCheckModel
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

func (r *WebJourneyCheckResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WebJourneyCheckResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, _ := strconv.Atoi(req.ID)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
