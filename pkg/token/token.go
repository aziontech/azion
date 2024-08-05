package token

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aziontech/azion-cli/messages/root"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/pelletier/go-toml/v2"
	"go.uber.org/zap"

	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/constants"
)

func New(c *Config) (*Token, error) {
	dir, err := config.Dir()
	if err != nil {
		return nil, err
	}

	return &Token{
		client:   c.Client,
		Endpoint: constants.AuthURL,
		filePath: filepath.Join(dir.Dir, dir.Settings),
		out:      c.Out,
	}, nil
}

func (t *Token) Validate(token *string) (bool, UserInfo, error) {
	logger.Debug("Validate token")

	req, err := http.NewRequest("GET", utils.Concat(t.Endpoint, "/user/me"), nil)
	if err != nil {
		return false, UserInfo{}, err
	}
	req.Header.Add("Accept", "application/json; version=3")
	req.Header.Add("Authorization", "token "+*token)

	resp, err := t.client.Do(req)
	if err != nil {
		return false, UserInfo{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return false, UserInfo{}, nil
	}

	var userInfo UserInfo
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return false, UserInfo{}, err
	}

	t.valid = true

	return true, userInfo, nil
}

func (t *Token) Save(b []byte) (string, error) {
	dir, err := config.Dir()
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(dir.Dir, os.ModePerm)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(t.filePath, b, 0777)
	if err != nil {
		return "", err
	}

	return t.filePath, nil
}

func (t *Token) Create(b64 string) (*Response, error) {
	logger.Debug("Create token", zap.Any("base64", b64))
	req, err := http.NewRequest(http.MethodPost, utils.Concat(t.Endpoint, "/tokens"), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json; version=3")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", utils.Concat("Basic ", b64))

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp != nil {
		logger.Debug("Error while creating token", zap.Error(err))
		err := utils.LogAndRewindBody(resp)
		if err != nil {
			return nil, err
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func WriteSettings(settings Settings) error {
	dir, err := config.Dir()
	if err != nil {
		return fmt.Errorf("Failed to get token dir: %w", err)
	}

	b, err := toml.Marshal(settings)
	if err != nil {
		return err
	}

	// Check if the directory exists, create it if not
	if err := os.MkdirAll(dir.Dir, 0777); err != nil {
		return fmt.Errorf("Error creating directory: %w", err)
	}

	if err := os.WriteFile(filepath.Join(dir.Dir, dir.Settings), b, 0777); err != nil {
		return fmt.Errorf(utils.ErrorWriteSettings.Error(), err)
	}

	return nil
}

func ReadSettings() (Settings, error) {
	dir, err := config.Dir()
	if err != nil {
		return Settings{}, fmt.Errorf("failed to get token dir: %w", err)
	}

	filePath := filepath.Join(dir.Dir, dir.Settings)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File does not exist, create it with default settings
		if config.GetPath() == config.DEFAULT_DIR {
			defaultSettings := Settings{}
			err := WriteSettings(defaultSettings)
			if err != nil {
				return Settings{}, fmt.Errorf("failed to create settings file: %w", err)
			}
			return defaultSettings, nil
		}

		return Settings{}, root.ErrorReadFileSettingsToml
	}

	// Read the file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return Settings{}, err
	}

	var settings Settings
	err = toml.Unmarshal(fileData, &settings)
	if err != nil {
		return Settings{}, fmt.Errorf("Failed parse byte to struct settings: %w", err)
	}

	return settings, nil
}
