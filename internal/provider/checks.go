package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Check struct {
	Id                  int64       `json:"id"`
	Name                string      `json:"name"`
	Description         string      `json:"description"`
	Enabled             bool        `json:"enabled"`
	MaintenanceOverride bool        `json:"maintenanceOverride"`
	CheckType           string      `json:"checkType"`
	CheckFrequency      int32       `json:"checkFrequency"`
	TriggerCount        int32       `json:"triggerCount"`
	ResultRetentionDays int32       `json:"resultRetentionDays"`
	CheckHost           *CheckHost  `json:"checkHost"`
	HostGroup           *HostGroup  `json:"hostGroup"`
	CheckGroup          *CheckGroup `json:"checkGroup"`
	ProxyHost           *ProxyHost  `json:"proxyHost"`
}

type CertificateCheck struct {
	Check
	AlertDaysRemaining   int32  `json:"alertDaysRemaining"`
	WarningDaysRemaining int32  `json:"warningDaysRemaining"`
	Url                  string `json:"url"`
	CheckDatesOnly       bool   `json:"checkDatesOnly"`
	CheckFullChain       bool   `json:"checkFullChain"`
}

type DnsCheck struct {
	Check
	Hostname          string   `json:"hostname"`
	ExpectedAddresses []string `json:"expectedAddresses"`
}

type PingCheck struct {
	Check
	Hostname            string `json:"hostname"`
	Timeout             int32  `json:"timeout"`
	WarningResponseTime int32  `json:"warningResponseTime"`
}

type SocketCheck struct {
	Check
	Hostname string `json:"hostname"`
	Port     int32  `json:"port"`
}

type UrlCheck struct {
	Check
	Url                  string              `json:"url"`
	RequestMethod        string              `json:"requestMethod"`
	ExpectedResponseCode int32               `json:"expectedResponseCode"`
	WarningRepsonseTime  int32               `json:"warningResponseTime"`
	AlertResponseTime    int32               `json:"alertResponseTime"`
	Timeout              int32               `json:"timeout"`
	AllowRedirects       bool                `json:"allowRedirects"`
	RequestBody          string              `json:"requestBody"`
	RequestHeaders       []RequestHeader     `json:"requestHeaders"`
	ResponseBodyChecks   []ResponseBodyCheck `json:"responseCheckStrings"`
}

type WebJourneyCheck struct {
	Check
	StartURL       string           `json:"startUrl"`
	WindowHeight   int              `json:"windowHeight"`
	WindowWidth    int              `json:"windowWidth"`
	MonitorDomains []MonitorDomain  `json:"monitorDomains"`
	Steps          []WebJourneyStep `json:"steps"`
}

type RequestHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ResponseBodyCheck struct {
	String     string `json:"string"`
	Comparator string `json:"comparator"`
}

type MonitorDomain struct {
	Domain            string `json:"domain"`
	IncludeSubDomains bool   `json:"includeSubDomains"`
}

type WebJourneyStep struct {
	Id                  int64                         `json:"id"`
	Sequence            int                           `json:"sequence"`
	Type                string                        `json:"type"`
	Name                string                        `json:"name"`
	CommonId            *int64                        `json:"commonId"`
	WaitTime            int                           `json:"waitTime"`
	WarningPageLoadTime int                           `json:"warningPageLoadTime"`
	AlertPageLoadTime   int                           `json:"alertPageLoadTime"`
	PageChecks          []*WebJourneyPageCheck        `json:"pageChecks,omitempty"`
	AlertSuppressions   []*WebJourneyAlertSuppression `json:"alertSuppressions,omitempty"`
	Actions             []*WebJourneyAction           `json:"actions,omitempty"`
}

type WebJourneyCommonStep struct {
	Id                  int                           `json:"id"`
	Name                string                        `json:"name"`
	Description         string                        `json:"description"`
	WaitTime            int                           `json:"waitTime"`
	WarningPageLoadTime int                           `json:"warningPageLoadTime"`
	AlertPageLoadTime   int                           `json:"alertPageLoadTime"`
	PageChecks          []*WebJourneyPageCheck        `json:"pageChecks,omitempty"`
	AlertSuppressions   []*WebJourneyAlertSuppression `json:"alertSuppressions,omitempty"`
	Actions             []*WebJourneyAction           `json:"actions,omitempty"`
}

type WebJourneyPageCheck struct {
	Id                   int                   `json:"id"`
	Description          string                `json:"description"`
	WarningOnly          bool                  `json:"warningOnly"`
	Type                 string                `json:"type"`
	PageCheckForText     *PageCheckForText     `json:"pageCheckForText"`
	PageCheckForElement  *PageCheckForElement  `json:"pageCheckForElement"`
	PageCheckCurrentURL  *PageCheckCurrentURL  `json:"pageCheckCurrentURL"`
	PageCheckURLResponse *PageCheckURLResponse `json:"pageCheckURLResponse"`
	PageCheckConsoleLog  *PageCheckConsoleLog  `json:"pageCheckConsoleLog"`
}

type PageCheckForText struct {
	Id          int     `json:"id"`
	TextToFind  string  `json:"textToFind"`
	ElementType *string `json:"elementType"`
	State       string  `json:"state"`
}

type PageCheckForElement struct {
	Id             int     `json:"id"`
	ElementId      *string `json:"elementId"`
	ElementName    *string `json:"elementName"`
	State          string  `json:"state"`
	AttributeName  *string `json:"attributeName"`
	AttributeValue *string `json:"attributeValue"`
	ElementConent  *string `json:"elementContent"`
}

type PageCheckCurrentURL struct {
	Id         int    `json:"id"`
	Url        string `json:"url"`
	Comparison string `json:"comparison"`
}

type PageCheckURLResponse struct {
	Id                     int    `json:"id"`
	Url                    string `json:"url"`
	Comparison             string `json:"comparison"`
	WarningResponseTime    int    `json:"warningResponseTime"`
	AlertResponseTime      int    `json:"alertResponseTime"`
	ResponseCode           int    `json:"responseCode"`
	AnyInfoResponse        bool   `json:"anyInfoResponseCode"`
	AnySuccessReponse      bool   `json:"anySuccessResponseCode"`
	AnyRedirectResponse    bool   `json:"anyRedirectResponseCode"`
	AnyClientErrorResponse bool   `json:"anyClientErrorResponseCode"`
	AnyServerErrorResponse bool   `json:"anyServerErrorResponseCode"`
}

type PageCheckConsoleLog struct {
	Id         int    `json:"id"`
	LogLevel   string `json:"logLevel"`
	Message    string `json:"message"`
	Comparison string `json:"comparison"`
}

type WebJourneyAlertSuppression struct {
	Id                 int                 `json:"id"`
	Description        string              `json:"description"`
	NetworkSuppression *NetworkSuppression `json:"networkSuppression"`
	ConsoleSuppression *ConsoleSuppression `json:"consoleSuppression"`
}

type NetworkSuppression struct {
	Id             int    `json:"id"`
	Url            string `json:"url"`
	Comparison     string `json:"comparison"`
	ResponseCode   *int32 `json:"responseCode"`
	AnyClientError bool   `json:"anyClientError"`
	AnyServerError bool   `json:"anyServerError"`
}

type ConsoleSuppression struct {
	Id         int    `json:"id"`
	LogLevel   string `json:"logLevel"`
	Message    string `json:"message"`
	Comparison string `json:"comparison"`
}

type WebJourneyAction struct {
	Id                            *int64                         `json:"id"`
	Sequence                      int                            `json:"sequence"`
	Description                   string                         `json:"description"`
	AlwaysRequired                bool                           `json:"alwaysRequired"`
	Type                          string                         `json:"type"`
	WebJourneyClickAction         *WebJourneyClickAction         `json:"webJourneyClickAction"`
	WebJourneyDoubleClickAction   *WebJourneyClickAction         `json:"webJourneyDoubleClickAction"`
	WebJourneyRightClickAction    *WebJourneyClickAction         `json:"webJourneyRightClickAction"`
	WebJourneyTextInputAction     *WebJourneyTextInputAction     `json:"webJourneyTextInputAction"`
	WebJourneyPasswordInputAction *WebJourneyPasswordInputAction `json:"webJourneyPasswordInputAction"`
	WebJourneyChangeWindowByOrder *WebJourneyChangeWindowByOrder `json:"webJourneyChangeWindowByOrder"`
	WebJourneyChangeWindowByTitle *WebJourneyChangeWindowByTitle `json:"webJourneyChangeWindowByTitle"`
	WebJourneyNavigateToUrl       *WebJourneyNavigateToUrl       `json:"webJourneyNavigateToUrl"`
	WebJourneyWait                *WebJourneyWait                `json:"webJourneyWait"`
	WebJourneySelectIframeByOrder *WebJourneySelectIframeByOrder `json:"webJourneySelectIframeByOrder"`
	WebJourneySelectIframeByXpath *WebJourneySelectIframeByXpath `json:"webJourneySelectIframeByXpath"`
	WebJourneyScrollToElement     *WebJourneyScrollToElement     `json:"webJourneyScrollToElement"`
	WebJourneySelectOption        *WebJourneySelectOption        `json:"webJourneySelectOption"`
}

type WebJourneyClickAction struct {
	Xpath       *string `json:"xpath"`
	SearchText  string  `json:"searchText"`
	ElementType *string `json:"elementType"`
}

type WebJourneyTextInputAction struct {
	Xpath       *string `json:"xpath"`
	ElementId   *string `json:"elementId"`
	ElementName *string `json:"elementName"`
	InputText   string  `json:"inputText"`
}

