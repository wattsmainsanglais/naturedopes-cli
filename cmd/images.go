package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/api"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/config"
)

var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "Get Images command",
}

var listImagesCmd = &cobra.Command{
	Use:   "list",
	Short: "Get list of images",
	Args:  cobra.ExactArgs(0),
	Run: func(command *cobra.Command, args []string) {
		baseUrl, _ := config.Get("api-url")
		key, _ := config.Get("api-key")
		client := api.NewClient(baseUrl, key)

		resp, err := client.ListImages()
		if err != nil {
			fmt.Errorf("could not retrieve images: %w", err)
			return
		}

		for _, image := range resp {
			fmt.Printf("name: %s, gps_long: %f, gps_lat: %f, image_path: %s\n", image.SpeciesName, image.GpsLong, image.GpsLat, image.ImagePath)
		}

	},
}

var getImageCmd = &cobra.Command{
	Use:   "get --id <id>",
	Short: "Get individual image",
	Args:  cobra.ExactArgs(1),
	Run: func(command *cobra.Command, args []string) {
		id := args[0]

	},
}
