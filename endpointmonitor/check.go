package endpointmonitor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Check struct {
	Id                  int        `json:"id"`
	Name                string     `json:"name"`
	Description         string     `json:"description"`
	Enabled             bool       `json:"enabled"`
	MaintenanceOverride bool       `json:"maintenanceOverride"`
	CheckType           string     `json:"checkType"`
	CheckHost           CheckHost  `json:"checkHost"`
	CheckGroup          CheckGroup `json:"checkGroup"`
}

type URLCheck struct {
	Id                   int             `json:"id"`
	Name                 string          `json:"name"`
	Description          string          `json:"description"`
	Enabled              bool            `json:"enabled"`
	CheckType            string          `json:"checkType"`
	MaintenanceOverride  bool            `json:"maintenanceOverride"`
	URL                  string          `json:"url"`
	TriggerCount         int             `json:"triggerCount"`
	ResultRetentionDays  int             `json:"resultRetentionDays"`
	RequestMethod        string          `json:"requestMethod"`
	ExpectedResponseCode int             `json:"expectedResponseCode"`
	WarningRepsonseTime  int             `json:"warningResponseTime"`
	AlertResponseTime    int             `json:"alertResponseTime"`
	Timeout              int             `json:"timeout"`
	AllowRedirects       bool            `json:"allowRedirects"`
	RequestBody          string          `json:"requestBody"`
	RequestHeaders       []RequestHeader `json:"requestHeaders"`
	CheckStrings         []CheckString   `json:"checkStrings"`
	CheckHost            CheckHost       `json:"checkHost"`
	CheckGroup           CheckGroup      `json:"checkGroup"`
}

type CheckString struct {
	String     string `json:"string"`
	Comparator string `json:"comparator"`
}

type SocketCheck struct {
	Id                  int        `json:"id"`
	Name                string     `json:"name"`
	Description         string     `json:"description"`
	Enabled             bool       `json:"enabled"`
	CheckType           string     `json:"checkType"`
	MaintenanceOverride bool       `json:"maintenanceOverride"`
	Hostname            string     `json:"hostname"`
	Port                int        `json:"port"`
	TriggerCount        int        `json:"triggerCount"`
	ResultRetentionDays int        `json:"resultRetentionDays"`
	CheckHost           CheckHost  `json:"checkHost"`
	CheckGroup          CheckGroup `json:"checkGroup"`
}

type DNSCheck struct {
	Id                  int        `json:"id"`
	Name                string     `json:"name"`
	Description         string     `json:"description"`
	Enabled             bool       `json:"enabled"`
	CheckType           string     `json:"checkType"`
	MaintenanceOverride bool       `json:"maintenanceOverride"`
	Hostname            string     `json:"hostname"`
	ExpectedAddresses   []string   `json:"expectedAddresses"`
	TriggerCount        int        `json:"triggerCount"`
	ResultRetentionDays int        `json:"resultRetentionDays"`
	CheckHost           CheckHost  `json:"checkHost"`
	CheckGroup          CheckGroup `json:"checkGroup"`
}

type PingCheck struct {
	Id                  int        `json:"id"`
	Name                string     `json:"name"`
	Description         string     `json:"description"`
	Enabled             bool       `json:"enabled"`
	CheckType           string     `json:"checkType"`
	MaintenanceOverride bool       `json:"maintenanceOverride"`
	Hostname            string     `json:"hostname"`
	WarningRepsonseTime int        `json:"warningResponseTime"`
	Timeout             int        `json:"timeout"`
	TriggerCount        int        `json:"triggerCount"`
	ResultRetentionDays int        `json:"resultRetentionDays"`
	CheckHost           CheckHost  `json:"checkHost"`
	CheckGroup          CheckGroup `json:"checkGroup"`
}

type WebJourneyCheck struct {
	Id                  int              `json:"id"`
	Name                string           `json:"name"`
	Description         string           `json:"description"`
	Enabled             bool             `json:"enabled"`
	CheckType           string           `json:"checkType"`
	MaintenanceOverride bool             `json:"maintenanceOverride"`
	StartURL            string           `json:"startUrl"`
	TriggerCount        int              `json:"triggerCount"`
	ResultRetentionDays int              `json:"resultRetentionDays"`
	WindowHeight        int              `json:"windowHeight"`
	WindowWidth         int              `json:"windowWidth"`
	MonitorDomains      []MonitorDomain  `json:"monitorDomains"`
	Steps               []WebJourneyStep `json:"steps"`
	CheckHost           CheckHost        `json:"checkHost"`
	CheckGroup          CheckGroup       `json:"checkGroup"`
}

