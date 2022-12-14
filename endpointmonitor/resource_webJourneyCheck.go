package endpointmonitor

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func webJourneyCheck() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage web journey checks that can be set up to navigate through a website and perform period checks to ensure page elements, network calls and console logs are there or not as expected.",
		CreateContext: resourceWebJourneyCheckCreate,
		ReadContext:   resourceWebJourneyCheckRead,
		UpdateContext: resourceWebJourneyCheckUpdate,
		DeleteContext: resourceWebJourneyCheckDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "A name to describe in the check, used throughout EndPoint Monitor to describe this check, including in notifications.",
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A space to provide a longer description of the check if needed. Will default to the name if not set.",
				Optional:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Allows the enabling/disabling of the check from executing.",
				Optional:    true,
				Default:     true,
			},
			"maintenance_override": {
				Type:        schema.TypeBool,
				Description: "If set true then notifications and alerts will be suppressed for the check.",
				Optional:    true,
				Default:     true,
			},
			"start_url": {
				Type:         schema.TypeString,
				Description:  "The URL to load start the journey at.",
				Required:     true,
				ValidateFunc: validateUrl(),
			},
			"trigger_count": {
				Type:         schema.TypeInt,
				Description:  "The sequential number of failures that need to occur for a check to trigger an alert or notification.",
				Required:     true,
				ValidateFunc: validatePositiveInt(),
			},
			"result_retention": {
				Type:         schema.TypeInt,
				Description:  "The number of days to store historic results of the check.",
				Optional:     true,
				Default:      366,
				ValidateFunc: validatePositiveInt(),
			},
			"window_height": {
				Type:         schema.TypeInt,
				Description:  "The height of the browser window used for the check.",
				Optional:     true,
				Default:      1080,
				ValidateFunc: validatePositiveInt(),
			},
			"window_width": {
				Type:         schema.TypeInt,
				Description:  "The width of the browser window used for the check.",
				Optional:     true,
				Default:      1920,
				ValidateFunc: validatePositiveInt(),
			},
			"monitor_domain": {
				Type:        schema.TypeList,
				Description: "Define a domain to monitor network calls from during the check. If no monitor_domain's are defined, then all calls will be monitored.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:         schema.TypeString,
							Description:  "The domain to monitor.",
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"include_sub_domains": {
							Type:        schema.TypeBool,
							Description: "If true, all sub-domains of the domain will be monitored too. If false, just the given domain will be monitored.",
							Required:    true,
						},
					},
				},
			},
			"step": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sequence": {
							Type:         schema.TypeInt,
							Description:  "This indicates the order in which the steps will executing during the check.",
							Required:     true,
							ValidateFunc: validatePositiveInt(),
						},
						"name": {
							Type:         schema.TypeString,
							Description:  "A name to describe what the step is doing. This will be included in any alerts and notifications.",
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"type": {
							Type:         schema.TypeString,
							Description:  "Should be COMMON or CUSTOM. COMMON allows use of a pre-defined Web Journey step, common_step_id must be set when using this option. CUSTOM allows a custom one to be defined for this check.",
							Required:     true,
							ValidateFunc: validateWebJourneyStepType(),
						},
						"common_step_id": {
							Type:         schema.TypeInt,
							Description:  "If type is set to COMMON, then this should be set. The id of the Common Web Journey Step to use.",
							Optional:     true,
							ValidateFunc: validatePositiveInt(),
						},
						"wait_time": {
							Type:         schema.TypeInt,
							Description:  "The number of milliseconds to wait for any page load / actions on the page to complete before any checks on this step are started.",
							Optional:     true,
							ValidateFunc: validatePositiveInt(),
						},
						"page_load_time_warning": {
							Type:         schema.TypeInt,
							Description:  "The maximum number of milliseconds that any discovered network call can take before a warning is created for it and the check is set to a warning status.",
							Optional:     true,
							ValidateFunc: validatePositiveInt(),
						},
						"page_load_time_alert": {
							Type:         schema.TypeInt,
							Description:  "The maximum number of milliseconds that any discovered network call can take before an alert is created for it, and the check is set to a failed status.",
							Optional:     true,
							ValidateFunc: validatePositiveInt(),
						},
						"page_check": {
							Type:        schema.TypeSet,
							Description: "The set of checks to run against the currently loaded content.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"description": {
										Type:        schema.TypeString,
										Description: "A description of what this check is doing. This will be used in alerts and notifications.",
										Required:    true,
									},
									"warning_only": {
										Type:        schema.TypeBool,
										Description: "If true then if this check fails, then it will only produce a warning, not a full check failure. Default is false.",
										Optional:    true,
										Default:     false,
									},
									"type": {
										Type:         schema.TypeString,
										Description:  "The type of check to execute. Options are: CHECK_FOR_TEXT - Check for any string on or not on the current page. CHECK_FOR_ELEMENT - Check for an element and it's properties on the current page. CHECK_CURRENT_URL - Check the current url. CHECK_URL_RESPONSE - Check for specific network calls made after the last step. CHECK_CONSOLE_LOG - Check for console logs made after the last step.",
										Required:     true,
										ValidateFunc: validateWebJourneyPageCheckType(),
									},
									"check_for_text": {
										Type:        schema.TypeSet,
										Description: "Check a specific stirng is present or absent on the current page.",
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"text_to_find": {
													Type:         schema.TypeString,
													Description:  "The string to search for for on the page.",
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"element_type": {
													Type:        schema.TypeString,
													Description: "Limit the search to specific elements.",
													Optional:    true,
												},
												"state": {
													Type:         schema.TypeString,
													Description:  "Must be either PRESENT or ABSENT. PRESENT means the text_to_find must be found on the page for the check to succeed. ABSENT mesns the text_to_find must not be on the page for the check to succeed.",
													Required:     true,
													ValidateFunc: validateWebJourneyState(),
												},
											},
										},
									},
									"check_element_on_page": {
										Type:        schema.TypeSet,
										Description: "Check for a specific element and it's attributes on the current page.",
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"elemenet_id": {
													Type:        schema.TypeString,
													Description: "The id of the element to check.",
													Optional:    true,
												},
												"elemenet_name": {
													Type:     schema.TypeString,
													Default:  "The name of the element to check.",
													Optional: true,
												},
												"state": {
													Type:         schema.TypeString,
													Description:  "Must be either PRESENT or ABSENT. PRESENT means the element must be found oth epage for the check to succeed. ABSNET means the element must not be on the page for the check to succeed.",
													Required:     true,
													ValidateFunc: validateWebJourneyState(),
												},
												"attribute_name": {
													Type:        schema.TypeString,
													Description: "Filter element matches out by those only containing a given attribute name.",
													Optional:    true,
												},
												"attribute_value": {
													Type:        schema.TypeString,
													Description: "Further filter element matches out by having a given attribute value too.",
													Optional:    true,
												},
												"element_content": {
													Type:        schema.TypeString,
													Description: "Filter element matches out by their content.",
													Optional:    true,
												},
											},
										},
									},
									"check_current_url": {
										Type:        schema.TypeSet,
										Description: "Check the URL of the current page.",
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"url": {
													Type:         schema.TypeString,
													Description:  "The URL to compare against the current URL of the page.",
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"comparison": {
													Type:         schema.TypeString,
													Description:  "Must be one of EQUALS, DOES_NOT_EQUAL, STARTS_WITH, ENDS_WITH, CONTAINS or DOES_NOT_CONTAIN. The way to compare the given url against the current URL of the page.",
													Required:     true,
													ValidateFunc: validateWebJourneyCommonComparitor(),
												},
											},
										},
									},
									"check_url_response": {
										Type:        schema.TypeSet,
										Description: "Check a network request made after the previous step.",
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"url": {
													Type:         schema.TypeString,
													Description:  "The URL to search for.",
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"comparison": {
													Type:         schema.TypeString,
													Description:  "Must be one of EQUALS, DOES_NOT_EQUAL, STARTS_WITH, ENDS_WITH, CONTAINS or DOES_NOT_CONTAIN. The way to compare the given url against the network requets made.",
													Required:     true,
													ValidateFunc: validateWebJourneyCommonComparitor(),
												},
												"warning_response_time": {
													Type:         schema.TypeInt,
													Description:  "The response time in milliseconds that will trigger a warning.",
													Required:     true,
													ValidateFunc: validatePositiveInt(),
												},
												"alert_response_time": {
													Type:         schema.TypeInt,
													Description:  "The response time in milliseconds that will trigger the check to fail.",
													Required:     true,
													ValidateFunc: validatePositiveInt(),
												},
												"response_code": {
													Type:         schema.TypeInt,
													Description:  "The response code required for the check to be successful.",
													Optional:     true,
													ValidateFunc: validatePositiveInt(),
												},
												"any_info_response": {
													Type:        schema.TypeBool,
													Description: "Accept any response code from 100-199.",
													Optional:    true,
												},
												"any_success_response": {
													Type:        schema.TypeBool,
													Description: "Accept any response code from 200-299.",
													Optional:    true,
												},
												"any_redirect_response": {
													Type:        schema.TypeBool,
													Description: "Accept any response code from 300-399.",
													Optional:    true,
												},
												"any_client_error_response": {
													Type:        schema.TypeBool,
													Description: "Accept any response code from 400-499.",
													Optional:    true,
												},
												"any_server_error_response": {
													Type:        schema.TypeBool,
													Description: "Accept any response code from 500-599.",
													Optional:    true,
												},
											},
										},
									},
									"check_console_log": {
										Type:        schema.TypeSet,
										Description: "Check for a log entry made after the past step.",
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"log_level": {
													Type:         schema.TypeString,
													Description:  "Must be one of: ANY, NORMAL, WARNING or ERROR. The level of the log to check for.",
													Required:     true,
													ValidateFunc: validateWebJourneyLogLevel(),
												},
												"message": {
													Type:         schema.TypeString,
													Description:  "The full or partial log message to check for.",
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"comparison": {
													Type:         schema.TypeString,
													Description:  "Must be one of EQUALS, DOES_NOT_EQUAL, STARTS_WITH, ENDS_WITH, CONTAINS or DOES_NOT_CONTAIN. The way to compare the given message against the console logs.",
													Required:     true,
													ValidateFunc: validateWebJourneyCommonComparitor(),
												},
											},
										},
									},
								},
							},
						},
						"network_suppression": {
							Type:        schema.TypeSet,
							Description: "Suppress one or more network calls from causing any warnings or failures.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"description": {
										Type:         schema.TypeString,
										Description:  "Space for a description of what this is supressing.",
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"url": {
										Type:         schema.TypeString,
										Description:  "The full or part URL to suppress.",
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"comparison": {
										Type:         schema.TypeString,
										Description:  "Must be EQUALS, STARTS_WITH, ENDS_WITH or CONTAINS. The way to compare the given url to the network calls made after the last step.",
										Required:     true,
										ValidateFunc: validateWebJourneyPositiveComparitor(),
									},
									"response_code": {
										Type:         schema.TypeInt,
										Description:  "The response code for the given url that is to be suppressed for warnings or alerts.",
										Required:     true,
										ValidateFunc: validatePositiveInt(),
									},
									"any_client_error": {
										Type:        schema.TypeBool,
										Description: "Suppress any 400-499 response code for the given url.",
										Optional:    true,
									},
									"any_server_error": {
										Type:        schema.TypeBool,
										Description: "Suppress any 500-599 response code for the given url.",
										Optional:    true,
									},
								},
							},
						},
						"console_message_suppression": {
							Type:        schema.TypeSet,
							Description: "Suppress one or more cosole log messages from creating a warning or failure for a Web Journey Step.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"description": {
										Type:         schema.TypeString,
										Description:  "Space for a description of what this is supressing.",
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"log_level": {
										Type:         schema.TypeString,
										Description:  "The log level to suppress. Must be ANY, WARNING or ERROR.",
										Required:     true,
										ValidateFunc: validateWebJourneyLogLevel(),
									},
									"message": {
										Type:         schema.TypeString,
										Description:  "The full log message or part of the log message to suppress.",
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"comparison": {
										Type:         schema.TypeString,
										Description:  "Must be EQUALS, STARTS_WITH, ENDS_WITH or CONTAINS. The way to compare the given full or part message to the console logs made after the previous step.",
										Required:     true,
										ValidateFunc: validateWebJourneyPositiveComparitor(),
									},
								},
							},
						},
						"action": {
							Type:        schema.TypeSet,
							Description: "The set of actions to perform at the end of the step such as clicking on elements or enterting text.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sequence": {
										Type:         schema.TypeInt,
										Description:  "This defines the order that actions will be taken, from number lowest first to highest number last.",
										Required:     true,
										ValidateFunc: validatePositiveInt(),
									},
									"description": {
										Type:         schema.TypeString,
										Description:  "Space for a description of what this action does.",
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"always_required": {
										Type:        schema.TypeBool,
										Description: "If true the the action given must be able to be completed against the current page, and if it can't the check will be marked as failed. If false, and the action can't complete, for example because the element is missing, the step will continue onto the next action regardless.",
										Optional:    true,
										Default:     false,
									},
									"type": {
										Type:         schema.TypeString,
										Description:  "The type of action to perform. Options are: CLICK, DOUBLE_CLICK, RIGHT_CLICK, TEXT_INPUT, PASSWORD_INPUT, CHANGE_WINDOW_BY_ORDER, CHANGE_WINDOW_BY_TITLE, NAVIGATE_URL, WAIT, REFRESH_PAGE, CLOSE_WINDOW, CHANGE_IFRAME_BY_ORDER, CHANGE_IFRAME_BY_XPATH, SCROLL_TO_ELEMENT, TAKE_SCREENSHOT or SAVE_DOM.",
										Required:     true,
										ValidateFunc: validateWebJourneyStepActionType(),
									},
									"click_action": {
										Type:        schema.TypeSet,
										Description: "The additional details needed for a CLICK, DOUBLE_CLICK or RIGHT_CLICK action type.",
										Optional:    true,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"xpath": {
													Type:         schema.TypeString,
													Description:  "The xpath of the element to click on. If multiple matches, the first will be used. Can not be used with search_text.",
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"search_text": {
													Type:         schema.TypeString,
													Description:  "The text on the page to click on. If this has multiple matches then the first will be used. Can not be used with xpath.",
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"element_type": {
													Type:         schema.TypeString,
													Description:  "Only to be used alongside search_text. The element type/name to help target the given search_text.",
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
									},
									"text_input_action": {
										Type:        schema.TypeSet,
										Description: "The additional details needed for a TEXT_INPUT action type.",
										Optional:    true,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"xpath": {
													Type:         schema.TypeString,
													Description:  "The xpath of the element to input text into. If multiple matches, the first will be used. Not to be used with element_id or element_name.",
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"element_id": {
													Type:         schema.TypeString,
													Description:  "The id of the element to input text into. Not to be used with xapth or element_name.",
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"element_name": {
													Type:         schema.TypeString,
													Description:  "The name of the element to input text into. Not to be used with xapth or element_id.",
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"input_text": {
													Type:         schema.TypeString,
													Description:  "The text to input.",
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
									},
									"password_input_action": {
										Type:        schema.TypeSet,
										Description: "The additional details needed for a PASSWORD_INPUT action type.",
										Optional:    true,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"xpath": {
													Type:         schema.TypeString,
													Description:  "The xpath of the element to input the password into. If multiple matches, the first will be used. Not to be used with element_id or element_name.",
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"element_id": {
													Type:         schema.TypeString,
													Description:  "The id of the element to input the password into. Not to be used with xapth or element_name.",
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"element_name": {
													Type:         schema.TypeString,
													Description:  "The name of the element to input the password into. Not to be used with xapth or element_id.",
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"input_password": {
													Type:         schema.TypeString,
													Description:  "The password to input. This will not be stored in your Terraform state and ideally should be passed in to your Terraform as a environment variable rather than statically stored in your Terraform code.",
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
													Sensitive:    true,
												},
											},
										},
									},
									"window_id": {
										Type:         schema.TypeInt,
										Description:  "The opening order number of the window to change focus to for CHANGE_WINDOW_BY_ORDER action types.",
										Optional:     true,
										ValidateFunc: validatePositiveInt(),
									},
									"window_title": {
										Type:         schema.TypeString,
										Description:  "The title of the window to change focus to for CHANGE_WINDOW_BY_TITLE action types.",
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"navigate_url": {
										Type:         schema.TypeString,
										Description:  "The URL to navigate to for the NAVIGATE_URL action type.",
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"wait_time": {
										Type:         schema.TypeInt,
										Description:  "The number of milliseconds to wait for the WAIT action type.",
										Optional:     true,
										ValidateFunc: validatePositiveInt(),
									},
									"iframe_id": {
										Type:         schema.TypeInt,
										Description:  "The order number of the iframe to set focus to for the CHANGE_IFRAME_BY_ORDER action type. Set to 0 if you need to move focus back to the main page.",
										Optional:     true,
										ValidateFunc: validatePositiveInt(),
									},
									"iframe_xpath": {
										Type:         schema.TypeString,
										Description:  "The xpath of the iframe to set focus to for the CHANGE_IFRAME_BY_XPATH action type.",
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"scroll_to_element": {
										Type:        schema.TypeSet,
										Description: "The additional details needed for the SCROLL_TO_ELEMENT action type.",
										Optional:    true,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"xpath": {
													Type:         schema.TypeString,
													Description:  "The xpath of the element to scroll to. If multiple matches, the first will be used. Can not be used with search_text.",
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"search_text": {
													Type:         schema.TypeString,
													Description:  "The text on the page to scroll to. If this has multiple matches then the first will be used. Can not be used with xpath.",
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"element_type": {
													Type:         schema.TypeString,
													Description:  "Only to be used alongside search_text. The element type/name to help target the given search_text.",
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
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
			"check_host_id": {
				Type:         schema.TypeInt,
				Description:  "The id of the Check Host to run the check on.",
				Required:     true,
				ValidateFunc: validatePositiveInt(),
			},
			"check_group_id": {
				Type:         schema.TypeInt,
				Description:  "The id of the Check Group the check belongs to. This also determines check frequency.",
				Required:     true,
				ValidateFunc: validatePositiveInt(),
			},
			"proxy_host_id": {
				Type:         schema.TypeInt,
				Description:  "The id of the Proxy Host the check should use for a HTTP proxy if needed.",
				Optional:     true,
				Default:      nil,
				ValidateFunc: validatePositiveInt(),
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceWebJourneyCheckRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkId := d.Id()

	check, err := c.GetWebJourneyCheck(checkId)
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.IsNewResource() && check == nil {
		d.SetId("")
		return nil
	}

	mapWebJourneyCheckSchema(*check, d)

	return diags
}

func resourceWebJourneyCheckCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	check := mapWebJourneyCheck(d)

	o, err := c.CreateWebJourneyCheck(check, ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(o.Id)))

	resourceWebJourneyCheckRead(ctx, d, m)

	return diags
}

func resourceWebJourneyCheckUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	_, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChangesExcept() {
		check := mapWebJourneyCheck(d)

		_, err := c.UpdateWebJourneyCheck(check)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceWebJourneyCheckRead(ctx, d, m)
}

func resourceWebJourneyCheckDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteCheck(checkId)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func mapWebJourneyCheck(d *schema.ResourceData) WebJourneyCheck {
	checkId, err := strconv.Atoi(d.Id())
	if err != nil {
		checkId = 0
	}

	check := WebJourneyCheck{
		Id:                  checkId,
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Enabled:             d.Get("enabled").(bool),
		CheckType:           "WEB_JOURNEY",
		MaintenanceOverride: d.Get("maintenance_override").(bool),
		StartURL:            d.Get("start_url").(string),
		TriggerCount:        d.Get("trigger_count").(int),
		ResultRetentionDays: d.Get("result_retention").(int),
		WindowHeight:        d.Get("window_height").(int),
		WindowWidth:         d.Get("window_width").(int),
		MonitorDomains:      mapMonitorDomains(d),
		Steps:               mapWebJourneySteps(d),
		CheckHost: CheckHost{
			Id: d.Get("check_host_id").(int),
		},
		CheckGroup: CheckGroup{
			Id: d.Get("check_group_id").(int),
		},
	}

	if d.Get("proxy_host_id").(int) != 0 {
		check.ProxyHost = &ProxyHost{
			Id: d.Get("proxy_host_id").(int),
		}
	}

	return check
}

func mapMonitorDomains(d *schema.ResourceData) []MonitorDomain {
	domains := []MonitorDomain{}
	monitor_domains := d.Get("monitor_domain").([]interface{})

	for _, domainEntry := range monitor_domains {
		de := domainEntry.(map[string]interface{})
		domains = append(domains, MonitorDomain{de["domain"].(string), de["include_sub_domains"].(bool)})
	}

	return domains
}

func mapWebJourneySteps(d *schema.ResourceData) []WebJourneyStep {
	steps := []WebJourneyStep{}
	resourceSteps := d.Get("step").(*schema.Set).List()
	stepId, _ := strconv.Atoi(d.Id())

	for _, rawStep := range resourceSteps {
		resourceStep := rawStep.(map[string]interface{})

		step := WebJourneyStep{
			Id:                  stepId,
			Sequence:            resourceStep["sequence"].(int),
			Type:                resourceStep["type"].(string),
			Name:                resourceStep["name"].(string),
			CommonId:            resourceStep["common_step_id"].(int),
			WaitTime:            resourceStep["wait_time"].(int),
			WarningPageLoadTime: resourceStep["page_load_time_warning"].(int),
			AlertPageLoadTime:   resourceStep["page_load_time_alert"].(int),
			PageChecks:          mapWebJourneyPageChecks(resourceStep["page_check"].(*schema.Set)),
			AlertSuppressions: append(
				mapWebJourneyNetworkSuppressions(resourceStep["network_suppression"].(*schema.Set)),
				mapWebJourneyConsoleSuppressions(resourceStep["console_message_suppression"].(*schema.Set))...),
			Actions: mapWebJourneyActions(resourceStep["action"].(*schema.Set)),
		}

		steps = append(steps, step)
	}

	return steps
}

func mapWebJourneyPageChecks(pageChecks *schema.Set) []*WebJourneyPageCheck {
	checks := []*WebJourneyPageCheck{}

	for _, rawPageCheck := range pageChecks.List() {
		pageCheck := rawPageCheck.(map[string]interface{})
		checks = append(checks, mapWebJourneyPageCheck(pageCheck))
	}

	return checks
}

func mapWebJourneyPageCheck(rawPageCheck map[string]interface{}) *WebJourneyPageCheck {
	pageCheck := WebJourneyPageCheck{
		Id:          rawPageCheck["id"].(int),
		Description: rawPageCheck["description"].(string),
		WarningOnly: rawPageCheck["warning_only"].(bool),
		Type:        rawPageCheck["type"].(string),
	}

	checkForText := rawPageCheck["check_for_text"].(*schema.Set).List()
	checkForElement := rawPageCheck["check_element_on_page"].(*schema.Set).List()
	checkCurrentUrl := rawPageCheck["check_current_url"].(*schema.Set).List()
	checkUrlResponse := rawPageCheck["check_url_response"].(*schema.Set).List()
	checkConsoleLog := rawPageCheck["check_console_log"].(*schema.Set).List()

	if len(checkForText) > 0 {
		pageCheck.PageCheckForText = mapWebJourneyCheckForText(checkForText[0].(map[string]interface{}))
	}

	if len(checkForElement) > 0 {
		pageCheck.PageCheckForElement = mapWebJourneyCheckForElement(checkForElement[0].(map[string]interface{}))
	}

	if len(checkCurrentUrl) > 0 {
		pageCheck.PageCheckCurrentURL = mapWebJourneyCheckCurrentUrl(checkCurrentUrl[0].(map[string]interface{}))
	}

	if len(checkUrlResponse) > 0 {
		pageCheck.PageCheckURLResponse = mapWebJourneyCheckUrlResponse(checkUrlResponse[0].(map[string]interface{}))
	}

	if len(checkConsoleLog) > 0 {
		pageCheck.PageCheckConsoleLog = mapWebJourneyCheckConsoleLog(checkConsoleLog[0].(map[string]interface{}))
	}

	return &pageCheck
}

func mapWebJourneyCheckForText(pageCheck map[string]interface{}) *PageCheckForText {

	return &PageCheckForText{
		Id:          pageCheck["id"].(int),
		TextToFind:  pageCheck["text_to_find"].(string),
		ElementType: pageCheck["element_type"].(string),
		State:       pageCheck["state"].(string),
	}
}

func mapWebJourneyCheckForElement(pageCheck map[string]interface{}) *PageCheckForElement {
	return &PageCheckForElement{
		Id:             pageCheck["id"].(int),
		ElementId:      pageCheck["elemenet_id"].(string),
		ElementName:    pageCheck["elemenet_name"].(string),
		State:          pageCheck["state"].(string),
		AttributeName:  pageCheck["attribute_name"].(string),
		AttributeValue: pageCheck["attribute_value"].(string),
		ElementConent:  pageCheck["element_content"].(string),
	}
}

func mapWebJourneyCheckCurrentUrl(pageCheck map[string]interface{}) *PageCheckCurrentURL {
	return &PageCheckCurrentURL{
		Id:         pageCheck["id"].(int),
		Url:        pageCheck["url"].(string),
		Comparison: pageCheck["comparison"].(string),
	}
}

func mapWebJourneyCheckUrlResponse(pageCheck map[string]interface{}) *PageCheckURLResponse {
	return &PageCheckURLResponse{
		Id:                     pageCheck["id"].(int),
		Url:                    pageCheck["url"].(string),
		Comparison:             pageCheck["comparison"].(string),
		WarningResponseTime:    pageCheck["warning_response_time"].(int),
		AlertResponseTime:      pageCheck["alert_response_time"].(int),
		ResponseCode:           pageCheck["response_code"].(int),
		AnyInfoResponse:        pageCheck["any_info_response"].(bool),
		AnySuccessReponse:      pageCheck["any_success_response"].(bool),
		AnyRedirectResponse:    pageCheck["any_redirect_response"].(bool),
		AnyClientErrorResponse: pageCheck["any_client_error_response"].(bool),
		AnyServerErrorResponse: pageCheck["any_server_error_response"].(bool),
	}
}

func mapWebJourneyCheckConsoleLog(pageCheck map[string]interface{}) *PageCheckConsoleLog {
	return &PageCheckConsoleLog{
		Id:       pageCheck["id"].(int),
		LogLevel: pageCheck["log_level"].(string),
		Message:  pageCheck["message"].(string),
	}
}

func mapWebJourneyNetworkSuppressions(networkSuppressions *schema.Set) []*WebJourneyAlertSuppression {
	suppressions := []*WebJourneyAlertSuppression{}

	for _, networkSuppression := range networkSuppressions.List() {
		suppression := networkSuppression.(map[string]interface{})
		suppressions = append(suppressions, mapWebJourneyNetworkSuppression(suppression))
	}

	return suppressions
}

func mapWebJourneyNetworkSuppression(suppression map[string]interface{}) *WebJourneyAlertSuppression {
	return &WebJourneyAlertSuppression{
		Id:          suppression["id"].(int),
		Description: suppression["description"].(string),
		NetworkSuppression: NetworkSuppression{
			Url:            suppression["url"].(string),
			Comparison:     suppression["comparison"].(string),
			ResponseCode:   suppression["response_code"].(int),
			AnyClientError: suppression["any_client_error"].(bool),
			AnyServerError: suppression["any_server_error"].(bool),
		},
	}
}

func mapWebJourneyConsoleSuppressions(consoleSuppressions *schema.Set) []*WebJourneyAlertSuppression {
	suppressions := []*WebJourneyAlertSuppression{}

	for _, consoleSuppression := range consoleSuppressions.List() {
		suppression := consoleSuppression.(map[string]interface{})
		suppressions = append(suppressions, mapWebJourneyConsoleSuppression(suppression))
	}

	return suppressions
}

func mapWebJourneyConsoleSuppression(suppression map[string]interface{}) *WebJourneyAlertSuppression {
	return &WebJourneyAlertSuppression{
		Id:          suppression["id"].(int),
		Description: suppression["description"].(string),
		ConsoleSuppression: ConsoleSuppression{
			LogLevel:   suppression["log_level"].(string),
			Message:    suppression["message"].(string),
			Comparison: suppression["comparison"].(string),
		},
	}
}

func mapWebJourneyActions(rawActions *schema.Set) []*WebJourneyAction {
	actions := make([]*WebJourneyAction, rawActions.Len())

	for _, rawAction := range rawActions.List() {
		action := rawAction.(map[string]interface{})
		actions = append(actions, &WebJourneyAction{
			Sequence:                      action["sequence"].(int),
			Description:                   action["description"].(string),
			AlwaysRequired:                action["always_required"].(bool),
			Type:                          action["type"].(string),
			WebJourneyClickAction:         mapWebJourneyClickAction(action["click_action"].(*schema.Set)),
			WebJourneyTextInputAction:     mapWebJourneyTextInputAction(action["text_input_action"].(*schema.Set)),
			WebJourneyPasswordInputAction: mapWebJourneyPasswordInputAction(action["text_input_action"].(*schema.Set)),
			WebJourneyChangeWindowByOrder: &WebJourneyChangeWindowByOrder{WindowId: action["window_id"].(int)},
			WebJourneyChangeWindowByTitle: &WebJourneyChangeWindowByTitle{Title: action["window_title"].(string)},
			WebJourneyNavigateToUrl:       &WebJourneyNavigateToUrl{action["navigate_url"].(string)},
			WebJourneyWait:                &WebJourneyWait{action["wait_time"].(int)},
			WebJourneySelectIframeByOrder: &WebJourneySelectIframeByOrder{action["iframe_id"].(int)},
			WebJourneySelectIframeByXpath: &WebJourneySelectIframeByXpath{action["iframe_xpath"].(string)},
			WebJourneyScrollToElement:     mapWebJourneyScrollToElementAction(action["scroll_to_element"].(*schema.Set)),
		})
	}

	return actions
}

func mapWebJourneyClickAction(rawAction *schema.Set) *WebJourneyClickAction {
	action := rawAction.List()[0].(map[string]interface{})

	return &WebJourneyClickAction{
		Xpath:       action["xpath"].(string),
		SearchText:  action["search_text"].(string),
		ElementType: action["element_type"].(string),
	}
}

func mapWebJourneyTextInputAction(rawAction *schema.Set) *WebJourneyTextInputAction {
	action := rawAction.List()[0].(map[string]interface{})

	return &WebJourneyTextInputAction{
		Xpath:       action["xpath"].(string),
		ElementId:   action["element_id"].(string),
		ElementName: action["element_name"].(string),
		InputText:   action["input_text"].(string),
	}
}

func mapWebJourneyPasswordInputAction(rawAction *schema.Set) *WebJourneyPasswordInputAction {
	action := rawAction.List()[0].(map[string]interface{})

	return &WebJourneyPasswordInputAction{
		Xpath:       action["xpath"].(string),
		ElementId:   action["element_id"].(string),
		ElementName: action["element_name"].(string),
		Password:    action["input_password"].(string),
	}
}

func mapWebJourneyScrollToElementAction(rawAction *schema.Set) *WebJourneyScrollToElement {
	action := rawAction.List()[0].(map[string]interface{})

	return &WebJourneyScrollToElement{
		Xpath:       action["xpath"].(string),
		SearchText:  action["search_text"].(string),
		ElementType: action["element_type"].(string),
	}
}

func mapWebJourneyCheckSchema(check WebJourneyCheck, d *schema.ResourceData) {
	if check.ProxyHost != nil {
		d.Set("proxy_host_id", check.ProxyHost.Id)
	}

	d.SetId(strconv.Itoa(check.Id))
	d.Set("name", check.Name)
	d.Set("description", check.Description)
	d.Set("enabled", check.Enabled)
	d.Set("mainteance_override", check.MaintenanceOverride)
	d.Set("startUrl", check.StartURL)
	d.Set("trigger_count", check.TriggerCount)
	d.Set("result_retention", check.ResultRetentionDays)
	d.Set("window_height", check.WindowHeight)
	d.Set("window_width", check.WindowWidth)
	d.Set("check_host_id", check.CheckHost.Id)
	d.Set("check_group_id", check.CheckGroup.Id)
}

func validateWebJourneyStepType() schema.SchemaValidateFunc {
	types := []string{
		"COMMON",
		"CUSTOM",
	}
	return validation.StringInSlice(types, false)
}

func validateWebJourneyPageCheckType() schema.SchemaValidateFunc {
	types := []string{
		"CHECK_FOR_TEXT",
		"CHECK_FOR_ELEMENT",
		"CHECK_CURRENT_URL",
		"CHECK_URL_RESPONSE",
		"CHECK_CONSOLE_LOG",
	}
	return validation.StringInSlice(types, false)
}

func validateWebJourneyState() schema.SchemaValidateFunc {
	types := []string{
		"ABSENT",
		"PRESENT",
	}
	return validation.StringInSlice(types, false)
}

func validateWebJourneyCommonComparitor() schema.SchemaValidateFunc {
	types := []string{
		"EQUALS",
		"DOES_NOT_EQUAL",
		"STARTS_WITH",
		"ENDS_WITH",
		"CONTAINS",
		"DOES_NOT_CONTAIN",
	}
	return validation.StringInSlice(types, false)
}

func validateWebJourneyPositiveComparitor() schema.SchemaValidateFunc {
	types := []string{
		"EQUALS",
		"STARTS_WITH",
		"ENDS_WITH",
		"CONTAINS",
	}
	return validation.StringInSlice(types, false)
}

func validateWebJourneyLogLevel() schema.SchemaValidateFunc {
	types := []string{
		"ANY",
		"NORMAL",
		"WARNING",
		"ERROR",
	}
	return validation.StringInSlice(types, false)
}

func validateWebJourneyStepActionType() schema.SchemaValidateFunc {
	types := []string{
		"CLICK",
		"DOUBLE_CLICK",
		"RIGHT_CLICK",
		"TEXT_INPUT",
		"PASSWORD_INPUT",
		"CHANGE_WINDOW_BY_ORDER",
		"CHANGE_WINDOW_BY_TITLE",
		"NAVIGATE_URL",
		"WAIT",
		"REFRESH_PAGE",
		"CLOSE_WINDOW",
		"CHANGE_IFRAME_BY_ORDER",
		"CHANGE_IFRAME_BY_XPATH",
		"SCROLL_TO_ELEMENT",
		"TAKE_SCREENSHOT",
		"SAVE_DOM",
	}
	return validation.StringInSlice(types, false)
}
