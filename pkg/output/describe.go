package output

import (
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/fatih/color"
	"github.com/maxwelbm/tablecli"
)

type DescribeOutput struct {
	GeneralOutput `json:"-" yaml:"-" toml:"-"`
	Fields        [][]string `json:"fields" yaml:"fields" toml:"fields"`
}

func (d *DescribeOutput) Format() (bool, error) {
	formated := false
	if len(d.FlagFormat) > 0 || len(d.FlagOutPath) > 0 {
		formated = true
		err := format(d, d.GeneralOutput)
		if err != nil {
			return formated, err
		}
	}
	return formated, nil
}

func (c *DescribeOutput) Output() {
	sliceNull := make([]string, 2)
	tbl := tablecli.NewTable(sliceNull)
	tbl.WithFirstColumnFormatter(color.New(color.FgBlue).SprintfFunc())

	for _, v := range c.Fields {
		tbl.AddRows(v)
	}

	logger.FInfo(c.Out, string(tbl.GetByteFormat()))
}
