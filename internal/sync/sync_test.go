package sync

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/konkasidiaris/gitvault/internal/github"
	"github.com/stretchr/testify/assert"
)

// mockGitOps stores calls made to clone/update for assertions.
type mockGitOps struct {
	cloneCalls             []cloneCall
	updateCalls            []string
	cloneErr               error
	updateErr              error
	cloneErrForRepository  map[string]error
	updateErrForRepository map[string]error
	createDirectoryOnClone bool
}

type cloneCall struct {
	sshURL          string
	targetDirectory string
}

func newMockGitOps() *mockGitOps {
	return &mockGitOps{
		cloneErrForRepository:  make(map[string]error),
		updateErrForRepository: make(map[string]error),
		createDirectoryOnClone: true,
	}
}

func (m *mockGitOps) clone(sshURL, targetDirectory string) error {
	m.cloneCalls = append(m.cloneCalls, cloneCall{sshURL, targetDirectory})

	if err, ok := m.cloneErrForRepository[targetDirectory]; ok {
		return err
	}
	if m.cloneErr != nil {
		return m.cloneErr
	}

	// Simulate git clone by creating the directory
	if m.createDirectoryOnClone {
		return os.MkdirAll(targetDirectory, 0755)
	}
	return nil
}

func (m *mockGitOps) update(repoDir string) error {
	m.updateCalls = append(m.updateCalls, repoDir)

	if err, ok := m.updateErrForRepository[repoDir]; ok {
		return err
	}
	return m.updateErr
}

func setupMocks(t *testing.T, repos []github.Repository, fetchErr error, ops *mockGitOps) {
	t.Helper()

	originalFetch := fetchGithubRepositories
	originalClone := cloneMirrorFn
	originalUpdate := remoteUpdateFn

	fetchGithubRepositories = func() ([]github.Repository, error) {
		return repos, fetchErr
	}
	cloneMirrorFn = ops.clone
	remoteUpdateFn = ops.update

	t.Cleanup(func() {
		fetchGithubRepositories = originalFetch
		cloneMirrorFn = originalClone
		remoteUpdateFn = originalUpdate
	})
}

func TestRepoName_Success(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "owner slash", input: "owner/my-repo", expected: "my-repo"},
		{name: "no slash", input: "my-repo", expected: "my-repo"},
		{name: "multiple slashes", input: "owner/repo/extra", expected: "repo/extra"},
		{name: "empty string", input: "", expected: ""},
		{name: "only slash", input: "/", expected: ""},
		{name: "trailing slash", input: "owner/repo", expected: "repo"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, repositoryName(tc.input))
		})
	}
}

func TestRun_FetchError(t *testing.T) {
	ops := newMockGitOps()
	setupMocks(t, nil, errors.New("API error"), ops)

	err := run(t.TempDir())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch repositories from GitHub")
	assert.Empty(t, ops.cloneCalls)
	assert.Empty(t, ops.updateCalls)
}

func TestRun_EmptyRepos(t *testing.T) {
	ops := newMockGitOps()
	setupMocks(t, []github.Repository{}, nil, ops)

	err := run(t.TempDir())

	assert.NoError(t, err)
	assert.Empty(t, ops.cloneCalls)
	assert.Empty(t, ops.updateCalls)
}

func TestRun_ClonesNewRepos(t *testing.T) {
	dir := t.TempDir()
	repos := []github.Repository{
		{ID: 1, FullName: "user/repo1", SSHURL: "git@github.com:user/repo1.git"},
		{ID: 2, FullName: "user/repo2", SSHURL: "git@github.com:user/repo2.git"},
	}

	ops := newMockGitOps()
	setupMocks(t, repos, nil, ops)

	err := run(dir)

	assert.NoError(t, err)
	assert.Len(t, ops.cloneCalls, 2)
	assert.Equal(t, "git@github.com:user/repo1.git", ops.cloneCalls[0].sshURL)
	assert.Equal(t, filepath.Join(dir, "repo1.git"), ops.cloneCalls[0].targetDirectory)
	assert.Equal(t, "git@github.com:user/repo2.git", ops.cloneCalls[1].sshURL)
	assert.Equal(t, filepath.Join(dir, "repo2.git"), ops.cloneCalls[1].targetDirectory)
	assert.Empty(t, ops.updateCalls)
}

