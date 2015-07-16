package gosquared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	baseUrl = "https://api.gosquared.com/tracking/v1"
)

type Client struct {
	APIKey     string
	SiteToken  string
	HTTPClient http.Client
}

func NewClient(apiKey string, siteToken string) *Client {
	return &Client{
		APIKey:    apiKey,
		SiteToken: siteToken,
	}
}

type EventRequest struct {
	PersonId string `json:"person_id,omitempty"`
	Event    `json:"event"`
}

type Event struct {
	Name string
	Data map[string]interface{}
}

func (c *Client) Event(name string, data map[string]interface{}, personId string) error {

	req := &EventRequest{
		PersonId: personId,
		Event: Event{
			Name: name,
			Data: data,
		},
	}

	resp, err := c.request("POST", "/event", req)
	if err != nil {
		return err
	}

	return c.respToError(resp)
}

func (c *Client) respToError(resp *http.Response) error {
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return fmt.Errorf("Non 200 reply from gosquared: %s", data)
}

func (c *Client) request(method, path string, payload interface{}) (*http.Response, error) {
	// serialize payload
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// construct url
	url := baseUrl + path + fmt.Sprintf("?api_key=%s&site_token=%s", url.QueryEscape(c.APIKey), url.QueryEscape(c.SiteToken))

	// new request
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// set length/content-type
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
		req.ContentLength = int64(len(body))
	}

	return c.HTTPClient.Do(req)
}
