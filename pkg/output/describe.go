package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/fatih/color"
	"github.com/maxwelbm/tablecli"
)

type DescribeOutput struct {
	FlagOutPath string
	FlagFormat  string
	Msg         string
	Fields      [][]string
	Out         io.Writer
}

func (d *DescribeOutput) Format() (bool, error) {
	format := false
	if len(d.FlagFormat) > 0 || len(d.FlagOutPath) > 0 {
		format = true

		b, err := json.MarshalIndent(d, "", " ")
		if err != nil {
			return format, err
		}

		if len(d.FlagOutPath) > 0 {
			err = cmdutil.WriteDetailsToFile(b, d.FlagOutPath, d.Out)
			if err != nil {
				return format, fmt.Errorf("%s: %w", utils.ErrorWriteFile, err)
			}
			logger.FInfo(d.Out, fmt.Sprintf(WRITE_SUCCESS, d.FlagOutPath))
			return format, nil
		}

		logger.FInfo(d.Out, string(b))
		return format, nil
	}
	return format, nil
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
