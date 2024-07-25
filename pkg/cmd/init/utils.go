package init

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/logger"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
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

	vul := vulcanPkg.NewVulcan()
	err = vul.CheckVulcanMajor(vulcanVer, cmd.f, vul)
	if err != nil {
		return err
	}

	cmdVulcanInit := "init"
	if len(cmd.preset) > 0 {
		cmdVulcanInit = fmt.Sprintf("%s --preset '%s'", cmdVulcanInit, cmd.preset)
	}
	if len(cmd.mode) > 0 {
		cmdVulcanInit = fmt.Sprintf("%s --mode '%s'", cmdVulcanInit, cmd.mode)
	}
	if len(cmd.pathWorkingDir) > 0 {
		cmdVulcanInit = fmt.Sprintf("%s --scope '%s'", cmdVulcanInit, cmd.pathWorkingDir)
	}

	command := vul.Command("", cmdVulcanInit, cmd.f)

	err = cmd.commandRunInteractive(cmd.f, command)
	if err != nil {
		return err
	}

	preset, mode, err := getVulcanEnvInfo(cmd)
	if err != nil {
		return err
	}

	if preset == strings.ToLower("vite") {
		preset = "vue"
	}

	cmd.preset = strings.ToLower(preset)
	cmd.mode = strings.ToLower(mode)
	return nil
}

func depsInstall(cmd *initCmd, packageManager string) error {
	logger.FInfoFlags(cmd.io.Out, msg.InitInstallDeps, cmd.f.Format, cmd.f.Out)
	command := fmt.Sprintf("%s install", packageManager)
	err := cmd.commandRunInteractive(cmd.f, command)
	if err != nil {
		logger.Debug("Error while running command with simultaneous output", zap.Error(err))
		return msg.ErrorDeps
	}

	return nil
}

func getVulcanEnvInfo(cmd *initCmd) (string, string, error) {
	err := godotenv.Load(cmd.pathWorkingDir + "/.vulcan")
	if err != nil {
		logger.Debug("Error loading .vulcan file", zap.Error(err))
		return "", "", err
	}

	// Access environment variables
	preset := os.Getenv("preset")
	mode := os.Getenv("mode")
	return preset, mode, nil
}
