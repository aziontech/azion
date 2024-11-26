package command

import (
	"bytes"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestCommandRunInteractive(t *testing.T) {
	tests := []struct {
		name    string
		factory *cmdutil.Factory
		command string
		wantErr bool
	}{
		{
			name: "Successful command with Silent false and no Format or Out",
			factory: &cmdutil.Factory{
				IOStreams: &iostreams.IOStreams{
					In:  os.Stdin,
					Out: os.Stdout,
					Err: os.Stderr,
				},
				Flags: cmdutil.Flags{
					Logger: logger.Logger{
						Silent: false,
					},
				},
			},
			command: "echo 'Hello, World!'",
			wantErr: false,
		},
		{
			name: "Successful command with Silent true",
			factory: &cmdutil.Factory{
				IOStreams: &iostreams.IOStreams{
					In:  os.Stdin,
					Out: os.Stdout,
					Err: os.Stderr,
				},
				Flags: cmdutil.Flags{
					Logger: logger.Logger{
						Silent: true,
					},
				},
			},
			command: "echo 'Hello, World!'",
			wantErr: false,
		},
		{
			name: "Successful command with Format flag",
			factory: &cmdutil.Factory{
				IOStreams: &iostreams.IOStreams{
					In:  os.Stdin,
					Out: os.Stdout,
					Err: os.Stderr,
				},
				Flags: cmdutil.Flags{
					Logger: logger.Logger{
						Silent: false,
					},
					Format: "json",
				},
			},
			command: "echo 'Hello, World!'",
			wantErr: false,
		},
		{
			name: "Successful command with Out flag",
			factory: &cmdutil.Factory{
				IOStreams: &iostreams.IOStreams{
					In:  os.Stdin,
					Out: os.Stdout,
					Err: os.Stderr,
				},
				Flags: cmdutil.Flags{
					Logger: logger.Logger{
						Silent: false,
					},
					Out: "output.txt",
				},
			},
			command: "echo 'Hello, World!'",
			wantErr: false,
		},
		{
			name: "Failed command",
			factory: &cmdutil.Factory{
				IOStreams: &iostreams.IOStreams{
					In:  os.Stdin,
					Out: os.Stdout,
					Err: os.Stderr,
				},
				Flags: cmdutil.Flags{
					Logger: logger.Logger{
						Silent: false,
					},
					Out: "output.txt",
				},
			},
			command: "invalidcommand",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CommandRunInteractive(tt.factory, tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommandRunInteractive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunCommandStreamOutput(t *testing.T) {
	tests := []struct {
		name       string
		envVars    []string
		comm       string
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "Successful command without env vars",
			envVars:    nil,
			comm:       "echo 'Hello, World!'",
			wantOutput: "Hello, World!\n",
			wantErr:    false,
		},
		{
			name:       "Successful command with env vars",
			envVars:    []string{"FOO=bar"},
			comm:       "echo $FOO",
			wantOutput: "bar\n",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			err := RunCommandStreamOutput(&out, tt.envVars, tt.comm)

			if (err != nil) != tt.wantErr {
				t.Errorf("RunCommandStreamOutput() error = %v, wantErr %v", err, tt.wantErr)
			}

			if gotOutput := out.String(); gotOutput != tt.wantOutput {
				t.Errorf("RunCommandStreamOutput() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestCommandRunInteractiveWithOutput(t *testing.T) {
	tests := []struct {
		name       string
		factory    *cmdutil.Factory
		command    string
		envVars    []string
		wantOutput string
		wantErr    bool
	}{
		{
			name: "Successful command without env vars and silent false",
			factory: &cmdutil.Factory{
				IOStreams: &iostreams.IOStreams{
					In:  os.Stdin,
					Out: os.Stdout,
					Err: os.Stderr,
				},
				Flags: cmdutil.Flags{
					Logger: logger.Logger{
						Silent: false,
					},
				},
			},
			command:    "echo 'Hello, World!'",
			envVars:    nil,
			wantOutput: "Hello, World!\n",
			wantErr:    false,
		},
		{
			name: "Successful command with env vars and silent false",
			factory: &cmdutil.Factory{
				IOStreams: &iostreams.IOStreams{
					In:  os.Stdin,
					Out: os.Stdout,
					Err: os.Stderr,
				},
				Flags: cmdutil.Flags{
					Logger: logger.Logger{
						Silent: false,
					},
				},
			},
			command:    "echo $FOO",
			envVars:    []string{"FOO=bar"},
			wantOutput: "bar\n",
			wantErr:    false,
		},
		{
			name: "Failed command without env vars and silent false",
			factory: &cmdutil.Factory{
				IOStreams: &iostreams.IOStreams{
					In:  os.Stdin,
					Out: os.Stdout,
					Err: os.Stderr,
				},
				Flags: cmdutil.Flags{
					Logger: logger.Logger{
						Silent: false,
					},
				},
			},
			command:    "invalidcommand",
			envVars:    nil,
			wantOutput: "",
			wantErr:    true,
		},
		{
			name: "Failed command with env vars and silent false",
			factory: &cmdutil.Factory{
				IOStreams: &iostreams.IOStreams{
					In:  os.Stdin,
					Out: os.Stdout,
					Err: os.Stderr,
				},
				Flags: cmdutil.Flags{
					Logger: logger.Logger{
						Silent: false,
					},
				},
			},
			command:    "invalidcommand",
			envVars:    []string{"FOO=bar"},
			wantOutput: "",
			wantErr:    true,
		},
		{
			name: "Successful command with env vars and silent true",
			factory: &cmdutil.Factory{
				IOStreams: &iostreams.IOStreams{
					In:  os.Stdin,
					Out: os.Stdout,
					Err: os.Stderr,
				},
				Flags: cmdutil.Flags{
					Logger: logger.Logger{
						Silent: true,
					},
				},
			},
			command:    "echo $FOO",
			envVars:    []string{"FOO=bar"},
			wantOutput: "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := CommandRunInteractiveWithOutput(tt.factory, tt.command, tt.envVars)

			if (err != nil) != tt.wantErr {
				t.Errorf("CommandRunInteractiveWithOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.wantOutput, output)
		})
	}
}

func TestRunCommandWithOutput(t *testing.T) {
	tests := []struct {
		name       string
		envVars    []string
		command    string
		wantOutput string
		wantCode   int
		wantErr    bool
	}{
		{
			name:       "Successful command without env vars",
			envVars:    nil,
			command:    "echo 'Hello, World!'",
			wantOutput: "Hello, World!\n",
			wantCode:   0,
			wantErr:    false,
		},
		{
			name:       "Successful command with env vars",
			envVars:    []string{"FOO=bar"},
			command:    "echo $FOO",
			wantOutput: "bar\n",
			wantCode:   0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, exitCode, err := RunCommandWithOutput(tt.envVars, tt.command)

			if (err != nil) != tt.wantErr {
				t.Errorf("RunCommandWithOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.wantOutput, output)
			require.Equal(t, tt.wantCode, exitCode)
		})
	}
}
