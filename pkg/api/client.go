package api

import (
	"bytes"
	"context"
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

func (c *Client) doRequest(ctx context.Context, method string, path string, body []byte) ([]byte, error) {

	url := c.BaseUrl + path
	var reqBody io.Reader = nil
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("could not create http request err: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.APIKey != "" {
		req.Header.Set("X-API-Key", c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("request cancelled, %w", ctx.Err())
		}
		return nil, fmt.Errorf("could not send request err: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response, err: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("returned status code: %d , message: %s ", resp.StatusCode, string(respBody))
	}

	return respBody, nil

}
