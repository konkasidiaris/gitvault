package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupFile(t *testing.T, filePath string, data []byte) {
	t.Helper()

	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	t.Cleanup(func() {
		os.Remove(filePath)
	})
}

func TestLoadConfig_Success(t *testing.T) {
	path := "/tmp/test_config.json"
	token := "test-github-token"
	username := "GitHub_username"
	jsonData := fmt.Sprintf(`{"github_token": "%s", "github_username": "%s"}`, token, username)
	setupFile(t, path, []byte(jsonData))

	cfg, err := LoadConfig(path)

	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, cfg.GitHubToken, token)
	assert.Equal(t, cfg.GitHubUsername, username)
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	cfg, err := LoadConfig("/tmp/non_existent_file.json")

	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	path := "/tmp/invalid.json"
	setupFile(t, path, []byte(`{invalid json}`))

	cfg, err := LoadConfig(path)

	assert.Nil(t, cfg)
	assert.Error(t, err)
}

func TestLoadConfig_EmptyFile(t *testing.T) {
	path := "/tmp/empty.json"
	setupFile(t, path, []byte(``))

	cfg, err := LoadConfig(path)

	assert.Nil(t, cfg)
	assert.Error(t, err)
}

func TestLoadConfig_WrongTypeOfGitHubToken(t *testing.T) {
	path := "/tmp/wrong-github-token-type.json"
	setupFile(t, path, []byte(`{ "github_token": true}`))

	cfg, err := LoadConfig(path)

	assert.Nil(t, cfg)
	assert.Error(t, err)

	var typeErr *json.UnmarshalTypeError
	assert.True(t, errors.As(err, &typeErr))
	assert.Equal(t, "github_token", typeErr.Field)
}

func TestLoadConfig_WrongTypeOfGitHubUsername(t *testing.T) {
	path := "/tmp/wrong-github-username-type.json"
	setupFile(t, path, []byte(`{ "github_username": 1}`))

	cfg, err := LoadConfig(path)

	assert.Nil(t, cfg)
	assert.Error(t, err)

	var typeErr *json.UnmarshalTypeError
	assert.True(t, errors.As(err, &typeErr))
	assert.Equal(t, "github_username", typeErr.Field)
}

func TestLoadConfig_TrimsWhitespace(t *testing.T) {
	path := "/tmp/test_config_whitespaces.json"
	token := "test-github-token"
	username := "GitHub_username"
	jsonData := fmt.Sprintf(`{"github_token": " %s\t", "github_username": "\t%s \n"}`, token, username)
	setupFile(t, path, []byte(jsonData))

	cfg, err := LoadConfig(path)

	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, cfg.GitHubToken, token)
	assert.Equal(t, cfg.GitHubUsername, username)
}

func TestLoadConfig_EmptyValues(t *testing.T) {
	path := "/tmp/empty_values.json"
	jsonData := `{"github_token": "", "github_username": ""}`
	setupFile(t, path, []byte(jsonData))

	cfg, err := LoadConfig(path)

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "", cfg.GitHubToken)
	assert.Equal(t, "", cfg.GitHubUsername)
}

func TestLoadConfig_MissingFields(t *testing.T) {
	path := "/tmp/missing_fields.json"
	jsonData := `{}`
	setupFile(t, path, []byte(jsonData))

	cfg, err := LoadConfig(path)

	assert.NoError(t, err)
	assert.Equal(t, "", cfg.GitHubToken)
	assert.Equal(t, "", cfg.GitHubUsername)
}

func TestLoadConfig_ExtraFieldsIgnored(t *testing.T) {
	path := "/tmp/extra-fields-ignored.json"
	token := "test-github-token"
	username := "GitHub_username"
	jsonData := fmt.Sprintf(`{"github_token": "%s", "github_username": "%s", "extra_field": false}`, token, username)
	setupFile(t, path, []byte(jsonData))

	cfg, err := LoadConfig(path)

	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, cfg.GitHubToken, token)
	assert.Equal(t, cfg.GitHubUsername, username)
}
