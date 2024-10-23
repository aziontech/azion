package root

import (
	"bytes"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type mockCommand struct {
	executeFunc func() error
}

func (m *mockCommand) Execute() error {
	return m.executeFunc()
}

func (m *mockCommand) InitDefaultHelpCmd() {}

func TestExecute(t *testing.T) {
	tests := []struct {
		name         string
		f            *factoryRoot
		exitCodeWant int
		exitCode     int
	}{
		{
			name: "success",
			f: &factoryRoot{
				factory: &cmdutil.Factory{
					HttpClient: &http.Client{
						Timeout: 50 * time.Second,
					},
					IOStreams: &iostreams.IOStreams{
						Out: &bytes.Buffer{},
						Err: &bytes.Buffer{},
					},
					Config: viper.New(),
				},
				command: &mockCommand{
					func() error { return nil },
				},
				doPreCommandCheck: func(cmd *cobra.Command, fact *factoryRoot) error {
					return nil
				},
				execSchedules: func(factory *cmdutil.Factory) {},
			},
			exitCode:     0,
			exitCodeWant: 0,
		},
		{
			name: "failure",
			f: &factoryRoot{
				factory: &cmdutil.Factory{
					HttpClient: &http.Client{
						Timeout: 50 * time.Second,
					},
					IOStreams: &iostreams.IOStreams{
						Out: &bytes.Buffer{},
						Err: &bytes.Buffer{},
					},
					Config: viper.New(),
				},
				command: &mockCommand{
					func() error {
						return errors.New("error")
					},
				},
				doPreCommandCheck: func(cmd *cobra.Command, fact *factoryRoot) error {
					return nil
				},
				execSchedules: func(factory *cmdutil.Factory) {},
			},
			exitCode:     1,
			exitCodeWant: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.osExit = func(code int) {
				tt.exitCode = code
			}

			Execute(tt.f)
			assert.Equal(t, tt.exitCodeWant, tt.exitCode)
		})
	}
}
