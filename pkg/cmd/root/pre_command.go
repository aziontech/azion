package root

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/aziontech/azion-cli/pkg/github"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/root"
	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/metric"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

type OSInfo struct {
	OS struct {
		Vendor string `json:"vendor"`
	} `json:"os"`
}

// doPreCommandCheck carry out all pre-cmd checks needed
func doPreCommandCheck(cmd *cobra.Command, fact *factoryRoot) error {

	// Apply timeout based on the parsed flags
	timeout, err := cmd.Flags().GetInt("timeout")
	if err != nil {
		return msg.ErrorParseTimeout
	}

	fact.factory.HttpClient.Timeout = time.Duration(timeout) * time.Second

	// get full command run and rewrite with our metrics pattern
	fact.commandName = cmd.CommandPath()
	rewrittenCommand := strings.ReplaceAll(strings.TrimPrefix(fact.commandName, "azion "), " ", "-")
	fact.commandName = rewrittenCommand

	if cmd.Flags().Changed("config") {
		if err := config.SetPath(fact.configFlag); err != nil {
			return err
		}
	}

	settings, err := token.ReadSettings()
	if err != nil {
		return err
	}
	fact.globalSettings = &settings

	t := token.New(&token.Config{
		Client: fact.factory.HttpClient,
		Out:    fact.factory.IOStreams.Out,
	})

	if cmd.Flags().Changed("token") {
		if err := checkTokenSent(fact, fact.globalSettings, t); err != nil {
			return err
		}
	}

	if err := checkAuthorizeMetricsCollection(cmd, fact.factory.GlobalFlagAll, fact.globalSettings); err != nil {
		return err
	}

	//both verifications occurs if 24 hours have passed since the last execution
	if err := checkForUpdateAndMetrics(version.BinVersion, fact.factory, fact.globalSettings); err != nil {
		return err
	}

	return nil
}

func checkTokenSent(fact *factoryRoot, settings *token.Settings, tokenStr *token.Token) error {
	if fact.tokenFlag == "" {
		return utils.ErrorTokenNotProvided
	}

	valid, user, err := tokenStr.Validate(&fact.tokenFlag)
	if err != nil {
		return err
	}

	if !valid {
		return utils.ErrorInvalidToken
	}

	strToken := token.Settings{
		Token:                      fact.tokenFlag,
		ClientId:                   user.Results.ClientID,
		Email:                      user.Results.Email,
		AuthorizeMetricsCollection: settings.AuthorizeMetricsCollection,
		S3AccessKey:                "",
		S3SecretKey:                "",
		S3Bucket:                   "",
	}

	bStrToken, err := toml.Marshal(strToken)
	if err != nil {
		return err
	}

	filePath, err := tokenStr.Save(bStrToken)
	if err != nil {
		return err
	}

	fact.globalSettings = &strToken

	logger.FInfo(fact.factory.IOStreams.Out, fmt.Sprintf(msg.TokenSavedIn, filePath))
	logger.FInfo(fact.factory.IOStreams.Out, msg.TokenUsedIn+"\n")
	return nil
}

func checkForUpdateAndMetrics(cVersion string, f *cmdutil.Factory, settings *token.Settings) error {
	logger.Debug("Verifying if an update is required")
	// checks if 24 hours have passed since the last check
	if time.Since(settings.LastCheck) < 24*time.Hour && !settings.LastCheck.IsZero() {
		return nil
	}

	// checks if user is Logged in before sending metrics
	if verifyUserInfo(settings) {
		metric.Send(settings)
	}

	git := github.NewGithub()

	tagName, publishedAt, err := git.GetVersionGitHub("azion")
	if err != nil {
		return err
	}

	logger.Debug("Current version: " + cVersion)
	logger.Debug("Latest version: " + tagName)
	logger.Debug("Published at: " + publishedAt)

	latestVersion, err := format(tagName)
	if err != nil {
		return err
	}
	currentVersion, err := format(cVersion)
	if err != nil {
		return err
	}

	if latestVersion > currentVersion {
		// Parse the published_at date
		publishedTime, err := time.Parse(time.RFC3339, publishedAt)
		if err != nil {
			logger.Debug("Failed to parse published_at date", zap.Error(err))
			// If we can't parse the date, fall back to showing the update message
			err := showUpdateMessage(f, tagName)
			if err != nil {
				return err
			}
		} else {
			// Only show update message if at least 24 hours have passed since publishing
			if time.Since(publishedTime) >= 24*time.Hour {
				err := showUpdateMessage(f, tagName)
				if err != nil {
					return err
				}
			}
		}
	}

	// Update the last update check time
	settings.LastCheck = time.Now()
	if err := token.WriteSettings(*settings); err != nil {
		return err
	}

	return nil
}

func showUpdateMessage(f *cmdutil.Factory, vNumber string) error {
	logger.FInfo(f.IOStreams.Out, msg.NewVersion)

	err := showUpdadeMessageSystem(f, vNumber)
	if err != nil {
		return err
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

// Mock for utils.Confirm
var confirmFn = utils.Confirm

// 0 = authorization was not asked yet, 1 = accepted, 2 = denied
func checkAuthorizeMetricsCollection(cmd *cobra.Command, globalFlagAll bool, settings *token.Settings) error {
	if settings.AuthorizeMetricsCollection > 0 || cmd.Name() == "completion" {
		return nil
	}

	authorize := confirmFn(globalFlagAll, msg.AskCollectMetrics, true)
	if authorize {
		settings.AuthorizeMetricsCollection = 1
	} else {
		settings.AuthorizeMetricsCollection = 2
	}

	if err := token.WriteSettings(*settings); err != nil {
		return err
	}

	return nil
}

func verifyUserInfo(settings *token.Settings) bool {
	if settings.ClientId != "" && settings.Email != "" {
		return true
	}

	return false
}
