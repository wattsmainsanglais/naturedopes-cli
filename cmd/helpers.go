package cmd

import (
	"fmt"
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
