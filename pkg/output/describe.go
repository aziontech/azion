package output

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/tablecli"
	"github.com/fatih/color"
)

type DescribeOutput struct {
	GeneralOutput `json:"-" yaml:"-" toml:"-"`
	Fields        map[string]string
	Values        interface{}
	Field         string // Used for large character values like codes or scripts that break the table.
}

func (d *DescribeOutput) Format() (bool, error) {
	formatted := false
	if len(d.Flags.Format) > 0 || len(d.Flags.Out) > 0 {
		formatted = true
		err := format(d.Values, d.GeneralOutput)
		if err != nil {
			return formatted, err
		}
	}
	return formatted, nil
}

func (c *DescribeOutput) Output() {
	tbl := tablecli.New("", "")

	if !c.Flags.NoColor {
		tbl.WithFirstColumnFormatter(color.New(color.FgBlue).SprintfFunc())
	}

	values := reflect.ValueOf(c.Values)
	interfaceValue := values.Elem()

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
			} else {
				dereferencedValue = reflect.ValueOf(fieldValue).Interface()
			}

			fieldName := interfaceValue.Type().Field(i).Name
			if vl, ok := c.Fields[fieldName]; ok {
				tbl.AddRow(fmt.Sprintf("%s: ", vl), checkPrimitiveType(dereferencedValue))
			}
		}

	}

	logger.FInfo(c.Out, string(tbl.GetByteFormat()))
	if len(c.Field) > 0 {
		format := fmt.Sprintf
		if !c.Flags.NoColor {
			format = color.New(color.FgGreen).SprintfFunc()
		}
		logger.FInfo(c.Out, format("\nCode: %s", c.Field))
	}
}

func checkPrimitiveType(value any) any {
	valueType := reflect.TypeOf(value)
	if valueType == nil {
		jsonValue, _ := json.Marshal(value)
		return string(jsonValue)
	}

	switch valueType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Bool, reflect.String:
		return value
	default:
		jsonValue, _ := json.Marshal(value)
		return string(jsonValue)
	}
}
