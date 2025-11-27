package api

import (
	"encoding/json"
	"fmt"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/models"
	"net/url"
	"strconv"
)

func (c *Client) ListImages() ([]models.Image, error) {

	var images []models.Image

	resp, err := c.doRequest("GET", "/images", nil)
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

	resp, err := c.doRequest("GET", fmt.Sprintf("/images/%d", id), nil)
	if err != nil {
		return nil, fmt.Errorf("could not obtain image: %w", err)
	}

	err = json.Unmarshal(resp, &images)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall to json: %w", err)
	}

	return &images, nil
}

func (c *Client) SearchImages(species string, userID int) ([]models.Image, error) {

	var images []models.Image

	path := "/images"
	params := url.Values{}
	if species != "" {
		params.Add("species_name", species)
	}
	if userID > 0 {
		params.Add("user_id", strconv.Itoa(userID))
	}

	if len(params) > 0 {
		path = path + "?" + params.Encode()
	}

	resp, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("could  not search images: %w", err)
	}

	err = json.Unmarshal(resp, &images)
	if err != nil {
		return nil, fmt.Errorf("could not read response from server: %w", err)
	}

	return images, nil

}
