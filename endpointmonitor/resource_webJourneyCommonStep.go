package endpointmonitor

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func webJourneyCommonStep() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage web journey common steps which are used to provide common checks and actions to take for web journey checks.",
		CreateContext: resourceWebJourneyCommonStepCreate,
		ReadContext:   resourceWebJourneyCommonStepRead,
		UpdateContext: resourceWebJourneyCommonStepUpdate,
		DeleteContext: resourceWebJourneyCommonStepDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "A name to describe what the step is doing. This will be included in any alerts and notifications.",
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Space to provide a longer description of what this common step can be used for.",
				Optional:    true,
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
							Optional:     true,
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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceWebJourneyCommonStepRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	stepId := d.Id()

	step, err := c.GetWebJourneyCommonStep(stepId)
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.IsNewResource() && step == nil {
		d.SetId("")
		return nil
	}

	mapWebJourneyCommonStepSchema(*step, d)

	return diags
}

func resourceWebJourneyCommonStepCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	step := mapWebJourneyCommonStep(d)

	o, err := c.CreateWebJourneyCommonStep(step, ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(o.Id)))

	resourceWebJourneyCommonStepRead(ctx, d, m)

	return diags
}

func resourceWebJourneyCommonStepUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	_, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChangesExcept() {
		step := mapWebJourneyCommonStep(d)

		_, err := c.UpdateWebJourneyCommonStep(step)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceWebJourneyCommonStepRead(ctx, d, m)
}

func resourceWebJourneyCommonStepDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	stepId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteCommonStep(stepId)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
