package endpointmonitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type MaintenancePeriod struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	DayOfWeek   string `json:"dayOfWeek"`
	Checks      []int  `json:"checks"`
	CheckGroups []int  `json:"checkGroups"`
	AppGroups   []int  `json:"appGroups"`
}

func (c *Client) SearchMaintenancePeriods(search string) (*[]MaintenancePeriod, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/maintenancePeriods/list?page=0&search=%s", c.HostURL, url.QueryEscape(search)), nil)

	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	maintenancePeriods := []MaintenancePeriod{}

	if body != nil {
		err = json.Unmarshal(body, &maintenancePeriods)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &maintenancePeriods, nil
}

func (c *Client) GetMaintenancePeriod(id string) (*MaintenancePeriod, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/maintenancePeriods/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	maintenancePeriod := MaintenancePeriod{}

	if body != nil {
		err = json.Unmarshal(body, &maintenancePeriod)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &maintenancePeriod, nil
}

func (c *Client) CreateMaintenancePeriod(maintenancePeriod MaintenancePeriod) (*MaintenancePeriod, error) {
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

	return &newMaintenancePeriod, nil
}

func (c *Client) UpdateMaintenancePeriod(maintenancePeriod MaintenancePeriod) (*MaintenancePeriod, error) {
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

	updateMaintenancePeriod := MaintenancePeriod{}
	err = json.Unmarshal(body, &updateMaintenancePeriod)
	if err != nil {
		return nil, err
	}

	return &updateMaintenancePeriod, nil
}

func (c *Client) DeleteMaintenancePeriod(maintenancePeriodId int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/maintenancePeriods/remove/%d", c.HostURL, maintenancePeriodId), nil)
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
