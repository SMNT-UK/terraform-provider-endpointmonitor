package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type CheckGroupModel struct {
	Id             types.Int32  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	DashboardGroup types.Int32  `tfsdk:"dashboard_group_id"`
}

type CheckCommonModel struct {
	Id                  types.Int64  `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	MaintenanceOverride types.Bool   `tfsdk:"maintenance_override"`
	CheckType           types.String `tfsdk:"-"`
	CheckFrequency      types.Int32  `tfsdk:"check_frequency"`
	TriggerCount        types.Int32  `tfsdk:"trigger_count"`
	ResultRetentionDays types.Int32  `tfsdk:"result_retention"`
	CheckHostId         types.Int32  `tfsdk:"check_host_id"`
	HostGroupId         types.Int32  `tfsdk:"check_host_group_id"`
	CheckGroupId        types.Int32  `tfsdk:"check_group_id"`
	ProxyHostId         types.Int32  `tfsdk:"proxy_host_id"`
}

type CheckHostModel struct {
	Id                  types.Int32  `tfsdk:"id"`
	Hostname            types.String `tfsdk:"hostname"`
	Description         types.String `tfsdk:"description"`
	Type                types.String `tfsdk:"type"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	MaxWebJourneyChecks types.Int32  `tfsdk:"max_checks"`
	SendCheckFiles      types.Bool   `tfsdk:"send_check_files"`
}

type CertificateCheckModel struct {
	CheckCommonModel
	AlertDaysRemaining   types.Int32  `tfsdk:"alert_days_remaining"`
	WarningDaysRemaining types.Int32  `tfsdk:"warning_days_remaining"`
	Url                  types.String `tfsdk:"url"`
	CheckDateOnly        types.Bool   `tfsdk:"check_date_only"`
	CheckFullChain       types.Bool   `tfsdk:"check_full_chain"`
}

