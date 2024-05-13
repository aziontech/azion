package output

import (
	"strings"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/fatih/color"
	"github.com/aziontech/tablecli"
)

type ListOutput struct {
	GeneralOutput `json:"-" yaml:"-" toml:"-"`
	Columns       []string   `json:"columns" yaml:"columns" toml:"columns"`
	Lines         [][]string `json:"lines" yaml:"lines" toml:"lines"`
}

func (l *ListOutput) Format() (bool, error) {
	formated := false
	if len(l.FlagFormat) > 0 || len(l.FlagOutPath) > 0 {
		formated = true
		err := format(l, l.GeneralOutput)
		if err != nil {
			return formated, err
		}
	}
	return formated, nil
}

func (c *ListOutput) Output() {
	tbl := tablecli.NewTable(c.Columns)
	tbl.WithWriter(c.Out)

	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgWhite).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, ln := range c.Lines {
		tbl.AddRows(ln)
	}

	format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
	tbl.CalculateWidths([]string{})

	logger.PrintHeader(tbl, format)
	for _, row := range tbl.GetRows() {
		logger.PrintRow(tbl, format, row)
	}
}
