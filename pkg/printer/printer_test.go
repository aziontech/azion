package printer

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Dummy struct {
	ID   int
	Name string
	Date time.Time
}

type DummyWithMethod struct {
	name  string
	price float64
}

func (d DummyWithMethod) GetName() string {
	return d.name
}

func (d DummyWithMethod) GetPrice() float64 {
	return d.price
}

func TestTab(t *testing.T) {
	t.Run("simple print", func(t *testing.T) {
		out := bytes.NewBuffer(nil)
		tp := NewTab(out)

		values := []Dummy{
			{1, "MyMock", time.Now()},
			{2, "MySuperMock", time.Now()},
		}

		tp.PrintWithHeaders(values, []string{"ID", "Name"}, []string{"ID", "Name"})

		assert.Equal(t, `ID    Name
1     MyMock
2     MySuperMock
`, out.String())
	})

	t.Run("fields that are methods", func(t *testing.T) {
		out := bytes.NewBuffer(nil)
		tp := NewTab(out)

		values := []DummyWithMethod{
			{name: "COSCARQUE H3", price: 59.99},
			{name: "CALPIS", price: 9.99},
		}

		tp.PrintWithHeaders(values, []string{"GetName()", "GetPrice()"}, []string{"Name", "Price"})

		assert.Equal(t, `Name            Price
COSCARQUE H3    59.99
CALPIS          9.99
`, out.String())
	})

	t.Run("not passing a slice", func(t *testing.T) {
		out := bytes.NewBuffer(nil)
		tp := NewTab(out)

		tp.PrintWithHeaders(1, []string{"ID", "Name"}, []string{"ID", "Name"})

		assert.Equal(t, "ID    Name\n", out.String())
	})

	t.Run("printing a datetime", func(t *testing.T) {
		out := bytes.NewBuffer(nil)
		tp := NewTab(out)

		values := []Dummy{
			{1, "MyMock", time.Date(2012, 12, 12, 9, 30, 15, 0, time.UTC)},
			{2, "MySuperMock", time.Date(2006, 01, 02, 23, 59, 32, 0, time.UTC)},
		}

		tp.PrintWithHeaders(values, []string{"ID", "Name", "Date"}, []string{"ID", "Name", "Date"})

		assert.Equal(t, `ID    Name           Date
1     MyMock         12 Dec 12 09:30 +0000
2     MySuperMock    02 Jan 06 23:59 +0000
`, out.String())
	})

	t.Run("header and fields with different lengths", func(t *testing.T) {
		out := bytes.NewBuffer(nil)
		tp := NewTab(out)

		values := []Dummy{
			{1, "MyMock", time.Date(2012, 12, 12, 9, 30, 15, 0, time.UTC)},
			{2, "MySuperMock", time.Date(2006, 01, 02, 23, 59, 32, 0, time.UTC)},
		}

		tp.PrintWithHeaders(values, []string{"ID", "Name"}, []string{"ID", "Name", "Date"})

		assert.Equal(t, `ID    Name    Date
1     MyMock
2     MySuperMock
`, out.String())
	})

}
