package output

import (
	"bytes"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestGeneralOutput_Format(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	outBuffer := &bytes.Buffer{}

	tests := []struct {
		name          string
		generalOutput *GeneralOutput
		wantFormatted bool
		wantErr       bool
	}{

		{
			name: "should format with no errors",
			generalOutput: &GeneralOutput{
				Out:   outBuffer,
				Flags: cmdutil.Flags{Format: "json", Out: ""},
			},
			wantFormatted: true,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted, err := tt.generalOutput.Format()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantFormatted, formatted)
		})
	}
}

func TestGeneralOutput_Output(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		noColor        bool
		expectedOutput string
	}{
		{
			name:           "With Color",
			noColor:        false,
			expectedOutput: color.New(color.FgGreen).Sprintf("%s", "Test message"),
		},
		{
			name:           "Without Color",
			noColor:        true,
			expectedOutput: "Test message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			g := &GeneralOutput{
				Flags: cmdutil.Flags{
					NoColor: tt.noColor,
				},
				Msg: "Test message",
				Out: &buf,
			}
			g.Output()
			assert.Equal(t, tt.expectedOutput, buf.String())
		})
	}
}