type WebJourneyPasswordInputAction struct {
	Xpath       *string `json:"xpath"`
	ElementId   *string `json:"elementId"`
	ElementName *string `json:"elementName"`
	NewPassword string  `json:"newInputPassword"`
}

type WebJourneyChangeWindowByOrder struct {
	WindowId int `json:"windowId"`
}

type WebJourneyChangeWindowByTitle struct {
	Title string `json:"title"`
}

type WebJourneyNavigateToUrl struct {
	Url string `json:"url"`
}

type WebJourneyWait struct {
	WaitTime int `json:"waitTime"`
}

type WebJourneySelectIframeByOrder struct {
	IframeId int `json:"iframeId"`
}

type WebJourneySelectIframeByXpath struct {
	Xpath string `json:"xpath"`
}

type WebJourneyScrollToElement struct {
	Xpath       *string `json:"xpath"`
	SearchText  *string `json:"searchText"`
	ElementType *string `json:"elementType"`
}

type WebJourneySelectOption struct {
	ElementId   *string `json:"elementId"`
	Xpath       *string `json:"xapth"`
	OptionIndex *int32  `json:"optionIndex"`
	OptionName  *string `json:"optionName"`
	OptionValue *string `json:"optionValue"`
}

// Probably nicer ways to map between the API objects and the Terraform models using reflection or something,
// but not familiar enough with Go yet to want to try that.

func mapToCertificateCheck(checkModel CertificateCheckModel) CertificateCheck {
	check := CertificateCheck{}
	check.Id = checkModel.Id.ValueInt64()
	check.Name = checkModel.Name.ValueString()
	check.Description = checkModel.Description.ValueString()
	check.Enabled = checkModel.Enabled.ValueBool()
	check.MaintenanceOverride = checkModel.MaintenanceOverride.ValueBool()
	check.CheckType = "TLS_CERTIFICATE"
	check.CheckFrequency = checkModel.CheckFrequency.ValueInt32()
	check.TriggerCount = checkModel.TriggerCount.ValueInt32()
	check.ResultRetentionDays = checkModel.ResultRetentionDays.ValueInt32()
	check.CheckGroup = &CheckGroup{Id: int(checkModel.CheckGroupId.ValueInt32())}

	if !checkModel.CheckHostId.IsNull() {
		check.CheckHost = &CheckHost{Id: int(checkModel.CheckHostId.ValueInt32())}
	}

	if !checkModel.CheckGroupId.IsNull() {
		check.HostGroup = &HostGroup{Id: int(checkModel.HostGroupId.ValueInt32())}
	}

	if !checkModel.ProxyHostId.IsNull() {
		check.ProxyHost = &ProxyHost{Id: int(checkModel.ProxyHostId.ValueInt32())}
	}

	check.Url = checkModel.Url.ValueString()
	check.AlertDaysRemaining = checkModel.AlertDaysRemaining.ValueInt32()
	check.WarningDaysRemaining = checkModel.WarningDaysRemaining.ValueInt32()
	check.CheckDatesOnly = checkModel.CheckDateOnly.ValueBool()
	check.CheckFullChain = checkModel.CheckFullChain.ValueBool()

	return check
}

func mapToDnsCheck(checkModel DnsCheckModel) DnsCheck {
	check := DnsCheck{}
	check.Id = checkModel.Id.ValueInt64()
	check.Name = checkModel.Name.ValueString()
	check.Description = checkModel.Description.ValueString()
	check.Enabled = checkModel.Enabled.ValueBool()
	check.MaintenanceOverride = checkModel.MaintenanceOverride.ValueBool()
	check.CheckType = "DNS"
	check.CheckFrequency = checkModel.CheckFrequency.ValueInt32()
	check.TriggerCount = checkModel.TriggerCount.ValueInt32()
	check.ResultRetentionDays = checkModel.ResultRetentionDays.ValueInt32()
	check.CheckGroup = &CheckGroup{Id: int(checkModel.CheckGroupId.ValueInt32())}

	if !checkModel.CheckHostId.IsNull() {
		check.CheckHost = &CheckHost{Id: int(checkModel.CheckHostId.ValueInt32())}
	}

	if !checkModel.CheckGroupId.IsNull() {
		check.HostGroup = &HostGroup{Id: int(checkModel.HostGroupId.ValueInt32())}
	}

	if !checkModel.ProxyHostId.IsNull() {
		check.ProxyHost = &ProxyHost{Id: int(checkModel.ProxyHostId.ValueInt32())}
	}

	check.Hostname = checkModel.Hostname.ValueString()
	check.ExpectedAddresses = make([]string, 0)

	for _, address := range checkModel.ExpectedAddresses {
		check.ExpectedAddresses = append(check.ExpectedAddresses, address.ValueString())
	}

	return check
}

func mapToPingCheck(checkModel PingCheckModel) PingCheck {
	check := PingCheck{}
	check.Id = checkModel.Id.ValueInt64()
	check.Name = checkModel.Name.ValueString()
	check.Description = checkModel.Description.ValueString()
	check.Enabled = checkModel.Enabled.ValueBool()
	check.MaintenanceOverride = checkModel.MaintenanceOverride.ValueBool()
	check.CheckType = "PING"
	check.CheckFrequency = checkModel.CheckFrequency.ValueInt32()
	check.TriggerCount = checkModel.TriggerCount.ValueInt32()
	check.ResultRetentionDays = checkModel.ResultRetentionDays.ValueInt32()
	check.CheckGroup = &CheckGroup{Id: int(checkModel.CheckGroupId.ValueInt32())}

	if !checkModel.CheckHostId.IsNull() {
		check.CheckHost = &CheckHost{Id: int(checkModel.CheckHostId.ValueInt32())}
	}

	if !checkModel.CheckGroupId.IsNull() {
		check.HostGroup = &HostGroup{Id: int(checkModel.HostGroupId.ValueInt32())}
	}

	if !checkModel.ProxyHostId.IsNull() {
		check.ProxyHost = &ProxyHost{Id: int(checkModel.ProxyHostId.ValueInt32())}
	}

	check.Hostname = checkModel.Hostname.ValueString()
	check.Timeout = checkModel.TimeoutTime.ValueInt32()
	check.WarningResponseTime = checkModel.WarningResponseTime.ValueInt32()

	return check
}

func mapToSocketCheck(checkModel SocketCheckModel) SocketCheck {
	check := SocketCheck{}
	check.Id = checkModel.Id.ValueInt64()
	check.Name = checkModel.Name.ValueString()
	check.Description = checkModel.Description.ValueString()
	check.Enabled = checkModel.Enabled.ValueBool()
	check.MaintenanceOverride = checkModel.MaintenanceOverride.ValueBool()
	check.CheckType = "SOCKET"
	check.CheckFrequency = checkModel.CheckFrequency.ValueInt32()
	check.TriggerCount = checkModel.TriggerCount.ValueInt32()
	check.ResultRetentionDays = checkModel.ResultRetentionDays.ValueInt32()
	check.CheckGroup = &CheckGroup{Id: int(checkModel.CheckGroupId.ValueInt32())}

	if !checkModel.CheckHostId.IsNull() {
		check.CheckHost = &CheckHost{Id: int(checkModel.CheckHostId.ValueInt32())}
	}

	if !checkModel.CheckGroupId.IsNull() {
		check.HostGroup = &HostGroup{Id: int(checkModel.HostGroupId.ValueInt32())}
	}

	if !checkModel.ProxyHostId.IsNull() {
		check.ProxyHost = &ProxyHost{Id: int(checkModel.ProxyHostId.ValueInt32())}
	}

	check.Hostname = checkModel.Hostname.ValueString()
	check.Port = checkModel.Port.ValueInt32()

	return check
}

func mapToUrlCheck(checkModel UrlCheckModel) UrlCheck {
	check := UrlCheck{}
	check.Id = checkModel.Id.ValueInt64()
	check.Name = checkModel.Name.ValueString()
	check.Description = checkModel.Description.ValueString()
	check.Enabled = checkModel.Enabled.ValueBool()
	check.MaintenanceOverride = checkModel.MaintenanceOverride.ValueBool()
	check.CheckType = "URL"
	check.CheckFrequency = checkModel.CheckFrequency.ValueInt32()
	check.TriggerCount = checkModel.TriggerCount.ValueInt32()
	check.ResultRetentionDays = checkModel.ResultRetentionDays.ValueInt32()
	check.CheckGroup = &CheckGroup{Id: int(checkModel.CheckGroupId.ValueInt32())}

	if !checkModel.CheckHostId.IsNull() {
		check.CheckHost = &CheckHost{Id: int(checkModel.CheckHostId.ValueInt32())}
	}

	if !checkModel.CheckGroupId.IsNull() {
		check.HostGroup = &HostGroup{Id: int(checkModel.HostGroupId.ValueInt32())}
	}

	if !checkModel.ProxyHostId.IsNull() {
		check.ProxyHost = &ProxyHost{Id: int(checkModel.ProxyHostId.ValueInt32())}
	}

	check.Url = checkModel.URL.ValueString()
	check.RequestMethod = checkModel.RequestMethod.ValueString()
	check.ExpectedResponseCode = checkModel.ExpectedResponseCode.ValueInt32()
	check.WarningRepsonseTime = checkModel.WarningRepsonseTime.ValueInt32()
	check.AlertResponseTime = checkModel.AlertResponseTime.ValueInt32()
	check.Timeout = checkModel.Timeout.ValueInt32()
	check.AllowRedirects = checkModel.AllowRedirects.ValueBool()
	check.RequestBody = checkModel.RequestBody.ValueString()
	check.RequestHeaders = []RequestHeader{}
	check.ResponseBodyChecks = []ResponseBodyCheck{}

	for _, header := range checkModel.RequestHeader {
		check.RequestHeaders = append(check.RequestHeaders, RequestHeader{Name: header.Name.ValueString(), Value: header.Value.ValueString()})
	}

	for _, bodyCheck := range checkModel.ResponseBodyCheck {
		check.ResponseBodyChecks = append(check.ResponseBodyChecks, ResponseBodyCheck{String: bodyCheck.String.ValueString(), Comparator: bodyCheck.Comparator.ValueString()})
	}

	return check
}

