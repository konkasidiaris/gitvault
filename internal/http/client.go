package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/konkasidiaris/gitvault/internal/config"
)

var UserAgent = fmt.Sprintf("GitVault/%s", config.GetGitVaultVersion())

type Client struct {
	Client *http.Client
}

func NewClient() *Client {
	return &Client{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", UserAgent)
	return c.Client.Do(req)
}
