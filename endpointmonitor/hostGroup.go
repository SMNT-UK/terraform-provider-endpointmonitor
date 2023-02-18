package endpointmonitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type HostGroup struct {
	Id          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Enabled     bool        `json:"enabled"`
	Hosts       []CheckHost `json:"checkHosts"`
}

func (c *Client) SearchHostGroups(search string) (*[]HostGroup, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hostGroups/list?page=0&search=%s", c.HostURL, url.QueryEscape(search)), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	hostGroup := []HostGroup{}

	if body != nil {
		err = json.Unmarshal(body, &hostGroup)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &hostGroup, nil
}

func (c *Client) GetHostGroup(id string) (*HostGroup, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hostGroups/%s", c.HostURL, id), nil)
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

	return &hostGroup, nil
}

func (c *Client) CreateHostGroup(hostGroup HostGroup) (*HostGroup, error) {
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

	return &newHostGroup, nil
}

func (c *Client) UpdateHostGroup(hostGroup HostGroup) (*HostGroup, error) {
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

	updatedHostGroup := HostGroup{}
	err = json.Unmarshal(body, &updatedHostGroup)
	if err != nil {
		return nil, err
	}

	return &updatedHostGroup, nil
}

func (c *Client) DeleteHostGroup(hostGroupId int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/hostGroups/remove/%d", c.HostURL, hostGroupId), nil)
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
