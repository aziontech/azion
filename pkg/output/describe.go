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
	Values        interface{}
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

	interfaceValue := reflect.ValueOf(c.Values).Elem()

	for i := 0; i < interfaceValue.NumField(); i++ {
		field := interfaceValue.Field(i)

		var dereferencedValue any
		if field.CanInterface() {
			fieldValue := field.Interface()

			if reflect.TypeOf(fieldValue).Kind() == reflect.Ptr {
				ptrValue := reflect.ValueOf(fieldValue)
				if !ptrValue.IsNil() {
					dereferencedValue = ptrValue.Elem().Interface()
				}
			}

			fieldName := interfaceValue.Type().Field(i).Name
			if vl, ok := c.Fields[fieldName]; ok {
				tbl.AddRow(fmt.Sprintf("%s: ", vl), dereferencedValue)
			}
		}
	}

	logger.FInfo(c.Out, string(tbl.GetByteFormat()))
}
