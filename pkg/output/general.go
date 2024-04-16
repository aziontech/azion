package output

import (
	"io"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/fatih/color"
)

type GeneralOutput struct {
	Msg         string    `json:"message" yaml:"message" toml:"message"`
	Out         io.Writer `json:"-" yaml:"-" toml:"-"`
	FlagOutPath string    `json:"-" yaml:"-" toml:"-"`
	FlagFormat  string    `json:"-" yaml:"-" toml:"-"`
}

func (g *GeneralOutput) Format() (bool, error) {
	formated := false
	if len(g.FlagFormat) > 0 || len(g.FlagOutPath) > 0 {
		formated = true
		err := format(g, *g)
		if err != nil {
			return formated, err
		}
	}
	return formated, nil
}

func (g *GeneralOutput) Output() {
	format := color.New(color.FgGreen).SprintfFunc()
	logger.FInfo(g.Out, format("%s", g.Msg))
}
