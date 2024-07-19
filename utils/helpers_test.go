package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/stretchr/testify/assert"
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

func TestGetWorkingDir(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "write json content",
			test: func(t *testing.T) {
				path, _ := GetWorkingDir()
				jsonConf := filepath.Join(path, "azion", "azion.json")
				err := os.MkdirAll(filepath.Dir(jsonConf), os.ModePerm)
				require.NoError(t, err)
				var azJsonData contracts.AzionApplicationOptions
				azJsonData.Name = "Test01"
				azJsonData.Function.Name = "MyFunc"
				azJsonData.Function.File = "myfile.js"
				azJsonData.Function.ID = 476
				err = WriteAzionJsonContent(&azJsonData, "azion")
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestResponseToBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		err      error
	}{
		{"yes", true, nil},
		{"Yes", true, nil},
		{"YES", true, nil},
		{" no", false, nil},
		{"no", false, nil},
		{"NO", false, nil},
		{"  ", false, nil},
		{"", false, nil},
		{"  maybe  ", false, ErrorInvalidOption},
		{"anything else", false, ErrorInvalidOption},
	}

	for _, test := range tests {
		result, err := ResponseToBool(test.input)
		assert.Equal(t, test.expected, result)
		assert.Equal(t, test.err, err)
	}
}

func TestGetAzionJsonContent(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "read json content",
			test: func(t *testing.T) {
				path, _ := GetWorkingDir()
				jsonConf := filepath.Join(path, "azion", "azion.json")
				_ = os.MkdirAll(filepath.Dir(jsonConf), os.ModePerm)
				azJsonData, err := GetAzionJsonContent("azion")
				require.NoError(t, err)
				require.Contains(t, azJsonData.Name, "Test01")
				require.Contains(t, azJsonData.Function.Name, "MyFunc")
				require.Contains(t, azJsonData.Function.File, "myfile.js")
				require.EqualValues(t, azJsonData.Function.ID, 476)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, test.test)
	}
}

func TestErrorPerStatusCode(t *testing.T) {
	tests := []struct {
		name          string
		httpResp      *http.Response
		err           error
		expectedError error
	}{
		{
			name: "status code 401",
			httpResp: &http.Response{
				StatusCode: 401,
			},
			err:           nil,
			expectedError: ErrorToken401,
		},
		{
			name: "status code 403",
			httpResp: &http.Response{
				StatusCode: 403,
			},
			err:           nil,
			expectedError: ErrorForbidden403,
		},
		{
			name: "status code 404",
			httpResp: &http.Response{
				StatusCode: 404,
			},
			err:           nil,
			expectedError: ErrorNotFound404,
		},
		{
			name: "status code 409",
			httpResp: &http.Response{
				StatusCode: 409,
			},
			err:           nil,
			expectedError: ErrorNameInUse,
		},
		{
			name: "status code 200 with error",
			httpResp: &http.Response{
				StatusCode: 200,
			},
			err:           errors.New("some error"),
			expectedError: errors.New("some error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ErrorPerStatusCode(test.httpResp, test.err)
			assert.Equal(t, test.expectedError, result)
		})
	}
}

func TestCheckStatusCode500Error(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		expectedError error
	}{
		{
			name:          "error contains 'Client.Timeout'",
			err:           errors.New("Client.Timeout: request timed out"),
			expectedError: ErrorTimeoutAPICall,
		},
		{
			name:          "error does not contain 'Client.Timeout'",
			err:           errors.New("some other error"),
			expectedError: ErrorInternalServerError,
		},
		{
			name:          "error contains 'Client.Timeout' but case is different",
			err:           errors.New("client.timeout: request timed out"),
			expectedError: ErrorInternalServerError,
		},
		{
			name:          "error is empty string",
			err:           errors.New(""),
			expectedError: ErrorInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := checkStatusCode500Error(test.err)
			assert.Equal(t, test.expectedError, result)
		})
	}
}

