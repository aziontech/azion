package output

import (
	"fmt"
	"reflect"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/fatih/color"
	"github.com/maxwelbm/tablecli"
)

type DescribeOutput struct {
	GeneralOutput `json:"-" yaml:"-" toml:"-"`
	Fields        map[string]string
	Values        interface{} `json:"fields" yaml:"fields" toml:"fields"`
}

func (d *DescribeOutput) Format() (bool, error) {
	formated := false
	if len(d.FlagFormat) > 0 || len(d.FlagOutPath) > 0 {
		formated = true
		err := format(d.Values, d.GeneralOutput)
		if err != nil {
			return formated, err
		}
	}
	return formated, nil
}

func (c *DescribeOutput) Output() {
	tbl := tablecli.New("", "")
	tbl.WithFirstColumnFormatter(color.New(color.FgBlue).SprintfFunc())

	interfaceFields := reflect.ValueOf(c.Values)
	for i := 0; i < interfaceFields.NumField(); i++ {
		field := interfaceFields.Type().Field(i)
		fieldValue := interfaceFields.Field(i).Interface()

		if vl, ok := c.Fields[field.Name]; ok {
			tbl.AddRow(fmt.Sprintf("%s: ", vl), fieldValue)
		}
	}

	logger.FInfo(c.Out, string(tbl.GetByteFormat()))
}
