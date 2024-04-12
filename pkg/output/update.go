package output

import (
	"io"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/fatih/color"
)

type UpdateOutput struct {
	Msg string
	Out io.Writer
}

func (c *UpdateOutput) Output() {
	format := color.New(color.FgGreen).SprintfFunc()
	logger.FInfo(c.Out, format("%s", c.Msg))
}
