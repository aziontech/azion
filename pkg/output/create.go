package output

import (
	"io"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/fatih/color"
)

type CreateOutput struct {
	Msg string
	Out io.Writer
}

func (c *CreateOutput) Output() {
	format := color.New(color.FgGreen).SprintfFunc()
	logger.FInfo(c.Out, format("%s", c.Msg))
}
