package output

import (
	"fmt"
	"io"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/fatih/color"
)

type GeneralOutput struct {
	Msg   string        `json:"message,omitempty" yaml:"message,omitempty" toml:"message,omitempty"`
	Out   io.Writer     `json:"-" yaml:"-" toml:"-"`
	Flags cmdutil.Flags `json:"-" yaml:"-" toml:"-"`
}

func (g *GeneralOutput) Format() (bool, error) {
	formatted := false
	if len(g.Flags.Format) > 0 || len(g.Flags.Out) > 0 {
		formatted = true
		err := format(g, *g)
		if err != nil {
			return formatted, err
		}
	}
	return formatted, nil
}

func (g *GeneralOutput) Output() {
	format := fmt.Sprintf
	if !g.Flags.NoColor {
		format = color.New(color.FgGreen).SprintfFunc()
	}
	logger.FInfo(g.Out, format("%s", g.Msg))
}
