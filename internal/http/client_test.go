package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	assert.NotNil(t, client.Client)
	assert.Equal(t, 30*time.Second, client.Client.Timeout)
}

func TestClientDo_AddsUserAgent(t *testing.T) {
	tests := []struct {
		name      string
		versionFn func() string
		expected  string
	}{
		{
			name:      "custom version",
			versionFn: func() string { return "mocked version" },
			expected:  "GitVault/mocked version",
		},
		{
			name:      "empty version",
			versionFn: func() string { return "" },
			expected:  "GitVault/",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expected, r.Header.Get("User-Agent"))
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			client := newClientWithVersion(tc.versionFn)
			req, _ := http.NewRequest("GET", server.URL, nil)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}
