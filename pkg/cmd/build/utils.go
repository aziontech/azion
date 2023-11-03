package build

import (
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/build"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

func runCommand(cmd *BuildCmd, command string) error {
	logger.FInfo(cmd.Io.Out, build.BuildStart)

	logger.FInfo(cmd.Io.Out, build.BuildRunningCmd)
	logger.FInfo(cmd.Io.Out, fmt.Sprintf("$ %s\n", command))

	err := cmd.CommandRunInteractive(cmd.f, command)
	if err != nil {
		logger.Debug("Error while running command with simultaneous output", zap.Error(err))
		return build.ErrFailedToRunBuildCommand
	}

	logger.FInfo(cmd.Io.Out, build.BuildSuccessful)
	return nil
}