func mapToWebJourneyCheck(checkModel WebJourneyCheckModel) WebJourneyCheck {
	check := WebJourneyCheck{}
	check.Id = checkModel.Id.ValueInt64()
	check.Name = checkModel.Name.ValueString()
	check.Description = checkModel.Description.ValueString()
	check.Enabled = checkModel.Enabled.ValueBool()
	check.MaintenanceOverride = checkModel.MaintenanceOverride.ValueBool()
	check.CheckType = "WEB_JOURNEY"
	check.CheckFrequency = checkModel.CheckFrequency.ValueInt32()
	check.TriggerCount = checkModel.TriggerCount.ValueInt32()
	check.ResultRetentionDays = checkModel.ResultRetentionDays.ValueInt32()
	check.CheckGroup = &CheckGroup{Id: int(checkModel.CheckGroupId.ValueInt32())}

	if !checkModel.CheckHostId.IsNull() {
		check.CheckHost = &CheckHost{Id: int(checkModel.CheckHostId.ValueInt32())}
	}

	if !checkModel.CheckGroupId.IsNull() {
		check.HostGroup = &HostGroup{Id: int(checkModel.HostGroupId.ValueInt32())}
	}

	if !checkModel.ProxyHostId.IsNull() {
		check.ProxyHost = &ProxyHost{Id: int(checkModel.ProxyHostId.ValueInt32())}
	}

	check.StartURL = checkModel.StartUrl.ValueString()
	check.WindowHeight = int(checkModel.WindowHeight.ValueInt32())
	check.WindowWidth = int(checkModel.WindowWidth.ValueInt32())

	for _, monitorDomain := range checkModel.MonitorDomains {
		check.MonitorDomains = append(check.MonitorDomains, MonitorDomain{
			Domain:            monitorDomain.Domain.ValueString(),
			IncludeSubDomains: monitorDomain.IncludeSubDomains.ValueBool(),
		})
	}

	for _, stepModel := range checkModel.Steps {
		step := WebJourneyStep{}
		step.Id = stepModel.Id.ValueInt64()
		step.Sequence = int(stepModel.Sequence.ValueInt32())
		step.Type = stepModel.Type.ValueString()
		step.Name = stepModel.Name.ValueString()
		step.CommonId = stepModel.CommonId.ValueInt64Pointer()
		step.WaitTime = int(stepModel.WaitTime.ValueInt32())
		step.WarningPageLoadTime = int(stepModel.WarningPageLoadTime.ValueInt32())
		step.AlertPageLoadTime = int(stepModel.AlertPageLoadTime.ValueInt32())

		for _, pageCheckModel := range stepModel.PageChecks {
			pageCheck := WebJourneyPageCheck{}
			pageCheck.Id = int(pageCheckModel.Id.ValueInt64())
			pageCheck.Description = pageCheckModel.Description.ValueString()
			pageCheck.WarningOnly = pageCheckModel.WarningOnly.ValueBool()
			pageCheck.Type = pageCheckModel.Type.ValueString()

			if pageCheckModel.PageCheckForText != nil {
				pageCheck.PageCheckForText = &PageCheckForText{
					Id:          int(pageCheckModel.PageCheckForText.Id.ValueInt64()),
					TextToFind:  pageCheckModel.PageCheckForText.TextToFind.ValueString(),
					ElementType: pageCheckModel.PageCheckForText.ElementType.ValueStringPointer(),
					State:       pageCheckModel.PageCheckForText.State.ValueString(),
				}
			}

			if pageCheckModel.PageCheckForElement != nil {
				pageCheck.PageCheckForElement = &PageCheckForElement{
					Id:             int(pageCheckModel.PageCheckForElement.Id.ValueInt64()),
					ElementId:      pageCheckModel.PageCheckForElement.ElementId.ValueStringPointer(),
					ElementName:    pageCheckModel.PageCheckForElement.ElementName.ValueStringPointer(),
					State:          pageCheckModel.PageCheckForElement.State.ValueString(),
					AttributeName:  pageCheckModel.PageCheckForElement.AttributeName.ValueStringPointer(),
					AttributeValue: pageCheckModel.PageCheckForElement.AttributeValue.ValueStringPointer(),
					ElementConent:  pageCheckModel.PageCheckForElement.ElementConent.ValueStringPointer(),
				}
			}

			if pageCheckModel.PageCheckCurrentURL != nil {
				pageCheck.PageCheckCurrentURL = &PageCheckCurrentURL{
					Id:         int(pageCheckModel.PageCheckCurrentURL.Id.ValueInt64()),
					Url:        pageCheckModel.PageCheckCurrentURL.Url.ValueString(),
					Comparison: pageCheckModel.PageCheckCurrentURL.Comparison.ValueString(),
				}
			}

			if pageCheckModel.PageCheckURLResponse != nil {
				pageCheck.PageCheckURLResponse = &PageCheckURLResponse{
					Id:                     int(pageCheckModel.PageCheckURLResponse.Id.ValueInt64()),
					Url:                    pageCheckModel.PageCheckURLResponse.Url.ValueString(),
					Comparison:             pageCheckModel.PageCheckURLResponse.Comparison.ValueString(),
					WarningResponseTime:    int(pageCheckModel.PageCheckURLResponse.WarningResponseTime.ValueInt32()),
					AlertResponseTime:      int(pageCheckModel.PageCheckURLResponse.AlertResponseTime.ValueInt32()),
					ResponseCode:           int(pageCheckModel.PageCheckURLResponse.ResponseCode.ValueInt32()),
					AnyInfoResponse:        pageCheckModel.PageCheckURLResponse.AnyInfoResponse.ValueBool(),
					AnySuccessReponse:      pageCheckModel.PageCheckURLResponse.AnySuccessReponse.ValueBool(),
					AnyRedirectResponse:    pageCheckModel.PageCheckURLResponse.AnyRedirectResponse.ValueBool(),
					AnyClientErrorResponse: pageCheckModel.PageCheckURLResponse.AnyClientErrorResponse.ValueBool(),
					AnyServerErrorResponse: pageCheckModel.PageCheckURLResponse.AnyServerErrorResponse.ValueBool(),
				}
			}

			if pageCheckModel.PageCheckConsoleLog != nil {
				pageCheck.PageCheckConsoleLog = &PageCheckConsoleLog{
					Id:         int(pageCheckModel.PageCheckConsoleLog.Id.ValueInt64()),
					LogLevel:   pageCheckModel.PageCheckConsoleLog.LogLevel.ValueString(),
					Message:    pageCheckModel.PageCheckConsoleLog.Message.ValueString(),
					Comparison: pageCheckModel.PageCheckConsoleLog.Comparison.ValueString(),
				}
			}

			step.PageChecks = append(step.PageChecks, &pageCheck)
		}

		for _, networkSuppressionModel := range stepModel.NetworkSuppressions {
			step.AlertSuppressions = append(step.AlertSuppressions, &WebJourneyAlertSuppression{
				Id:          int(networkSuppressionModel.Id.ValueInt64()),
				Description: networkSuppressionModel.Description.ValueString(),
				NetworkSuppression: &NetworkSuppression{
					Id:             int(networkSuppressionModel.Id.ValueInt64()),
					Url:            networkSuppressionModel.Url.ValueString(),
					Comparison:     networkSuppressionModel.Comparison.ValueString(),
					ResponseCode:   networkSuppressionModel.ResponseCode.ValueInt32Pointer(),
					AnyClientError: networkSuppressionModel.AnyClientError.ValueBool(),
					AnyServerError: networkSuppressionModel.AnyServerError.ValueBool(),
				},
			})
		}

		for _, consoleSuppressionModel := range stepModel.ConsoleMessageSuppressions {
			step.AlertSuppressions = append(step.AlertSuppressions, &WebJourneyAlertSuppression{
				Id:          int(consoleSuppressionModel.Id.ValueInt64()),
				Description: consoleSuppressionModel.Description.ValueString(),
				ConsoleSuppression: &ConsoleSuppression{
					Id:         int(consoleSuppressionModel.Id.ValueInt64()),
					LogLevel:   consoleSuppressionModel.LogLevel.ValueString(),
					Message:    consoleSuppressionModel.Message.ValueString(),
					Comparison: consoleSuppressionModel.Comparison.ValueString(),
				},
			})
		}

		for _, actionModel := range stepModel.Actions {
			action := &WebJourneyAction{
				Id:             actionModel.Id.ValueInt64Pointer(),
				Sequence:       int(actionModel.Sequence.ValueInt32()),
				Description:    actionModel.Description.ValueString(),
				AlwaysRequired: actionModel.AlwaysRequired.ValueBool(),
				Type:           actionModel.Type.ValueString(),
			}

			switch action.Type {
			case "CLICK":
				action.WebJourneyClickAction = &WebJourneyClickAction{
					Xpath:       actionModel.Click.Xpath.ValueStringPointer(),
					ElementType: actionModel.Click.ElementType.ValueStringPointer(),
					SearchText:  actionModel.Click.SearchText.ValueString(),
				}

			case "DOUBLE_CLICK":
				action.WebJourneyDoubleClickAction = &WebJourneyClickAction{
					Xpath:       actionModel.Click.Xpath.ValueStringPointer(),
					ElementType: actionModel.Click.ElementType.ValueStringPointer(),
					SearchText:  actionModel.Click.SearchText.ValueString(),
				}

			case "RIGHT_CLICK":
				action.WebJourneyRightClickAction = &WebJourneyClickAction{
					Xpath:       actionModel.Click.Xpath.ValueStringPointer(),
					ElementType: actionModel.Click.ElementType.ValueStringPointer(),
					SearchText:  actionModel.Click.SearchText.ValueString(),
				}

			case "TEXT_INPUT":
				action.WebJourneyTextInputAction = &WebJourneyTextInputAction{
					Xpath:       actionModel.TextInput.Xpath.ValueStringPointer(),
					ElementId:   actionModel.TextInput.ElementId.ValueStringPointer(),
					ElementName: actionModel.TextInput.ElementName.ValueStringPointer(),
					InputText:   actionModel.TextInput.InputText.ValueString(),
				}

			case "PASSWORD_INPUT":
				action.WebJourneyPasswordInputAction = &WebJourneyPasswordInputAction{
					Xpath:       actionModel.PasswordInput.Xpath.ValueStringPointer(),
					ElementId:   actionModel.PasswordInput.ElementId.ValueStringPointer(),
					ElementName: actionModel.PasswordInput.ElementName.ValueStringPointer(),
					NewPassword: actionModel.PasswordInput.InputPassword.ValueString(),
				}

			case "CHANGE_WINDOW_BY_ORDER":
				action.WebJourneyChangeWindowByOrder = &WebJourneyChangeWindowByOrder{
					WindowId: int(actionModel.WindowId.ValueInt32()),
				}

			case "CHANGE_WINDOW_BY_TITLE":
				action.WebJourneyChangeWindowByTitle = &WebJourneyChangeWindowByTitle{
					Title: actionModel.WindowTitle.ValueString(),
				}

			case "NAVIGATE_URL":
				action.WebJourneyNavigateToUrl = &WebJourneyNavigateToUrl{
					Url: actionModel.NavigateUrl.ValueString(),
				}

			case "WAIT":
				action.WebJourneyWait = &WebJourneyWait{
					WaitTime: int(actionModel.WaitTime.ValueInt32()),
				}

			case "CHANGE_IFRAME_BY_ORDER":
				action.WebJourneySelectIframeByOrder = &WebJourneySelectIframeByOrder{
					IframeId: int(actionModel.IframeId.ValueInt32()),
				}

			case "CHANGE_IFRAME_BY_XPATH":
				action.WebJourneySelectIframeByXpath = &WebJourneySelectIframeByXpath{
					Xpath: actionModel.IframeXpath.ValueString(),
				}

			case "SCROLL_TO_ELEMENT":
				action.WebJourneyScrollToElement = &WebJourneyScrollToElement{
					Xpath:       actionModel.ScrollToElement.Xpath.ValueStringPointer(),
					SearchText:  actionModel.ScrollToElement.SearchText.ValueStringPointer(),
					ElementType: actionModel.ScrollToElement.ElementType.ValueStringPointer(),
				}

			case "SELECT_OPTION":
				action.WebJourneySelectOption = &WebJourneySelectOption{
					ElementId:   actionModel.SelectOption.ElementId.ValueStringPointer(),
					Xpath:       actionModel.SelectOption.Xpath.ValueStringPointer(),
					OptionIndex: actionModel.SelectOption.OptionIndex.ValueInt32Pointer(),
					OptionName:  actionModel.SelectOption.OptionName.ValueStringPointer(),
					OptionValue: actionModel.SelectOption.OptionName.ValueStringPointer(),
				}
			}

			step.Actions = append(step.Actions, action)
		}

		check.Steps = append(check.Steps, step)
	}

	return check
}

