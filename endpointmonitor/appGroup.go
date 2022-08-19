package endpointmonitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type AppGroup struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *Client) SearchAppGroups(search string) (*[]AppGroup, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/appGroups/list?page=0&search=%s", c.HostURL, search), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	appGroups := []AppGroup{}

	if body != nil {
		err = json.Unmarshal(body, &appGroups)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &appGroups, nil
}

func (c *Client) GetAppGroup(id string) (*AppGroup, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/appGroups/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	appGroup := AppGroup{}

	if body != nil {
		err = json.Unmarshal(body, &appGroup)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &appGroup, nil
}

func (c *Client) CreateAppGroup(appGroup AppGroup) (*AppGroup, error) {
	rb, err := json.Marshal(appGroup)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/appGroups/add", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newAppGroup := AppGroup{}

	err = json.Unmarshal(body, &newAppGroup)
	if err != nil {
		return nil, err
	}

	return &newAppGroup, nil
}

func (c *Client) UpdateAppGroup(appGroup AppGroup) (*AppGroup, error) {
	rb, err := json.Marshal(appGroup)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/appGroups/update", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedAppGroup := AppGroup{}
	err = json.Unmarshal(body, &updatedAppGroup)
	if err != nil {
		return nil, err
	}

	return &updatedAppGroup, nil
}

func (c *Client) DeleteAppGroup(groupId int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/appGroups/remove/%d", c.HostURL, groupId), nil)
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
