package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupFile(t *testing.T, filepath string, data []byte) {
	t.Helper()

	err := os.WriteFile(filepath, data, 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	t.Cleanup(func() {
		os.Remove(filepath)
	})
}

func TestGetGitHubRepositories_Success(t *testing.T) {
	path := "/tmp/db.json"
	repositoriesCount := 15
	repositories := make([]string, repositoriesCount)
	for index := range repositories {
		repositories[index] = fmt.Sprintf("repository-%d", index)
	}
	stringifiedRepositories, err := json.Marshal(repositories)
	jsonData := fmt.Sprintf(`{"github":{"repositories": %s}}`, stringifiedRepositories)
	setupFile(t, path, []byte(jsonData))

	ghRepositories, err := getGitHubRepositories(path)

	assert.NoError(t, err)
	assert.NotNil(t, ghRepositories)
	assert.Equal(t, repositoriesCount, len(ghRepositories))
}

func TestGetGitHubRepositories_FileNotFound(t *testing.T) {
	ghRepositories, err := getGitHubRepositories("/tmp/non_existent_db.json")

	assert.Nil(t, ghRepositories)
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestGetGitHubRepositories_InvalidJSON(t *testing.T) {
	path := "/tmp/invalid.json"
	setupFile(t, path, []byte(`{invalid json}`))

	ghRepositories, err := getGitHubRepositories(path)

	assert.Nil(t, ghRepositories)
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid character 'i' looking for beginning of object key string")
}

func TestGetGitHubRepositories_EmptyFile(t *testing.T) {
	path := "/tmp/empty.json"
	setupFile(t, path, []byte(``))

	ghRepositories, err := getGitHubRepositories(path)

	assert.Nil(t, ghRepositories)
	assert.Error(t, err)
	assert.EqualError(t, err, "EOF")
}

func TestGetGitHubRepositories_EmptyGithubObject(t *testing.T) {
	path := "/tmp/empty_github_object.json"
	setupFile(t, path, []byte(`{"github": {}}`))

	ghRepositories, err := getGitHubRepositories(path)

	assert.NotNil(t, ghRepositories)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(ghRepositories))

	assert.Equal(t, []string{}, ghRepositories)
}

func TestGetGitHubRepositories_WrongTypeRepositories(t *testing.T) {
	path := "/tmp/wrong_type_repositories.json"
	setupFile(t, path, []byte(`{"github":{"repositories": true}}`))

	ghRepositories, err := getGitHubRepositories(path)

	assert.Nil(t, ghRepositories)
	assert.Error(t, err)

	var typeErr *json.UnmarshalTypeError
	assert.True(t, errors.As(err, &typeErr))
	assert.Equal(t, "github.repositories", typeErr.Field)
}

func TestGetGitHubRepositories_TrimsWhitespaceOnRepositories(t *testing.T) {
	path := "/tmp/repositories_with_whitespace.json"
	repositories := []string{
		"repo-1",
		"  repo-2  ",
		"\trepo-3\t",
	}
	stringifiedRepositories, err := json.Marshal(repositories)
	if err != nil {
		log.Fatal("Could not create repositories with whitespace")
	}
	jsonData := fmt.Sprintf(`{"github":{"repositories": %s}}`, stringifiedRepositories)
	setupFile(t, path, []byte(jsonData))

	ghRepositories, err := getGitHubRepositories(path)

	assert.NoError(t, err)
	assert.NotNil(t, ghRepositories)
	assert.Equal(t, len(repositories), len(ghRepositories))
	for index := range repositories {
		assert.Equal(t, strings.TrimSpace(repositories[index]), ghRepositories[index])
	}
}

func TestInitializeDB_DBExists(t *testing.T) {
	path := "/tmp/db.json"
	jsonData := `{"github":{"repositories":[]}}`
	setupFile(t, path, []byte(jsonData))

	err := initializeDB(path)

	assert.NoError(t, err)
}

func TestInitializeDB_Success(t *testing.T) {
	path := "/tmp/new_db.json"

	err := initializeDB(path)
	assert.NoError(t, err)

	ghRepositories, err := getGitHubRepositories(path)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(ghRepositories))
	os.Remove(path)
}

func TestUpdateGithubRepositories_EmptyRepositoriesSuccess(t *testing.T) {
	path := "/tmp/existing_db.json"
	jsonData := `{"github":{"repositories":[]}}`
	setupFile(t, path, []byte(jsonData))

	err := initializeDB(path)
	assert.NoError(t, err)

	repositories := []string{"repo"}
	err = updateGithubRepositories(repositories, path)
	assert.NoError(t, err)

	db, err := getDB(path)
	assert.NoError(t, err)

	assert.Equal(t, repositories, db.Github.Repositories)
}

func TestUpdateGithubRepositories_ExistingRepositoriesSuccess(t *testing.T) {
	path := "/tmp/existing_db.json"
	jsonData := `{"github":{"repositories":["repo1"]}}`
	setupFile(t, path, []byte(jsonData))

	err := initializeDB(path)
	assert.NoError(t, err)

	repositories := []string{"repo1", "repo2"}
	err = updateGithubRepositories(repositories, path)
	assert.NoError(t, err)

	db, err := getDB(path)
	assert.NoError(t, err)

	assert.Equal(t, repositories, db.Github.Repositories)
}
