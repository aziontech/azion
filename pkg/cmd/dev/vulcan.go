package dev

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/dev"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	vul "github.com/aziontech/azion-cli/pkg/vulcan"
	"go.uber.org/zap"
)

func vulcan(f *cmdutil.Factory, cmd *DevCmd, isFirewall bool) error {

	command := vul.Command("", "dev")
	if isFirewall {
		command = vul.Command("", "dev --firewall")
	}

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
