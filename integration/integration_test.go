//go:build integration

package integration

import (
	"context"
	"os"
	"testing"

	"github.com/konkasidiaris/gitvault/integration/helpers"
)

func TestMain(m *testing.M) {
	// Disable sidecar container for cleanups
	_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	os.Exit(m.Run())
}

func TestWireMock(t *testing.T) {
	ctx := context.Background()

	container, err := helpers.CreateWiremockContainer(t, ctx)
	if err != nil {
		t.Fatal(err)
	}

	statusCode, out, err := container.SendHttpGet("/users/unauthorizeduser/repos")
	if err != nil {
		t.Fatal(err, "Failed to get a response")
	}

	if statusCode != 401 {
		t.Fatalf("expected HTTP-401 but got %d", statusCode)
	}

	expectedBody := `{"message":"Bad credentials"}`
	if string(out) != expectedBody {
		t.Fatalf("expected '%s' but got '%s'", expectedBody, string(out))
	}
}

// // binaryPath holds the path to the compiled gitvault binary.
// var binaryPath string

// func TestMain(m *testing.M) {
// 	_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
// 	_ = os.Setenv("TESTCONTAINERS_CHECKS_DISABLE", "true")
// 	// Build the gitvault binary once for all integration tests.
// 	tmpDir, err := os.MkdirTemp("", "gitvault-integration-*")
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "failed to create temp dir: %v\n", err)
// 		os.Exit(1)
// 	}
// 	defer os.RemoveAll(tmpDir)

