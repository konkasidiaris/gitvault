package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestClient(baseURL, token, username string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:  baseURL,
		token:    token,
		username: username,
		http:     httpClient,
	}
}

func TestGetUserRepos_Success(t *testing.T) {
	expected := []Repository{
		{ID: 1, FullName: "user/repo1", SSHURL: "git@github.com:user/repo1.git"},
		{ID: 2, FullName: "user/repo2", SSHURL: "git@github.com:user/repo2.git"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Path; got != "/users/testuser/repos" {
			t.Fatalf("path = %s, want /users/testuser/repos", got)
		}
		if r.Header.Get("Accept") != "application/vnd.github+json" ||
			r.Header.Get("X-GitHub-Api-Version") != "2022-11-28" ||
			r.Header.Get("Authorization") != "Bearer test-token" {
			t.Fatalf("headers not set as expected: %+v", r.Header)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	client := newTestClient(server.URL, "test-token", "testuser", server.Client())

	repos, err := client.GetUserRepos()
	if err != nil {
		t.Fatalf("got err: %v", err)
	}
	if len(repos) != len(expected) {
		t.Fatalf("len(repos) = %d, want %d", len(repos), len(expected))
	}
	for i := range repos {
		if repos[i] != expected[i] {
			t.Fatalf("repo[%d] = %+v, want %+v", i, repos[i], expected[i])
		}
	}
}

func TestGetUserRepos_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Repository{})
	}))
	defer server.Close()

	client := newTestClient(server.URL, "test-token", "testuser", server.Client())

	repos, err := client.GetUserRepos()
	if err != nil {
		t.Fatalf("got err: %v", err)
	}
	if len(repos) != 0 {
		t.Fatalf("len(repos) = %d, want 0", len(repos))
	}
}

func TestGetUserRepos_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message": "Bad credentials"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL, "invalid-token", "testuser", server.Client())

	_, err := client.GetUserRepos()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedErr := "unexpected status code: 401"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestGetUserRepos_Forbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"message": "API rate limit exceeded"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL, "test-token", "testuser", server.Client())

	_, err := client.GetUserRepos()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedErr := "unexpected status code: 403"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestGetUserRepos_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Not Found"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL, "test-token", "nonexistent-user", server.Client())

	_, err := client.GetUserRepos()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedErr := "unexpected status code: 404"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestGetUserRepos_InternalServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Internal Server Error"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL, "test-token", "testuser", server.Client())

	_, err := client.GetUserRepos()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedErr := "unexpected status code: 500"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestGetUserRepos_ServiceUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"message": "Service Unavailable"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL, "test-token", "testuser", server.Client())

	_, err := client.GetUserRepos()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedErr := "unexpected status code: 503"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestGetUserRepos_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid json`))
	}))
	defer server.Close()

	client := newTestClient(server.URL, "test-token", "testuser", server.Client())

	_, err := client.GetUserRepos()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to decode response") {
		t.Errorf("expected error to start with 'failed to decode response', got '%s'", err.Error())
	}
}

func TestGetUserRepos_DifferentConfigs(t *testing.T) {
	testCases := []struct {
		name     string
		token    string
		username string
	}{
		{name: "standard config", token: "ghp_xxxxxxxxxxxx", username: "johndoe"},
		{name: "special characters in username", token: "token123", username: "user-name_123"},
		{name: "long token", token: "ghp_" + strings.Repeat("abc", 100), username: "user"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/users/" + tc.username + "/repos"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path '%s', got '%s'", expectedPath, r.URL.Path)
				}

				expectedAuth := "Bearer " + tc.token
				if r.Header.Get("Authorization") != expectedAuth {
					t.Errorf("expected Authorization '%s', got '%s'", expectedAuth, r.Header.Get("Authorization"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]Repository{})
			}))
			defer server.Close()

			client := newTestClient(server.URL, tc.token, tc.username, server.Client())
			_, err := client.GetUserRepos()
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

// TestGetUserRepos_LargeResponse tests handling of large repository lists
func TestGetUserRepos_LargeResponse(t *testing.T) {
	var largeRepoList []Repository
	for i := 0; i < 1000; i++ {
		largeRepoList = append(largeRepoList, Repository{
			ID:       int64(i),
			FullName: "user/repo" + string(rune(i)),
			SSHURL:   "git@github.com:user/repo" + string(rune(i)) + ".git",
		})
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(largeRepoList)
	}))
	defer server.Close()

	client := newTestClient(server.URL, "test-token", "testuser", server.Client())

	repos, err := client.GetUserRepos()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repos) != 1000 {
		t.Fatalf("expected 1000 repos, got %d", len(repos))
	}
}