func mapToWebJourneyCommonStep(stepModel WebJourneyCommonStepModel) WebJourneyCommonStep {
	step := WebJourneyCommonStep{}
	step.Id = int(stepModel.Id.ValueInt64())
	step.Name = stepModel.Name.ValueString()
	step.Description = stepModel.Description.ValueString()
	step.WaitTime = int(stepModel.WaitTime.ValueInt32())
	step.WarningPageLoadTime = int(stepModel.WarningPageLoadTime.ValueInt32())
	step.AlertPageLoadTime = int(stepModel.AlertPageLoadTime.ValueInt32())

	for _, pageCheckModel := range stepModel.PageChecks {
		pageCheck := WebJourneyPageCheck{}
		pageCheck.Id = int(pageCheckModel.Id.ValueInt64())
		pageCheck.Description = pageCheckModel.Description.ValueString()
		pageCheck.WarningOnly = pageCheckModel.WarningOnly.ValueBool()
		pageCheck.Type = pageCheckModel.Type.ValueString()

		if pageCheckModel.PageCheckForText != nil {
			pageCheck.PageCheckForText = &PageCheckForText{
				Id:          int(pageCheckModel.PageCheckForText.Id.ValueInt64()),
				TextToFind:  pageCheckModel.PageCheckForText.TextToFind.ValueString(),
				ElementType: pageCheckModel.PageCheckForText.ElementType.ValueStringPointer(),
				State:       pageCheckModel.PageCheckForText.State.ValueString(),
			}
		}

		if pageCheck.PageCheckForElement != nil {
			pageCheck.PageCheckForElement = &PageCheckForElement{
				Id:             int(pageCheckModel.PageCheckForElement.Id.ValueInt64()),
				ElementId:      pageCheckModel.PageCheckForElement.ElementId.ValueStringPointer(),
				ElementName:    pageCheckModel.PageCheckForElement.ElementName.ValueStringPointer(),
				State:          pageCheckModel.PageCheckForElement.State.ValueString(),
				AttributeName:  pageCheckModel.PageCheckForElement.AttributeName.ValueStringPointer(),
				AttributeValue: pageCheckModel.PageCheckForElement.AttributeValue.ValueStringPointer(),
				ElementConent:  pageCheckModel.PageCheckForElement.ElementConent.ValueStringPointer(),
			}
		}

		if pageCheck.PageCheckCurrentURL != nil {
			pageCheck.PageCheckCurrentURL = &PageCheckCurrentURL{
				Id:         int(pageCheckModel.PageCheckCurrentURL.Id.ValueInt64()),
				Url:        pageCheckModel.PageCheckCurrentURL.Url.ValueString(),
				Comparison: pageCheckModel.PageCheckCurrentURL.Comparison.ValueString(),
			}
		}

		if pageCheck.PageCheckURLResponse != nil {
			pageCheck.PageCheckURLResponse = &PageCheckURLResponse{
				Id:                     int(pageCheckModel.PageCheckURLResponse.Id.ValueInt64()),
				Url:                    pageCheckModel.PageCheckURLResponse.Url.ValueString(),
				Comparison:             pageCheckModel.PageCheckURLResponse.Comparison.ValueString(),
				WarningResponseTime:    int(pageCheckModel.PageCheckURLResponse.WarningResponseTime.ValueInt32()),
				AlertResponseTime:      int(pageCheckModel.PageCheckURLResponse.AlertResponseTime.ValueInt32()),
				ResponseCode:           int(pageCheckModel.PageCheckURLResponse.ResponseCode.ValueInt32()),
				AnyInfoResponse:        pageCheckModel.PageCheckURLResponse.AnyInfoResponse.ValueBool(),
				AnySuccessReponse:      pageCheckModel.PageCheckURLResponse.AnySuccessReponse.ValueBool(),
				AnyRedirectResponse:    pageCheckModel.PageCheckURLResponse.AnyRedirectResponse.ValueBool(),
				AnyClientErrorResponse: pageCheckModel.PageCheckURLResponse.AnyClientErrorResponse.ValueBool(),
				AnyServerErrorResponse: pageCheckModel.PageCheckURLResponse.AnyServerErrorResponse.ValueBool(),
			}
		}

		if pageCheck.PageCheckConsoleLog != nil {
			pageCheck.PageCheckConsoleLog = &PageCheckConsoleLog{
				Id:         int(pageCheckModel.PageCheckConsoleLog.Id.ValueInt64()),
				LogLevel:   pageCheckModel.PageCheckConsoleLog.LogLevel.ValueString(),
				Message:    pageCheckModel.PageCheckConsoleLog.Message.ValueString(),
				Comparison: pageCheckModel.PageCheckConsoleLog.Comparison.ValueString(),
			}
		}

		step.PageChecks = append(step.PageChecks, &pageCheck)
	}

	for _, networkSuppressionModel := range stepModel.NetworkSuppressions {
		step.AlertSuppressions = append(step.AlertSuppressions, &WebJourneyAlertSuppression{
			Id:          int(networkSuppressionModel.Id.ValueInt64()),
			Description: networkSuppressionModel.Description.ValueString(),
			NetworkSuppression: &NetworkSuppression{
				Id:             int(networkSuppressionModel.Id.ValueInt64()),
				Url:            networkSuppressionModel.Url.ValueString(),
				Comparison:     networkSuppressionModel.Comparison.ValueString(),
				ResponseCode:   networkSuppressionModel.ResponseCode.ValueInt32Pointer(),
				AnyClientError: networkSuppressionModel.AnyClientError.ValueBool(),
				AnyServerError: networkSuppressionModel.AnyServerError.ValueBool(),
			},
		})
	}

	for _, consoleSuppressionModel := range stepModel.ConsoleMessageSuppressions {
		step.AlertSuppressions = append(step.AlertSuppressions, &WebJourneyAlertSuppression{
			Id:          int(consoleSuppressionModel.Id.ValueInt64()),
			Description: consoleSuppressionModel.Description.ValueString(),
			ConsoleSuppression: &ConsoleSuppression{
				Id:         int(consoleSuppressionModel.Id.ValueInt64()),
				LogLevel:   consoleSuppressionModel.LogLevel.ValueString(),
				Message:    consoleSuppressionModel.Message.ValueString(),
				Comparison: consoleSuppressionModel.Comparison.ValueString(),
			},
		})
	}

	for _, actionModel := range stepModel.Actions {
		action := &WebJourneyAction{
			Id:             actionModel.Id.ValueInt64Pointer(),
			Sequence:       int(actionModel.Sequence.ValueInt32()),
			Description:    actionModel.Description.ValueString(),
			AlwaysRequired: actionModel.AlwaysRequired.ValueBool(),
			Type:           actionModel.Type.ValueString(),
		}

		switch action.Type {
		case "CLICK":
			action.WebJourneyClickAction = &WebJourneyClickAction{
				Xpath:       actionModel.Click.Xpath.ValueStringPointer(),
				ElementType: actionModel.Click.ElementType.ValueStringPointer(),
				SearchText:  actionModel.Click.SearchText.ValueString(),
			}

		case "DOUBLE_CLICK":
			action.WebJourneyDoubleClickAction = &WebJourneyClickAction{
				Xpath:       actionModel.Click.Xpath.ValueStringPointer(),
				ElementType: actionModel.Click.ElementType.ValueStringPointer(),
				SearchText:  actionModel.Click.SearchText.ValueString(),
			}

		case "RIGHT_CLICK":
			action.WebJourneyRightClickAction = &WebJourneyClickAction{
				Xpath:       actionModel.Click.Xpath.ValueStringPointer(),
				ElementType: actionModel.Click.ElementType.ValueStringPointer(),
				SearchText:  actionModel.Click.SearchText.ValueString(),
			}

		case "TEXT_INPUT":
			action.WebJourneyTextInputAction = &WebJourneyTextInputAction{
				Xpath:       actionModel.TextInput.Xpath.ValueStringPointer(),
				ElementId:   actionModel.TextInput.ElementId.ValueStringPointer(),
				ElementName: actionModel.TextInput.ElementName.ValueStringPointer(),
				InputText:   actionModel.TextInput.InputText.ValueString(),
			}

		case "PASSWORD_INPUT":
			action.WebJourneyPasswordInputAction = &WebJourneyPasswordInputAction{
				Xpath:       actionModel.PasswordInput.Xpath.ValueStringPointer(),
				ElementId:   actionModel.PasswordInput.ElementId.ValueStringPointer(),
				ElementName: actionModel.PasswordInput.ElementName.ValueStringPointer(),
				NewPassword: actionModel.PasswordInput.InputPassword.ValueString(),
			}

		case "CHANGE_WINDOW_BY_ORDER":
			action.WebJourneyChangeWindowByOrder = &WebJourneyChangeWindowByOrder{
				WindowId: int(actionModel.WindowId.ValueInt32()),
			}

		case "CHANGE_WINDOW_BY_TITLE":
			action.WebJourneyChangeWindowByTitle = &WebJourneyChangeWindowByTitle{
				Title: actionModel.WindowTitle.ValueString(),
			}

		case "NAVIGATE_URL":
			action.WebJourneyNavigateToUrl = &WebJourneyNavigateToUrl{
				Url: actionModel.NavigateUrl.ValueString(),
			}

		case "WAIT":
			action.WebJourneyWait = &WebJourneyWait{
				WaitTime: int(actionModel.WaitTime.ValueInt32()),
			}

		case "CHANGE_IFRAME_BY_ORDER":
			action.WebJourneySelectIframeByOrder = &WebJourneySelectIframeByOrder{
				IframeId: int(actionModel.IframeId.ValueInt32()),
			}

		case "CHANGE_IFRAME_BY_XPATH":
			action.WebJourneySelectIframeByXpath = &WebJourneySelectIframeByXpath{
				Xpath: actionModel.IframeXpath.ValueString(),
			}

		case "SCROLL_TO_ELEMENT":
			action.WebJourneyScrollToElement = &WebJourneyScrollToElement{
				Xpath:       actionModel.ScrollToElement.Xpath.ValueStringPointer(),
				SearchText:  actionModel.ScrollToElement.SearchText.ValueStringPointer(),
				ElementType: actionModel.ScrollToElement.ElementType.ValueStringPointer(),
			}

		case "SELECT_OPTION":
			action.WebJourneySelectOption = &WebJourneySelectOption{
				ElementId:   actionModel.SelectOption.ElementId.ValueStringPointer(),
				Xpath:       actionModel.SelectOption.Xpath.ValueStringPointer(),
				OptionIndex: actionModel.SelectOption.OptionIndex.ValueInt32Pointer(),
				OptionName:  actionModel.SelectOption.OptionName.ValueStringPointer(),
				OptionValue: actionModel.SelectOption.OptionName.ValueStringPointer(),
			}
		}

		step.Actions = append(step.Actions, action)
	}

	return step
}

