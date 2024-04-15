package output

import (
	"os"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/fatih/color"
)

type ErrorOutput struct {
	FlagOutPath string
	FlagFormat  string
	Err         error
}

func (e *ErrorOutput) Format() (bool, error) {
	return false, nil
}

func (e *ErrorOutput) Output() {
	if e.Err != nil {
		format := color.New(color.FgRed).SprintfFunc()
		logger.FInfo(os.Stderr, format("Error: %s", e.Err.Error()))
		os.Exit(1)
	}
}
