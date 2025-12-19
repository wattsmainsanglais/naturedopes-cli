package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/api"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/config"
	"strconv"
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

		if !checkApiKey(key) {
			return
		}

		client := api.NewClient(baseUrl, key)

		resp, err := client.ListImages()
		if err != nil {
			fmt.Printf("could not retrieve images: %v\n", err)
			return
		}

		for _, image := range resp {
			fmt.Printf("name: %s, gps_long: %f, gps_lat: %f, image_path: %s\n", image.SpeciesName, image.GpsLong, image.GpsLat, image.ImagePath)
		}

	},
}

var getImageCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get individual image",
	Args:  cobra.ExactArgs(1),
	Run: func(command *cobra.Command, args []string) {
		id := args[0]
		integer, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("Error, invalid ID, please check you've supplied an integer as argument: %v\n", err)
			return
		}

		baseUrl, _ := config.Get("api-url")
		key, _ := config.Get("api-key")

		if !checkApiKey(key) {
			return
		}
		client := api.NewClient(baseUrl, key)

		image, err := client.GetImage(integer)
		if err != nil {
			fmt.Printf("could not retrieve image data: %v\n", err)
			return
		}

		fmt.Printf("id:%d name: %s, gps_long: %f, gps_lat: %f, image_path: %s, user_id: %d\n", image.ID, image.SpeciesName, image.GpsLong, image.GpsLat, image.ImagePath, image.UserID)

	},
}

var searchImagesCmd = &cobra.Command{
	Use:   "search <species_name> <id>",
	Short: "Enter the name of the species that you'd like to search for, you can optionally add the user Id as a filter, enter 0 for this argumnet if you don't ",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		id := args[1]

		idInt, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("could not convert to integer: %s\n ", err)
			return
		}

		baseUrl, _ := config.Get("api-url")
		key, _ := config.Get("api-key")

		if !checkApiKey(key) {
			return
		}
		client := api.NewClient(baseUrl, key)

		images, err := client.SearchImages(name, idInt)
		if err != nil {
			fmt.Printf("could not return images: %s\n", err)
			return
		}

		for _, i := range images {
			fmt.Printf("id:%d species_name: %s, gps_long: %f, gps_lat: %f, image_path: %s user_id: %d\n", i.ID, i.SpeciesName, i.GpsLong, i.GpsLat, i.ImagePath, i.UserID)
		}

	},
}

func init() {
	rootCmd.AddCommand(imagesCmd)
	imagesCmd.AddCommand(listImagesCmd)
	imagesCmd.AddCommand(getImageCmd)
	imagesCmd.AddCommand(searchImagesCmd)

}