func mapToCertificateCheckModel(check CertificateCheck) CertificateCheckModel {
	checkModel := CertificateCheckModel{}
	checkModel.Id = types.Int64Value(check.Id)
	checkModel.Name = types.StringValue(check.Name)
	checkModel.Description = types.StringValue(check.Description)
	checkModel.Enabled = types.BoolValue(check.Enabled)
	checkModel.MaintenanceOverride = types.BoolValue(check.MaintenanceOverride)
	checkModel.CheckFrequency = types.Int32Value(check.CheckFrequency)
	checkModel.TriggerCount = types.Int32Value(check.TriggerCount)
	checkModel.ResultRetentionDays = types.Int32Value(check.ResultRetentionDays)

	if check.CheckHost != nil {
		checkModel.CheckHostId = types.Int32Value(int32(check.CheckHost.Id))
	}

	if check.HostGroup != nil {
		checkModel.HostGroupId = types.Int32Value(int32(check.HostGroup.Id))
	}

	if check.CheckGroup != nil {
		checkModel.CheckGroupId = types.Int32Value(int32(check.CheckGroup.Id))
	}

	if check.ProxyHost != nil {
		checkModel.ProxyHostId = types.Int32Value(int32(check.ProxyHost.Id))
	}

	checkModel.Url = types.StringValue(check.Url)
	checkModel.AlertDaysRemaining = types.Int32Value(check.AlertDaysRemaining)
	checkModel.WarningDaysRemaining = types.Int32Value(check.WarningDaysRemaining)
	checkModel.CheckDateOnly = types.BoolValue(check.CheckDatesOnly)
	checkModel.CheckFullChain = types.BoolValue(check.CheckFullChain)

	return checkModel
}

func mapToDnsCheckModel(check DnsCheck) DnsCheckModel {
	checkModel := DnsCheckModel{}
	checkModel.Id = types.Int64Value(check.Id)
	checkModel.Name = types.StringValue(check.Name)
	checkModel.Description = types.StringValue(check.Description)
	checkModel.Enabled = types.BoolValue(check.Enabled)
	checkModel.MaintenanceOverride = types.BoolValue(check.MaintenanceOverride)
	checkModel.CheckFrequency = types.Int32Value(check.CheckFrequency)
	checkModel.TriggerCount = types.Int32Value(check.TriggerCount)
	checkModel.ResultRetentionDays = types.Int32Value(check.ResultRetentionDays)

	if check.CheckHost != nil {
		checkModel.CheckHostId = types.Int32Value(int32(check.CheckHost.Id))
	}

	if check.HostGroup != nil {
		checkModel.HostGroupId = types.Int32Value(int32(check.HostGroup.Id))
	}

	if check.CheckGroup != nil {
		checkModel.CheckGroupId = types.Int32Value(int32(check.CheckGroup.Id))
	}

	if check.ProxyHost != nil {
		checkModel.ProxyHostId = types.Int32Value(int32(check.ProxyHost.Id))
	}

	checkModel.Hostname = types.StringValue(check.Hostname)

	for _, address := range check.ExpectedAddresses {
		checkModel.ExpectedAddresses = append(checkModel.ExpectedAddresses, types.StringValue(address))
	}

	return checkModel
}

func mapToPingCheckModel(check PingCheck) PingCheckModel {
	checkModel := PingCheckModel{}
	checkModel.Id = types.Int64Value(check.Id)
	checkModel.Name = types.StringValue(check.Name)
	checkModel.Description = types.StringValue(check.Description)
	checkModel.Enabled = types.BoolValue(check.Enabled)
	checkModel.MaintenanceOverride = types.BoolValue(check.MaintenanceOverride)
	checkModel.CheckFrequency = types.Int32Value(check.CheckFrequency)
	checkModel.TriggerCount = types.Int32Value(check.TriggerCount)
	checkModel.ResultRetentionDays = types.Int32Value(check.ResultRetentionDays)

	if check.CheckHost != nil {
		checkModel.CheckHostId = types.Int32Value(int32(check.CheckHost.Id))
	}

	if check.HostGroup != nil {
		checkModel.HostGroupId = types.Int32Value(int32(check.HostGroup.Id))
	}

	if check.CheckGroup != nil {
		checkModel.CheckGroupId = types.Int32Value(int32(check.CheckGroup.Id))
	}

	if check.ProxyHost != nil {
		checkModel.ProxyHostId = types.Int32Value(int32(check.ProxyHost.Id))
	}

	checkModel.Hostname = types.StringValue(check.Hostname)
	checkModel.TimeoutTime = types.Int32Value(check.Timeout)
	checkModel.WarningResponseTime = types.Int32Value(check.WarningResponseTime)

	return checkModel
}

