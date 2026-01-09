package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const version = "v0.0.1"

var secretsBasePath = "/secrets" // declaring this here, to be able to change it in tests

type Config struct {
	Version     string
	GitHubToken string
}

var (
	instance *Config
	once     sync.Once
	loadErr  error
)

func readSecretFromFile(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	trimmedString := strings.TrimSpace(string(data))

	return trimmedString, nil
}

func Get() (*Config, error) {
	once.Do(
		func() {
			path := filepath.Join(secretsBasePath, "github-token.txt")
			githubToken, err := readSecretFromFile(path)

			if err != nil || githubToken == "" {
				loadErr = errors.New("github-token.txt is missing OR github token cannot be an empty string")
			}

			instance = &Config{
				Version:     version,
				GitHubToken: githubToken,
			}
		},
	)

	return instance, loadErr
}

func GetGitVaultVersion() string {
	if instance == nil {
		Get()
	}

	return instance.Version
}

func GetGitHubToken() string {
	if instance == nil {
		Get()
	}

	return instance.GitHubToken
}
