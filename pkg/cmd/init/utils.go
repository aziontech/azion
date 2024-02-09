package init

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/logger"
	vul "github.com/aziontech/azion-cli/pkg/vulcan"
	helpers "github.com/aziontech/azion-cli/utils"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func shouldDevDeploy(info *InitInfo, msg string, defaultYes bool) bool {
	if info.GlobalFlagAll {
		return true
	}
	return helpers.Confirm(msg, defaultYes)
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

	// checking if vulcan major is correct
	vulcanVer, err := cmd.CommandRunnerOutput(cmd.F, "npm show edge-functions version", []string{})
	if err != nil {
		return err
	}

	err = vul.CheckVulcanMajor(vulcanVer, cmd.F)
	if err != nil {
		return err
	}

	logger.FInfo(cmd.Io.Out, msg.InitGettingVulcan)

	command := vul.Command("", "init --name "+info.Name)

	err = cmd.CommandRunInteractive(cmd.F, command)
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

	command = vul.Command("", "presets ls --preset "+preset)
	output, _, err := cmd.CommandRunner(command, []string{"CLEAN_OUTPUT_MODE=true"})
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
		info.Template = strings.ToLower(preset)
		info.Mode = strings.ToLower(answer)
		return nil
	}

	if len(newLineSplit) < 1 {
		logger.Debug("No mode found for the selected preset: "+preset, zap.Error(err))
		return msg.ErrorModeNotFound
	}

	var mds string = newLineSplit[0]

	info.Template = strings.ToLower(preset)
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