func mapToSocketCheckModel(check SocketCheck) SocketCheckModel {
	checkModel := SocketCheckModel{}
	checkModel.Id = types.Int64Value(check.Id)
	checkModel.Name = types.StringValue(check.Name)
	checkModel.Description = types.StringValue(check.Description)
	checkModel.Enabled = types.BoolValue(check.Enabled)
	checkModel.MaintenanceOverride = types.BoolValue(check.MaintenanceOverride)
	checkModel.CheckFrequency = types.Int32Value(check.CheckFrequency)
	checkModel.TriggerCount = types.Int32Value(check.TriggerCount)
	checkModel.ResultRetentionDays = types.Int32Value(check.ResultRetentionDays)

	if check.CheckHost != nil {
		checkModel.CheckHostId = types.Int32Value(int32(check.CheckHost.Id))
	}

	if check.HostGroup != nil {
		checkModel.HostGroupId = types.Int32Value(int32(check.HostGroup.Id))
	}

	if check.CheckGroup != nil {
		checkModel.CheckGroupId = types.Int32Value(int32(check.CheckGroup.Id))
	}

	if check.ProxyHost != nil {
		checkModel.ProxyHostId = types.Int32Value(int32(check.ProxyHost.Id))
	}

	checkModel.Hostname = types.StringValue(check.Hostname)
	checkModel.Port = types.Int32Value(check.Port)

	return checkModel
}

func mapToUrlCheckModel(check UrlCheck) UrlCheckModel {
	checkModel := UrlCheckModel{}
	checkModel.Id = types.Int64Value(check.Id)
	checkModel.Name = types.StringValue(check.Name)
	checkModel.Description = types.StringValue(check.Description)
	checkModel.Enabled = types.BoolValue(check.Enabled)
	checkModel.MaintenanceOverride = types.BoolValue(check.MaintenanceOverride)
	checkModel.CheckFrequency = types.Int32Value(check.CheckFrequency)
	checkModel.TriggerCount = types.Int32Value(check.TriggerCount)
	checkModel.ResultRetentionDays = types.Int32Value(check.ResultRetentionDays)

	if check.CheckHost != nil {
		checkModel.CheckHostId = types.Int32Value(int32(check.CheckHost.Id))
	}

	if check.HostGroup != nil {
		checkModel.HostGroupId = types.Int32Value(int32(check.HostGroup.Id))
	}

	if check.CheckGroup != nil {
		checkModel.CheckGroupId = types.Int32Value(int32(check.CheckGroup.Id))
	}

	if check.ProxyHost != nil {
		checkModel.ProxyHostId = types.Int32Value(int32(check.ProxyHost.Id))
	}

	checkModel.URL = types.StringValue(check.Url)
	checkModel.RequestMethod = types.StringValue(check.RequestMethod)
	checkModel.ExpectedResponseCode = types.Int32Value(check.ExpectedResponseCode)
	checkModel.WarningRepsonseTime = types.Int32Value(check.WarningRepsonseTime)
	checkModel.AlertResponseTime = types.Int32Value(check.AlertResponseTime)
	checkModel.Timeout = types.Int32Value(check.Timeout)
	checkModel.AllowRedirects = types.BoolValue(check.AllowRedirects)
	checkModel.RequestBody = types.StringValue(check.RequestBody)

	for _, header := range check.RequestHeaders {
		checkModel.RequestHeader = append(checkModel.RequestHeader, struct {
			Name  basetypes.StringValue "tfsdk:\"name\""
			Value basetypes.StringValue "tfsdk:\"value\""
		}{Name: types.StringValue(header.Name), Value: types.StringValue(header.Value)})
	}

	for _, bodyCheck := range check.ResponseBodyChecks {
		checkModel.ResponseBodyCheck = append(checkModel.ResponseBodyCheck, struct {
			String     basetypes.StringValue "tfsdk:\"string\""
			Comparator basetypes.StringValue "tfsdk:\"comparator\""
		}{
			String:     types.StringValue(bodyCheck.String),
			Comparator: types.StringValue(bodyCheck.Comparator),
		})
	}

	return checkModel
}

