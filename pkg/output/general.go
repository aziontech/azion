package output

import (
	"io"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/fatih/color"
)

type GeneralOutput struct {
	FlagOutPath string
	FlagFormat  string
	Msg         string
	Out         io.Writer
}

func (g *GeneralOutput) Format() (bool, error) {
	return false, nil
}

func (g *GeneralOutput) Output() {
	format := color.New(color.FgGreen).SprintfFunc()
	logger.FInfo(g.Out, format("%s", g.Msg))
}
