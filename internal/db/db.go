package db

import (
	"encoding/json"
	"os"
	"strings"
)

const lockfilePath = "gitvault.lock.json"

type DB struct {
	Github struct {
		Repositories []string `json:"repositories"`
	} `json:"github"`
}

func getDB(filepath string) (*DB, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var db DB
	if err := json.NewDecoder(file).Decode(&db); err != nil {
		return nil, err
	}

	if db.Github.Repositories == nil {
		db.Github.Repositories = []string{}
	}

	for index := range db.Github.Repositories {
		db.Github.Repositories[index] = strings.TrimSpace(db.Github.Repositories[index])
	}

	return &db, nil
}

func GetGitHubRepositories() ([]string, error) {
	return getGitHubRepositories(lockfilePath)
}

func getGitHubRepositories(filepath string) ([]string, error) {
	db, err := getDB(filepath)
	if err != nil {
		return nil, err
	}
	return db.Github.Repositories, nil
}

func InitializeDB() error {
	return initializeDB(lockfilePath)
}

func initializeDB(filepath string) error {
	if _, err := os.Stat(filepath); err == nil {
		return nil
	}

	initialDB := DB{
		Github: struct {
			Repositories []string `json:"repositories"`
		}{
			Repositories: []string{},
		},
	}

	data, err := json.MarshalIndent(initialDB, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}

func UpdateGithubRepositories(repositories []string) error {
	return updateGithubRepositories(repositories, lockfilePath)
}

func updateGithubRepositories(repositories []string, filepath string) error {
	db, err := getDB(filepath)
	if err != nil {
		return err
	}

	db.Github.Repositories = repositories
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}