func mapToWebJourneyCheckModel(check WebJourneyCheck) WebJourneyCheckModel {
	checkModel := WebJourneyCheckModel{}
	checkModel.Id = types.Int64Value(check.Id)
	checkModel.Name = types.StringValue(check.Name)
	checkModel.Description = types.StringValue(check.Description)
	checkModel.Enabled = types.BoolValue(check.Enabled)
	checkModel.MaintenanceOverride = types.BoolValue(check.MaintenanceOverride)
	checkModel.CheckFrequency = types.Int32Value(check.CheckFrequency)
	checkModel.TriggerCount = types.Int32Value(check.TriggerCount)
	checkModel.ResultRetentionDays = types.Int32Value(check.ResultRetentionDays)

	if check.CheckHost != nil {
		checkModel.CheckHostId = types.Int32Value(int32(check.CheckHost.Id))
	}

	if check.HostGroup != nil {
		checkModel.HostGroupId = types.Int32Value(int32(check.HostGroup.Id))
	}

	if check.CheckGroup != nil {
		checkModel.CheckGroupId = types.Int32Value(int32(check.CheckGroup.Id))
	}

	if check.ProxyHost != nil {
		checkModel.ProxyHostId = types.Int32Value(int32(check.ProxyHost.Id))
	}

	checkModel.StartUrl = types.StringValue(check.StartURL)
	checkModel.WindowHeight = types.Int32Value(int32(check.WindowHeight))
	checkModel.WindowWidth = types.Int32Value(int32(check.WindowWidth))

	for _, monitorDomain := range check.MonitorDomains {
		checkModel.MonitorDomains = append(checkModel.MonitorDomains,
			struct {
				Domain            basetypes.StringValue "tfsdk:\"domain\""
				IncludeSubDomains basetypes.BoolValue   "tfsdk:\"include_sub_domains\""
			}{
				Domain:            types.StringValue(monitorDomain.Domain),
				IncludeSubDomains: types.BoolValue(monitorDomain.IncludeSubDomains),
			})
	}

	for _, step := range check.Steps {
		stepModel := WebJourneyStepModel{
			Id:                  types.Int64Value(step.Id),
			Sequence:            types.Int32Value(int32(step.Sequence)),
			Type:                types.StringValue(step.Type),
			Name:                types.StringValue(step.Name),
			CommonId:            types.Int64PointerValue(step.CommonId),
			WaitTime:            types.Int32Value(int32(step.WaitTime)),
			WarningPageLoadTime: types.Int32Value(int32(step.WarningPageLoadTime)),
			AlertPageLoadTime:   types.Int32Value(int32(step.AlertPageLoadTime)),
		}

		for _, pageCheck := range step.PageChecks {
			pageCheckModel := WebJourneyPageCheckModel{
				Id:          types.Int64Value(int64(pageCheck.Id)),
				Description: types.StringValue(pageCheck.Description),
				WarningOnly: types.BoolValue(pageCheck.WarningOnly),
				Type:        types.StringValue(pageCheck.Type),
			}

			if pageCheck.PageCheckForText != nil {
				pageCheckModel.PageCheckForText = &PageCheckForTextModel{
					Id:          types.Int64Value(int64(pageCheck.PageCheckForText.Id)),
					TextToFind:  types.StringValue(pageCheck.PageCheckForText.TextToFind),
					ElementType: types.StringPointerValue(pageCheck.PageCheckForText.ElementType),
					State:       types.StringValue(pageCheck.PageCheckForText.State),
				}
			}

			if pageCheck.PageCheckForElement != nil {
				pageCheckModel.PageCheckForElement = &PageCheckForElementModel{
					Id:             types.Int64Value(int64(pageCheck.PageCheckForElement.Id)),
					ElementId:      types.StringPointerValue(pageCheck.PageCheckForElement.ElementId),
					ElementName:    types.StringPointerValue(pageCheck.PageCheckForElement.ElementName),
					State:          types.StringValue(pageCheck.PageCheckForElement.State),
					AttributeName:  types.StringPointerValue(pageCheck.PageCheckForElement.AttributeName),
					AttributeValue: types.StringPointerValue(pageCheck.PageCheckForElement.AttributeValue),
					ElementConent:  types.StringPointerValue(pageCheck.PageCheckForElement.ElementConent),
				}
			}

			if pageCheck.PageCheckCurrentURL != nil {
				pageCheckModel.PageCheckCurrentURL = &PageCheckCurrentURLModel{
					Id:         types.Int64Value(int64(pageCheck.PageCheckCurrentURL.Id)),
					Url:        types.StringValue(pageCheck.PageCheckCurrentURL.Url),
					Comparison: types.StringValue(pageCheck.PageCheckCurrentURL.Comparison),
				}
			}

			if pageCheck.PageCheckURLResponse != nil {
				pageCheckModel.PageCheckURLResponse = &PageCheckURLResponseModel{
					Id:                     types.Int64Value(int64(pageCheck.PageCheckURLResponse.Id)),
					Url:                    types.StringValue(pageCheck.PageCheckURLResponse.Url),
					Comparison:             types.StringValue(pageCheck.PageCheckURLResponse.Comparison),
					WarningResponseTime:    types.Int32Value(int32(pageCheck.PageCheckURLResponse.WarningResponseTime)),
					AlertResponseTime:      types.Int32Value(int32(pageCheck.PageCheckURLResponse.AlertResponseTime)),
					ResponseCode:           types.Int32Value(int32(pageCheck.PageCheckURLResponse.ResponseCode)),
					AnyInfoResponse:        types.BoolValue(pageCheck.PageCheckURLResponse.AnyInfoResponse),
					AnySuccessReponse:      types.BoolValue(pageCheck.PageCheckURLResponse.AnySuccessReponse),
					AnyRedirectResponse:    types.BoolValue(pageCheck.PageCheckURLResponse.AnyRedirectResponse),
					AnyClientErrorResponse: types.BoolValue(pageCheck.PageCheckURLResponse.AnyClientErrorResponse),
					AnyServerErrorResponse: types.BoolValue(pageCheck.PageCheckURLResponse.AnyServerErrorResponse),
				}
			}

			if pageCheck.PageCheckConsoleLog != nil {
				pageCheckModel.PageCheckConsoleLog = &PageCheckConsoleLogModel{
					Id:         types.Int64Value(int64(pageCheck.PageCheckConsoleLog.Id)),
					LogLevel:   types.StringValue(pageCheck.PageCheckConsoleLog.LogLevel),
					Message:    types.StringValue(pageCheck.PageCheckConsoleLog.Message),
					Comparison: types.StringValue(pageCheck.PageCheckConsoleLog.Comparison),
				}
			}

			stepModel.PageChecks = append(stepModel.PageChecks, pageCheckModel)
		}

		for _, suppression := range step.AlertSuppressions {
			if suppression.ConsoleSuppression != nil {
				stepModel.ConsoleMessageSuppressions = append(stepModel.ConsoleMessageSuppressions, ConsoleMessageSuppressionModel{
					Id:          types.Int64Value(int64(suppression.Id)),
					Description: types.StringValue(suppression.Description),
					LogLevel:    types.StringValue(suppression.ConsoleSuppression.LogLevel),
					Message:     types.StringValue(suppression.ConsoleSuppression.Message),
					Comparison:  types.StringValue(suppression.ConsoleSuppression.Comparison),
				})
			}

			if suppression.NetworkSuppression != nil {
				stepModel.NetworkSuppressions = append(stepModel.NetworkSuppressions, NetworkSuppressionModel{
					Id:             types.Int64Value(int64(suppression.Id)),
					Description:    types.StringValue(suppression.Description),
					Url:            types.StringValue(suppression.NetworkSuppression.Url),
					Comparison:     types.StringValue(suppression.NetworkSuppression.Comparison),
					AnyClientError: types.BoolValue(suppression.NetworkSuppression.AnyClientError),
					AnyServerError: types.BoolValue(suppression.NetworkSuppression.AnyServerError),
					ResponseCode:   types.Int32PointerValue(suppression.NetworkSuppression.ResponseCode),
				})
			}
		}

		for _, action := range step.Actions {
			actionModel := WebJourneyActionModel{
				Id:             types.Int64PointerValue(action.Id),
				Sequence:       types.Int32Value(int32(action.Sequence)),
				Description:    types.StringValue(action.Description),
				AlwaysRequired: types.BoolValue(action.AlwaysRequired),
				Type:           types.StringValue(action.Type),
			}

			switch action.Type {
			case "CLICK":
				if action.WebJourneyClickAction != nil {
					actionModel.Click = &WebJourneyActionClickModel{
						Xpath:       types.StringPointerValue(action.WebJourneyClickAction.Xpath),
						ElementType: types.StringPointerValue(action.WebJourneyClickAction.ElementType),
						SearchText:  types.StringValue(action.WebJourneyClickAction.SearchText),
					}
				}

			case "DOUBLE_CLICK":
				if action.WebJourneyDoubleClickAction != nil {
					actionModel.Click = &WebJourneyActionClickModel{
						Xpath:       types.StringPointerValue(action.WebJourneyDoubleClickAction.Xpath),
						ElementType: types.StringPointerValue(action.WebJourneyDoubleClickAction.ElementType),
						SearchText:  types.StringValue(action.WebJourneyDoubleClickAction.SearchText),
					}
				}

			case "RIGHT_CLICK":
				if action.WebJourneyRightClickAction != nil {
					actionModel.Click = &WebJourneyActionClickModel{
						Xpath:       types.StringPointerValue(action.WebJourneyRightClickAction.Xpath),
						ElementType: types.StringPointerValue(action.WebJourneyRightClickAction.ElementType),
						SearchText:  types.StringValue(action.WebJourneyRightClickAction.SearchText),
					}
				}

			case "TEXT_INPUT":
				if action.WebJourneyTextInputAction != nil {
					actionModel.TextInput = &WebJourneyActionTextInputModel{
						Xpath:       types.StringPointerValue(action.WebJourneyTextInputAction.Xpath),
						ElementId:   types.StringPointerValue(action.WebJourneyTextInputAction.ElementId),
						ElementName: types.StringPointerValue(action.WebJourneyTextInputAction.ElementName),
						InputText:   types.StringValue(action.WebJourneyTextInputAction.InputText),
					}
				}

			case "PASSWORD_INPUT":
				if action.WebJourneyPasswordInputAction != nil {
					actionModel.PasswordInput = &WebJourneyActionPasswordInputModel{
						Xpath:         types.StringPointerValue(action.WebJourneyPasswordInputAction.Xpath),
						ElementId:     types.StringPointerValue(action.WebJourneyPasswordInputAction.ElementId),
						ElementName:   types.StringPointerValue(action.WebJourneyPasswordInputAction.ElementName),
						InputPassword: types.StringValue(action.WebJourneyPasswordInputAction.NewPassword),
					}
				}

			case "CHANGE_WINDOW_BY_ORDER":
				if action.WebJourneyChangeWindowByOrder != nil {
					actionModel.WindowId = types.Int32Value(int32(action.WebJourneyChangeWindowByOrder.WindowId))
				}

			case "CHANGE_WINDOW_BY_TITLE":
				if action.WebJourneyChangeWindowByTitle != nil {
					actionModel.WindowTitle = types.StringValue(action.WebJourneyChangeWindowByTitle.Title)
				}

			case "NAVIGATE_URL":
				if action.WebJourneyNavigateToUrl != nil {
					actionModel.NavigateUrl = types.StringValue(action.WebJourneyNavigateToUrl.Url)
				}

			case "WAIT":
				if action.WebJourneyWait != nil {
					actionModel.WaitTime = types.Int32Value(int32(action.WebJourneyWait.WaitTime))
				}

			case "CHANGE_IFRAME_BY_ORDER":
				if action.WebJourneySelectIframeByOrder != nil {
					actionModel.IframeId = types.Int32Value(int32(action.WebJourneySelectIframeByOrder.IframeId))
				}

			case "CHANGE_IFRAME_BY_XPATH":
				if action.WebJourneySelectIframeByXpath != nil {
					actionModel.IframeXpath = types.StringValue(action.WebJourneySelectIframeByXpath.Xpath)
				}

			case "SCROLL_TO_ELEMENT":
				if action.WebJourneyScrollToElement != nil {
					actionModel.ScrollToElement = &WebJourneyActionScrollToElementModel{
						Xpath:       types.StringPointerValue(action.WebJourneyScrollToElement.Xpath),
						SearchText:  types.StringPointerValue(action.WebJourneyScrollToElement.SearchText),
						ElementType: types.StringPointerValue(action.WebJourneyScrollToElement.ElementType),
					}
				}

			case "SELECT_OPTION":
				if action.WebJourneySelectOption != nil {
					actionModel.SelectOption = &WebJourneyActionSelectOptionModel{
						ElementId:   types.StringPointerValue(action.WebJourneySelectOption.ElementId),
						Xpath:       types.StringPointerValue(action.WebJourneySelectOption.Xpath),
						OptionIndex: types.Int32PointerValue(action.WebJourneySelectOption.OptionIndex),
						OptionName:  types.StringPointerValue(action.WebJourneySelectOption.OptionName),
						OptionValue: types.StringPointerValue(action.WebJourneySelectOption.OptionValue),
					}
				}
			}

			stepModel.Actions = append(stepModel.Actions, actionModel)
		}

		checkModel.Steps = append(checkModel.Steps, stepModel)
	}

	return checkModel
}

