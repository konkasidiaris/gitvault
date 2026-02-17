package sync

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/konkasidiaris/gitvault/internal/github"
)

const backupDirectory = "/backup"

var (
	fetchGithubRepositories = getGithubRepositories
	cloneMirrorFn           = gitCloneMirror
	remoteUpdateFn          = gitRemoteUpdate
)

func getGithubRepositories() ([]github.Repository, error) {
	client := github.NewClient()
	return client.GetUserRepos()
}

func repositoryName(fullName string) string {
	parts := strings.SplitN(fullName, "/", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return fullName
}

func gitCloneMirror(sshURL, targetDirectory string) error {
	cmd := exec.Command("git", "clone", "--mirror", sshURL, targetDirectory)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func gitRemoteUpdate(repository string) error {
	cmd := exec.Command("git", "remote", "update")
	cmd.Dir = repository
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func run(dir string) error {
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create backup directory %s: %w", dir, err)
		}
	}

	repos, err := fetchGithubRepositories()
	if err != nil {
		return fmt.Errorf("failed to fetch repositories from GitHub: %w", err)
	}

	slog.Info(fmt.Sprintf("fetched %d repositories from GitHub", len(repos)))

	for _, repository := range repos {
		name := repositoryName(repository.FullName)
		repositoryDirectory := filepath.Join(dir, name+".git")

		if info, err := os.Stat(repositoryDirectory); err == nil && info.IsDir() {
			slog.Info("updating mirror", "repository", repository.FullName, "dir", repositoryDirectory)
			if err := remoteUpdateFn(repositoryDirectory); err != nil {
				slog.Error("failed to update mirror", "repo", repository.FullName, "error", err)
				continue
			}
		} else {
			slog.Info("cloning mirror", "repository", repository.FullName, "dir", repositoryDirectory)
			if err := cloneMirrorFn(repository.SSHURL, repositoryDirectory); err != nil {
				slog.Error("failed to clone mirror", "repository", repository.FullName, "error", err)
				continue
			}
		}
	}

	slog.Info("sync completed successfully", "total", len(repos))
	return nil
}

func Run() error {
	return run(backupDirectory)
}
