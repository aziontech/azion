package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
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

func Clone(url, path string) error {
	cmd := exec.Command("git", "clone", url)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(utils.ERROR_CLONE, err.Error())
	}
	if err := os.Chdir(path); err != nil {
		return fmt.Errorf(utils.ERROR_CDDIR, err.Error())
	}
	return nil
}