func TestRun_UpdatesExistingMirrors(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "repo1.git"), 0755)
	os.MkdirAll(filepath.Join(dir, "repo2.git"), 0755)

	repos := []github.Repository{
		{ID: 1, FullName: "user/repo1", SSHURL: "git@github.com:user/repo1.git"},
		{ID: 2, FullName: "user/repo2", SSHURL: "git@github.com:user/repo2.git"},
	}

	ops := newMockGitOps()
	setupMocks(t, repos, nil, ops)

	err := run(dir)

	assert.NoError(t, err)
	assert.Empty(t, ops.cloneCalls)
	assert.Len(t, ops.updateCalls, 2)
	assert.Equal(t, filepath.Join(dir, "repo1.git"), ops.updateCalls[0])
	assert.Equal(t, filepath.Join(dir, "repo2.git"), ops.updateCalls[1])
}

func TestRun_MixedCloneAndUpdate(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "repo1.git"), 0755)

	repos := []github.Repository{
		{ID: 1, FullName: "user/repo1", SSHURL: "git@github.com:user/repo1.git"},
		{ID: 2, FullName: "user/repo2", SSHURL: "git@github.com:user/repo2.git"},
	}

	ops := newMockGitOps()
	setupMocks(t, repos, nil, ops)

	err := run(dir)

	assert.NoError(t, err)
	assert.Len(t, ops.updateCalls, 1)
	assert.Equal(t, filepath.Join(dir, "repo1.git"), ops.updateCalls[0])
	assert.Len(t, ops.cloneCalls, 1)
	assert.Equal(t, "git@github.com:user/repo2.git", ops.cloneCalls[0].sshURL)
	assert.Equal(t, filepath.Join(dir, "repo2.git"), ops.cloneCalls[0].targetDirectory)
}

func TestRun_CloneErrorContinues(t *testing.T) {
	dir := t.TempDir()

	repos := []github.Repository{
		{ID: 1, FullName: "user/repo1", SSHURL: "git@github.com:user/repo1.git"},
		{ID: 2, FullName: "user/repo2", SSHURL: "git@github.com:user/repo2.git"},
	}

	ops := newMockGitOps()
	ops.cloneErrForRepository[filepath.Join(dir, "repo1.git")] = errors.New("clone failed")
	setupMocks(t, repos, nil, ops)

	err := run(dir)

	assert.NoError(t, err)
	assert.Len(t, ops.cloneCalls, 2)
}

func TestRun_UpdateErrorContinues(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "repo1.git"), 0755)
	os.MkdirAll(filepath.Join(dir, "repo2.git"), 0755)

	repos := []github.Repository{
		{ID: 1, FullName: "user/repo1", SSHURL: "git@github.com:user/repo1.git"},
		{ID: 2, FullName: "user/repo2", SSHURL: "git@github.com:user/repo2.git"},
	}

	ops := newMockGitOps()
	ops.updateErrForRepository[filepath.Join(dir, "repo1.git")] = errors.New("update failed")
	setupMocks(t, repos, nil, ops)

	err := run(dir)

	assert.NoError(t, err)
	assert.Len(t, ops.updateCalls, 2)
}

func TestRun_CreatesBackupDir(t *testing.T) {
	parent := t.TempDir()
	dir := filepath.Join(parent, "nested", "backup")

	ops := newMockGitOps()
	setupMocks(t, []github.Repository{}, nil, ops)

	err := run(dir)

	assert.NoError(t, err)
	info, statErr := os.Stat(dir)
	assert.NoError(t, statErr)
	assert.True(t, info.IsDir())
}

func TestRun_AllClonesFail(t *testing.T) {
	dir := t.TempDir()

	repos := []github.Repository{
		{ID: 1, FullName: "user/repo1", SSHURL: "git@github.com:user/repo1.git"},
		{ID: 2, FullName: "user/repo2", SSHURL: "git@github.com:user/repo2.git"},
	}

	ops := newMockGitOps()
	ops.cloneErr = errors.New("clone failed")
	setupMocks(t, repos, nil, ops)

	err := run(dir)

	assert.NoError(t, err)
	assert.Len(t, ops.cloneCalls, 2)
}

func TestRun_FileExistsButNotDir(t *testing.T) {
	dir := t.TempDir()

	filePath := filepath.Join(dir, "repo1.git")
	os.WriteFile(filePath, []byte("not a dir"), 0644)

	repos := []github.Repository{
		{ID: 1, FullName: "user/repo1", SSHURL: "git@github.com:user/repo1.git"},
	}

	ops := newMockGitOps()
	setupMocks(t, repos, nil, ops)

	err := run(dir)

	assert.NoError(t, err)
	assert.Len(t, ops.cloneCalls, 1)
	assert.Empty(t, ops.updateCalls)
}
