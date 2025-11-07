package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	apiKey string
	apiUrl string
)

var rootCmd = &cobra.Command{
	Use:   "naturedopes-cli",
	Short: "CLI tool for nature dopes Api",
	Long: `A command-line interface for interacting with the Nature Dopes API.

  Manage images, search for flora species, and work with API keys.

  Example usage:
    naturedopes-cli images list
    naturedopes-cli keys create --name "My Key"`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//global persistant flags
	rootCmd.PersistentFlags().StringVar(&apiUrl, "api-url", "http://localhost:8080", "API base Url")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key for auth")
}