// 	binaryPath = filepath.Join(tmpDir, "gitvault")
// 	cmd := exec.Command("go", "build", "-o", binaryPath, "../cmd/gitvault")
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	if err := cmd.Run(); err != nil {
// 		fmt.Fprintf(os.Stderr, "failed to build gitvault: %v\n", err)
// 		os.Exit(1)
// 	}

// 	os.Exit(m.Run())
// }

// func isDockerAvailable() bool {
// 	if os.Getenv("DOCKER_HOST") != "" {
// 		return true
// 	}

// 	if _, err := os.Stat("/var/run/docker.sock"); err == nil {
// 		return true
// 	}

// 	cmd := exec.Command("docker", "info")
// 	if err := cmd.Run(); err == nil {
// 		return true
// 	}

// 	return false
// }

// func wiremockFiles(t *testing.T) []testcontainers.ContainerFile {
// 	t.Helper()

// 	_, currentFile, _, ok := runtime.Caller(0)
// 	require.True(t, ok, "failed to resolve current test file path")

// 	mappingsDir := filepath.Join(filepath.Dir(currentFile), "wiremock", "mappings")

// 	entries, err := os.ReadDir(mappingsDir)
// 	require.NoError(t, err, "wiremock mappings directory not found: %s", mappingsDir)

// 	files := make([]testcontainers.ContainerFile, 0, len(entries))
// 	for _, entry := range entries {
// 		if entry.IsDir() {
// 			continue
// 		}
// 		files = append(files, testcontainers.ContainerFile{
// 			HostFilePath:      filepath.Join(mappingsDir, entry.Name()),
// 			ContainerFilePath: filepath.Join("/home/wiremock/mappings", entry.Name()),
// 			FileMode:          0o644,
// 		})
// 	}

// 	require.NotEmpty(t, files, "no wiremock mapping files found in %s", mappingsDir)

// 	return files
// }

// func startWireMock(t *testing.T) (string, testcontainers.Container) {
// 	t.Helper()
// 	ctx := context.Background()

// 	if !isDockerAvailable() {
// 		t.Skip("Docker daemon not available. Skipping integration test requiring testcontainers.")
// 	}

// 	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
// 		ContainerRequest: testcontainers.ContainerRequest{
// 			Image:        "wiremock/wiremock:3.13.0",
// 			ExposedPorts: []string{"8080/tcp"},
// 			WaitingFor:   wait.ForHTTP("/__admin/mappings").WithPort("8080/tcp").WithStartupTimeout(30 * time.Second),
// 			Files:        wiremockFiles(t),
// 		},
// 		Started: true,
// 	})
// 	require.NoError(t, err)

// 	t.Cleanup(func() {
// 		_ = container.Terminate(ctx)
// 	})

// 	host, err := container.Host(ctx)
// 	require.NoError(t, err)

// 	port, err := container.MappedPort(ctx, "8080/tcp")
// 	require.NoError(t, err)

// 	baseURL := fmt.Sprintf("http://%s:%s", host, port.Port())
// 	return baseURL, container
// }

// // writeConfigFile creates a temporary gitvault.json config file.
// func writeConfigFile(t *testing.T, token, username string) string {
// 	t.Helper()

// 	configDir := t.TempDir()
// 	configPath := filepath.Join(configDir, "gitvault.json")

// 	config := map[string]string{
// 		"github_token":    token,
// 		"github_username": username,
// 	}
// 	data, err := json.Marshal(config)
// 	require.NoError(t, err)
// 	require.NoError(t, os.WriteFile(configPath, data, 0644))

// 	t.Cleanup(func() {
// 		os.RemoveAll(configDir)
// 	})

// 	return configPath
// }

// // runGitVault executes the gitvault binary with the given environment variables.
// // Returns stdout+stderr combined output, and the error (if any).
// func runGitVault(t *testing.T, env map[string]string) (string, error) {
// 	t.Helper()

// 	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
// 	defer cancel()

// 	cmd := exec.CommandContext(ctx, binaryPath)
// 	cmd.Env = os.Environ()
// 	for k, v := range env {
// 		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
// 	}

// 	output, err := cmd.CombinedOutput()
// 	return string(output), err
// }

// func TestIntegration_SyncClonesRepositories(t *testing.T) {
// 	wiremockURL, _ := startWireMock(t)

// 	backupDir := t.TempDir()
// 	t.Cleanup(func() {
// 		os.RemoveAll(backupDir)
// 	})

// 	configPath := writeConfigFile(t, "fake-token", "testuser")

// 	output, err := runGitVault(t, map[string]string{
// 		"GITVAULT_GITHUB_BASE_URL": wiremockURL,
// 		"GITVAULT_CONFIG_PATH":     configPath,
// 		"GITVAULT_BACKUP_DIR":      backupDir,
// 	})

// 	// The clone will fail because the SSH URLs aren't real git repos,
// 	// but we verify gitvault attempted to clone the right repositories.
// 	// This validates the full pipeline: config → API call → clone attempt.
// 	t.Logf("gitvault output:\n%s", output)

// 	// The program should still exit (it continues on clone errors).
// 	// We check that the GitHub API was called and repos were discovered.
// 	assert.Contains(t, output, "fetched 2 repositories from GitHub")
// 	assert.Contains(t, output, "testuser/repo-alpha")
// 	assert.Contains(t, output, "testuser/repo-beta")

// 	// Even though clone fails, we verify it didn't crash with a config or API error.
// 	if err != nil {
// 		// gitvault logs clone errors but still exits 0 — if it exits non-zero,
// 		// it means the config or API layer failed, which would be a real bug.
// 		assert.Contains(t, output, "clone", "if exit non-zero, expected clone-related error, not config/API")
// 	}
// }

// func TestIntegration_SyncWithNoRepositories(t *testing.T) {
// 	wiremockURL, _ := startWireMock(t)

// 	backupDir := t.TempDir()
// 	t.Cleanup(func() {
// 		os.RemoveAll(backupDir)
// 	})

// 	configPath := writeConfigFile(t, "fake-token", "emptyuser")

// 	output, err := runGitVault(t, map[string]string{
// 		"GITVAULT_GITHUB_BASE_URL": wiremockURL,
// 		"GITVAULT_CONFIG_PATH":     configPath,
// 		"GITVAULT_BACKUP_DIR":      backupDir,
// 	})

// 	t.Logf("gitvault output:\n%s", output)

// 	assert.NoError(t, err, "gitvault should exit 0 with no repos")
// 	assert.Contains(t, output, "fetched 0 repositories from GitHub")
// 	assert.Contains(t, output, "sync completed successfully")

// 	// Verify backup dir is empty (no repos cloned).
// 	entries, readErr := os.ReadDir(backupDir)
// 	assert.NoError(t, readErr)
// 	assert.Empty(t, entries, "backup directory should be empty when there are no repos")
// }

// func TestIntegration_SyncFailsOnUnauthorized(t *testing.T) {
// 	wiremockURL, _ := startWireMock(t)

// 	backupDir := t.TempDir()
// 	t.Cleanup(func() {
// 		os.RemoveAll(backupDir)
// 	})

// 	configPath := writeConfigFile(t, "bad-token", "unauthorizeduser")

// 	output, err := runGitVault(t, map[string]string{
// 		"GITVAULT_GITHUB_BASE_URL": wiremockURL,
// 		"GITVAULT_CONFIG_PATH":     configPath,
// 		"GITVAULT_BACKUP_DIR":      backupDir,
// 	})

// 	t.Logf("gitvault output:\n%s", output)

// 	assert.Error(t, err, "gitvault should exit non-zero on API error")
// 	assert.Contains(t, output, "unexpected status code: 401")
// }

// func TestIntegration_SyncFailsWithMissingConfig(t *testing.T) {
// 	backupDir := t.TempDir()
// 	t.Cleanup(func() {
// 		os.RemoveAll(backupDir)
// 	})

// 	output, err := runGitVault(t, map[string]string{
// 		"GITVAULT_CONFIG_PATH": "/nonexistent/gitvault.json",
// 		"GITVAULT_BACKUP_DIR":  backupDir,
// 	})

// 	t.Logf("gitvault output:\n%s", output)

// 	assert.Error(t, err, "gitvault should exit non-zero with missing config")
// 	assert.Contains(t, output, "error")
// }

// func TestIntegration_SyncFailsWithEmptyToken(t *testing.T) {
// 	backupDir := t.TempDir()
// 	t.Cleanup(func() {
// 		os.RemoveAll(backupDir)
// 	})

// 	configPath := writeConfigFile(t, "", "someuser")

// 	output, err := runGitVault(t, map[string]string{
// 		"GITVAULT_CONFIG_PATH": configPath,
// 		"GITVAULT_BACKUP_DIR":  backupDir,
// 	})

// 	t.Logf("gitvault output:\n%s", output)

// 	assert.Error(t, err, "gitvault should exit non-zero with empty token")
// 	assert.Contains(t, output, "GitHub token")
// }

// func TestIntegration_SyncUpdatesExistingMirror(t *testing.T) {
// 	wiremockURL, _ := startWireMock(t)

// 	backupDir := t.TempDir()
// 	t.Cleanup(func() {
// 		os.RemoveAll(backupDir)
// 	})

// 	// Pre-create a bare git repo to simulate an existing mirror.
// 	repoDir := filepath.Join(backupDir, "repo-alpha.git")
// 	initCmd := exec.Command("git", "init", "--bare", repoDir)
// 	require.NoError(t, initCmd.Run(), "failed to init bare repo for test setup")

// 	configPath := writeConfigFile(t, "fake-token", "testuser")

// 	output, err := runGitVault(t, map[string]string{
// 		"GITVAULT_GITHUB_BASE_URL": wiremockURL,
// 		"GITVAULT_CONFIG_PATH":     configPath,
// 		"GITVAULT_BACKUP_DIR":      backupDir,
// 	})

// 	t.Logf("gitvault output:\n%s", output)

// 	// repo-alpha should be updated (not cloned), repo-beta should be cloned.
// 	assert.Contains(t, output, "updating mirror")
// 	assert.Contains(t, output, "repo-alpha")

// 	// repo-beta will be a new clone attempt.
// 	assert.Contains(t, output, "cloning mirror")
// 	assert.Contains(t, output, "repo-beta")

// 	_ = err // clone/update may fail since SSH URLs aren't real, but the logic paths are exercised.
// }

// func TestIntegration_BackupDirectoryCreated(t *testing.T) {
// 	wiremockURL, _ := startWireMock(t)

// 	parentDir := t.TempDir()
// 	backupDir := filepath.Join(parentDir, "deep", "nested", "backup")
// 	t.Cleanup(func() {
// 		os.RemoveAll(parentDir)
// 	})

// 	configPath := writeConfigFile(t, "fake-token", "emptyuser")

// 	output, err := runGitVault(t, map[string]string{
// 		"GITVAULT_GITHUB_BASE_URL": wiremockURL,
// 		"GITVAULT_CONFIG_PATH":     configPath,
// 		"GITVAULT_BACKUP_DIR":      backupDir,
// 	})

// 	t.Logf("gitvault output:\n%s", output)

// 	assert.NoError(t, err)

// 	info, statErr := os.Stat(backupDir)
// 	assert.NoError(t, statErr, "backup directory should have been created")
// 	assert.True(t, info.IsDir())
// }

// func TestIntegration_ConcurrentSyncIsolation(t *testing.T) {
// 	wiremockURL, _ := startWireMock(t)

// 	// Run two syncs in parallel with different backup dirs to prove isolation.
// 	t.Run("sync-a", func(t *testing.T) {
// 		t.Parallel()

// 		backupDir := t.TempDir()
// 		t.Cleanup(func() { os.RemoveAll(backupDir) })

// 		configPath := writeConfigFile(t, "fake-token", "testuser")

// 		output, _ := runGitVault(t, map[string]string{
// 			"GITVAULT_GITHUB_BASE_URL": wiremockURL,
// 			"GITVAULT_CONFIG_PATH":     configPath,
// 			"GITVAULT_BACKUP_DIR":      backupDir,
// 		})
// 		assert.Contains(t, output, "fetched 2 repositories")
// 	})

// 	t.Run("sync-b", func(t *testing.T) {
// 		t.Parallel()

// 		backupDir := t.TempDir()
// 		t.Cleanup(func() { os.RemoveAll(backupDir) })

// 		configPath := writeConfigFile(t, "fake-token", "emptyuser")

// 		output, err := runGitVault(t, map[string]string{
// 			"GITVAULT_GITHUB_BASE_URL": wiremockURL,
// 			"GITVAULT_CONFIG_PATH":     configPath,
// 			"GITVAULT_BACKUP_DIR":      backupDir,
// 		})
// 		assert.NoError(t, err)
// 		assert.Contains(t, output, "fetched 0 repositories")
// 	})
// }

// func TestIntegration_OutputIsStructuredJSON(t *testing.T) {
// 	wiremockURL, _ := startWireMock(t)

// 	backupDir := t.TempDir()
// 	t.Cleanup(func() {
// 		os.RemoveAll(backupDir)
// 	})

// 	configPath := writeConfigFile(t, "fake-token", "emptyuser")

// 	output, err := runGitVault(t, map[string]string{
// 		"GITVAULT_GITHUB_BASE_URL": wiremockURL,
// 		"GITVAULT_CONFIG_PATH":     configPath,
// 		"GITVAULT_BACKUP_DIR":      backupDir,
// 	})

// 	assert.NoError(t, err)

// 	// Verify each line of stderr output is valid JSON (structured logging).
// 	lines := strings.Split(strings.TrimSpace(output), "\n")
// 	for _, line := range lines {
// 		line = strings.TrimSpace(line)
// 		if line == "" {
// 			continue
// 		}
// 		var parsed map[string]interface{}
// 		jsonErr := json.Unmarshal([]byte(line), &parsed)
// 		assert.NoError(t, jsonErr, "each log line should be valid JSON, got: %s", line)
// 	}
// }
