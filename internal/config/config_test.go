package config

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockConfigLoader struct {
	config *GitVaultFileConfig
	err    error
}

func (m *mockConfigLoader) Load(filepath string) (*GitVaultFileConfig, error) {
	return m.config, m.err
}

func mockConfig(t *testing.T, c *GitVaultFileConfig, e error) {
	t.Helper()

	originalLoader := configLoader
	configLoader = &mockConfigLoader{
		config: c,
		err:    e,
	}

	t.Cleanup(func() {
		configLoader = originalLoader
		reset()
	})
}

func reset() {
	instance = nil
	once = sync.Once{}
	loadErr = nil
}

func TestGet_Success(t *testing.T) {
	mockGitVaultConfig := &GitVaultFileConfig{
		GitHubToken:    "test-github-token",
		GitHubUsername: "test-github-username",
	}
	mockConfig(t, mockGitVaultConfig, nil)

	assert.Nil(t, instance)
	assert.Nil(t, loadErr)
	cfg, err := Get()

	assert.NoError(t, err)
	assert.NotNil(t, instance)
	assert.Nil(t, loadErr)
	assert.Equal(t, version, cfg.Version)
	assert.Equal(t, "test-github-token", cfg.GitHubToken)
}

func TestGet_LoadingGitVaultConfigError(t *testing.T) {
	mockConfig(t, nil, errors.New("Missing file"))
	cfg, err := Get()

	assert.EqualError(t, err, "[Config] error while loading /secrets/gitvault.json: Missing file")
	assert.True(t, cfg == nil)
}

func TestGet_MissingGitHubToken(t *testing.T) {
	mockGitVaultConfig := &GitVaultFileConfig{
		GitHubToken:    "",
		GitHubUsername: "test-github-username",
	}
	mockConfig(t, mockGitVaultConfig, nil)
	cfg, err := Get()

	assert.EqualError(t, err, "[Config] GitHub token is either missing or empty")
	assert.True(t, cfg == nil)
}

func TestGet_MissingGitHubUsername(t *testing.T) {
	mockGitVaultConfig := &GitVaultFileConfig{
		GitHubToken:    "test-github-token",
		GitHubUsername: "",
	}
	mockConfig(t, mockGitVaultConfig, nil)
	cfg, err := Get()

	assert.EqualError(t, err, "[Config] GitHub Username is either missing or empty")
	assert.True(t, cfg == nil)
}

func TestGetGitVaultVersion_Success(t *testing.T) {
	mockGitVaultConfig := &GitVaultFileConfig{
		GitHubToken:    "test-github-token",
		GitHubUsername: "test-github-username",
	}
	mockConfig(t, mockGitVaultConfig, nil)

	assert.Nil(t, instance)
	assert.Nil(t, loadErr)
	localVersion := GetGitVaultVersion()
	assert.Equal(t, version, localVersion)

	assert.NotNil(t, instance)
	assert.Nil(t, loadErr)
	localVersion = GetGitVaultVersion()
	assert.Equal(t, version, localVersion)
}

func TestGetGitHubToken_Success(t *testing.T) {
	mockGitVaultConfig := &GitVaultFileConfig{
		GitHubToken:    "test-github-token",
		GitHubUsername: "test-github-username",
	}
	mockConfig(t, mockGitVaultConfig, nil)

	assert.Nil(t, instance)
	assert.Nil(t, loadErr)
	token := GetGitHubToken()
	assert.Equal(t, "test-github-token", token)

	assert.NotNil(t, instance)
	assert.Nil(t, loadErr)
	token = GetGitHubToken()
	assert.Equal(t, "test-github-token", token)
}

func TestGetGitHubUsername_Success(t *testing.T) {
	mockGitVaultConfig := &GitVaultFileConfig{
		GitHubToken:    "test-github-token",
		GitHubUsername: "test-github-username",
	}
	mockConfig(t, mockGitVaultConfig, nil)

	assert.Nil(t, instance)
	assert.Nil(t, loadErr)
	username := GetGitHubUsername()
	assert.Equal(t, "test-github-username", username)

	assert.NotNil(t, instance)
	assert.Nil(t, loadErr)
	username = GetGitHubUsername()
	assert.Equal(t, "test-github-username", username)
}
