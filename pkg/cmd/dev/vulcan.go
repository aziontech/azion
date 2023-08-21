package dev

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/dev"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

func vulcan(cmd *DevCmd) error {
	const command string = "npx --yes edge-functions@1.1.0 run"

	err := runCommand(cmd, command)
	if err != nil {
		return fmt.Errorf(msg.ErrorVulcanExecute.Error(), err.Error())
	}

	return nil
}

func runCommand(cmd *DevCmd, command string) error {
	logger.Debug("Running vulcan run command")
	logger.Debug(fmt.Sprintf("$ %s\n", command))

	logger.FInfo(cmd.Io.Out, msg.RunningDevCommand)

	err := cmd.CommandRunnerStream(cmd.Io.Out, command, []string{})
	if err != nil {
		logger.Debug("Error while running command with simultaneous output", zap.Error(err))
		return msg.ErrFailedToRunDevCommand
	}
	return nil
}
