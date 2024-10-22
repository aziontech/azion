package output

import (
	"strings"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/tablecli"
	"github.com/fatih/color"
)

type ListOutput struct {
	GeneralOutput `json:"-" yaml:"-" toml:"-"`
	Columns       []string   `json:"columns" yaml:"columns" toml:"columns"`
	Lines         [][]string `json:"lines" yaml:"lines" toml:"lines"`
}

func (l *ListOutput) Format() (bool, error) {
	formated := false
	if len(l.Flags.Format) > 0 || len(l.Flags.Out) > 0 {
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

	if !c.Flags.NoColor {
		headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgWhite).SprintfFunc()
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	}

	for _, ln := range c.Lines {
		tbl.AddRows(ln)
	}

	format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
	tbl.CalculateWidths(c.Columns)

	logger.PrintHeader(tbl, format)
	for _, row := range tbl.GetRows() {
		logger.PrintRow(tbl, format, row)
	}
}
