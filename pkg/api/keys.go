package api

import (
	"encoding/json"
	"fmt"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/models"
)

func (client *Client) GenerateKey(name string) (*models.ApiKey, error) {
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

	resp, err := client.doRequest("POST", "/api/keys", jsonData)
	if err != nil {
		return nil, fmt.Errorf("could not create api keys from naturedopesApi: %w", err)
	}

	err = json.Unmarshal(resp, &apiKey)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response: %w", err)
	}

	return &apiKey, nil

}

func (client *Client) ListKeys() ([]models.ApiKey, error) {
	var apiKeys []models.ApiKey

	resp, err := client.doRequest("GET", "/api/keys", nil)
	if err != nil {
		return nil, fmt.Errorf("could not get apikeys: %w", err)
	}

	err = json.Unmarshal(resp, &apiKeys)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall json: %w", err)
	}

	return apiKeys, nil

}

func (client *Client) RevokeKey(id int) error {
	_, err := client.doRequest("DELETE", fmt.Sprintf("/api/keys/%d", id), nil)
	if err != nil {
		return fmt.Errorf("could not delete api-key: %w", err)
	}

	return nil
}
