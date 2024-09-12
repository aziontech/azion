package root

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/aziontech/azion-cli/pkg/github"

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
		Client: fact.Factory.HttpClient,
		Out:    fact.Factory.IOStreams.Out,
	})

	if cmd.Flags().Changed("token") {
		if err := checkTokenSent(fact, fact.globalSettings, t); err != nil {
			return err
		}
	}

	if err := checkAuthorizeMetricsCollection(cmd, fact.Factory.GlobalFlagAll, fact.globalSettings); err != nil {
		return err
	}

	//both verifications occurs if 24 hours have passed since the last execution
	if err := checkForUpdateAndMetrics(version.BinVersion, fact.Factory, fact.globalSettings); err != nil {
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

	logger.FInfo(fact.Factory.IOStreams.Out, fmt.Sprintf(msg.TokenSavedIn, filePath))
	logger.FInfo(fact.Factory.IOStreams.Out, msg.TokenUsedIn+"\n")
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

	tagName, err := git.GetVersionGitHub("azion")
	if err != nil {
		return err
	}

	logger.Debug("Current version: " + cVersion)
	logger.Debug("Latest version: " + tagName)

	latestVersion, err := format(tagName)
	if err != nil {
		return err
	}
	currentVersion, err := format(cVersion)
	if err != nil {
		return err
	}

	logger.Debug("Formatted current version: " + fmt.Sprint(currentVersion))
	logger.Debug("Formatted latest version: " + fmt.Sprint(latestVersion))

	if latestVersion > currentVersion {
		err := showUpdateMessage(f, tagName)
		if err != nil {
			return err
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

// Mock para utils.Confirm
var confirmFn = func(globalFlagAll bool, msg string, defaultValue bool) bool {
	return defaultValue
}

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
