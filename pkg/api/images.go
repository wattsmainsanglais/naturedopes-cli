package api

import (
	"encoding/json"
	"fmt"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/models"
)

func (c *Client) ListImages() ([]models.Image, error) {

	var images []models.Image

	resp, err := c.doRequest("GET", "/images")
	if err != nil {
		return nil, fmt.Errorf("could not retrieve images: %w", err)
	}

	err = json.Unmarshal(resp, &images)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall to json: %w", err)
	}
	return images, nil

}

func (c *Client) GetImage(id int) (*models.Image, error) {

	var images models.Image

	resp, err := c.doRequest("GET", fmt.Sprintf("/images/%d", id))
	if err != nil {
		return nil, fmt.Errorf("Could not obtain image: %w", err)
	}

	err = json.Unmarshal(resp, &images)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall to json: %w", err)
	}

	return &images, nil
}
