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
		apiUrl := os.Getenv("API_URL")
		if apiUrl == "" {
			apiUrl = "http://localhost:8080"
		}

		return &Config{
			ApiURL: apiUrl,
			ApiKey: os.Getenv("API_KEY"),
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

func (config *Config) Save() error {

	path, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("couldn't get home directory: %w", err)
	}

	configDir := filepath.Dir(path)

	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("couldn't marshal JSON: %w", err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("could not write data to file: %w", err)
	}

	return nil

}

func Set(key, value string) error {

	currentConfig, err := Load()
	if err != nil {
		return fmt.Errorf("could not load config file: %w", err)
	}

	switch key {
	case "api-url":
		currentConfig.ApiURL = value
	case "api-key":
		currentConfig.ApiKey = value
	default:
		return fmt.Errorf("invalid key: %s", key)
	}

	err = currentConfig.Save()
	if err != nil {
		return fmt.Errorf("could not save config file: %w", err)
	}

	return nil

}

func Get(key string) (string, error) {

	currentConfig, err := Load()
	if err != nil {
		return "", fmt.Errorf("could not load config file: %w", err)
	}

	switch key {
	case "api-url":
		return currentConfig.ApiURL, nil
	case "api-key":
		return currentConfig.ApiKey, nil
	default:
		return "", fmt.Errorf("invalid key: %s", key)
	}

}