func mapToWebJourneyCommonStepModel(step WebJourneyCommonStep) WebJourneyCommonStepModel {
	stepModel := WebJourneyCommonStepModel{
		Id:                  types.Int64Value(int64(step.Id)),
		Name:                types.StringValue(step.Name),
		Description:         types.StringValue(step.Description),
		WaitTime:            types.Int32Value(int32(step.WaitTime)),
		WarningPageLoadTime: types.Int32Value(int32(step.WarningPageLoadTime)),
		AlertPageLoadTime:   types.Int32Value(int32(step.AlertPageLoadTime)),
	}

	for _, pageCheck := range step.PageChecks {
		pageCheckModel := WebJourneyPageCheckModel{
			Id:          types.Int64Value(int64(pageCheck.Id)),
			Description: types.StringValue(pageCheck.Description),
			WarningOnly: types.BoolValue(pageCheck.WarningOnly),
			Type:        types.StringValue(pageCheck.Type),
		}

		if pageCheck.PageCheckForText != nil {
			pageCheckModel.PageCheckForText = &PageCheckForTextModel{
				Id:          types.Int64Value(int64(pageCheck.PageCheckForText.Id)),
				TextToFind:  types.StringValue(pageCheck.PageCheckForText.TextToFind),
				ElementType: types.StringPointerValue(pageCheck.PageCheckForText.ElementType),
				State:       types.StringValue(pageCheck.PageCheckForText.State),
			}
		}

		if pageCheck.PageCheckForElement != nil {
			pageCheckModel.PageCheckForElement = &PageCheckForElementModel{
				Id:             types.Int64Value(int64(pageCheck.PageCheckForElement.Id)),
				ElementId:      types.StringPointerValue(pageCheck.PageCheckForElement.ElementId),
				ElementName:    types.StringPointerValue(pageCheck.PageCheckForElement.ElementName),
				State:          types.StringValue(pageCheck.PageCheckForElement.State),
				AttributeName:  types.StringPointerValue(pageCheck.PageCheckForElement.AttributeName),
				AttributeValue: types.StringPointerValue(pageCheck.PageCheckForElement.AttributeValue),
				ElementConent:  types.StringPointerValue(pageCheck.PageCheckForElement.ElementConent),
			}
		}

		if pageCheck.PageCheckCurrentURL != nil {
			pageCheckModel.PageCheckCurrentURL = &PageCheckCurrentURLModel{
				Id:         types.Int64Value(int64(pageCheck.PageCheckCurrentURL.Id)),
				Url:        types.StringValue(pageCheck.PageCheckCurrentURL.Url),
				Comparison: types.StringValue(pageCheck.PageCheckCurrentURL.Comparison),
			}
		}

		if pageCheck.PageCheckURLResponse != nil {
			pageCheckModel.PageCheckURLResponse = &PageCheckURLResponseModel{
				Id:                     types.Int64Value(int64(pageCheck.PageCheckURLResponse.Id)),
				Url:                    types.StringValue(pageCheck.PageCheckURLResponse.Url),
				Comparison:             types.StringValue(pageCheck.PageCheckURLResponse.Comparison),
				WarningResponseTime:    types.Int32Value(int32(pageCheck.PageCheckURLResponse.WarningResponseTime)),
				AlertResponseTime:      types.Int32Value(int32(pageCheck.PageCheckURLResponse.AlertResponseTime)),
				ResponseCode:           types.Int32Value(int32(pageCheck.PageCheckURLResponse.ResponseCode)),
				AnyInfoResponse:        types.BoolValue(pageCheck.PageCheckURLResponse.AnyInfoResponse),
				AnySuccessReponse:      types.BoolValue(pageCheck.PageCheckURLResponse.AnySuccessReponse),
				AnyRedirectResponse:    types.BoolValue(pageCheck.PageCheckURLResponse.AnyRedirectResponse),
				AnyClientErrorResponse: types.BoolValue(pageCheck.PageCheckURLResponse.AnyClientErrorResponse),
				AnyServerErrorResponse: types.BoolValue(pageCheck.PageCheckURLResponse.AnyServerErrorResponse),
			}
		}

		if pageCheck.PageCheckConsoleLog != nil {
			pageCheckModel.PageCheckConsoleLog = &PageCheckConsoleLogModel{
				Id:         types.Int64Value(int64(pageCheck.PageCheckConsoleLog.Id)),
				LogLevel:   types.StringValue(pageCheck.PageCheckConsoleLog.LogLevel),
				Message:    types.StringValue(pageCheck.PageCheckConsoleLog.Message),
				Comparison: types.StringValue(pageCheck.PageCheckConsoleLog.Comparison),
			}
		}

		stepModel.PageChecks = append(stepModel.PageChecks, pageCheckModel)
	}

	for _, suppression := range step.AlertSuppressions {
		if suppression.ConsoleSuppression != nil {
			stepModel.ConsoleMessageSuppressions = append(stepModel.ConsoleMessageSuppressions, ConsoleMessageSuppressionModel{
				Id:          types.Int64Value(int64(suppression.Id)),
				Description: types.StringValue(suppression.Description),
				LogLevel:    types.StringValue(suppression.ConsoleSuppression.LogLevel),
				Message:     types.StringValue(suppression.ConsoleSuppression.Message),
				Comparison:  types.StringValue(suppression.ConsoleSuppression.Comparison),
			})
		}

		if suppression.NetworkSuppression != nil {
			stepModel.NetworkSuppressions = append(stepModel.NetworkSuppressions, NetworkSuppressionModel{
				Id:             types.Int64Value(int64(suppression.Id)),
				Description:    types.StringValue(suppression.Description),
				Url:            types.StringValue(suppression.NetworkSuppression.Url),
				Comparison:     types.StringValue(suppression.NetworkSuppression.Comparison),
				AnyClientError: types.BoolValue(suppression.NetworkSuppression.AnyClientError),
				AnyServerError: types.BoolValue(suppression.NetworkSuppression.AnyServerError),
				ResponseCode:   types.Int32PointerValue(suppression.NetworkSuppression.ResponseCode),
			})
		}
	}

	for _, action := range step.Actions {
		actionModel := WebJourneyActionModel{
			Id:             types.Int64PointerValue(action.Id),
			Sequence:       types.Int32Value(int32(action.Sequence)),
			Description:    types.StringValue(action.Description),
			AlwaysRequired: types.BoolValue(action.AlwaysRequired),
			Type:           types.StringValue(action.Type),
		}

		switch action.Type {
		case "CLICK":
			if action.WebJourneyClickAction != nil {
				actionModel.Click = &WebJourneyActionClickModel{
					Xpath:       types.StringPointerValue(action.WebJourneyClickAction.Xpath),
					ElementType: types.StringPointerValue(action.WebJourneyClickAction.ElementType),
					SearchText:  types.StringValue(action.WebJourneyClickAction.SearchText),
				}
			}

		case "DOUBLE_CLICK":
			if action.WebJourneyDoubleClickAction != nil {
				actionModel.Click = &WebJourneyActionClickModel{
					Xpath:       types.StringPointerValue(action.WebJourneyDoubleClickAction.Xpath),
					ElementType: types.StringPointerValue(action.WebJourneyDoubleClickAction.ElementType),
					SearchText:  types.StringValue(action.WebJourneyDoubleClickAction.SearchText),
				}
			}

		case "RIGHT_CLICK":
			if action.WebJourneyRightClickAction != nil {
				actionModel.Click = &WebJourneyActionClickModel{
					Xpath:       types.StringPointerValue(action.WebJourneyRightClickAction.Xpath),
					ElementType: types.StringPointerValue(action.WebJourneyRightClickAction.ElementType),
					SearchText:  types.StringValue(action.WebJourneyRightClickAction.SearchText),
				}
			}

		case "TEXT_INPUT":
			if action.WebJourneyTextInputAction != nil {
				actionModel.TextInput = &WebJourneyActionTextInputModel{
					Xpath:       types.StringPointerValue(action.WebJourneyTextInputAction.Xpath),
					ElementId:   types.StringPointerValue(action.WebJourneyTextInputAction.ElementId),
					ElementName: types.StringPointerValue(action.WebJourneyTextInputAction.ElementName),
					InputText:   types.StringValue(action.WebJourneyTextInputAction.InputText),
				}
			}

		case "PASSWORD_INPUT":
			if action.WebJourneyPasswordInputAction != nil {
				actionModel.PasswordInput = &WebJourneyActionPasswordInputModel{
					Xpath:         types.StringPointerValue(action.WebJourneyPasswordInputAction.Xpath),
					ElementId:     types.StringPointerValue(action.WebJourneyPasswordInputAction.ElementId),
					ElementName:   types.StringPointerValue(action.WebJourneyPasswordInputAction.ElementName),
					InputPassword: types.StringValue(action.WebJourneyPasswordInputAction.NewPassword),
				}
			}

		case "CHANGE_WINDOW_BY_ORDER":
			if action.WebJourneyChangeWindowByOrder != nil {
				actionModel.WindowId = types.Int32Value(int32(action.WebJourneyChangeWindowByOrder.WindowId))
			}

		case "CHANGE_WINDOW_BY_TITLE":
			if action.WebJourneyChangeWindowByTitle != nil {
				actionModel.WindowTitle = types.StringValue(action.WebJourneyChangeWindowByTitle.Title)
			}

		case "NAVIGATE_URL":
			if action.WebJourneyNavigateToUrl != nil {
				actionModel.NavigateUrl = types.StringValue(action.WebJourneyNavigateToUrl.Url)
			}

		case "WAIT":
			if action.WebJourneyWait != nil {
				actionModel.WaitTime = types.Int32Value(int32(action.WebJourneyWait.WaitTime))
			}

		case "CHANGE_IFRAME_BY_ORDER":
			if action.WebJourneySelectIframeByOrder != nil {
				actionModel.IframeId = types.Int32Value(int32(action.WebJourneySelectIframeByOrder.IframeId))
			}

		case "CHANGE_IFRAME_BY_XPATH":
			if action.WebJourneySelectIframeByXpath != nil {
				actionModel.IframeXpath = types.StringValue(action.WebJourneySelectIframeByXpath.Xpath)
			}

		case "SCROLL_TO_ELEMENT":
			if action.WebJourneyScrollToElement != nil {
				actionModel.ScrollToElement = &WebJourneyActionScrollToElementModel{
					Xpath:       types.StringPointerValue(action.WebJourneyScrollToElement.Xpath),
					SearchText:  types.StringPointerValue(action.WebJourneyScrollToElement.SearchText),
					ElementType: types.StringPointerValue(action.WebJourneyScrollToElement.ElementType),
				}
			}

		case "SELECT_OPTION":
			if action.WebJourneySelectOption != nil {
				actionModel.SelectOption = &WebJourneyActionSelectOptionModel{
					ElementId:   types.StringPointerValue(action.WebJourneySelectOption.ElementId),
					Xpath:       types.StringPointerValue(action.WebJourneySelectOption.Xpath),
					OptionIndex: types.Int32PointerValue(action.WebJourneySelectOption.OptionIndex),
					OptionName:  types.StringPointerValue(action.WebJourneySelectOption.OptionName),
					OptionValue: types.StringPointerValue(action.WebJourneySelectOption.OptionValue),
				}
			}
		}

		stepModel.Actions = append(stepModel.Actions, actionModel)
	}

	return stepModel
}
