package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

var ghFile *os.File

const githubTestToken = "test-github-token"

func reset() {
	instance = nil
	once = sync.Once{}
	loadErr = nil
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) SetupTest() {
	reset()
	var err error
	secretsBasePath = s.T().TempDir()
	path := filepath.Join(secretsBasePath, "github-token.txt")
	ghFile, err = os.Create(path)
	if err != nil {
		log.Fatalf("Could not create github token text file, e:%e", err)
	}

	data := []byte("test-github-token\n")
	err = os.WriteFile(path, data, 0644)

	if err != nil {
		log.Fatalf("Could not write to test file, e:%e", err)
	}
}

func (s *ConfigTestSuite) TearDownTest() {
	reset()

	os.Remove(ghFile.Name())
}

func (s *ConfigTestSuite) TestGet() {
	cfg, err := Get()

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), version, cfg.Version)
	assert.Equal(s.T(), githubTestToken, cfg.GitHubToken)
}

func (s *ConfigTestSuite) TestGet_MissingGitHubTokenFile() {
	os.Remove(ghFile.Name())
	cfg, err := Get()

	assert.EqualError(s.T(), err, "github-token.txt is missing OR github token cannot be an empty string")
	assert.Equal(s.T(), version, cfg.Version)
	assert.Equal(s.T(), "", cfg.GitHubToken)
}

func (s *ConfigTestSuite) TestGet_EmptyGitHubTokenFile() {
	ghFile.Truncate(0)
	cfg, err := Get()

	assert.EqualError(s.T(), err, "github-token.txt is missing OR github token cannot be an empty string")
	assert.Equal(s.T(), version, cfg.Version)
	assert.Equal(s.T(), "", cfg.GitHubToken)
}

func (s *ConfigTestSuite) TestGetGitVaultVersion() {
	localVersion := GetGitVaultVersion()

	assert.Equal(s.T(), version, localVersion)
}

func (s *ConfigTestSuite) TestGetGitHubToken() {
	token := GetGitHubToken()

	assert.Equal(s.T(), githubTestToken, token)
}
