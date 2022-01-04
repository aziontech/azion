package printer

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
	"time"
)

const padding = 4 // spaces

// TabPrinter represents a printer that uses tab separation
type TabPrinter struct {
	writer *tabwriter.Writer
}

// NewTab creates a printer for tab separated values
func NewTab(w io.Writer) *TabPrinter {
	return &TabPrinter{
		writer: tabwriter.NewWriter(w, 0, 0, padding, ' ', 0),
	}
}

// PrintWitHeaders prints using Print but with an additional header line
// Each element of the `headers` slice will be used to name a given tab-separated column
// `headers` and `fields` should be match one to one i.e. for each header there should be a field
func (p *TabPrinter) PrintWithHeaders(elems interface{}, fields []string, headers []string) {
	fmt.Fprintf(p.writer, "%s", buildLine(headers))
	p.Print(elems, fields)
}

// Print can be used to print a some fields of a struct slice, all tab separated.
// The `fields` slice should contain the names of the struct fields that should be printed.
// If a field does not exist in a given element of `elems`, it will panic
// If `elems` is not of a slice type, nothing will be printed.
func (p *TabPrinter) Print(elems interface{}, fields []string) {
	rows := buildRows(elems, fields)
	for _, row := range rows {
		fmt.Fprintf(p.writer, "%s", buildLine(row))
	}

	p.writer.Flush()
}

func buildLine(args []string) string {
	var b strings.Builder

	for idx, arg := range args {
		b.WriteString(arg)
		if idx+1 != len(args) {
			b.WriteString("\t")
		}
	}
	b.WriteString("\n")

	return b.String()
}

// buildRows returns a slice of string slice
// Each top-level slice represents a 'row' and the values inside a given
// row are the columns
// NOTE: Not tested for non-basic types (structs,slices,maps)
func buildRows(elems interface{}, fields []string) [][]string {
	if reflect.TypeOf(elems).Kind() != reflect.Slice {
		// ignore if not a slice
		return nil
	}

	var rows [][]string

	slice := reflect.ValueOf(elems)

	for i := 0; i < slice.Len(); i++ {
		var columns []string

		for _, name := range fields {
			value := slice.Index(i).FieldByName(name).Interface()
			columns = append(columns, toString(value))
		}

		rows = append(rows, columns)
	}

	return rows
}

func toString(value interface{}) string {
	switch v := value.(type) {
	case time.Time:
		return v.Format(time.RFC822Z)
	default:
		return fmt.Sprintf("%v", v)
	}
}
