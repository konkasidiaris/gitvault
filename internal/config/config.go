package config

import (
	"fmt"
	"sync"
)

const version = "v0.0.1"

type ConfigLoader interface {
	Load(filepath string) (*GitVaultFileConfig, error)
}

type FileConfigLoader struct{}

func (f *FileConfigLoader) Load(filepath string) (*GitVaultFileConfig, error) {
	return LoadConfig(filepath)
}

var configLoader ConfigLoader = &FileConfigLoader{}

type Config struct {
	Version        string
	GitHubToken    string
	GitHubUsername string
}

var (
	instance *Config
	once     sync.Once
	loadErr  error
)

func Get() (*Config, error) {
	once.Do(
		func() {
			fileConfig, err := configLoader.Load("/secrets/gitvault.json")
			if err != nil {
				loadErr = fmt.Errorf("[Config] error while loading /secrets/gitvault.json: %w", err)
				return
			}

			if fileConfig.GitHubToken == "" {
				loadErr = fmt.Errorf("[Config] GitHub token is either missing or empty")
				return
			}

			if fileConfig.GitHubUsername == "" {
				loadErr = fmt.Errorf("[Config] GitHub Username is either missing or empty")
				return
			}

			instance = &Config{
				Version:        version,
				GitHubToken:    fileConfig.GitHubToken,
				GitHubUsername: fileConfig.GitHubUsername,
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

func GetGitHubUsername() string {
	if instance == nil {
		Get()
	}

	return instance.GitHubUsername
}
