package api

import (
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseUrl    string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(BaseUrl string, APIKey string) *Client {
	return &Client{
		BaseUrl:    BaseUrl,
		APIKey:     APIKey,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) doRequest(method string, path string) ([]byte, error) {

	url := c.BaseUrl + path

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create http request err: %w", err)
	}

	if c.APIKey != "" {
		req.Header.Set("X-API-Key", c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request err: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response, err: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("returned status code: %d , message: %s ", resp.StatusCode, string(body))
	}

	return body, nil

}
