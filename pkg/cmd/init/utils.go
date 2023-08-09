package init

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

func yesNoFlagToResponse(info *InitInfo) bool {
	if info.YesOption {
		return info.YesOption
	}

	return false
}

func shouldFetch(cmd *InitCmd, info *InitInfo) (bool, error) {
	var err error
	var shouldFetchTemplates bool
	if empty, _ := cmd.IsDirEmpty("./azion"); !empty {
		if info.NoOption || info.YesOption {
			shouldFetchTemplates = yesNoFlagToResponse(info)
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
	output, _, err := cmd.CommandRunner("npx --yes edge-functions@1.0.0 presets ls", []string{"CLEAN_OUTPUT_MODE=true"})
	if err != nil {
		return err
	}

	newLineSplit := strings.Split(output, "\n")
	newLineSplit[len(newLineSplit)-1] = "nextjs (faststore)"

	answer := ""
	template := ""
	mode := ""
	prompt := &survey.Select{
		Message: "Choose a template:",
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
