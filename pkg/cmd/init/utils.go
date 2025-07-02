package init

import (
	"fmt"
	"os"
	"path"
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
	err := cmd.askOne(prompt, &userInput, survey.WithKeepFilter(true), survey.WithStdio(os.Stdin, os.Stderr, os.Stdout))
	if err != nil {
		return "", err
	}

	return userInput, nil
}

func (cmd *initCmd) selectVulcanTemplates(vul *vulcanPkg.VulcanPkg) error {
	logger.Debug("Running bundler store init")
	// checking if vulcan major is correct
	vulcanVer, err := cmd.commandRunnerOutput(cmd.f, "npm show edge-functions version", []string{})
	if err != nil {
		return err
	}

	err = vul.CheckVulcanMajor(vulcanVer, cmd.f, vul)
	if err != nil {
		return err
	}

	preset, err := cmd.getVulcanInfo()
	if err != nil {
		return err
	}

	if preset == strings.ToLower("vite") {
		preset = "vue"
	}

	cmd.preset = strings.ToLower(preset)
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

func (cmd *initCmd) getVulcanInfo() (string, error) {

	fileContent, err := cmd.fileReader(path.Join(cmd.pathWorkingDir, "info.json"))
	if err != nil {
		logger.Debug("Error reading template info", zap.Error(err))
		return "", err
	}

	var infoJson map[string]string
	err = cmd.unmarshal(fileContent, &infoJson)
	if err != nil {
		logger.Debug("Error unmarshalling template info", zap.Error(err))
		return "", err
	}

	logger.Debug("Information about the template:", zap.Any("preset", infoJson["preset"]))
	return infoJson["preset"], nil
}
