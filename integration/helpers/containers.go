package helpers

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type WireMockContainer struct {
	testcontainers.Container
}

func (c *WireMockContainer) BaseURL(ctx context.Context) (string, error) {
	ip, err := c.ContainerIP(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get container IP: %w", err)
	}
	return fmt.Sprintf("http://%s:8080", ip), nil
}

func (c *WireMockContainer) SendHttpGet(path string) (int, []byte, error) {
	ctx := context.Background()
	base, err := c.BaseURL(ctx)
	if err != nil {
		return 0, nil, err
	}
	resp, err := http.Get(base + path)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	body := make([]byte, 0)
	buf := make([]byte, 512)
	for {
		n, err := resp.Body.Read(buf)
		body = append(body, buf[:n]...)
		if err != nil {
			break
		}
	}
	return resp.StatusCode, body, nil
}

func CreateWiremockContainer(t *testing.T, ctx context.Context) (*WireMockContainer, error) {

	_, thisFile, _, _ := runtime.Caller(0)
	mappingsDir := filepath.Join(filepath.Dir(thisFile), "..", "wiremock", "mappings")
	entries, err := os.ReadDir(mappingsDir)
	if err != nil {
		t.Fatal(err)
	}

	files := make([]testcontainers.ContainerFile, 0, len(entries))

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".json") {
			file := testcontainers.ContainerFile{
				HostFilePath:      filepath.Join(mappingsDir, entry.Name()),
				ContainerFilePath: fmt.Sprintf("/home/wiremock/mappings/%s", entry.Name()),
				FileMode:          0o644,
			}

			files = append(files, file)
		}
	}

	req := testcontainers.ContainerRequest{
		Image:      "docker.io/wiremock/wiremock:3.13.2",
		Name:       fmt.Sprintf("wiremock-testcontainer-%d", rand.Int()),
		Files:      files,
		WaitingFor: wait.ForLog("extensions:").WithStartupTimeout(60 * time.Second),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	wm := &WireMockContainer{Container: c}
	t.Cleanup(func() { _ = c.Terminate(ctx) })
	return wm, nil
}
