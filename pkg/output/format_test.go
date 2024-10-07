package output

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestFormat(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	// Mock for WriteDetailsToFile (simulates error or success when writing to a file)
	WriteDetailsToFile = func(_ []byte, filename string) error {
		if filename == "error" {
			return fmt.Errorf("write error")
		}
		return nil
	}

	tests := []struct {
		name        string
		v           interface{}
		g           GeneralOutput
		expectError bool
		expectLog   string
	}{
		{
			name: "valid JSON format",
			v:    map[string]string{"key": "value"},
			g: GeneralOutput{
				Flags: cmdutil.Flags{Format: "json", Out: ""},
				Out:   &bytes.Buffer{},
			},
			expectError: false,
			expectLog:   "{\n \"key\": \"value\"\n}",
		},
		{
			name: "valid YAML format",
			v:    map[string]string{"key": "value"},
			g: GeneralOutput{
				Flags: cmdutil.Flags{Format: "yaml", Out: ""},
				Out:   &bytes.Buffer{},
			},
			expectError: false,
			expectLog:   "key: value\n",
		},
		{
			name: "valid TOML format",
			v:    map[string]string{"key": "value"},
			g: GeneralOutput{
				Flags: cmdutil.Flags{Format: "toml", Out: ""},
				Out:   &bytes.Buffer{},
			},
			expectError: false,
			expectLog:   "key = \"value\"\n",
		},
		{
			name: "invalid type for JSON format",
			v:    make(chan int), // Tipo não serializável
			g: GeneralOutput{
				Flags: cmdutil.Flags{Format: "json", Out: ""},
				Out:   &bytes.Buffer{},
			},
			expectError: true,
		},
		{
			name: "write to file success",
			v:    map[string]string{"key": "value"},
			g: GeneralOutput{
				Flags: cmdutil.Flags{Format: "json", Out: "output.txt"},
				Out:   &bytes.Buffer{},
			},
			expectError: false,
			expectLog:   fmt.Sprintf(WRITE_SUCCESS, "output.txt"),
		},
		{
			name: "write to file failure",
			v:    map[string]string{"key": "value"},
			g: GeneralOutput{
				Flags: cmdutil.Flags{Format: "json", Out: "error"},
				Out:   &bytes.Buffer{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := format(tt.v, tt.g)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, tt.g.Out.(*bytes.Buffer).String(), tt.expectLog)
			}
		})
	}
}
