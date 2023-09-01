package init

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/logger"
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
	logger.FInfo(cmd.Io.Out, msg.InitGettingTemplates)

	err := cmd.CommandRunInteractive(cmd.F, []string{}, "npx --yes edge-functions@1.5.0 init --name "+info.Name)
	if err != nil {
		return err
	}

	output, _, err := cmd.CommandRunner("npx --yes edge-functions@1.5.0 presets ls", []string{"CLEAN_OUTPUT_MODE=true"})
	if err != nil {
		return err
	}

	newLineSplit := strings.Split(output, "\n")
	newLineSplit[len(newLineSplit)-1] = "static (azion)"

	answer := ""
	template := ""
	mode := ""
	prompt := &survey.Select{
		Message: "Choose a mode:",
		Options: newLineSplit,
	}
	err = survey.AskOne(prompt, &answer)
	if err != nil {
		return err
	}

	modeSplit := strings.Split(answer, " ")
	template = modeSplit[0]
	mode = strings.Replace(strings.Replace(modeSplit[1], "(", "", -1), ")", "", -1)

	info.Template = template
	info.Mode = mode

	return nil
}

func yarnInstall(cmd *InitCmd) error {

	logger.FInfo(cmd.Io.Out, msg.InitInstallDeps)

	err := cmd.CommandRunInteractive(cmd.F, []string{}, "yarn install")
	if err != nil {
		logger.Debug("Error while running command with simultaneous output", zap.Error(err))
		return msg.ErrorDeps
	}

	return nil
}
