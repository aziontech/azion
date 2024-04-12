package output

import (
	"io"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/fatih/color"
)

type DeleteOutput struct {
	Msg string
	Out io.Writer
}

func (d *DeleteOutput) Output() {
	format := color.New(color.FgGreen).SprintfFunc()
	logger.FInfo(d.Out, format("%s", d.Msg))
}
