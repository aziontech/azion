package build

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/build"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

func runCommand(cmd *BuildCmd, command string) error {
	logger.FInfo(cmd.Io.Out, msg.BuildStart)

	logger.FInfo(cmd.Io.Out, msg.BuildRunningCmd)
	logger.FInfo(cmd.Io.Out, fmt.Sprintf("$ %s\n", command))

	err := cmd.CommandRunInteractive(cmd.f, []string{}, command)
	if err != nil {
		logger.Debug("Error while running command with simultaneous output", zap.Error(err))
		return msg.ErrFailedToRunBuildCommand
	}

	logger.FInfo(cmd.Io.Out, msg.BuildSuccessful)
	return nil
}
