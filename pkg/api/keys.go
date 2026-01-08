package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wattsmainsanglais/naturedopes-cli/pkg/models"
)

func (client *Client) GenerateKey(ctx context.Context, name string) (*models.ApiKey, error) {
	var apiKey models.ApiKey

	requestBody := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("could not create jsonData: %w", err)
	}

	resp, err := client.doRequest(ctx, "POST", "/api/keys", jsonData)
	if err != nil {
		return nil, fmt.Errorf("could not create api keys from naturedopesApi: %w", err)
	}

	err = json.Unmarshal(resp, &apiKey)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response: %w", err)
	}

	return &apiKey, nil

}

func (client *Client) ListKeys(ctx context.Context) ([]models.ApiKey, error) {
	var apiKeys []models.ApiKey

	resp, err := client.doRequest(ctx, "GET", "/api/keys/list", nil)
	if err != nil {
		return nil, fmt.Errorf("could not get apikeys: %w", err)
	}

	err = json.Unmarshal(resp, &apiKeys)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall json: %w", err)
	}

	return apiKeys, nil

}

func (client *Client) GetKeyInfo(ctx context.Context, key string) (*models.ApiKey, error) {
	var apiKey models.ApiKey

	requestBody := struct {
		ApiKey string `json:"api-key"`
	}{
		ApiKey: key,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("could not create jsonData: %w", err)
	}

	resp, err := client.doRequest(ctx, "GET", "/api/keys/get", jsonData)
	if err != nil {
		return nil, fmt.Errorf("could not get api key; %w", err)
	}

	err = json.Unmarshal(resp, &apiKey)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response: %w", err)
	}

	return &apiKey, nil
}

func (client *Client) RevokeKey(ctx context.Context) error {
	_, err := client.doRequest(ctx, "DELETE", "/api/keys", nil)
	if err != nil {
		return fmt.Errorf("could not delete api-key: %w", err)
	}

	return nil
}