func TestCheckStatusCode400Error(t *testing.T) {
	cases := []struct {
		name         string
		responseBody string
		expectedErr  string
	}{
		{
			name:         "All checks pass",
			responseBody: `{"key": "value"}`,
			expectedErr:  "\"key\": \"value\"",
		},
		{
			name:         "checkNoProduct fails",
			responseBody: `{"no_product": true}`,
			expectedErr:  "\"no_product\": true",
		},
		{
			name:         "checkTlsVersion fails",
			responseBody: `{"tls_version": "1.0"}`,
			expectedErr:  "\"tls_version\": \"1.0\"",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock para o http.Response
			mockResp := &http.Response{
				Body: io.NopCloser(strings.NewReader(tc.responseBody)),
			}

			err := checkStatusCode400Error(mockResp)
			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckNoProduct(t *testing.T) {
	cases := []struct {
		name        string
		body        string
		expectedErr error
	}{
		{
			name:        "No user_has_no_product key",
			body:        `{"some_key": "some_value"}`,
			expectedErr: nil,
		},
		{
			name:        "user_has_no_product key exists",
			body:        `{"user_has_no_product": "product123"}`,
			expectedErr: fmt.Errorf("%w: %s", ErrorProductNotOwned, "product123"),
		},
		{
			name:        "user_has_no_product key with empty value",
			body:        `{"user_has_no_product": ""}`,
			expectedErr: fmt.Errorf("%w: %s", ErrorProductNotOwned, ""),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := checkNoProduct(tc.body)
			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckOriginlessCacheSettings(t *testing.T) {
	cases := []struct {
		name        string
		body        string
		expectedErr error
	}{
		{
			name:        "No originless_cache_settings key",
			body:        `{"some_key": "some_value"}`,
			expectedErr: nil,
		},
		{
			name:        "originless_cache_settings key with value",
			body:        `{"originless_cache_settings": "cache_setting_123"}`,
			expectedErr: fmt.Errorf("cache_setting_123"),
		},
		{
			name:        "originless_cache_settings key with empty value",
			body:        `{"originless_cache_settings": ""}`,
			expectedErr: fmt.Errorf(""),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := checkOriginlessCacheSettings(tc.body)
			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckTlsVersion(t *testing.T) {
	cases := []struct {
		name        string
		body        string
		expectedErr error
	}{
		{
			name:        "No minimum_tls_version key",
			body:        `{"some_key": "some_value"}`,
			expectedErr: nil,
		},
		{
			name:        "minimum_tls_version key present",
			body:        `{"minimum_tls_version": "1.2"}`,
			expectedErr: ErrorMinTlsVersion,
		},
		{
			name:        "minimum_tls_version key with empty value",
			body:        `{"minimum_tls_version": ""}`,
			expectedErr: ErrorMinTlsVersion,
		},
		{
			name:        "minimum_tls_version key with unexpected JSON structure",
			body:        `{"minimum_tls_version": {"nested_key": "nested_value"}}`,
			expectedErr: ErrorMinTlsVersion,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := checkTlsVersion(tc.body)
			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckNameInUse(t *testing.T) {
	cases := []struct {
		name        string
		body        string
		expectedErr error
	}{
		{
			name:        "No name_in_use related keys",
			body:        `{"some_key": "some_value"}`,
			expectedErr: nil,
		},
		{
			name:        "name_already_in_use key present",
			body:        `{"name_already_in_use": "true"}`,
			expectedErr: ErrorNameInUse,
		},
		{
			name:        "bucket name is already in use message",
			body:        `{"error": "bucket name is already in use"}`,
			expectedErr: ErrorNameInUse,
		},
		{
			name:        "name taken error message",
			body:        `{"error": "name taken"}`,
			expectedErr: ErrorNameInUse,
		},
		{
			name:        "Mixed error messages",
			body:        `{"error": "bucket name is already in use", "name_already_in_use": "true"}`,
			expectedErr: ErrorNameInUse,
		},
		{
			name:        "name_already_in_use key with empty value",
			body:        `{"name_already_in_use": ""}`,
			expectedErr: ErrorNameInUse,
		},
		{
			name:        "bucket name is already in use key with empty value",
			body:        `{"bucket name is already in use": ""}`,
			expectedErr: ErrorNameInUse,
		},
		{
			name:        "name taken key with empty value",
			body:        `{"name taken": ""}`,
			expectedErr: ErrorNameInUse,
		},
		{
			name:        "Unknown key that is not related",
			body:        `{"other_key": "other_value"}`,
			expectedErr: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := checkNameInUse(tc.body)
			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckDetail(t *testing.T) {
	cases := []struct {
		name        string
		body        string
		expectedErr error
	}{
		{
			name:        "No detail key",
			body:        `{"some_key": "some_value"}`,
			expectedErr: nil,
		},
		{
			name:        "detail key with string value",
			body:        `{"detail": "detailed error message"}`,
			expectedErr: fmt.Errorf("detailed error message"),
		},
		{
			name:        "detail key with empty value",
			body:        `{"detail": ""}`,
			expectedErr: fmt.Errorf(""),
		},
		{
			name:        "detail key with number value",
			body:        `{"detail": 12345}`,
			expectedErr: fmt.Errorf("12345"),
		},
		{
			name:        "detail key with boolean value",
			body:        `{"detail": true}`,
			expectedErr: fmt.Errorf("true"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := checkDetail(tc.body)
			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckOrderField(t *testing.T) {
	cases := []struct {
		name        string
		body        string
		expectedErr error
	}{
		{
			name:        "No invalid_order_field key",
			body:        `{"some_key": "some_value"}`,
			expectedErr: nil,
		},
		{
			name:        "invalid_order_field key with string value",
			body:        `{"invalid_order_field": "invalid field error"}`,
			expectedErr: fmt.Errorf("invalid field error"),
		},
		{
			name:        "invalid_order_field key with empty value",
			body:        `{"invalid_order_field": ""}`,
			expectedErr: fmt.Errorf(""),
		},
		{
			name:        "invalid_order_field key with number value",
			body:        `{"invalid_order_field": 6789}`,
			expectedErr: fmt.Errorf("6789"),
		},
		{
			name:        "invalid_order_field key with boolean value",
			body:        `{"invalid_order_field": false}`,
			expectedErr: fmt.Errorf("false"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := checkOrderField(tc.body)
			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "String with less than 30 characters",
			input:    "short string",
			expected: "short string",
		},
		{
			name:     "String with exactly 30 characters",
			input:    "123456789012345678901234567890",
			expected: "123456789012345678901234567890",
		},
		{
			name:     "String with more than 30 characters",
			input:    "12345678901234567890123456789012345",
			expected: "123456789012345678901234567890...",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "String with exactly 31 characters",
			input:    "1234567890123456789012345678901",
			expected: "123456789012345678901234567890...",
		},
		{
			name:     "String with special characters",
			input:    "123456789012345678901234567890!@#",
			expected: "123456789012345678901234567890...",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := TruncateString(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsEmpty(t *testing.T) {
	emptyString := ""
	nonEmptyString := "non-empty"
	emptyIntSlice := []int{}
	nonEmptyIntSlice := []int{1}
	emptyStringSlice := []string{}
	nonEmptyStringSlice := []string{"non-empty"}
	emptyIntMap := map[string]int{}
	nonEmptyIntMap := map[string]int{"key": 1}
	emptyStringMap := map[string]string{}
	nonEmptyStringMap := map[string]string{"key": "value"}
	nilStringPtr := (*string)(nil)
	nonEmptyStringPtr := &nonEmptyString
	nilIntPtr := (*int)(nil)
	intVal := 1
	nonEmptyIntPtr := &intVal
	nilBoolPtr := (*bool)(nil)
	nilFloat64Ptr := (*float64)(nil)
	nilIntSlicePtr := (*[]int)(nil)
	nonEmptyIntSlicePtr := &nonEmptyIntSlice
	nilStringSlicePtr := (*[]string)(nil)
	nonEmptyStringSlicePtr := &nonEmptyStringSlice
	nilIntMapPtr := (*map[string]int)(nil)
	nonEmptyIntMapPtr := &nonEmptyIntMap
	nilStringMapPtr := (*map[string]string)(nil)
	nonEmptyStringMapPtr := &nonEmptyStringMap

	cases := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{"nil value", nil, true},
		{"empty string", "", true},
		{"non-empty string", "non-empty", false},
		{"empty []int", emptyIntSlice, true},
		{"non-empty []int", nonEmptyIntSlice, false},
		{"empty []string", emptyStringSlice, true},
		{"non-empty []string", nonEmptyStringSlice, false},
		{"empty map[string]int", emptyIntMap, true},
		{"non-empty map[string]int", nonEmptyIntMap, false},
		{"empty map[string]string", emptyStringMap, true},
		{"non-empty map[string]string", nonEmptyStringMap, false},
		{"nil *string", nilStringPtr, true},
		{"empty *string", &emptyString, true},
		{"non-empty *string", nonEmptyStringPtr, false},
		{"nil *int", nilIntPtr, true},
		{"non-nil *int", nonEmptyIntPtr, false},
		{"nil *bool", nilBoolPtr, true},
		{"nil *float64", nilFloat64Ptr, true},
		{"nil *[]int", nilIntSlicePtr, true},
		{"empty *[]int", &emptyIntSlice, true},
		{"non-empty *[]int", nonEmptyIntSlicePtr, false},
		{"nil *[]string", nilStringSlicePtr, true},
		{"empty *[]string", &emptyStringSlice, true},
		{"non-empty *[]string", nonEmptyStringSlicePtr, false},
		{"nil *map[string]int", nilIntMapPtr, true},
		{"empty *map[string]int", &emptyIntMap, true},
		{"non-empty *map[string]int", nonEmptyIntMapPtr, false},
		{"nil *map[string]string", nilStringMapPtr, true},
		{"empty *map[string]string", &emptyStringMap, true},
		{"non-empty *map[string]string", nonEmptyStringMapPtr, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsEmpty(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// Mock para survey.AskOne
type mockAskOne func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error

func (m mockAskOne) AskOne(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	return m(p, response, opts...)
}

func TestGetPackageManager(t *testing.T) {
	tests := []struct {
		name         string
		mockResponse mockAskOne
		expected     string
		expectedErr  error
	}{
		{
			name: "User selects npm",
			mockResponse: func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
				*(response.(*string)) = "npm"
				return nil
			},
			expected:    "npm",
			expectedErr: nil,
		},
		{
			name: "User selects yarn",
			mockResponse: func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
				*(response.(*string)) = "yarn"
				return nil
			},
			expected:    "yarn",
			expectedErr: nil,
		},
		{
			name: "User exits without selecting",
			mockResponse: func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
				return errors.New("interrupted")
			},
			expected:    "",
			expectedErr: errors.New("interrupted"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AskOne = tt.mockResponse.AskOne
			defer func() { AskOne = survey.AskOne }() // Restore original function after test

			result, err := GetPackageManager()
			assert.Equal(t, tt.expected, result)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type mockAsk func(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error

func (m mockAsk) Ask(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error {
	return m(qs, response, opts...)
}

// Testa a função AskInputEmpty
func TestAskInputEmpty(t *testing.T) {

	tests := []struct {
		name         string
		mockResponse mockAsk
		expected     string
		expectedErr  error
		exitExpected bool
	}{
		{
			name: "User enters a non-empty string",
			mockResponse: func(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error {
				*(response.(*string)) = "non-empty input"
				return nil
			},
			expected:    "non-empty input",
			expectedErr: nil,
		},
		// {
		// 	name: "User enters an empty string",
		// 	mockResponse: func(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error {
		// 		*(response.(*string)) = ""
		// 		return nil
		// 	},
		// 	expected:    "",
		// 	expectedErr: nil,
		// },
		// {
		// 	name: "User interrupts the input",
		// 	mockResponse: func(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error {
		// 		return terminal.InterruptErr
		// 	},
		// 	expected:     "",
		// 	expectedErr:  ErrorCancelledContextInput,
		// 	exitExpected: true,
		// },
		// {
		// 	name: "Error while parsing answer",
		// 	mockResponse: func(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error {
		// 		return errors.New("parse error")
		// 	},
		// 	expected:    "",
		// 	expectedErr: ErrorParseResponse,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ask = tt.mockResponse.Ask
			defer func() { ask = survey.Ask }() // Restore original function after test

			if tt.exitExpected {
				// Handle os.Exit call
				assert.PanicsWithValue(t, 0, func() { AskInputEmpty("Test") })
			} else {
				result, err := AskInputEmpty("Test")
				assert.Equal(t, tt.expected, result)
				if tt.expectedErr != nil {
					assert.EqualError(t, err, tt.expectedErr.Error())
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}

// Testa a função AskInput
func TestAskInput(t *testing.T) {

	tests := []struct {
		name         string
		mockResponse mockAsk
		expected     string
		expectedErr  error
		exitExpected bool
	}{
		{
			name: "User enters a valid string",
			mockResponse: func(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error {
				*(response.(*string)) = "valid input"
				return nil
			},
			expected:    "valid input",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ask = tt.mockResponse.Ask
			defer func() { ask = survey.Ask }() // Restore original function after test

			if tt.exitExpected {
				// Handle os.Exit call
				assert.PanicsWithValue(t, 0, func() { AskInput("Test") })
			} else {
				result, err := AskInput("Test")
				assert.Equal(t, tt.expected, result)
				if tt.expectedErr != nil {
					assert.EqualError(t, err, tt.expectedErr.Error())
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestAskPassword(t *testing.T) {
	tests := []struct {
		name         string
		mockResponse mockAsk
		expected     string
		expectedErr  error
		exitExpected bool
	}{
		{
			name: "User enters a valid password",
			mockResponse: func(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error {
				*(response.(*string)) = "validpassword"
				return nil
			},
			expected:    "validpassword",
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ask = tt.mockResponse.Ask
			defer func() { ask = survey.Ask }() // Restore original function after test

			if tt.exitExpected {
				// Handle os.Exit call
				assert.PanicsWithValue(t, 0, func() { AskPassword("Test") })
			} else {
				result, err := AskPassword("Test")
				assert.Equal(t, tt.expected, result)
				if tt.expectedErr != nil {
					assert.EqualError(t, err, tt.expectedErr.Error())
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}