type DashboardGroupModel struct {
	Id          types.Int32  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type DnsCheckModel struct {
	CheckCommonModel
	Hostname          types.String   `tfsdk:"hostname"`
	ExpectedAddresses []types.String `tfsdk:"expected_addresses"`
}

type GenericSingleDataSource struct {
	Search types.String `tfsdk:"search"`
	Id     types.Int32  `tfsdk:"id"`
}

type GenericMultipleDataSource struct {
	Search types.String  `tfsdk:"search"`
	Ids    []types.Int32 `tfsdk:"ids"`
}

type GenericSingleDataSource64 struct {
	Search types.String `tfsdk:"search"`
	Id     types.Int64  `tfsdk:"id"`
}

type GenericMultipleDataSource64 struct {
	Search types.String  `tfsdk:"search"`
	Ids    []types.Int64 `tfsdk:"ids"`
}

type HostGroupModel struct {
	Id          types.Int32   `tfsdk:"id"`
	Name        types.String  `tfsdk:"name"`
	Description types.String  `tfsdk:"description"`
	Enabled     types.Bool    `tfsdk:"enabled"`
	Hosts       []types.Int32 `tfsdk:"check_host_ids"`
}

type MaintenancePeriodModel struct {
	Id              types.Int32   `tfsdk:"id"`
	Description     types.String  `tfsdk:"description"`
	Enabled         types.Bool    `tfsdk:"enabled"`
	StartTime       types.String  `tfsdk:"start_time"`
	EndTime         types.String  `tfsdk:"end_time"`
	DayOfWeek       types.String  `tfsdk:"day_of_week"`
	Checks          []types.Int32 `tfsdk:"check_ids"`
	CheckGroups     []types.Int32 `tfsdk:"check_group_ids"`
	DashboardGroups []types.Int32 `tfsdk:"dashboard_group_ids"`
}

type PingCheckModel struct {
	CheckCommonModel
	Hostname            types.String `tfsdk:"hostname"`
	TimeoutTime         types.Int32  `tfsdk:"timeout_time"`
	WarningResponseTime types.Int32  `tfsdk:"warning_response_time"`
}

type ProxyHostModel struct {
	Id          types.Int32  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Hostname    types.String `tfsdk:"hostname"`
	Port        types.Int32  `tfsdk:"port"`
}

type SocketCheckModel struct {
	CheckCommonModel
	Hostname types.String `tfsdk:"hostname"`
	Port     types.Int32  `tfsdk:"port"`
}

type UrlCheckModel struct {
	CheckCommonModel
	URL                  types.String `tfsdk:"url"`
	RequestMethod        types.String `tfsdk:"request_method"`
	ExpectedResponseCode types.Int32  `tfsdk:"expected_response_code"`
	WarningRepsonseTime  types.Int32  `tfsdk:"warning_response_time"`
	AlertResponseTime    types.Int32  `tfsdk:"alert_response_time"`
	Timeout              types.Int32  `tfsdk:"timeout"`
	AllowRedirects       types.Bool   `tfsdk:"allow_redirects"`
	RequestBody          types.String `tfsdk:"request_body"`
	RequestHeader        []struct {
		Name  types.String `tfsdk:"name"`
		Value types.String `tfsdk:"value"`
	} `tfsdk:"request_header"`
	ResponseBodyCheck []struct {
		String     types.String `tfsdk:"string"`
		Comparator types.String `tfsdk:"comparator"`
	} `tfsdk:"response_body_check"`
}

type AndroidJourneyCheckModel struct {
	CheckCommonModel
	Apk                  types.String                        `tfsdk:"apk"`
	ApkChecksum          types.String                        `tfsdk:"apk_checksum"`
	ScreenOrientation    types.String                        `tfsdk:"screen_orientation"`
	OverridePackageName  types.String                        `tfsdk:"override_package_name"`
	OverrideMainActivity types.String                        `tfsdk:"override_main_activity"`
	CommonSteps          []AndroidJourneyCommonStepStepModel `tfsdk:"common_step"`
	CustomSteps          []AndroidJourneyCustomStepModel     `tfsdk:"custom_step"`
}

type AndroidJourneyCustomStepModel struct {
	Id               types.Int64                   `tfsdk:"id"`
	Sequence         types.Int32                   `tfsdk:"sequence"`
	Name             types.String                  `tfsdk:"name"`
	WaitTime         types.Int32                   `tfsdk:"wait_time"`
	StepChecks       []AndroidStepCheckModel       `tfsdk:"step_check"`
	StepInteractions []AndroidStepInteractionModel `tfsdk:"step_interaction"`
}

type AndroidJourneyCommonStepStepModel struct {
	Id           types.Int64 `tfsdk:"id"`
	Sequence     types.Int32 `tfsdk:"sequence"`
	CommonStepId types.Int64 `tfsdk:"common_step_id"`
}

type AndroidJourneyCommonStepModel struct {
	Id               types.Int64                   `tfsdk:"id"`
	Name             types.String                  `tfsdk:"name"`
	Description      types.String                  `tfsdk:"description"`
	WaitTime         types.Int32                   `tfsdk:"wait_time"`
	StepChecks       []AndroidStepCheckModel       `tfsdk:"step_check"`
	StepInteractions []AndroidStepInteractionModel `tfsdk:"step_interaction"`
}

type AndroidStepCheckModel struct {
	Id              types.Int64                  `tfsdk:"id"`
	Description     types.String                 `tfsdk:"description"`
	WarningOnly     types.Bool                   `tfsdk:"warning_only"`
	Type            types.String                 `tfsdk:"type"`
	CheckForText    *AndroidCheckForTextModel    `tfsdk:"check_for_text"`
	CheckForElement *AndroidCheckForElementModel `tfsdk:"check_for_element"`
}

type AndroidCheckForTextModel struct {
	Id         types.Int64  `tfsdk:"id"`
	TextToFind types.String `tfsdk:"text_to_find"`
	State      types.String `tfsdk:"state"`
}

type AndroidCheckForElementModel struct {
	Id             types.Int64  `tfsdk:"id"`
	ComponentId    types.String `tfsdk:"component_id"`
	ComponentType  types.String `tfsdk:"component_type"`
	Xpath          types.String `tfsdk:"xpath"`
	State          types.String `tfsdk:"state"`
	AttributeName  types.String `tfsdk:"attribute_name"`
	AttributeValue types.String `tfsdk:"attribute_value"`
}

type AndroidStepInteractionModel struct {
	Id                  types.Int64                            `tfsdk:"id"`
	Sequence            types.Int32                            `tfsdk:"sequence"`
	Description         types.String                           `tfsdk:"description"`
	AlwaysRequired      types.Bool                             `tfsdk:"always_required"`
	Type                types.String                           `tfsdk:"type"`
	WaitTime            types.Int32                            `tfsdk:"wait_time"`
	Click               *AndroidClickActionModel               `tfsdk:"click"`
	TextInput           *AndroidInputTextActionModel           `tfsdk:"text_input"`
	PasswordInput       *AndroidInputPasswordActionModel       `tfsdk:"password_input"`
	RotateDisplay       *AndroidRotateDisplayActionModel       `tfsdk:"rotate_display"`
	SelectSpinnerOption *AndroidSelectSpinnerOptionActionModel `tfsdk:"select_spinner_option"`
	Swipe               *AndroidSwipeActionModel               `tfsdk:"swipe"`
}

type AndroidClickActionModel struct {
	ComponentId types.String `tfsdk:"component_id"`
	Xpath       types.String `tfsdk:"xpath"`
	SearchText  types.String `tfsdk:"search_text"`
}

type AndroidInputTextActionModel struct {
	ComponentId types.String `tfsdk:"component_id"`
	Xpath       types.String `tfsdk:"xpath"`
	InputText   types.String `tfsdk:"input_text"`
}

type AndroidInputPasswordActionModel struct {
	ComponentId   types.String `tfsdk:"component_id"`
	Xpath         types.String `tfsdk:"xpath"`
	InputPassword types.String `tfsdk:"input_password"`
}

type AndroidRotateDisplayActionModel struct {
	Orientation types.String `tfsdk:"orientation"`
}

type AndroidSelectSpinnerOptionActionModel struct {
	ComponentId        types.String `tfsdk:"component_id"`
	Xpath              types.String `tfsdk:"xpath"`
	SearchText         types.String `tfsdk:"search_text"`
	OptionListPosition types.Int32  `tfsdk:"option_list_position"`
	OptionListText     types.String `tfsdk:"option_list_text"`
}

type AndroidSwipeActionModel struct {
	ComponentId           types.String `tfsdk:"component_id"`
	Xpath                 types.String `tfsdk:"xpath"`
	StartSwipeCoordinates types.String `tfsdk:"start_swipe_coordinates"`
	SwipeDirection        types.String `tfsdk:"swipe_direction"`
	SwipeLength           types.Int32  `tfsdk:"swipe_length"`
}

type WebJourneyCheckModel struct {
	CheckCommonModel
	StartUrl       types.String `tfsdk:"start_url"`
	WindowHeight   types.Int32  `tfsdk:"window_height"`
	WindowWidth    types.Int32  `tfsdk:"window_width"`
	MonitorDomains []struct {
		Domain            types.String `tfsdk:"domain"`
		IncludeSubDomains types.Bool   `tfsdk:"include_sub_domains"`
	} `tfsdk:"monitor_domain"`
	Steps []WebJourneyStepModel `tfsdk:"step"`
}

type WebJourneyStepModel struct {
	Id                         types.Int64                      `tfsdk:"id"`
	Sequence                   types.Int32                      `tfsdk:"sequence"`
	Type                       types.String                     `tfsdk:"type"`
	Name                       types.String                     `tfsdk:"name"`
	CommonId                   types.Int64                      `tfsdk:"common_step_id"`
	WaitTime                   types.Int32                      `tfsdk:"wait_time"`
	WarningPageLoadTime        types.Int32                      `tfsdk:"page_load_time_warning"`
	AlertPageLoadTime          types.Int32                      `tfsdk:"page_load_time_alert"`
	PageChecks                 []WebJourneyPageCheckModel       `tfsdk:"page_check"`
	ConsoleMessageSuppressions []ConsoleMessageSuppressionModel `tfsdk:"console_message_suppression"`
	NetworkSuppressions        []NetworkSuppressionModel        `tfsdk:"network_suppression"`
	Actions                    []WebJourneyActionModel          `tfsdk:"action"`
}

type WebJourneyCommonStepModel struct {
	Id                         types.Int64                      `tfsdk:"id"`
	Name                       types.String                     `tfsdk:"name"`
	Description                types.String                     `tfsdk:"description"`
	WaitTime                   types.Int32                      `tfsdk:"wait_time"`
	WarningPageLoadTime        types.Int32                      `tfsdk:"page_load_time_warning"`
	AlertPageLoadTime          types.Int32                      `tfsdk:"page_load_time_alert"`
	PageChecks                 []WebJourneyPageCheckModel       `tfsdk:"page_check"`
	ConsoleMessageSuppressions []ConsoleMessageSuppressionModel `tfsdk:"console_message_suppression"`
	NetworkSuppressions        []NetworkSuppressionModel        `tfsdk:"network_suppression"`
	Actions                    []WebJourneyActionModel          `tfsdk:"action"`
}

type WebJourneyPageCheckModel struct {
	Id                   types.Int64                `tfsdk:"id"`
	Description          types.String               `tfsdk:"description"`
	WarningOnly          types.Bool                 `tfsdk:"warning_only"`
	Type                 types.String               `tfsdk:"type"`
	PageCheckForText     *PageCheckForTextModel     `tfsdk:"check_for_text"`
	PageCheckForElement  *PageCheckForElementModel  `tfsdk:"check_element_on_page"`
	PageCheckCurrentURL  *PageCheckCurrentURLModel  `tfsdk:"check_current_url"`
	PageCheckURLResponse *PageCheckURLResponseModel `tfsdk:"check_url_response"`
	PageCheckConsoleLog  *PageCheckConsoleLogModel  `tfsdk:"check_console_log"`
}

type PageCheckForTextModel struct {
	Id          types.Int64  `tfsdk:"id"`
	TextToFind  types.String `tfsdk:"text_to_find"`
	ElementType types.String `tfsdk:"element_type"`
	State       types.String `tfsdk:"state"`
}

type PageCheckForElementModel struct {
	Id             types.Int64  `tfsdk:"id"`
	ElementId      types.String `tfsdk:"element_id"`
	ElementName    types.String `tfsdk:"element_name"`
	State          types.String `tfsdk:"state"`
	AttributeName  types.String `tfsdk:"attribute_name"`
	AttributeValue types.String `tfsdk:"attribute_value"`
	ElementConent  types.String `tfsdk:"element_content"`
}

type PageCheckCurrentURLModel struct {
	Id         types.Int64  `tfsdk:"id"`
	Url        types.String `tfsdk:"url"`
	Comparison types.String `tfsdk:"comparison"`
}

type PageCheckURLResponseModel struct {
	Id                     types.Int64  `tfsdk:"id"`
	Url                    types.String `tfsdk:"url"`
	Comparison             types.String `tfsdk:"comparison"`
	WarningResponseTime    types.Int32  `tfsdk:"warning_response_time"`
	AlertResponseTime      types.Int32  `tfsdk:"alert_response_time"`
	ResponseCode           types.Int32  `tfsdk:"response_code"`
	AnyInfoResponse        types.Bool   `tfsdk:"any_info_response"`
	AnySuccessReponse      types.Bool   `tfsdk:"any_success_response"`
	AnyRedirectResponse    types.Bool   `tfsdk:"any_redirect_response"`
	AnyClientErrorResponse types.Bool   `tfsdk:"any_client_error_response"`
	AnyServerErrorResponse types.Bool   `tfsdk:"any_server_error_response"`
}

type PageCheckConsoleLogModel struct {
	Id         types.Int64  `tfsdk:"id"`
	LogLevel   types.String `tfsdk:"log_level"`
	Message    types.String `tfsdk:"message"`
	Comparison types.String `tfsdk:"comparison"`
}

type ConsoleMessageSuppressionModel struct {
	Id          types.Int64  `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	LogLevel    types.String `tfsdk:"log_level"`
	Message     types.String `tfsdk:"message"`
	Comparison  types.String `tfsdk:"comparison"`
}

type NetworkSuppressionModel struct {
	Id             types.Int64  `tfsdk:"id"`
	Description    types.String `tfsdk:"description"`
	Url            types.String `tfsdk:"url"`
	Comparison     types.String `tfsdk:"comparison"`
	AnyClientError types.Bool   `tfsdk:"any_client_error"`
	AnyServerError types.Bool   `tfsdk:"any_server_error"`
	ResponseCode   types.Int32  `tfsdk:"response_code"`
}

type WebJourneyActionModel struct {
	Id              types.Int64                           `tfsdk:"id"`
	Sequence        types.Int32                           `tfsdk:"sequence"`
	Description     types.String                          `tfsdk:"description"`
	AlwaysRequired  types.Bool                            `tfsdk:"always_required"`
	Type            types.String                          `tfsdk:"type"`
	Click           *WebJourneyActionClickModel           `tfsdk:"click"`
	IframeId        types.Int32                           `tfsdk:"iframe_id"`
	IframeXpath     types.String                          `tfsdk:"iframe_xpath"`
	NavigateUrl     types.String                          `tfsdk:"navigate_url"`
	PasswordInput   *WebJourneyActionPasswordInputModel   `tfsdk:"password_input"`
	ScrollToElement *WebJourneyActionScrollToElementModel `tfsdk:"scroll_to_element"`
	SelectOption    *WebJourneyActionSelectOptionModel    `tfsdk:"select_option"`
	TextInput       *WebJourneyActionTextInputModel       `tfsdk:"text_input"`
	WaitTime        types.Int32                           `tfsdk:"wait_time"`
	WindowId        types.Int32                           `tfsdk:"window_id"`
	WindowTitle     types.String                          `tfsdk:"window_title"`
}

type WebJourneyActionClickModel struct {
	ElementType types.String `tfsdk:"element_type"`
	SearchText  types.String `tfsdk:"search_text"`
	Xpath       types.String `tfsdk:"xpath"`
}

type WebJourneyActionPasswordInputModel struct {
	InputPassword types.String `tfsdk:"input_password"`
	ElementId     types.String `tfsdk:"element_id"`
	ElementName   types.String `tfsdk:"element_name"`
	Xpath         types.String `tfsdk:"xpath"`
}

type WebJourneyActionScrollToElementModel struct {
	ElementType types.String `tfsdk:"element_type"`
	SearchText  types.String `tfsdk:"search_text"`
	Xpath       types.String `tfsdk:"xpath"`
}

type WebJourneyActionSelectOptionModel struct {
	ElementId   types.String `tfsdk:"element_id"`
	OptionIndex types.Int32  `tfsdk:"option_index"`
	OptionName  types.String `tfsdk:"option_name"`
	OptionValue types.String `tfsdk:"option_value"`
	Xpath       types.String `tfsdk:"xpath"`
}

type WebJourneyActionTextInputModel struct {
	InputText   types.String `tfsdk:"input_text"`
	ElementId   types.String `tfsdk:"element_id"`
	ElementName types.String `tfsdk:"element_name"`
	Xpath       types.String `tfsdk:"xpath"`
}
