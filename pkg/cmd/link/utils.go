package link

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/vulcan"
	vul "github.com/aziontech/azion-cli/pkg/vulcan"
	helpers "github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func shouldConfigure(info *LinkInfo) bool {
	if info.GlobalFlagAll || info.Auto {
		return true
	}
	msg := fmt.Sprintf("Do you want to link %s to Azion? (y/N)", info.PathWorkingDir)
	return helpers.Confirm(msg, false)
}

func shouldDevDeploy(info *LinkInfo, msg string, defaultYes bool) bool {
	if info.GlobalFlagAll {
		return true
	}
	return helpers.Confirm(msg, defaultYes)
}

func shouldFetch(cmd *LinkCmd, info *LinkInfo) (bool, error) {
	var err error
	var shouldFetchTemplates bool
	if empty, _ := cmd.IsDirEmpty("./azion"); !empty {
		if info.GlobalFlagAll || info.Auto {
			shouldFetchTemplates = true
		} else {
			return helpers.Confirm("This project was already configured. Do you want to override the previous configuration? (y/N)", false), nil
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

func (cmd *LinkCmd) selectVulcanMode(info *LinkInfo) error {
	if info.Preset == "nextjs" {
		return nil
	}

	// checking is vulcan major is correct
	vulcanVer, err := cmd.CommandRunner(cmd.F, "npm show edge-functions version", []string{})
	if err != nil {
		return err
	}

	err = vulcan.CheckVulcanMajor(vulcanVer, cmd.F)
	if err != nil {
		return err
	}

	logger.FInfo(cmd.Io.Out, msg.InitGettingTemplates)

	command := vul.Command("--loglevel=error --no-update-notifier", "presets ls")

	output, err := cmd.CommandRunner(cmd.F, command, []string{"CLEAN_OUTPUT_MODE=true"})
	if err != nil {
		return err
	}

	newLineSplit := strings.Split(output, "\n")
	newLineSplit[len(newLineSplit)-1] = "static (azion)"

	answer := ""
	template := ""
	mode := ""
	prompt := &survey.Select{
		Message: "Choose a preset and mode:",
		Options: newLineSplit,
	}
	err = survey.AskOne(prompt, &answer)
	if err != nil {
		return err
	}

	modeSplit := strings.Split(answer, " ")
	template = modeSplit[0]
	mode = strings.Replace(strings.Replace(modeSplit[1], "(", "", -1), ")", "", -1)

	info.Preset = template
	info.Mode = mode

	return nil
}

func depsInstall(cmd *LinkCmd, packageManager string) error {
	logger.FInfo(cmd.Io.Out, msg.InitInstallDeps)
	command := fmt.Sprintf("%s install", packageManager)
	err := cmd.CommandRunInteractive(cmd.F, command)
	if err != nil {
		logger.Debug("Error while running command with simultaneous output", zap.Error(err))
		return msg.ErrorDeps
	}

	return nil
}
