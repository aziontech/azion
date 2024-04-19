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

func shouldDevDeploy(msg string, globalFlagAll bool, defaultYes bool) bool {
	return helpers.Confirm(globalFlagAll, msg, defaultYes)
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

func (cmd *initCmd) selectVulcanTemplates() error {
	// checking if vulcan major is correct
	vulcanVer, err := cmd.commandRunnerOutput(cmd.f, "npm show edge-functions version", []string{})
	if err != nil {
		return err
	}

	err = vul.CheckVulcanMajor(vulcanVer, cmd.f)
	if err != nil {
		return err
	}

	logger.FInfo(cmd.io.Out, msg.InitGettingVulcan)

	cmdVulcanInit := fmt.Sprintf("init --name %s", cmd.name)
	if len(cmd.preset) > 0 {
		cmdVulcanInit = fmt.Sprintf("%s --preset '%s'", cmdVulcanInit, cmd.preset)
	}
	if len(cmd.template) > 0 {
		cmdVulcanInit = fmt.Sprintf("%s --template '%s'", cmdVulcanInit, cmd.template)
	}

	command := vul.Command("", cmdVulcanInit)

	err = cmd.commandRunInteractive(cmd.f, command)
	if err != nil {
		return err
	}

	preset, err := getVulcanEnvInfo(cmd)
	if err != nil {
		return err
	}

	if preset == strings.ToLower("vite") {
		preset = "vue"
	}

	command = vul.Command("", "presets ls --preset "+preset)
	output, _, err := cmd.commandRunner(command, []string{"CLEAN_OUTPUT_MODE=true"})
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
		cmd.preset = strings.ToLower(preset)
		cmd.mode = strings.ToLower(answer)
		return nil
	}

	if len(newLineSplit) < 1 {
		logger.Debug("No mode found for the selected preset: "+preset, zap.Error(err))
		return msg.ErrorModeNotFound
	}

	var mds string = newLineSplit[0]

	cmd.preset = strings.ToLower(preset)
	cmd.mode = strings.ToLower(mds)
	logger.FInfo(cmd.io.Out, fmt.Sprintf(msg.ModeAutomatic, mds, preset))
	return nil
}

func depsInstall(cmd *initCmd, packageManager string) error {
	logger.FInfo(cmd.io.Out, msg.InitInstallDeps)
	command := fmt.Sprintf("%s install", packageManager)
	err := cmd.commandRunInteractive(cmd.f, command)
	if err != nil {
		logger.Debug("Error while running command with simultaneous output", zap.Error(err))
		return msg.ErrorDeps
	}

	return nil
}

func getVulcanEnvInfo(cmd *initCmd) (string, error) {
	err := godotenv.Load(cmd.pathWorkingDir + "/.vulcan")
	if err != nil {
		logger.Debug("Error loading .vulcan file", zap.Error(err))
		return "", err
	}

	// Access environment variables
	preset := os.Getenv("preset")
	return preset, nil
}
