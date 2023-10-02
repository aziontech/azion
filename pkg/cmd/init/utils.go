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

	err := cmd.CommandRunInteractive(cmd.F, "npx --yes edge-functions@1.7.0 init --name "+info.Name)
	if err != nil {
		return err
	}

	preset, err := getVulcanEnvInfo(info)
	if err != nil {
		return err
	}

	if preset == strings.ToLower("vite") {
		preset = "vue"
	}

	output, _, err := cmd.CommandRunner("npx --yes edge-functions@1.7.0 presets ls --preset "+preset, []string{"CLEAN_OUTPUT_MODE=true"})
	if err != nil {
		return err
	}

	newLineSplit := strings.Split(output, "\n")
	if newLineSplit[len(newLineSplit)-1] == "" {
		newLineSplit = newLineSplit[:len(newLineSplit)-1]
	}

	answer := ""
	if len(newLineSplit) > 1 {
		prompt := &survey.Select{
			Message: "Choose a mode:",
			Options: newLineSplit,
		}
		err = survey.AskOne(prompt, &answer)
		if err != nil {
			return err
		}
		info.Template = preset
		info.Mode = answer
		return nil
	}

	if len(newLineSplit) < 1 {
		logger.Debug("No mode found for the selected preset: "+preset, zap.Error(err))
		return msg.ErrorModeNotFound
	}

	var mds string = newLineSplit[0]

	info.Template = preset
	info.Mode = strings.ToLower(mds)
	logger.FInfo(cmd.Io.Out, fmt.Sprintf(msg.ModeAutomatic, mds, preset))

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
