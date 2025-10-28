package token

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/pelletier/go-toml/v2"
	"go.uber.org/zap"

	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/constants"
)

type TokenInterface interface {
	Validate(token *string) (bool, UserInfo, error)
	Save(b []byte) (string, error)
	Create(b64 string) (*Response, error)
}

func New(c *Config) *Token {
	dir := config.Dir()

	return &Token{
		client:   c.Client,
		Endpoint: constants.AuthURL,
		filePath: filepath.Join(dir.Dir, dir.Settings), //TODO: here
		out:      c.Out,
	}
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
	dir := config.Dir()
	err := os.MkdirAll(dir.Dir, os.ModePerm)
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

func WriteSettings(settings Settings, subdir string) error {
	dir := config.Dir()
	b, err := toml.Marshal(settings)
	if err != nil {
		return err
	}

	if subdir != "" {
		dir.Dir = filepath.Join(dir.Dir, subdir)
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

func ReadProfiles() (Profile, string, error) {
	dir := config.Dir()
	filePath := filepath.Join(dir.Dir, dir.Profiles)

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return Profile{}, "", err
	}

	var profile Profile
	err = toml.Unmarshal(fileData, &profile)
	if err != nil {
		return Profile{}, "", fmt.Errorf("Failed parse byte to struct profile: %w", err)
	}

	settingsPath := filepath.Join(dir.Dir, dir.Profiles)

	return profile, settingsPath, nil
}

func WriteProfiles(profile Profile) error {
	dir := config.Dir()
	b, err := toml.Marshal(profile)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir.Dir, 0777); err != nil {
		return fmt.Errorf("Error creating directory: %w", err)
	}

	if err := os.WriteFile(filepath.Join(dir.Dir, dir.Profiles), b, 0777); err != nil {
		return fmt.Errorf(utils.ErrorWriteSettings.Error(), err)
	}

	return nil
}

func ReadSettings(path string) (Settings, error) {
	dir := config.Dir()
	if path != "" {
		dir.Dir = filepath.Join(dir.Dir, path)
	}
	filePath := filepath.Join(dir.Dir, dir.Settings)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File does not exist, create it with default settings for the profile
		defaultSettings := Settings{}
		err := WriteSettings(defaultSettings, path)
		if err != nil {
			return Settings{}, fmt.Errorf(utils.ErrorWriteSettings.Error(), err)
		}
		return defaultSettings, nil
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
