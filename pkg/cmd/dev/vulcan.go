package dev

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/dev"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

func vulcan(cmd *DevCmd, isFirewall bool) error {

	vul := cmd.Vulcan()
	command := vul.Command("", "dev", cmd.F)
	if isFirewall {
		command = vul.Command("", "dev --firewall", cmd.F)
	}

	err := runCommand(cmd, command)
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
