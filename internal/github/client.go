package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/konkasidiaris/gitvault/internal/config"
)

const (
	baseURL = "https://api.github.com"
)

type Repository struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	SSHURL   string `json:"ssh_url"`
}

type Client struct {
	baseURL  string
	token    string
	username string
	http     *http.Client
}

func NewClient() *Client {
	return &Client{
		baseURL:  baseURL,
		token:    config.GetGitHubToken(),
		username: config.GetGitHubUsername(),
		http:     &http.Client{},
	}
}

func (c *Client) GetUserRepos() ([]Repository, error) {
	var repos []Repository

	url := fmt.Sprintf("%s/users/%s/repos", c.baseURL, c.username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return repos, nil
}
