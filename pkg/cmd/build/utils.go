package build

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/build"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"go.uber.org/zap"
)

func runCommand(cmd *BuildCmd, command string) error {
	msgs := []string{}
	logger.FInfoFlags(cmd.Io.Out, msg.BuildStart, cmd.f.Format, cmd.f.Out)
	msgs = append(msgs, msg.BuildStart)

	logger.FInfoFlags(cmd.Io.Out, msg.BuildRunningCmd, cmd.f.Format, cmd.f.Out)
	msgs = append(msgs, msg.BuildRunningCmd)
	logger.FInfoFlags(cmd.Io.Out, fmt.Sprintf("$ %s\n", command), cmd.f.Format, cmd.f.Out)
	msgs = append(msgs, fmt.Sprintf("$ %s\n", command))

	err := cmd.CommandRunInteractive(cmd.f, command)
	if err != nil {
		logger.Debug("Error while running command with simultaneous output", zap.Error(err))
		return msg.ErrFailedToRunBuildCommand
	}

	logger.FInfoFlags(cmd.Io.Out, msg.BuildSuccessful, cmd.f.Format, cmd.f.Out)
	msgs = append(msgs, msg.BuildSuccessful)

	outSlice := output.SliceOutput{
		Messages: msgs,
		GeneralOutput: output.GeneralOutput{
			Out:   cmd.f.IOStreams.Out,
			Flags: cmd.f.Flags,
		},
	}

	return output.Print(&outSlice)
}
