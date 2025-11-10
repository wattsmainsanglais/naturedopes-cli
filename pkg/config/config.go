package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	ApiURL string `json:"api_url"`
	ApiKey string `json:"api_key"`
}

func getConfigFilePath() (string, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("user home directory not found: %w", err)
	}

	fullPath := filepath.Join(homeDir, ".naturedopes-cli", "config.json")

	return fullPath, nil

}

func Load() (*Config, error) {

	path, err := getConfigFilePath()
	if err != nil {
		return nil, fmt.Errorf("couldn't get home directory: %w", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{
			ApiURL: "http://localhost:8080",
			ApiKey: "",
		}, nil
	}

	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file: %w", err)
	}

	var config Config

	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal JSON: %w", err)
	}

	return &config, nil

}
