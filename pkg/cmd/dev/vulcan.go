package dev

import (
	"fmt"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/dev"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

func vulcan(cmd *DevCmd, port int) error {

	vul := cmd.Vulcan()
	baseCommand := vul.Command("", "dev", cmd.F)

	var commandBuilder strings.Builder
	commandBuilder.WriteString(baseCommand)

	if port > 0 {
		commandBuilder.WriteString(" --port ")
		commandBuilder.WriteString(fmt.Sprintf("%d", port))
	}

	if SkipFramework {
		commandBuilder.WriteString(" --skip-framework-build")
	}

	err := runCommand(cmd, commandBuilder.String())
	if err != nil {
		return fmt.Errorf(msg.ErrorVulcanExecute.Error(), err.Error())
	}

	return nil
}

func runCommand(cmd *DevCmd, command string) error {
	logger.Debug("Running Bundler run command")
	logger.Debug(fmt.Sprintf("$ %s\n", command))

	err := cmd.CommandRunInteractive(cmd.F, command)
	if err != nil {
		logger.Debug("Error while running command with simultaneous output", zap.Error(err))
		return msg.ErrFailedToRunDevCommand
	}
	return nil
}
