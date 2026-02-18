package config

import (
	"fmt"
	"os"
	"sync"
)

const version = "v0.0.1"
const defaultConfigPath = "/secrets/gitvault.json"

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

func getConfigPath() string {
	if path := os.Getenv("GITVAULT_CONFIG_PATH"); path != "" {
		return path
	}
	return defaultConfigPath
}

func Get() (*Config, error) {
	once.Do(
		func() {
			configPath := getConfigPath()
			fileConfig, err := configLoader.Load(configPath)
			if err != nil {
				loadErr = fmt.Errorf("[Config] error while loading %s: %w", configPath, err)
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
