package github

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/go-git/go-git/v5"
	"go.uber.org/zap"

	gitignore "github.com/sabhiram/go-gitignore"
)

type Github struct {
	GetVersionGitHub  func(name string) (string, string, error)
	Clone             func(cloneOptions *git.CloneOptions, url, path string) error
	GetNameRepo       func(url string) string
	CheckGitignore    func(path string) (bool, error)
	WriteGitignore    func(path string) error
	CheckWorkflowFile func(path string) (bool, error)
	WriteWorkflowFile func(path string) error
}

type Release struct {
	TagName     string `json:"tag_name"`
	PublishedAt string `json:"published_at"`
}

var (
	ApiURL string
)

func NewGithub() *Github {
	return &Github{
		GetVersionGitHub:  getVersionGitHub,
		Clone:             clone,
		GetNameRepo:       getNameRepo,
		CheckGitignore:    checkGitignore,
		WriteGitignore:    writeGitignore,
		CheckWorkflowFile: checkWorkflowFile,
		WriteWorkflowFile: writeWorkflowFile,
	}
}

func getVersionGitHub(name string) (string, string, error) {
	ApiURL = fmt.Sprintf("https://api.github.com/repos/aziontech/%s/releases/latest", name)

	response, err := http.Get(ApiURL)
	if err != nil {
		logger.Debug("Failed to get latest version of "+name, zap.Error(err))
		return "", "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logger.Debug("Failed to get latest version of "+name, zap.Error(err))
		return "", "", nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Debug("Failed to read response body", zap.Error(err))
		return "", "", err
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		logger.Debug("Failed to unmarshal response body", zap.Error(err))
		return "", "", err
	}

	return release.TagName, release.PublishedAt, nil
}

func clone(cloneOptions *git.CloneOptions, url, path string) error {
	cloneOptions.URL = url
	_, err := git.PlainClone(path, false, cloneOptions)
	if err != nil {
		return err
	}
	if err := os.Chdir(path); err != nil {
		return fmt.Errorf(utils.ERROR_CDDIR, err.Error())
	}
	return nil
}

func getNameRepo(url string) string {
	parts := strings.Split(url, "/")
	repoPart := parts[len(parts)-1]
	repoPart = strings.TrimSuffix(repoPart, ".git")
	return repoPart
}

func checkGitignore(path string) (bool, error) {
	logger.Debug("Checking .gitignore file for existence of Vulcan files")
	path = filepath.Join(path, ".gitignore")

	object, err := gitignore.CompileIgnoreFile(path)
	if err != nil {
		// if the error is "no such file or directory" we can return false and nil for error, because the code that called this func will create
		// the .gitignore file
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	if !object.MatchesPath(".edge/") || !object.MatchesPath(".vulcan") || !object.MatchesPath(".open-next") {
		return false, nil
	}

	return true, nil
}

func writeGitignore(path string) error {
	logger.Debug("Writing .gitignore file")
	path = filepath.Join(path, ".gitignore")

	linesToAdd := []string{"#Paths added by Azion CLI", ".edge/", ".vulcan", ".open-next"}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error("Error opening file", zap.Error(err))
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, line := range linesToAdd {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			logger.Error("Error writing to .gitignore file", zap.Error(err))
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		logger.Error("Error flushing writer", zap.Error(err))
		return err
	}

	return nil
}

func checkWorkflowFile(path string) (bool, error) {
	logger.Debug("Checking for GitHub Actions workflow file")
	workflowPath := filepath.Join(path, ".github", "workflows", "azion-deploy.yml")

	_, err := os.Stat(workflowPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func writeWorkflowFile(path string) error {
	logger.Debug("Writing GitHub Actions workflow file")
	workflowDir := filepath.Join(path, ".github", "workflows")
	workflowPath := filepath.Join(workflowDir, "azion-deploy.yml")

	// Create .github/workflows directory if it doesn't exist
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		logger.Error("Error creating workflow directory", zap.Error(err))
		return err
	}

	// Write the workflow file
	if err := os.WriteFile(workflowPath, []byte(workflowContent), 0644); err != nil {
		logger.Error("Error writing workflow file", zap.Error(err))
		return err
	}

	return nil
}
