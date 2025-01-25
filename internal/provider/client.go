package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EndPointMonitorClient struct {
	HTTPClient *http.Client
	HostURL    string
	ApiKey     string
}

func NewEPMClient(hostUrl string, apiKey *string) (*EndPointMonitorClient, error) {
	c := EndPointMonitorClient{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		HostURL:    hostUrl,
		ApiKey:     *apiKey,
	}

	return &c, nil
}

func (c *EndPointMonitorClient) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("x-epm-auth", c.ApiKey)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, err
	}

	return body, err
}

func (c *EndPointMonitorClient) GetCheckGroup(id int32) (*CheckGroupModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checkGroups/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	checkGroup := CheckGroup{}

	if body != nil {
		err = json.Unmarshal(body, &checkGroup)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	checkGroupModel := mapToCheckGroupModel(checkGroup)

	return &checkGroupModel, nil
}

func (c *EndPointMonitorClient) GetCheckHost(id int32) (*CheckHostModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hosts/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	checkHost := CheckHost{}

	if body != nil {
		err = json.Unmarshal(body, &checkHost)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	checkHostModel := mapToCheckHostModel(checkHost)

	return &checkHostModel, nil
}

func (c *EndPointMonitorClient) GetDashboardGroup(id int32) (*DashboardGroupModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dashboardGroups/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	dashboardGroup := DashboardGroup{}

	if body != nil {
		err = json.Unmarshal(body, &dashboardGroup)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	dashboardGroupModel := mapToDashboardGroupModel(dashboardGroup)

	return &dashboardGroupModel, nil
}

func (c *EndPointMonitorClient) GetCertificateCheck(id int64) (*CertificateCheckModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	check := CertificateCheck{}

	if body != nil {
		err = json.Unmarshal(body, &check)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	checkModel := mapToCertificateCheckModel(check)

	return &checkModel, nil
}

func (c *EndPointMonitorClient) GetDnsCheck(id int64) (*DnsCheckModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	check := DnsCheck{}

	if body != nil {
		err = json.Unmarshal(body, &check)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	checkModel := mapToDnsCheckModel(check)

	return &checkModel, nil
}

func (c *EndPointMonitorClient) GetHostGroup(id int32) (*HostGroupModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hostGroups/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	hostGroup := HostGroup{}

	if body != nil {
		err = json.Unmarshal(body, &hostGroup)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	hostGroupModel := mapToHostGroupModel(hostGroup)

	return &hostGroupModel, nil
}

func (c *EndPointMonitorClient) GetMaintenancePeriod(id int32) (*MaintenancePeriodModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/maintenancePeriods/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	maintenancePereiod := MaintenancePeriod{}

	if body != nil {
		err = json.Unmarshal(body, &maintenancePereiod)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	maintenancePereiodModel := mapToMaintenancePeriodModel(maintenancePereiod)

	return &maintenancePereiodModel, nil
}

func (c *EndPointMonitorClient) GetPingCheck(id int64) (*PingCheckModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	check := PingCheck{}

	if body != nil {
		err = json.Unmarshal(body, &check)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	checkModel := mapToPingCheckModel(check)

	return &checkModel, nil
}

func (c *EndPointMonitorClient) GetProxyHost(id int32) (*ProxyHostModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/proxies/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	proxyHost := ProxyHost{}

	if body != nil {
		err = json.Unmarshal(body, &proxyHost)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	proxyHostModel := mapToProxyHostModel(proxyHost)

	return &proxyHostModel, nil
}

func (c *EndPointMonitorClient) GetSocketCheck(id int64) (*SocketCheckModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	check := SocketCheck{}

	if body != nil {
		err = json.Unmarshal(body, &check)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	checkModel := mapToSocketCheckModel(check)

	return &checkModel, nil
}

func (c *EndPointMonitorClient) GetUrlCheck(id int64) (*UrlCheckModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	check := UrlCheck{}

	if body != nil {
		err = json.Unmarshal(body, &check)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	checkModel := mapToUrlCheckModel(check)

	return &checkModel, nil
}

func (c *EndPointMonitorClient) GetAndroidJourneyCheck(id int64) (*AndroidJourneyCheckModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	c.HTTPClient.Timeout = time.Minute * 5

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	check := AndroidJourneyCheck{}

	if body != nil {
		err = json.Unmarshal(body, &check)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	checkModel := mapToAndroidJourneyCheckModel(check)

	return &checkModel, nil
}

func (c *EndPointMonitorClient) GetWebJourneyCheck(id int64) (*WebJourneyCheckModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	check := WebJourneyCheck{}

	if body != nil {
		err = json.Unmarshal(body, &check)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	checkModel := mapToWebJourneyCheckModel(check)

	return &checkModel, nil
}

func (c *EndPointMonitorClient) GetCommonAndroidJourneyStep(id int64) (*AndroidJourneyCommonStepModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/commonSteps/android/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	step := AndroidJourneyCommonStep{}

	if body != nil {
		err = json.Unmarshal(body, &step)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	stepModel := mapToAndroidJourneyCommonStepModel(step)

	return &stepModel, nil
}

func (c *EndPointMonitorClient) GetCommonWebJourneyStep(id int64) (*WebJourneyCommonStepModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/commonSteps/web/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	step := WebJourneyCommonStep{}

	if body != nil {
		err = json.Unmarshal(body, &step)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	stepModel := mapToWebJourneyCommonStepModel(step)

	return &stepModel, nil
}

func (c *EndPointMonitorClient) CreateCheckGroup(checkGroupModel CheckGroupModel, ctx context.Context) (*CheckGroupModel, error) {
	checkGroup := mapToCheckGroup(checkGroupModel)

	rb, err := json.Marshal(checkGroup)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/checkGroups/add", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newCheckGroup := CheckGroup{}

	err = json.Unmarshal(body, &newCheckGroup)
	if err != nil {
		return nil, err
	}

	newCheckGroupModel := mapToCheckGroupModel(newCheckGroup)

	return &newCheckGroupModel, nil
}

func (c *EndPointMonitorClient) CreateCheckHost(checkHostModel CheckHostModel, ctx context.Context) (*CheckHostModel, error) {
	checkHost := mapToCheckHost(checkHostModel)

	rb, err := json.Marshal(checkHost)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/hosts/add", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newCheckHost := CheckHost{}

	err = json.Unmarshal(body, &newCheckHost)
	if err != nil {
		return nil, err
	}

	newCheckHostModel := mapToCheckHostModel(newCheckHost)

	return &newCheckHostModel, nil
}

func (c *EndPointMonitorClient) CreateCertificateCheck(checkModel CertificateCheckModel, ctx context.Context) (*CertificateCheckModel, error) {
	check := mapToCertificateCheck(checkModel)

	rb, err := json.Marshal(check)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/checks/add/certificate", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newCheck := CertificateCheck{}

	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	newCheckModel := mapToCertificateCheckModel(newCheck)

	return &newCheckModel, nil
}

func (c *EndPointMonitorClient) CreateDashboardGroup(dashboardGroupModel DashboardGroupModel, ctx context.Context) (*DashboardGroupModel, error) {
	dashboardGroup := mapToDashboardGroup(dashboardGroupModel)

	rb, err := json.Marshal(dashboardGroup)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/dashboardGroups/add", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newDashboardGroup := DashboardGroup{}

	err = json.Unmarshal(body, &newDashboardGroup)
	if err != nil {
		return nil, err
	}

	newDashboardGroupModel := mapToDashboardGroupModel(newDashboardGroup)

	return &newDashboardGroupModel, nil
}

func (c *EndPointMonitorClient) CreateDnsCheck(checkModel DnsCheckModel, ctx context.Context) (*DnsCheckModel, error) {
	check := mapToDnsCheck(checkModel)

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

	newCheck := DnsCheck{}

	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	newCheckModel := mapToDnsCheckModel(newCheck)

	return &newCheckModel, nil
}

func (c *EndPointMonitorClient) CreateHostGroup(hostGroupModel HostGroupModel, ctx context.Context) (*HostGroupModel, error) {
	hostGroup := mapToHostGroup(hostGroupModel)

	rb, err := json.Marshal(hostGroup)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/hostGroups/add", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newHostGroup := HostGroup{}

	err = json.Unmarshal(body, &newHostGroup)
	if err != nil {
		return nil, err
	}

	newHostGroupModel := mapToHostGroupModel(newHostGroup)

	return &newHostGroupModel, nil
}

func (c *EndPointMonitorClient) CreateMaintenancePeriod(maintenancePeriodModel MaintenancePeriodModel, ctx context.Context) (*MaintenancePeriodModel, error) {
	maintenancePeriod := mapToMaintenancePeriod(maintenancePeriodModel)

	rb, err := json.Marshal(maintenancePeriod)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/maintenancePeriods/add", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newMaintenancePeriod := MaintenancePeriod{}

	err = json.Unmarshal(body, &newMaintenancePeriod)
	if err != nil {
		return nil, err
	}

	newMaintenancePeriodModel := mapToMaintenancePeriodModel(newMaintenancePeriod)

	return &newMaintenancePeriodModel, nil
}

func (c *EndPointMonitorClient) CreatePingCheck(checkModel PingCheckModel, ctx context.Context) (*PingCheckModel, error) {
	check := mapToPingCheck(checkModel)

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

	newCheckModel := mapToPingCheckModel(newCheck)

	return &newCheckModel, nil
}

func (c *EndPointMonitorClient) CreateProxyHost(proxyHostModel ProxyHostModel, ctx context.Context) (*ProxyHostModel, error) {
	proxyHost := mapToProxyHost(proxyHostModel)

	rb, err := json.Marshal(proxyHost)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/proxies/add", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newProxyHost := ProxyHost{}

	err = json.Unmarshal(body, &newProxyHost)
	if err != nil {
		return nil, err
	}

	newProxyHostModel := mapToProxyHostModel(newProxyHost)

	return &newProxyHostModel, nil
}

func (c *EndPointMonitorClient) CreateSocketCheck(checkModel SocketCheckModel, ctx context.Context) (*SocketCheckModel, error) {
	check := mapToSocketCheck(checkModel)

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

	newCheckModel := mapToSocketCheckModel(newCheck)

	return &newCheckModel, nil
}

func (c *EndPointMonitorClient) CreateUrlCheck(checkModel UrlCheckModel, ctx context.Context) (*UrlCheckModel, error) {
	check := mapToUrlCheck(checkModel)

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

	newCheck := UrlCheck{}

	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	newCheckModel := mapToUrlCheckModel(newCheck)

	return &newCheckModel, nil
}

func (c *EndPointMonitorClient) CreateAndroidJourneyCheck(checkModel AndroidJourneyCheckModel, ctx context.Context) (*AndroidJourneyCheckModel, error) {
	check := mapToAndroidJourneyCheck(checkModel)

	rb, err := json.Marshal(check)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/checks/add/androidJourney", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	c.HTTPClient.Timeout = time.Minute * 5

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newCheck := AndroidJourneyCheck{}

	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	newCheckModel := mapToAndroidJourneyCheckModel(newCheck)

	return &newCheckModel, nil
}

func (c *EndPointMonitorClient) CreateWebJourneyCheck(checkModel WebJourneyCheckModel, ctx context.Context) (*WebJourneyCheckModel, error) {
	check := mapToWebJourneyCheck(checkModel)

	rb, err := json.Marshal(check)
	if err != nil {
		return nil, err
	}

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

	newCheckModel := mapToWebJourneyCheckModel(newCheck)

	return &newCheckModel, nil
}

func (c *EndPointMonitorClient) CreateAndroidJourneyCommonStep(stepModel AndroidJourneyCommonStepModel, ctx context.Context) (*AndroidJourneyCommonStepModel, error) {
	step := mapToAndroidJourneyCommonStep(stepModel)

	rb, err := json.Marshal(step)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/checks/commonSteps/android/add", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newStep := AndroidJourneyCommonStep{}

	err = json.Unmarshal(body, &newStep)
	if err != nil {
		return nil, err
	}

	newStepModel := mapToAndroidJourneyCommonStepModel(newStep)

	return &newStepModel, nil
}

func (c *EndPointMonitorClient) CreateWebJourneyCommonStep(stepModel WebJourneyCommonStepModel, ctx context.Context) (*WebJourneyCommonStepModel, error) {
	step := mapToWebJourneyCommonStep(stepModel)

	rb, err := json.Marshal(step)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/checks/commonSteps/web/add", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newStep := WebJourneyCommonStep{}

	err = json.Unmarshal(body, &newStep)
	if err != nil {
		return nil, err
	}

	newStepModel := mapToWebJourneyCommonStepModel(newStep)

	return &newStepModel, nil
}

func (c *EndPointMonitorClient) UpdateCheckGroup(checkGroupModel CheckGroupModel, ctx context.Context) (*CheckGroupModel, error) {
	checkGroup := mapToCheckGroup(checkGroupModel)

	rb, err := json.Marshal(checkGroup)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/checkGroups/update", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newCheckGroup := CheckGroup{}

	err = json.Unmarshal(body, &newCheckGroup)
	if err != nil {
		return nil, err
	}

	newCheckGroupModel := mapToCheckGroupModel(newCheckGroup)

	return &newCheckGroupModel, nil
}

func (c *EndPointMonitorClient) UpdateCheckHost(checkHostModel CheckHostModel, ctx context.Context) (*CheckHostModel, error) {
	checkHost := mapToCheckHost(checkHostModel)

	rb, err := json.Marshal(checkHost)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hosts/update", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newCheckHost := CheckHost{}

	err = json.Unmarshal(body, &newCheckHost)
	if err != nil {
		return nil, err
	}

	newCheckHostModel := mapToCheckHostModel(newCheckHost)

	return &newCheckHostModel, nil
}

func (c *EndPointMonitorClient) UpdateCertificateCheck(checkModel CertificateCheckModel, ctx context.Context) (*CertificateCheckModel, error) {
	check := mapToCertificateCheck(checkModel)

	rb, err := json.Marshal(check)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/checks/update/certificate", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newCheck := CertificateCheck{}

	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	newCheckModel := mapToCertificateCheckModel(newCheck)

	return &newCheckModel, nil
}

func (c *EndPointMonitorClient) UpdateDashboardGroup(dashboardGroupModel DashboardGroupModel, ctx context.Context) (*DashboardGroupModel, error) {
	dashboardGroup := mapToDashboardGroup(dashboardGroupModel)

	rb, err := json.Marshal(dashboardGroup)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/dashboardGroups/update", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newDashboardGroup := DashboardGroup{}

	err = json.Unmarshal(body, &newDashboardGroup)
	if err != nil {
		return nil, err
	}

	newDashboardGroupModel := mapToDashboardGroupModel(newDashboardGroup)

	return &newDashboardGroupModel, nil
}

func (c *EndPointMonitorClient) UpdateDnsCheck(checkModel DnsCheckModel, ctx context.Context) (*DnsCheckModel, error) {
	check := mapToDnsCheck(checkModel)

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

	newCheck := DnsCheck{}

	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	newCheckModel := mapToDnsCheckModel(newCheck)

	return &newCheckModel, nil
}

func (c *EndPointMonitorClient) UpdateHostGroup(hostGroupModel HostGroupModel, ctx context.Context) (*HostGroupModel, error) {
	hostGroup := mapToHostGroup(hostGroupModel)

	rb, err := json.Marshal(hostGroup)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/hostGroups/update", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newHostGroup := HostGroup{}

	err = json.Unmarshal(body, &newHostGroup)
	if err != nil {
		return nil, err
	}

	newHostGroupModel := mapToHostGroupModel(newHostGroup)

	return &newHostGroupModel, nil
}

func (c *EndPointMonitorClient) UpdateMaintenancePeriod(maintenancePeriodModel MaintenancePeriodModel, ctx context.Context) (*MaintenancePeriodModel, error) {
	maintenancePeriod := mapToMaintenancePeriod(maintenancePeriodModel)

	rb, err := json.Marshal(maintenancePeriod)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/maintenancePeriods/update", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newMaintenancePeriod := MaintenancePeriod{}

	err = json.Unmarshal(body, &newMaintenancePeriod)
	if err != nil {
		return nil, err
	}

	newMaintenancePeriodModel := mapToMaintenancePeriodModel(newMaintenancePeriod)

	return &newMaintenancePeriodModel, nil
}

func (c *EndPointMonitorClient) UpdatePingCheck(checkModel PingCheckModel, ctx context.Context) (*PingCheckModel, error) {
	check := mapToPingCheck(checkModel)

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

	newCheck := PingCheck{}

	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	newCheckModel := mapToPingCheckModel(newCheck)

	return &newCheckModel, nil
}

func (c *EndPointMonitorClient) UpdateProxyHost(proxyHostModel ProxyHostModel, ctx context.Context) (*ProxyHostModel, error) {
	proxyHost := mapToProxyHost(proxyHostModel)

	rb, err := json.Marshal(proxyHost)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/proxies/update", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newProxyHost := ProxyHost{}

	err = json.Unmarshal(body, &newProxyHost)
	if err != nil {
		return nil, err
	}

	newProxyHostModel := mapToProxyHostModel(newProxyHost)

	return &newProxyHostModel, nil
}

func (c *EndPointMonitorClient) UpdateSocketCheck(checkModel SocketCheckModel, ctx context.Context) (*SocketCheckModel, error) {
	check := mapToSocketCheck(checkModel)

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

	newCheck := SocketCheck{}

	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	newCheckModel := mapToSocketCheckModel(newCheck)

	return &newCheckModel, nil
}

func (c *EndPointMonitorClient) UpdateUrlCheck(checkModel UrlCheckModel, ctx context.Context) (*UrlCheckModel, error) {
	check := mapToUrlCheck(checkModel)

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

	newCheck := UrlCheck{}

	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	newCheckModel := mapToUrlCheckModel(newCheck)

	return &newCheckModel, nil
}

func (c *EndPointMonitorClient) UpdateAndroidJourneyCheck(checkModel AndroidJourneyCheckModel, ctx context.Context) (*AndroidJourneyCheckModel, error) {
	check := mapToAndroidJourneyCheck(checkModel)

	rb, err := json.Marshal(check)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/checks/update/androidJourney", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	c.HTTPClient.Timeout = time.Minute * 5

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newCheck := AndroidJourneyCheck{}

	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	newCheckModel := mapToAndroidJourneyCheckModel(newCheck)

	return &newCheckModel, nil
}

func (c *EndPointMonitorClient) UpdateWebJourneyCheck(checkModel WebJourneyCheckModel, ctx context.Context) (*WebJourneyCheckModel, error) {
	check := mapToWebJourneyCheck(checkModel)

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

	newCheck := WebJourneyCheck{}

	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	newCheckModel := mapToWebJourneyCheckModel(newCheck)

	return &newCheckModel, nil
}

func (c *EndPointMonitorClient) UpdateAndroidJourneyCommonStep(stepModel AndroidJourneyCommonStepModel, ctx context.Context) (*AndroidJourneyCommonStepModel, error) {
	step := mapToAndroidJourneyCommonStep(stepModel)

	rb, err := json.Marshal(step)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/checks/commonSteps/android/update", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newStep := AndroidJourneyCommonStep{}

	err = json.Unmarshal(body, &newStep)
	if err != nil {
		return nil, err
	}

	newStepModel := mapToAndroidJourneyCommonStepModel(newStep)

	return &newStepModel, nil
}

func (c *EndPointMonitorClient) UpdateWebJourneyCommonStep(stepModel WebJourneyCommonStepModel, ctx context.Context) (*WebJourneyCommonStepModel, error) {
	step := mapToWebJourneyCommonStep(stepModel)

	rb, err := json.Marshal(step)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/checks/commonSteps/web/update", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newStep := WebJourneyCommonStep{}

	err = json.Unmarshal(body, &newStep)
	if err != nil {
		return nil, err
	}

	newStepModel := mapToWebJourneyCommonStepModel(newStep)

	return &newStepModel, nil
}

func (c *EndPointMonitorClient) DeleteCheck(id int64) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/checks/remove/%d", c.HostURL, id), nil)
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

func (c *EndPointMonitorClient) DeleteCheckGroup(id int32) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/checkGroups/remove/%d", c.HostURL, id), nil)
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

func (c *EndPointMonitorClient) DeleteCheckHost(id int32) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/hosts/remove/%d", c.HostURL, id), nil)
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

func (c *EndPointMonitorClient) DeleteAndroidCommonStep(id int64) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/checks/commonSteps/android/remove/%d", c.HostURL, id), nil)
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

func (c *EndPointMonitorClient) DeleteWebCommonStep(id int64) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/checks/commonSteps/web/remove/%d", c.HostURL, id), nil)
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

func (c *EndPointMonitorClient) DeleteDashboardGroup(id int32) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/dashboardGroups/remove/%d", c.HostURL, id), nil)
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

func (c *EndPointMonitorClient) DeleteHostGroup(id int32) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/hostGroups/remove/%d", c.HostURL, id), nil)
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

func (c *EndPointMonitorClient) DeleteMaintenancePeriod(id int32) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/maintenancePeriods/remove/%d", c.HostURL, id), nil)
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

func (c *EndPointMonitorClient) DeleteProxyHost(id int32) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/proxies/remove/%d", c.HostURL, id), nil)
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

func (c *EndPointMonitorClient) SearchCheckGroups(search string) ([]types.Int32, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checkGroups/list?page=0&search=%s", c.HostURL, url.QueryEscape(search)), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	checkGroups := []CheckGroup{}

	if body != nil {
		err = json.Unmarshal(body, &checkGroups)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	ids := make([]types.Int32, 0, len(checkGroups))

	for _, check := range checkGroups {
		ids = append(ids, types.Int32Value(int32(check.Id)))
	}

	return ids, nil
}

func (c *EndPointMonitorClient) SearchCheckHosts(search string) ([]types.Int32, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hosts/list?page=0&search=%s", c.HostURL, url.QueryEscape(search)), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	checkHosts := []CheckHost{}

	if body != nil {
		err = json.Unmarshal(body, &checkHosts)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	ids := make([]types.Int32, 0, len(checkHosts))

	for _, check := range checkHosts {
		ids = append(ids, types.Int32Value(int32(check.Id)))
	}

	return ids, nil
}

func (c *EndPointMonitorClient) SearchChecks(search string) ([]types.Int64, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/list?page=0&search=%s", c.HostURL, url.QueryEscape(search)), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	checks := []Check{}

	if body != nil {
		err = json.Unmarshal(body, &checks)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	ids := make([]types.Int64, 0, len(checks))

	for _, check := range checks {
		ids = append(ids, types.Int64Value(int64(check.Id)))
	}

	return ids, nil
}

func (c *EndPointMonitorClient) SearchDashboardGroups(search string) ([]types.Int32, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dashboardGroups/list?page=0&search=%s", c.HostURL, url.QueryEscape(search)), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	dashboardGroups := []DashboardGroup{}

	if body != nil {
		err = json.Unmarshal(body, &dashboardGroups)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	ids := make([]types.Int32, 0, len(dashboardGroups))

	for _, check := range dashboardGroups {
		ids = append(ids, types.Int32Value(int32(check.Id)))
	}

	return ids, nil
}

func (c *EndPointMonitorClient) SearchHostGroups(search string) ([]types.Int32, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hostGroups/list?page=0&search=%s", c.HostURL, url.QueryEscape(search)), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	hostGroups := []HostGroup{}

	if body != nil {
		err = json.Unmarshal(body, &hostGroups)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	ids := make([]types.Int32, 0, len(hostGroups))

	for _, check := range hostGroups {
		ids = append(ids, types.Int32Value(int32(check.Id)))
	}

	return ids, nil
}

func (c *EndPointMonitorClient) SearchMaintenancePeriods(search string) ([]types.Int32, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/maintenancePeriods/list?page=0&search=%s", c.HostURL, url.QueryEscape(search)), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	maintenancePeriod := []MaintenancePeriod{}

	if body != nil {
		err = json.Unmarshal(body, &maintenancePeriod)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	ids := make([]types.Int32, 0, len(maintenancePeriod))

	for _, check := range maintenancePeriod {
		ids = append(ids, types.Int32Value(int32(check.Id)))
	}

	return ids, nil
}

func (c *EndPointMonitorClient) SearchProxyHosts(search string) ([]types.Int32, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/proxies/list?page=0&search=%s", c.HostURL, url.QueryEscape(search)), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	proxyHost := []ProxyHost{}

	if body != nil {
		err = json.Unmarshal(body, &proxyHost)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	ids := make([]types.Int32, 0, len(proxyHost))

	for _, check := range proxyHost {
		ids = append(ids, types.Int32Value(int32(check.Id)))
	}

	return ids, nil
}

func (c *EndPointMonitorClient) SearchAndroidJoureyCommonSteps(search string) ([]types.Int32, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/commonSteps/android/list?page=0&search=%s", c.HostURL, url.QueryEscape(search)), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	commonStep := []AndroidJourneyCommonStep{}

	if body != nil {
		err = json.Unmarshal(body, &commonStep)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	ids := make([]types.Int32, 0, len(commonStep))

	for _, check := range commonStep {
		ids = append(ids, types.Int32Value(int32(check.Id)))
	}

	return ids, nil
}

func (c *EndPointMonitorClient) SearchWebJoureyCommonSteps(search string) ([]types.Int32, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checks/commonSteps/web/list?page=0&search=%s", c.HostURL, url.QueryEscape(search)), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	commonStep := []WebJourneyCommonStep{}

	if body != nil {
		err = json.Unmarshal(body, &commonStep)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	ids := make([]types.Int32, 0, len(commonStep))

	for _, check := range commonStep {
		ids = append(ids, types.Int32Value(int32(check.Id)))
	}

	return ids, nil
}
