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

// Print prints a table using tab separation
// `headers` and `fields` should have the same length, representing each column of the table
// `elems` should be of a slice type else only the header will be printed.
func (p *TabPrinter) Print(headers []string, fields []string, elems interface{}) {
	fmt.Fprintf(p.writer, "%s", buildLine(headers))

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
	switch value.(type) {
	case time.Time:
		return value.(time.Time).Format(time.RFC822Z)
	default:
		return fmt.Sprintf("%v", value)
	}
}
