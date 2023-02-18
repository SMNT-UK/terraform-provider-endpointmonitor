package endpointmonitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type CheckHost struct {
	Id                  int     `json:"id"`
	Hostname            string  `json:"hostname"`
	Description         string  `json:"description"`
	Type                *string `json:"type"`
	Enabled             bool    `json:"enabled"`
	MaxWebJourneyChecks int     `json:"maxWebJourneyChecks"`
	SendCheckFiles      bool    `json:"sendCheckFiles"`
}

func (c *Client) SearchCheckHosts(search string) (*[]CheckHost, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hosts/list?page=0&search=%s", c.HostURL, url.QueryEscape(search)), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	checkHost := []CheckHost{}

	if body != nil {
		err = json.Unmarshal(body, &checkHost)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &checkHost, nil
}

func (c *Client) GetCheckHost(id string) (*CheckHost, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hosts/%s", c.HostURL, id), nil)
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

	return &checkHost, nil
}

func (c *Client) CreateCheckHost(checkHost CheckHost) (*CheckHost, error) {
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

	return &newCheckHost, nil
}

func (c *Client) UpdateCheckHost(checkHost CheckHost) (*CheckHost, error) {
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

	updatedCheckHost := CheckHost{}
	err = json.Unmarshal(body, &updatedCheckHost)
	if err != nil {
		return nil, err
	}

	return &updatedCheckHost, nil
}

func (c *Client) DeleteCheckHost(hostId int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/hosts/remove/%d", c.HostURL, hostId), nil)
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
