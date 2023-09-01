package dev

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/dev"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

func vulcan(cmd *DevCmd) error {
	const command string = "npx --yes edge-functions@1.5.0 dev"

	err := runCommand(f, cmd, command)
	if err != nil {
		return fmt.Errorf(msg.ErrorVulcanExecute.Error(), err.Error())
	}

	return nil
}

func runCommand(f *cmdutil.Factory, cmd *DevCmd, command string) error {
	logger.Debug("Running vulcan run command")
	logger.Debug(fmt.Sprintf("$ %s\n", command))

	err := cmd.CommandRunInteractive(f, command)
	if err != nil {
		logger.Debug("Error while running command with simultaneous output", zap.Error(err))
		return msg.ErrFailedToRunDevCommand
	}
	return nil
}
