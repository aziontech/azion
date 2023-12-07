package root

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/user"
	"runtime"
	"strconv"
	"time"
	"unicode"

	msg "github.com/aziontech/azion-cli/messages/root"
	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"github.com/zcalusic/sysinfo"
	"go.uber.org/zap"
)

type PreCmd struct {
	token  string
	config string
}

type OSInfo struct {
	OS struct {
		Vendor string `json:"vendor"`
	} `json:"os"`
}

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

// doPreCommandCheck carry out all pre-cmd checks needed
func doPreCommandCheck(cmd *cobra.Command, f *cmdutil.Factory, pre PreCmd) error {

	if err := setConfigPath(cmd, pre.config); err != nil {
		return err
	}

	if err := checkTokenSent(cmd, f, pre.token); err != nil {
		return err
	}

	if err := checkForUpdate(version.BinVersion, f); err != nil {
		return err
	}

	return nil
}

func setConfigPath(cmd *cobra.Command, cfg string) error {

	if cmd.Flags().Changed("config") {
		config.SetPath(cfg)
		return nil
	}

	return nil
}

func checkTokenSent(cmd *cobra.Command, f *cmdutil.Factory, configureToken string) error {

	// if global --token flag was sent, verify it and save it locally
	if cmd.Flags().Changed("token") {
		t, err := token.New(&token.Config{
			Client: f.HttpClient,
			Out:    f.IOStreams.Out,
		})
		if err != nil {
			return fmt.Errorf("%s: %w", utils.ErrorTokenManager, err)
		}

		if configureToken == "" {
			return utils.ErrorTokenNotProvided
		}

		valid, err := t.Validate(&configureToken)
		if err != nil {
			return err
		}

		if !valid {
			return utils.ErrorInvalidToken
		}

		strToken := token.Settings{Token: configureToken}
		bStrToken, err := toml.Marshal(strToken)
		if err != nil {
			return err
		}

		filePath, err := t.Save(bStrToken)
		if err != nil {
			return err
		}

		logger.LogSuccess(f.IOStreams.Out, fmt.Sprintf(msg.TokenSavedIn, filePath))
		logger.FInfo(f.IOStreams.Out, msg.TokenUsedIn)
	}

	return nil
}

func checkForUpdate(cVersion string, f *cmdutil.Factory) error {
	config, err := token.ReadSettings()
	if err != nil {
		return err
	}
	// Check if 12 hours have passed since the last update check
	if time.Since(config.LastUpdateCheck) < 12*time.Hour && !config.LastUpdateCheck.IsZero() {
		return nil
	}

	apiURL := "https://api.github.com/repos/aziontech/azion/releases/latest"

	response, err := http.Get(apiURL)
	if err != nil {
		logger.Debug("Failed to get latest version of Azion CLI", zap.Error(err))
		return nil
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logger.Debug("Failed to get latest version of Azion CLI", zap.Error(err))
		return nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Debug("Failed to read response body", zap.Error(err))
		return nil
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		logger.Debug("Failed to unmarshal response body", zap.Error(err))
		return nil
	}

	latestVersion, err := format(release.TagName)
	if err != nil {
		return err
	}
	currentVersion, err := format(cVersion)
	if err != nil {
		return err
	}

	if latestVersion >= currentVersion {
		err := showUpdateMessage(f)
		if err != nil {
			return err
		}
	}

	// Update the last update check time
	config.LastUpdateCheck = time.Now()
	if err := token.WriteSettings(config); err != nil {
		return err
	}

	return nil
}

func showUpdateMessage(f *cmdutil.Factory) error {
	logger.FInfo(f.IOStreams.Out, msg.NewVersion)
	os := runtime.GOOS

	switch os {
	case "darwin":
		logger.FInfo(f.IOStreams.Out, msg.BrewUpdate)
		return nil
	case "linux":
		err := linuxUpdateMessage(f)
		if err != nil {
			return err
		}
	default:
		logger.FInfo(f.IOStreams.Out, msg.UnsupportedOS)
	}

	return nil
}

func format(input string) (int, error) {
	numberString := ""
	for _, char := range input {
		if unicode.IsDigit(char) {
			numberString += string(char)
		}
	}

	// to avoid converting errors
	if numberString == "" {
		numberString = "0"
	}

	number, err := strconv.Atoi(numberString)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func linuxUpdateMessage(f *cmdutil.Factory) error {
	current, err := user.Current()
	if err != nil {
		logger.Debug("Error while getting current user's information", zap.Error(err))
		return msg.ErrorCurrentUser
	}

	if current.Uid != "0" {
		logger.FInfo(f.IOStreams.Out, msg.CouldNotGetUser)
		return nil
	}

	var si sysinfo.SysInfo

	si.GetSysInfo()

	data, err := json.MarshalIndent(&si, "", "  ")
	if err != nil {
		logger.Debug("Error while marshaling current user's information", zap.Error(err))
		return msg.ErrorMarshalUserInfo
	}

	var osInfo OSInfo
	err = json.Unmarshal(data, &osInfo)
	if err != nil {
		logger.Debug("Error while unmarshaling current user's information", zap.Error(err))
		return msg.ErrorUnmarshalUserInfo
	}

	switch osInfo.OS.Vendor {
	case "debian":
		logger.FInfo(f.IOStreams.Out, msg.DpkgUpdate)
	case "alpine":
		logger.FInfo(f.IOStreams.Out, msg.ApkUpdate)
	case "centos", "fedora", "opensuse", "mageia", "mandriva":
		logger.FInfo(f.IOStreams.Out, msg.RpmUpdate)
	}

	return nil
}
