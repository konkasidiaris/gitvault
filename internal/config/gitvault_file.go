package config

import (
	"encoding/json"
	"os"
	"strings"
)

type GitVaultFileConfig struct {
	GitHubToken    string `json:"github_token"`
	GitHubUsername string `json:"github_username"`
}

func LoadConfig(filepath string) (*GitVaultFileConfig, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg GitVaultFileConfig
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	cfg.GitHubToken = strings.TrimSpace(cfg.GitHubToken)
	cfg.GitHubUsername = strings.TrimSpace(cfg.GitHubUsername)

	return &cfg, nil
}
