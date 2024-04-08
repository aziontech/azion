package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/go-git/go-git/v5"
	"go.uber.org/zap"
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
	if strings.HasSuffix(repoPart, ".git") {
		repoPart = repoPart[:len(repoPart)-4]
	}
	return repoPart
}
