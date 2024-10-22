package output

import (
	"bytes"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestDescribeOutput_Format(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	outBuffer := &bytes.Buffer{}

	type ValuesTest struct {
		ID   string
		Name string
	}

	type BadStruct struct {
		// Field with a non-exportable type (non-exported field and a channel)
		InvalidField chan int
	}

	tests := []struct {
		name          string
		describeOut   *DescribeOutput
		wantFormatted bool
		wantErr       bool
	}{

		{
			name: "should format with no errors",
			describeOut: &DescribeOutput{
				GeneralOutput: GeneralOutput{
					Out:   outBuffer,
					Flags: cmdutil.Flags{Format: "json", Out: ""},
				},
				Fields: map[string]string{"ID": "ID", "Name": "NAME"},
				Values: ValuesTest{ID: "0123", Name: "Luffy"},
			},
			wantFormatted: true,
			wantErr:       false,
		},
		{
			name: "should return an error when formatting fails",
			describeOut: &DescribeOutput{
				GeneralOutput: GeneralOutput{
					Out:   outBuffer,
					Flags: cmdutil.Flags{Format: "json", Out: ""},
				},
				Values: BadStruct{
					InvalidField: make(chan int),
				}, // causing format error
			},
			wantFormatted: true,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted, err := tt.describeOut.Format()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantFormatted, formatted)
		})
	}
}

func TestDescribeOutput_Output(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	outBuffer := &bytes.Buffer{}

	describeOut := &DescribeOutput{
		GeneralOutput: GeneralOutput{Out: outBuffer},
		Fields: map[string]string{
			"Field1": "Field 1",
			"Field2": "Field 2",
		},
		Values: &struct {
			Field1 string
			Field2 int
		}{
			Field1: "test",
			Field2: 123,
		},
	}

	describeOut.Output()

	expectedOutput := "Field 1:   test  \nField 2:   123   \n"
	assert.Contains(t, outBuffer.String(), expectedOutput)
}

func TestCheckPrimitiveType(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{"int type", 42, 42},
		{"string type", "hello", "hello"},
		{"float64 type", 3.14, 3.14},
		{"complex type", struct{ Name string }{"test"}, `{"Name":"test"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkPrimitiveType(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
