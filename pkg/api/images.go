package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/wattsmainsanglais/naturedopes-cli/pkg/models"
)

func (c *Client) ListImages(ctx context.Context) ([]models.Image, error) {

	var images []models.Image

	resp, err := c.doRequest(ctx, "GET", "/images", nil)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve images: %w", err)
	}

	err = json.Unmarshal(resp, &images)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall to json: %w", err)
	}
	return images, nil

}

func (c *Client) GetImage(ctx context.Context, id int) (*models.Image, error) {

	var images models.Image

	resp, err := c.doRequest(ctx, "GET", fmt.Sprintf("/images/%d", id), nil)
	if err != nil {
		return nil, fmt.Errorf("could not obtain image: %w", err)
	}

	err = json.Unmarshal(resp, &images)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall to json: %w", err)
	}

	return &images, nil
}

func (c *Client) SearchImages(ctx context.Context, species string, userID int) ([]models.Image, error) {

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

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("could  not search images: %w", err)
	}

	err = json.Unmarshal(resp, &images)
	if err != nil {
		return nil, fmt.Errorf("could not read response from server: %w", err)
	}

	return images, nil

}
