package init

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func shouldDevDeploy(info *InitInfo, msg string) (bool, error) {
	if info.GlobalFlagAll {
		return true, nil
	}
	var shouldConfigure bool
	prompt := &survey.Confirm{
		Message: msg,
	}
	err := survey.AskOne(prompt, &shouldConfigure)
	if err != nil {
		return false, err
	}
	return shouldConfigure, nil
}

func shouldFetch(cmd *InitCmd, info *InitInfo) (bool, error) {
	var err error
	var shouldFetchTemplates bool
	if empty, _ := cmd.IsDirEmpty("./azion"); !empty {
		if info.GlobalFlagAll {
			shouldFetchTemplates = info.GlobalFlagAll
		} else {
			prompt := &survey.Confirm{
				Message: "This project was already configured. Do you want to override the previous configuration?",
			}
			err := survey.AskOne(prompt, &shouldFetchTemplates)
			if err != nil {
				return false, err
			}
		}

		if shouldFetchTemplates {
			err = cmd.CleanDir("./azion")
			if err != nil {
				logger.Debug("Error while trying to clean azion directory", zap.Error(err))
				return false, err
			}
		}
		return shouldFetchTemplates, nil
	}
	return true, nil
}

func askForInput(msg string, defaultIn string) (string, error) {
	var userInput string
	prompt := &survey.Input{
		Message: msg,
		Default: defaultIn,
	}

	// Prompt the user for input
	err := survey.AskOne(prompt, &userInput, survey.WithKeepFilter(true))
	if err != nil {
		return "", err
	}
	return userInput, nil
}

func (cmd *InitCmd) selectVulcanTemplates(info *InitInfo) error {
	logger.FInfo(cmd.Io.Out, msg.InitGettingVulcan)

	err := cmd.CommandRunInteractive(cmd.F, "npx --yes edge-functions@1.6.0 init --name "+info.Name)
	if err != nil {
		return err
	}

	if info.Template == "" || info.Mode == "" {
		output, _, err := cmd.CommandRunner("npx --yes edge-functions@1.6.0 presets ls", []string{"CLEAN_OUTPUT_MODE=true"})
		if err != nil {
			return err
		}

		newLineSplit := strings.Split(output, "\n")
		preset, err := getVulcanEnvInfo(info)
		if err != nil {
			return err
		}

		var modes []string

		for _, line := range newLineSplit {
			if strings.Contains(strings.ToLower(line), strings.ToLower(preset)) {
				modeSplit := strings.Split(line, " ")
				modes = append(modes, strings.ToLower(strings.Replace(strings.Replace(modeSplit[1], "(", "", -1), ")", "", -1)))
			}
		}

		answer := ""
		if len(modes) > 1 {
			prompt := &survey.Select{
				Message: "Choose a mode:",
				Options: modes,
			}
			err = survey.AskOne(prompt, &answer)
			if err != nil {
				return err
			}
			info.Template = preset
			info.Mode = answer
			return nil
		}

		info.Template = preset
		info.Mode = strings.ToLower(modes[0])
		logger.FInfo(cmd.Io.Out, fmt.Sprintf(msg.ModeAutomatic, modes[0]))
	}

	return nil
}

func depsInstall(cmd *InitCmd, packageManager string) error {
	logger.FInfo(cmd.Io.Out, msg.InitInstallDeps)
	command := fmt.Sprintf("%s install", packageManager)
	err := cmd.CommandRunInteractive(cmd.F, command)
	if err != nil {
		logger.Debug("Error while running command with simultaneous output", zap.Error(err))
		return msg.ErrorDeps
	}

	return nil
}

func getVulcanEnvInfo(info *InitInfo) (string, error) {
	err := godotenv.Load(info.PathWorkingDir + "/.vulcan")
	if err != nil {
		logger.Debug("Error loading .vulcan file", zap.Error(err))
		return "", err
	}

	// Access environment variables
	preset := os.Getenv("preset")
	return preset, nil
}
