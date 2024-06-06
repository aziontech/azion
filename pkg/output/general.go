package output

import (
	"fmt"
	"io"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/fatih/color"
)

type GeneralOutput struct {
	Msg   string    `json:"message" yaml:"message" toml:"message"`
	Out   io.Writer `json:"-" yaml:"-" toml:"-"`
	Flags cmdutil.Flags
}

func (g *GeneralOutput) Format() (bool, error) {
	formated := false
	if len(g.Flags.Format) > 0 || len(g.Flags.Out) > 0 {
		formated = true
		err := format(g, *g)
		if err != nil {
			return formated, err
		}
	}
	return formated, nil
}

func (g *GeneralOutput) Output() {
	format := fmt.Sprintf
	if !g.Flags.NoColor {
		format = color.New(color.FgGreen).SprintfFunc()
	}
	logger.FInfo(g.Out, format("%s", g.Msg))
}
