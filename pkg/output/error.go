package output

import (
	"os"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/fatih/color"
)

type ErrorOutput struct {
	GeneralOutput `json:"-" yaml:"-" toml:"-"`
	Err           error `json:"error"`
}

func (e *ErrorOutput) Format() (bool, error) {
	formated := false
	if len(e.FlagFormat) > 0 || len(e.FlagOutPath) > 0 {
		formated = true
		err := format(e, e.GeneralOutput)
		if err != nil {
			return formated, err
		}
	}
	return formated, nil
}

func (e *ErrorOutput) Output() {
	if e.Err != nil {
		format := color.New(color.FgRed).SprintfFunc()
		logger.FInfo(os.Stderr, format("Error: %s", e.Err.Error()))
		os.Exit(1)
	}
}
