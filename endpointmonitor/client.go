package endpointmonitor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const HostURL string = "http://localhost:19090"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	ApiKey     string
}

func NewClient(hostUrl string, apiKey *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		// Default Hashicups URL
		HostURL: hostUrl,
		ApiKey:  *apiKey,
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("x-epm-auth", c.ApiKey)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
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
