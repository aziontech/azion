package init

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/logger"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"go.uber.org/zap"
)

func (cmd *initCmd) askForInput(msg string, defaultIn string) (string, error) {
	var userInput string
	prompt := &survey.Input{
		Message: msg,
		Default: defaultIn,
	}

	// Prompt the user for input
	err := cmd.askOne(prompt, &userInput, survey.WithKeepFilter(true))
	if err != nil {
		return "", err
	}

	return userInput, nil
}

func (cmd *initCmd) selectVulcanTemplates(vul *vulcanPkg.VulcanPkg) error {
	// checking if vulcan major is correct
	vulcanVer, err := cmd.commandRunnerOutput(cmd.f, "npm show edge-functions version", []string{})
	if err != nil {
		return err
	}

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

	preset, mode, err := cmd.getVulcanEnvInfo()
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

func (cmd *initCmd) depsInstall() error {
	command := fmt.Sprintf("%s install", cmd.packageManager)
	err := cmd.commandRunInteractive(cmd.f, command)
	if err != nil {
		logger.Debug("Error while running command with simultaneous output", zap.Error(err))
		return msg.ErrorDeps
	}
	return nil
}

func (cmd *initCmd) getVulcanEnvInfo() (string, string, error) {
	err := cmd.load(cmd.pathWorkingDir + "/.vulcan")
	if err != nil {
		logger.Debug("Error loading .vulcan file", zap.Error(err))
		return "", "", err
	}

	// Access environment variables
	preset := os.Getenv("preset")
	mode := os.Getenv("mode")
	return preset, mode, nil
}
