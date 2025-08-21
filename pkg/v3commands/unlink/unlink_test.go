package unlink

import (
	"context"
	"errors"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	app "github.com/aziontech/azion-cli/pkg/v3commands/delete/edge_application"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func mockDelSuccess() *app.DeleteCmd {
	return &app.DeleteCmd{
		Cascade: func(ctx context.Context, del *app.DeleteCmd) error {
			return nil
		},
	}
}

func mockDelFail() *app.DeleteCmd {
	return &app.DeleteCmd{
		Cascade: func(ctx context.Context, del *app.DeleteCmd) error {
			return errors.New("Failed to cascade delete")
		},
	}
}

func TestUnlink(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name           string
		isDirEmpty     bool
		cleanDirError  error
		expectedOutput string
		expectedError  string
		mockDeleteCmd  *app.DeleteCmd
	}{
		{
			name:           "unlink - directory is empty",
			isDirEmpty:     true,
			cleanDirError:  nil,
			expectedOutput: "Unliked successfully",
			expectedError:  "",
			mockDeleteCmd:  mockDelSuccess(),
		},
		{
			name:           "unlink - clean directory successfully",
			isDirEmpty:     false,
			cleanDirError:  nil,
			expectedOutput: "Unliked successfully",
			expectedError:  "",
			mockDeleteCmd:  mockDelSuccess(),
		},
		{
			name:           "unlink - failed to clean directory",
			isDirEmpty:     false,
			cleanDirError:  errors.New("failed to clean directory"),
			expectedOutput: "",
			expectedError:  "failed to clean directory",
			mockDeleteCmd:  mockDelSuccess(),
		},
		{
			name:           "unlink - cascade delete fails",
			isDirEmpty:     false,
			cleanDirError:  nil,
			expectedOutput: "",
			expectedError:  "Failed to cascade delete",
			mockDeleteCmd:  mockDelFail(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockIsDirEmpty := func(dirpath string) (bool, error) {
				return tt.isDirEmpty, nil
			}

			mockCleanDir := func(dirpath string) error {
				return tt.cleanDirError
			}

			mockDeleteCmd := tt.mockDeleteCmd

			mock := &httpmock.Registry{}

			f, out, _ := testutils.NewFactory(mock)
			f.GlobalFlagAll = true // Simulate --yes flag

			unlinkCmd := &UnlinkCmd{
				F:           f,
				IsDirEmpty:  mockIsDirEmpty,
				CleanDir:    mockCleanDir,
				ShouldClean: shouldClean,
				Clean:       clean,
				DeleteCmd: func(f *cmdutil.Factory) *app.DeleteCmd {
					return mockDeleteCmd
				},
			}
			cmd := NewCobraCmd(unlinkCmd, f)

			_, err := cmd.ExecuteC()
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, out.String())
			}

		})
	}
}
