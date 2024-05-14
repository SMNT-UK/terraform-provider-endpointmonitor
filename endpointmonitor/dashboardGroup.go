package endpointmonitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type DashboardGroup struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *Client) SearchDashboardGroups(search string) (*[]DashboardGroup, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dashboardGroups/list?page=0&search=%s", c.HostURL, search), nil)
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

	return &dashboardGroups, nil
}

func (c *Client) GetDashboardGroup(id string) (*DashboardGroup, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dashboardGroups/%s", c.HostURL, id), nil)
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

	return &dashboardGroup, nil
}

func (c *Client) CreateDashboardGroup(dashboardGroup DashboardGroup) (*DashboardGroup, error) {
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

	return &newDashboardGroup, nil
}

func (c *Client) UpdateDashboardGroup(dashboardGroup DashboardGroup) (*DashboardGroup, error) {
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

	updatedDashboardGroup := DashboardGroup{}
	err = json.Unmarshal(body, &updatedDashboardGroup)
	if err != nil {
		return nil, err
	}

	return &updatedDashboardGroup, nil
}

func (c *Client) DeleteDashboardGroup(groupId int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/dashboardGroups/remove/%d", c.HostURL, groupId), nil)
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