type MonitorDomain struct {
	Domain            string `json:"domain"`
	IncludeSubDomains bool   `json:"inlcudeSubDomains"`
}

type WebJourneyStep struct {
	Id                  int                           `json:"id"`
	Sequence            int                           `json:"sequence"`
	Type                string                        `json:"type"`
	Name                string                        `json:"name"`
	CommonId            int                           `json:"commonId"`
	WaitTime            int                           `json:"waitTime"`
	WarningPageLoadTime int                           `json:"warningPageLoadTime"`
	AlertPageLoadTime   int                           `json:"alertPageLoadTime"`
	PageChecks          []*WebJourneyPageCheck        `json:"pageChecks"`
	AlertSuppressions   []*WebJourneyAlertSuppression `json:"alertSuppressions"`
	Actions             []*WebJourneyAction           `json:"actions"`
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
	Id          int    `json:"id"`
	TextToFind  string `json:"textToFind"`
	ElementType string `json:"elementType"`
	State       string `json:"state"`
}

type PageCheckForElement struct {
	Id             int    `json:"id"`
	ElementId      string `json:"elementId"`
	ElementName    string `json:"elementName"`
	State          string `json:"state"`
	AttributeName  string `json:"attributeName"`
	AttributeValue string `json:"attributeValue"`
	ElementConent  string `json:"elementContent"`
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
	WarningResponseTime    int    `json:"warningRepsonseTime"`
	AlertResponseTime      int    `json:"alertResponseTime"`
	ResponseCode           int    `json:"responseCode"`
	AnyInfoResponse        bool   `json:"anyInfoResponseCode"`
	AnySuccessReponse      bool   `json:"anySuccessResponseCode"`
	AnyRedirectResponse    bool   `json:"anyRedirectResponseCode"`
	AnyClientErrorResponse bool   `json:"anyClientErrorResponseCode"`
	AnyServerErrorResponse bool   `json:"anyServerErrorResponseCode"`
}

type PageCheckConsoleLog struct {
	Id       int    `json:"id"`
	LogLevel string `json:"logLevel"`
	Message  string `json:"message"`
}

type WebJourneyAlertSuppression struct {
	Id                 int                `json:"id"`
	Description        string             `json:"description"`
	NetworkSuppression NetworkSuppression `json:"networkSuppression"`
	ConsoleSuppression ConsoleSuppression `json:"consoleSuppression"`
}

type NetworkSuppression struct {
	Id             int    `json:"id"`
	Url            string `json:"url"`
	Comparison     string `json:"comparison"`
	ResponseCode   int    `json:"responseCode"`
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
}

type WebJourneyClickAction struct {
	Xpath       string `json:"xpath"`
	SearchText  string `json:"searchText"`
	ElementType string `json:"elementType"`
}

type WebJourneyTextInputAction struct {
	Xpath       string `json:"xpath"`
	ElementId   string `json:"elementId"`
	ElementName string `json:"elementName"`
	InputText   string `json:"inputText"`
}

type WebJourneyPasswordInputAction struct {
	Xpath       string `json:"xpath"`
	ElementId   string `json:"elementId"`
	ElementName string `json:"elementName"`
	Password    string `json:"newInputPassword"`
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
	Xpath       string `json:"xpath"`
	SearchText  string `json:"searchText"`
	ElementType string `json:"elementType"`
}

type RequestHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (c *Client) SearchWebJourneyCommonSteps(search string) (*[]WebJourneyStep, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/commonWebJourneySteps/list?page=0&search=%s", c.HostURL, url.QueryEscape(search)), nil)

	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	steps := []WebJourneyStep{}

	if body != nil {
		err = json.Unmarshal(body, &steps)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &steps, nil
}

func (c *Client) GetUrlCheck(id string) (*URLCheck, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	check := Check{}
	urlCheck := URLCheck{}

	if body != nil {
		err = json.Unmarshal(body, &check)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	if check.CheckType != "URL" {
		err = errors.New("Existing check is not a URL Check.")
		return nil, err
	}

	err = json.Unmarshal(body, &urlCheck)
	if err != nil {
		return nil, err
	}

	return &urlCheck, nil
}

func (c *Client) GetSocketCheck(id string) (*SocketCheck, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	check := Check{}
	socketCheck := SocketCheck{}

	if body != nil {
		err = json.Unmarshal(body, &check)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	if check.CheckType != "SOCKET" {
		err = errors.New("Existing check is not a Socket Check.")
		return nil, err
	}

	err = json.Unmarshal(body, &socketCheck)
	if err != nil {
		return nil, err
	}

	return &socketCheck, nil
}

func (c *Client) GetDNSCheck(id string) (*DNSCheck, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	check := Check{}
	dnsCheck := DNSCheck{}

	if body != nil {
		err = json.Unmarshal(body, &check)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	if check.CheckType != "DNS" {
		err = errors.New("Existing check is not a DNS Check.")
		return nil, err
	}

	err = json.Unmarshal(body, &dnsCheck)
	if err != nil {
		return nil, err
	}

	return &dnsCheck, nil
}

func (c *Client) GetPingCheck(id string) (*PingCheck, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	check := Check{}
	pingCheck := PingCheck{}

	if body != nil {
		err = json.Unmarshal(body, &check)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	if check.CheckType != "PING" {
		err = errors.New("Existing check is not a Ping Check.")
		return nil, err
	}

	err = json.Unmarshal(body, &pingCheck)
	if err != nil {
		return nil, err
	}

	return &pingCheck, nil
}

func (c *Client) GetWebJourneyCheck(id string) (*WebJourneyCheck, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	check := Check{}
	pingCheck := WebJourneyCheck{}

	if body != nil {
		err = json.Unmarshal(body, &check)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	if check.CheckType != "WEB_JOURNEY" {
		err = errors.New("Existing check is not a Web Journey Check.")
		return nil, err
	}

	err = json.Unmarshal(body, &pingCheck)
	if err != nil {
		return nil, err
	}

	return &pingCheck, nil
}

func (c *Client) CreateUrlCheck(check URLCheck) (*URLCheck, error) {
	rb, err := json.Marshal(check)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/checks/add/url", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newCheck := URLCheck{}

	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	return &newCheck, nil
}

func (c *Client) CreateSocketCheck(check SocketCheck) (*SocketCheck, error) {
	rb, err := json.Marshal(check)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/checks/add/socket", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newCheck := SocketCheck{}

	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	return &newCheck, nil
}

func (c *Client) CreateDNSCheck(check DNSCheck) (*DNSCheck, error) {
	rb, err := json.Marshal(check)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/checks/add/dns", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newCheck := DNSCheck{}

	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	return &newCheck, nil
}

func (c *Client) CreatePingCheck(check PingCheck) (*PingCheck, error) {
	rb, err := json.Marshal(check)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/checks/add/ping", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newCheck := PingCheck{}

	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	return &newCheck, nil
}

func (c *Client) CreateWebJourneyCheck(check WebJourneyCheck, ctx context.Context) (*WebJourneyCheck, error) {
	rb, err := json.Marshal(check)
	if err != nil {
		return nil, err
	}

	tflog.Error(ctx, string(rb))

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/checks/add/webJourney", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newCheck := WebJourneyCheck{}

	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	return &newCheck, nil
}

func (c *Client) UpdateUrlCheck(check URLCheck) (*URLCheck, error) {
	rb, err := json.Marshal(check)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/checks/update/url", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedCheck := URLCheck{}
	err = json.Unmarshal(body, &updatedCheck)
	if err != nil {
		return nil, err
	}

	return &updatedCheck, nil
}

func (c *Client) UpdateSocketCheck(check SocketCheck) (*SocketCheck, error) {
	rb, err := json.Marshal(check)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/checks/update/socket", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedCheck := SocketCheck{}
	err = json.Unmarshal(body, &updatedCheck)
	if err != nil {
		return nil, err
	}

	return &updatedCheck, nil
}

func (c *Client) UpdateDNSCheck(check DNSCheck) (*DNSCheck, error) {
	rb, err := json.Marshal(check)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/checks/update/dns", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedCheck := DNSCheck{}
	err = json.Unmarshal(body, &updatedCheck)
	if err != nil {
		return nil, err
	}

	return &updatedCheck, nil
}

func (c *Client) UpdatePingCheck(check PingCheck) (*PingCheck, error) {
	rb, err := json.Marshal(check)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/checks/update/ping", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedCheck := PingCheck{}
	err = json.Unmarshal(body, &updatedCheck)
	if err != nil {
		return nil, err
	}

	return &updatedCheck, nil
}

func (c *Client) UpdateWebJourneyCheck(check WebJourneyCheck) (*WebJourneyCheck, error) {
	rb, err := json.Marshal(check)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/checks/update/webJourney", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedCheck := WebJourneyCheck{}
	err = json.Unmarshal(body, &updatedCheck)
	if err != nil {
		return nil, err
	}

	return &updatedCheck, nil
}

func (c *Client) DeleteCheck(checkId int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/checks/remove/%d", c.HostURL, checkId), nil)
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	if string(body) != "{\"success\":true}" {
		return errors.New(string(body))
	}

	return nil
}
