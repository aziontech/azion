package output

import (
	"fmt"
	"os"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/fatih/color"
)

type ErrorOutput struct {
	GeneralOutput `json:"-" yaml:"-" toml:"-"`
	Err           error `json:"error"`
}

func (e *ErrorOutput) Format() (bool, error) {
	formatted := false
	if len(e.Flags.Format) > 0 || len(e.Flags.Out) > 0 {
		formatted = true
		err := format(e, e.GeneralOutput)
		if err != nil {
			return formatted, err
		}
	}
	return formatted, nil
}

func (e *ErrorOutput) Output() {
	if e.Err != nil {
		format := fmt.Sprintf
		if !e.Flags.NoColor {
			format = color.New(color.FgRed).SprintfFunc()
		}
		logger.FInfo(os.Stderr, format("Error: %s\n", e.Err.Error()))
		os.Exit(1)
	}
}
