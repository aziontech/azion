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

type Release struct {
	TagName string `json:"tag_name"`
}

func GetVersionGitHub(name string) (string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/aziontech/%s/releases/latest", name)

	response, err := http.Get(apiURL)
	if err != nil {
		logger.Debug("Failed to get latest version of "+name, zap.Error(err))
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logger.Debug("Failed to get latest version of "+name, zap.Error(err))
		return "", nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Debug("Failed to read response body", zap.Error(err))
		return "", err
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		logger.Debug("Failed to unmarshal response body", zap.Error(err))
		return "", err
	}

	return release.TagName, nil
}

// Clone clone the repository using git
func Clone(url, path string) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return err
	}
	if err := os.Chdir(path); err != nil {
		return fmt.Errorf(utils.ERROR_CDDIR, err.Error())
	}
	return nil
}

// GetNameRepoFunction to get the repository name from the URL
func GetNameRepo(url string) string {
	// Remove the initial part of the URL
	parts := strings.Split(url, "/")
	repoPart := parts[len(parts)-1]
	// Remove the .git folder if it exists.
	repoPart = strings.TrimSuffix(repoPart, ".git")
	return repoPart
}

func CheckGitignore(path string) (bool, error) {
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

	if !object.MatchesPath(".edge/") || !object.MatchesPath(".vulcan") {
		return false, nil
	}

	return true, nil
}

func WriteGitignore(path string) error {
	logger.Debug("Writing .gitignore file")
	path = filepath.Join(path, ".gitignore")

	// Lines to add to .gitignore
	linesToAdd := []string{"#Paths added by Azion CLI", ".edge/", ".vulcan"}

	// Open the file in append mode, create if not exists
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

	// Ensure all data is written to the file
	if err := writer.Flush(); err != nil {
		logger.Error("Error flushing writer", zap.Error(err))
		return err
	}

	return nil
}
