package endpointmonitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type ProxyHost struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Hostname    string `json:"hostname"`
	Port        int    `json:"port"`
}

func (c *Client) SearchProxyHosts(search string) (*[]ProxyHost, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/proxies/list?page=0&search=%s", c.HostURL, url.QueryEscape(search)), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	proxyHosts := []ProxyHost{}

	if body != nil {
		err = json.Unmarshal(body, &proxyHosts)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &proxyHosts, nil
}

func (c *Client) GetProxyHost(id string) (*ProxyHost, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/proxies/%s", c.HostURL, id), nil)
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

	return &proxyHost, nil
}

func (c *Client) CreateProxyHost(proxyHost ProxyHost) (*ProxyHost, error) {
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

	return &newProxyHost, nil
}

func (c *Client) UpdateProxyHost(proxyHost ProxyHost) (*ProxyHost, error) {
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

	updatedProxyHost := ProxyHost{}
	err = json.Unmarshal(body, &updatedProxyHost)
	if err != nil {
		return nil, err
	}

	return &updatedProxyHost, nil
}

func (c *Client) DeleteProxyHost(proxyId int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/proxies/remove/%d", c.HostURL, proxyId), nil)
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
