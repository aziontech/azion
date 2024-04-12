package output

import (
	"io"
	"strings"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/fatih/color"
	"github.com/maxwelbm/tablecli"
)

type ListOutput struct {
	Columns []string
	Lines   [][]string
	Page    int64
	Out     io.Writer
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

	// print the header only in the first flow
	if c.Page == 1 {
		logger.PrintHeader(tbl, format)
	}

	for _, row := range tbl.GetRows() {
		logger.PrintRow(tbl, format, row)
	}
}
