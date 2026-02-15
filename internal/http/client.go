package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/konkasidiaris/gitvault/internal/config"
)

type Client struct {
	Client    *http.Client
	versionFn func() string
}

func NewClient() *Client {
	return newClientWithVersion(config.GetGitVaultVersion)
}

func newClientWithVersion(versionFn func() string) *Client {

	return &Client{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		versionFn: versionFn,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	UserAgent := fmt.Sprintf("GitVault/%s", c.versionFn())
	req.Header.Set("User-Agent", UserAgent)
	return c.Client.Do(req)
}
