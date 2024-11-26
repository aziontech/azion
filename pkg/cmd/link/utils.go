package link

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/logger"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	helpers "github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func shouldConfigure(info *LinkInfo) bool {
	if info.Auto {
		return true
	}
	msg := fmt.Sprintf("Do you want to link %s to Azion? (y/N)", info.PathWorkingDir)
	return helpers.Confirm(info.GlobalFlagAll, msg, false)
}

func shouldDevDeploy(info *LinkInfo, msg string, defaultYes bool) bool {
	return helpers.Confirm(info.GlobalFlagAll, msg, defaultYes)
}

func shouldFetch(cmd *LinkCmd, info *LinkInfo) (bool, error) {
	var err error
	var shouldFetchTemplates bool
	if empty, _ := cmd.IsDirEmpty(info.projectPath); !empty {
		if info.GlobalFlagAll || info.Auto {
			shouldFetchTemplates = true
		} else {
			return helpers.Confirm(info.GlobalFlagAll, "This project was already configured. Do you want to override the previous configuration? (y/N)", false), nil
		}

		if shouldFetchTemplates {
			err = cmd.CleanDir(info.projectPath)
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
	err := survey.AskOne(prompt, &userInput, survey.WithKeepFilter(true), survey.WithStdio(os.Stdin, os.Stderr, os.Stdout))
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

	vul := vulcanPkg.NewVulcan()
	err = vul.CheckVulcanMajor(vulcanVer, cmd.F, vul)
	if err != nil {
		return err
	}
	logger.FInfo(cmd.Io.Out, msg.InitGettingTemplates)

	command := vul.Command("--loglevel=error --no-update-notifier", "presets ls", cmd.F)
	output, err := cmd.CommandRunner(cmd.F, command, []string{"CLEAN_OUTPUT_MODE=true"})
	if err != nil {
		return err
	}

	// The list that comes from Vulcan has a blank line that we should remove.
	outputInline := strings.Split(output, "\n")
	noLastItem := len(outputInline) - 1
	listPresets := make([]string, noLastItem)
	copy(listPresets, outputInline[:noLastItem])

	prompt := &survey.Select{
		Message:  "Choose a preset:",
		Options:  listPresets,
		PageSize: len(listPresets),
	}

	var answer string
	err = survey.AskOne(prompt, &answer)
	if err != nil {
		return err
	}

	info.Preset = answer

	return nil
}

func depsInstall(cmd *LinkCmd, packageManager string) error {
	command := fmt.Sprintf("%s install", packageManager)
	err := cmd.CommandRunInteractive(cmd.F, command)
	if err != nil {
		logger.Debug("Error while running command with simultaneous output", zap.Error(err))
		return msg.ErrorDeps
	}

	return nil
}
