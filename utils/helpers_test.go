package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestCleanDirectory(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "Successful cleaning",
			test: func(t *testing.T) {
				dir, err := os.MkdirTemp("", "testdir")
				if err != nil {
					t.Fatalf("Error creating temporary directory: %v", err)
				}
				defer os.RemoveAll(dir)

				tempFile1 := filepath.Join(dir, "file1.txt")
				tempFile2 := filepath.Join(dir, "file2.txt")
				if err := os.WriteFile(tempFile1, []byte("content1"), 0666); err != nil {
					t.Fatalf("Error creating temporary file: %v", err)
				}
				if err := os.WriteFile(tempFile2, []byte("content2"), 0666); err != nil {
					t.Fatalf("Error creating temporary file: %v", err)
				}

				if _, err := os.Stat(tempFile1); os.IsNotExist(err) {
					t.Fatalf("file %s not created", tempFile1)
				}
				if _, err := os.Stat(tempFile2); os.IsNotExist(err) {
					t.Fatalf("file %s not created", tempFile2)
				}

				if err := CleanDirectory(dir); err != nil {
					t.Errorf("CleanDirectory failed: %v", err)
				}

				if _, err := os.Stat(dir); !os.IsNotExist(err) {
					t.Errorf("Dir %s not removed", dir)
				}
			},
		},
		{
			name: "Error cleaning directory",
			test: func(t *testing.T) {
				nonExistentDir := "."
				err := CleanDirectory(nonExistentDir)
				errExpected := "Failed to clean the directory's contents because the directory is read-only and/or isn't accessible. Change the attributes of the directory to read/write and/or give access to it - ."
				if err != nil && err.Error() != errExpected {
					t.Errorf("Error not expected %q", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestIsDirEmpty(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() (string, error)
		cleanup   func(string)
		wantEmpty bool
		wantErr   bool
	}{
		{
			name: "directory does not exist",
			setup: func() (string, error) {
				return "/non/existent/directory", nil
			},
			cleanup:   func(string) {},
			wantEmpty: true,
			wantErr:   false,
		},
		{
			name: "directory is empty",
			setup: func() (string, error) {
				dir, err := os.MkdirTemp("", "emptydir")
				if err != nil {
					return "", err
				}
				return dir, nil
			},
			cleanup: func(dir string) {
				os.RemoveAll(dir)
			},
			wantEmpty: true,
			wantErr:   false,
		},
		{
			name: "directory is not empty",
			setup: func() (string, error) {
				dir, err := os.MkdirTemp("", "notemptydir")
				if err != nil {
					return "", err
				}
				f, err := os.CreateTemp(dir, "file")
				if err != nil {
					return "", err
				}
				f.Close()
				return dir, nil
			},
			cleanup: func(dir string) {
				os.RemoveAll(dir)
			},
			wantEmpty: false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := tt.setup()
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}
			defer tt.cleanup(dir)

			gotEmpty, err := IsDirEmpty(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsDirEmpty() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotEmpty != tt.wantEmpty {
				t.Errorf("IsDirEmpty() = %v, want %v", gotEmpty, tt.wantEmpty)
			}
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

func TestLoadEnvVarsFromFile(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		wantVars    []string
		wantErr     bool
	}{
		{
			name:        "Valid env file",
			fileContent: "VAR1=test1\nVAR2=test2",
			wantVars:    []string{"VAR1=test1", "VAR2=test2"},
			wantErr:     false,
		},
		{
			name:        "Empty env file",
			fileContent: "",
			wantVars:    []string{},
			wantErr:     false,
		},
		{
			name:        "Invalid env file path",
			fileContent: "",
			wantVars:    nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string
			var err error
			if tt.name != "Invalid env file path" {
				filePath = filepath.Join(os.TempDir(), "ThisIsAzionCliFileVarTest", "vars.txt")
				_ = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
				err = os.WriteFile(filePath, []byte(tt.fileContent), 0644)
				require.NoError(t, err)
			} else {
				filePath = "invalid/path/to/vars.txt"
			}

			envs, err := LoadEnvVarsFromFile(filePath)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.ElementsMatch(t, tt.wantVars, envs)
			}

			if tt.name != "Invalid env file path" {
				_ = os.RemoveAll(filepath.Dir(filePath))
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

// func TestCobraCmd(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		test func(t *testing.T)
// 	}{
// 		{
// 			name: "clean directory",
// 			test: func(t *testing.T) {
// 				_ = os.MkdirAll("/tmp/ThisIsAzionCliTestDir", os.ModePerm)
// 				err := CleanDirectory("/tmp/ThisIsAzionCliTestDir")
// 				require.NoError(t, err)
// 			},
// 		},
// 		{
// 			name: "response to bool yes",
// 			test: func(t *testing.T) {
// 				resp, err := ResponseToBool("yes")
// 				require.True(t, resp)
// 				require.NoError(t, err)
// 			},
// 		},
// 		{
// 			name: "response to bool no",
// 			test: func(t *testing.T) {
// 				resp, err := ResponseToBool("no")
// 				require.False(t, resp)
// 				require.NoError(t, err)
// 			},
// 		},
// 		{
// 			name: "is directory empty",
// 			test: func(t *testing.T) {
// 				_ = os.MkdirAll("/tmp/ThisIsAzionCliTestDir", os.ModePerm)
// 				isEmpty, err := IsDirEmpty("/tmp/ThisIsAzionCliTestDir")
// 				require.True(t, isEmpty)
// 				require.NoError(t, err)
// 			},
// 		},
// 		{
// 			name: "load env from file vars",
// 			test: func(t *testing.T) {
// 				_ = os.MkdirAll("/tmp/ThisIsAzionCliFileVarTest", os.ModePerm)
// 				data := []byte("VAR1=test1\nVAR2=test2")
// 				_ = os.WriteFile("/tmp/ThisIsAzionCliFileVarTest/vars.txt", data, 0644)
// 				envs, err := LoadEnvVarsFromFile("/tmp/ThisIsAzionCliFileVarTest/vars.txt")
// 				require.Contains(t, envs[0], "test1")
// 				require.Contains(t, envs[1], "test2")
// 				require.NoError(t, err)
// 			},
// 		},
// 		{
// 			name: "write json content",
// 			test: func(t *testing.T) {
// 				path, _ := GetWorkingDir()
// 				jsonConf := filepath.Join(path, "azion", "azion.json")
// 				err := os.MkdirAll(filepath.Dir(jsonConf), os.ModePerm)
// 				require.NoError(t, err)
// 				var azJsonData contracts.AzionApplicationOptions
// 				azJsonData.Name = "Test01"
// 				azJsonData.Function.Name = "MyFunc"
// 				azJsonData.Function.File = "myfile.js"
// 				azJsonData.Function.ID = 476
// 				err = WriteAzionJsonContent(&azJsonData, "azion")
// 				require.NoError(t, err)
// 			},
// 		},
// 		{
// 			name: "read json content",
// 			test: func(t *testing.T) {
// 				path, _ := GetWorkingDir()
// 				jsonConf := filepath.Join(path, "azion", "azion.json")
// 				_ = os.MkdirAll(filepath.Dir(jsonConf), os.ModePerm)
// 				azJsonData, err := GetAzionJsonContent("azion")
// 				require.NoError(t, err)
// 				require.Contains(t, azJsonData.Name, "Test01")
// 				require.Contains(t, azJsonData.Function.Name, "MyFunc")
// 				require.Contains(t, azJsonData.Function.File, "myfile.js")
// 				require.EqualValues(t, azJsonData.Function.ID, 476)
// 			},
// 		},
// 		{
// 			name: "returns invalid order_by",
// 			test: func(t *testing.T) {
// 				body := `{"invalid_order_field":"'edge_domain' is not a valid option for 'order_by'","available_order_fields":["id","name","cnames","cname_access_only","digital_certificate_id","edge_application_id","is_active"]}`
// 				err := checkOrderField(body)
// 				require.Equal(t, `'edge_domain' is not a valid option for 'order_by'`, err.Error())
// 			},
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, tt.test)
// 	}
// }

// func TestIsEmpty(t *testing.T) {
// 	type args struct {
// 		value interface{}
// 	}
//
// 	var str *string
// 	var num *int
//
// 	tests := []struct {
// 		value interface{}
// 		want  bool
// 	}{
// 		{value: "string", want: false},
// 		{value: "", want: true},
// 		{value: str, want: true},
// 		{value: 1, want: false},
// 		{value: 0, want: false},
// 		{value: num, want: true},
// 	}
// 	for _, tt := range tests {
// 		t.Run("", func(t *testing.T) {
// 			if got := IsEmpty(tt.value); got != tt.want {
// 				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
