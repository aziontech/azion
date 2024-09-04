package build

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/build"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"go.uber.org/zap"
)

func (b *BuildCmd) runCommand(command string, msgs *[]string) error {
	var hasDeployMessage bool
	if len(*msgs) > 0 {
		hasDeployMessage = true
	}

	logger.FInfoFlags(b.Io.Out, msg.BuildStart, b.f.Format, b.f.Out)
	*msgs = append(*msgs, msg.BuildStart)

	logger.FInfoFlags(b.Io.Out, msg.BuildRunningCmd, b.f.Format, b.f.Out)
	*msgs = append(*msgs, msg.BuildRunningCmd)
	logger.FInfoFlags(b.Io.Out, fmt.Sprintf("$ %s\n", command), b.f.Format, b.f.Out)
	*msgs = append(*msgs, fmt.Sprintf("$ %s\n", command))

	err := b.CommandRunInteractive(b.f, command)
	if err != nil {
		logger.Debug("Error while running command with simultaneous output", zap.Error(err))
		return msg.ErrFailedToRunBuildCommand
	}

	logger.FInfoFlags(b.Io.Out, msg.BuildSuccessful, b.f.Format, b.f.Out)
	*msgs = append(*msgs, msg.BuildSuccessful)

	if hasDeployMessage {
		return nil
	}

	outSlice := output.SliceOutput{
		Messages: *msgs,
		GeneralOutput: output.GeneralOutput{
			Out:   b.f.IOStreams.Out,
			Flags: b.f.Flags,
		},
	}

	return output.Print(&outSlice)
}
