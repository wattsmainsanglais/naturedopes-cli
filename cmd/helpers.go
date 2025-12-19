package cmd

import (
	"fmt"
	"net/url"
)

func checkApiKey(apiKey string) bool {
	if apiKey == "" {
		fmt.Println("Error: No API key configured.")
		fmt.Println("To get started, generate an API key:")
		fmt.Println("  naturedopes-cli keys generate <name>")
		return false
	}
	return true
}

func validatePositiveInt(id int) bool {
	if id <= 0 {
		fmt.Printf("Error: Id numbers must be greater than zero, got: %d\n", id)
		return false
	}
	return true
}

func validUrl(urlString string) bool {
	_, err := url.Parse(urlString)
	if err != nil {
		fmt.Printf("Error: please use a valid url for the api-url field, got: %s", urlString)
		return false
	}
	return true
}
