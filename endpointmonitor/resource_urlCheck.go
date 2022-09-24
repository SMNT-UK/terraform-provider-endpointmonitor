package endpointmonitor

import (
	"context"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func urlCheck() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage URL checks that test a given URL for an expected response.",
		CreateContext: resourceUrlCheckCreate,
		ReadContext:   resourceUrlCheckRead,
		UpdateContext: resourceUrlCheckUpdate,
		DeleteContext: resourceUrlCheckDelete,
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
			"url": {
				Type:         schema.TypeString,
				Description:  "The URL to check",
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
			"request_method": {
				Type:         schema.TypeString,
				Description:  "The HTTP verb used to send the request",
				Required:     true,
				ValidateFunc: validateHttpMethod(),
			},
			"expected_response_code": {
				Type:         schema.TypeInt,
				Description:  "The expected successful response code. Any code other than this will be considered a failure.",
				Required:     true,
				ValidateFunc: validatePositiveInt(),
			},
			"alert_response_time": {
				Type:         schema.TypeInt,
				Description:  "The alert response time threshold in milliseconds.",
				Required:     true,
				ValidateFunc: validatePositiveInt(),
			},
			"warning_response_time": {
				Type:         schema.TypeInt,
				Description:  "The warning response time threshold in milliseconds.",
				Required:     true,
				ValidateFunc: validatePositiveInt(),
			},
			"timeout": {
				Type:         schema.TypeInt,
				Description:  "The number of milliseconds to wait for a response before giving up.",
				Required:     true,
				ValidateFunc: validatePositiveInt(),
			},
			"allow_redirects": {
				Type:        schema.TypeBool,
				Description: "If true, the check will follow redirects. If false the initial response will be evaluated for the check.",
				Required:    true,
			},
			"request_body": {
				Type:        schema.TypeString,
				Description: "The body to send as part of the check.",
				Optional:    true,
			},
			"request_header": {
				Type:        schema.TypeList,
				Description: "Header to send as part of the check.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
			"response_body_checks": {
				Type:        schema.TypeList,
				Description: "A list of string checks to perform against the returned body from the URL.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"string": {
							Type:         schema.TypeString,
							Description:  "The string to used in this check.",
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"comparator": {
							Type:         schema.TypeString,
							Description:  "The comparison to use between the string given and the response body.",
							Required:     true,
							ValidateFunc: validateWebJourneyCommonComparitor(),
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
				ValidateFunc: validatePositiveInt(),
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceUrlCheckRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkId := d.Id()

	check, err := c.GetUrlCheck(checkId)
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.IsNewResource() && check == nil {
		d.SetId("")
		return nil
	}

	mapUrlCheckSchema(*check, d)

	return diags
}

func resourceUrlCheckCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	check := mapUrlCheck(d)

	o, err := c.CreateUrlCheck(check)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(o.Id)))

	resourceUrlCheckRead(ctx, d, m)

	return diags
}

func resourceUrlCheckUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	_, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChangesExcept() {
		check := mapUrlCheck(d)

		_, err := c.UpdateUrlCheck(check)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceUrlCheckRead(ctx, d, m)
}

func resourceUrlCheckDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func mapUrlCheck(d *schema.ResourceData) URLCheck {
	checkId, err := strconv.Atoi(d.Id())
	if err != nil {
		checkId = 0
	}

	headers := d.Get("request_header").([]interface{})
	requestHeaders := []RequestHeader{}

	for _, header := range headers {
		requestHeader := header.(map[string]interface{})
		requestHeaders = append(requestHeaders, RequestHeader{
			Name:  requestHeader["name"].(string),
			Value: requestHeader["value"].(string),
		})
	}

	bodyChecks := d.Get("response_body_checks").([]interface{})
	checkStrings := []CheckString{}

	for _, rawCheckString := range bodyChecks {
		checkString := rawCheckString.(map[string]interface{})
		checkStrings = append(checkStrings, CheckString{
			String:     checkString["string"].(string),
			Comparator: checkString["comparator"].(string),
		})
	}

	check := URLCheck{
		Id:                   checkId,
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		Enabled:              d.Get("enabled").(bool),
		CheckType:            "URL",
		MaintenanceOverride:  d.Get("maintenance_override").(bool),
		URL:                  d.Get("url").(string),
		TriggerCount:         d.Get("trigger_count").(int),
		ResultRetentionDays:  d.Get("result_retention").(int),
		RequestMethod:        d.Get("request_method").(string),
		ExpectedResponseCode: d.Get("expected_response_code").(int),
		CheckStrings:         checkStrings,
		WarningRepsonseTime:  d.Get("warning_response_time").(int),
		AlertResponseTime:    d.Get("alert_response_time").(int),
		Timeout:              d.Get("timeout").(int),
		AllowRedirects:       d.Get("allow_redirects").(bool),
		RequestBody:          d.Get("request_body").(string),
		RequestHeaders:       requestHeaders,
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

func mapUrlCheckSchema(check URLCheck, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(check.Id))
	d.Set("name", check.Name)
	d.Set("description", check.Description)
	d.Set("enabled", check.Enabled)
	d.Set("maintenance_override", check.MaintenanceOverride)
	d.Set("url", check.URL)
	d.Set("trigger_count", check.TriggerCount)
	d.Set("result_retention", check.ResultRetentionDays)
	d.Set("request_method", check.RequestMethod)
	d.Set("expected_response_code", check.ExpectedResponseCode)
	d.Set("response_body_checks", check.CheckStrings)
	d.Set("warning_response_time", check.WarningRepsonseTime)
	d.Set("alert_response_time", check.AlertResponseTime)
	d.Set("timeout", check.Timeout)
	d.Set("allow_redirects", check.AllowRedirects)
	d.Set("request_body", check.RequestBody)
	d.Set("check_host_id", check.CheckHost.Id)
	d.Set("check_group_id", check.CheckGroup.Id)
	d.Set("request_header", check.RequestHeaders)

	if check.ProxyHost != nil {
		d.Set("proxy_host_id", check.ProxyHost.Id)
	}
}

func validateUrl() schema.SchemaValidateFunc {
	check, _ := regexp.Compile(`(http)s?(:\/\/)`)
	return validation.StringMatch(check, "URL must start with http:// or https://.")
}

func validateHttpMethod() schema.SchemaValidateFunc {
	methods := []string{
		"GET",
		"PUT",
		"POST",
		"OPTIONS",
		"HEAD",
	}
	return validation.StringInSlice(methods, false)
}

func validatePositiveInt() schema.SchemaValidateFunc {
	return validation.IntAtLeast(0)
}
