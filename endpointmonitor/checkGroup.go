package endpointmonitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type CheckGroup struct {
	Id             int            `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	DashboardGroup DashboardGroup `json:"dashboardGroup"`
}

func (c *Client) SearchCheckGroups(search string) (*[]CheckGroup, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checkGroups/list?page=0&search=%s", c.HostURL, url.QueryEscape(search)), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	checkGroup := []CheckGroup{}

	if body != nil {
		err = json.Unmarshal(body, &checkGroup)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &checkGroup, nil
}

func (c *Client) GetCheckGroup(id string) (*CheckGroup, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checkGroups/%s", c.HostURL, id), nil)
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

	return &checkGroup, nil
}

func (c *Client) CreateCheckGroup(checkGroup CheckGroup) (*CheckGroup, error) {
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

	return &newCheckGroup, nil
}

func (c *Client) UpdateCheckGroup(checkGroup CheckGroup) (*CheckGroup, error) {
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

	updatedCheckGroup := CheckGroup{}
	err = json.Unmarshal(body, &updatedCheckGroup)
	if err != nil {
		return nil, err
	}

	return &updatedCheckGroup, nil
}

func (c *Client) DeleteCheckGroup(groupId int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/checkGroups/remove/%d", c.HostURL, groupId), nil)
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
