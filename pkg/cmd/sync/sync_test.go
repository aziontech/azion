package sync

import (
	"errors"
	"fmt"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/sync"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestSync(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name               string
		mockGetContentFunc func(confPath string) (*contracts.AzionApplicationOptions, error)
		mockWriteFunc      func(conf *contracts.AzionApplicationOptions, confPath string) error
		mockSyncResources  func(f *cmdutil.Factory, info contracts.SyncOpts, synch *SyncCmd) error
		expectedError      error
	}{
		{
			name: "sync - successful synchronization",
			mockGetContentFunc: func(confPath string) (*contracts.AzionApplicationOptions, error) {
				return &contracts.AzionApplicationOptions{}, nil
			},
			mockWriteFunc: func(conf *contracts.AzionApplicationOptions, confPath string) error {
				return nil
			},
			mockSyncResources: func(f *cmdutil.Factory, info contracts.SyncOpts, synch *SyncCmd) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name: "sync - failed to get content",
			mockGetContentFunc: func(confPath string) (*contracts.AzionApplicationOptions, error) {
				return nil, errors.New("Failed to synchronize local resources with remote resources: failed to get azion.json content")
			},
			mockWriteFunc: func(conf *contracts.AzionApplicationOptions, confPath string) error {
				return nil
			},
			mockSyncResources: func(f *cmdutil.Factory, info contracts.SyncOpts, synch *SyncCmd) error {
				return nil
			},
			expectedError: fmt.Errorf(msg.ERRORSYNC, "failed to get azion.json content"),
		},
		{
			name: "sync - failed to write content",
			mockGetContentFunc: func(confPath string) (*contracts.AzionApplicationOptions, error) {
				return &contracts.AzionApplicationOptions{
					// Mock relevant fields
				}, nil
			},
			mockWriteFunc: func(conf *contracts.AzionApplicationOptions, confPath string) error {
				return errors.New("failed to write azion.json content")
			},
			mockSyncResources: func(f *cmdutil.Factory, info contracts.SyncOpts, synch *SyncCmd) error {
				return errors.New("failed to write azion.json content")
			},
			expectedError: errors.New("failed to write azion.json content"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			f, _, _ := testutils.NewFactory(mock)

			syncCmd := NewSyncCmd(f)

			syncCmd.GetAzionJsonContent = tt.mockGetContentFunc
			syncCmd.WriteAzionJsonContent = tt.mockWriteFunc

			syncCmd.SyncResources = tt.mockSyncResources

			err := Run(syncCmd)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSyncFull(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("sync full - no items", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1000000/rules_engine/request/rules"),
			httpmock.JSONFromFile("./fixtures/rules.json"),
		)

		mock.Register(
			httpmock.REST("GET", "edge_applications/1000000/cache_settings"),
			httpmock.JSONFromFile("./fixtures/cache.json"),
		)

		mock.Register(
			httpmock.REST("GET", "edge_applications/1000000/origins"),
			httpmock.JSONFromFile("./fixtures/origins.json"),
		)

		mock.Register(
			httpmock.REST("GET", "variables"),
			httpmock.JSONFromFile("./fixtures/variables.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		syncCmd := NewSyncCmd(f)
		syncCmd.GetAzionJsonContent = func(confPath string) (*contracts.AzionApplicationOptions, error) {
			return &contracts.AzionApplicationOptions{
				Application: contracts.AzionJsonDataApplication{
					ID:   1000000,
					Name: "testename",
				},
			}, nil
		}
		syncCmd.WriteManifest = func(manifest *contracts.Manifest, pathMan string) error {
			return nil
		}
		syncCmd.CommandRunInteractive = func(f *cmdutil.Factory, comm string) error {
			return nil
		}
		syncCmd.ReadEnv = func(filenames ...string) (envMap map[string]string, err error) {
			return nil, nil
		}

		cmd := NewCobraCmd(syncCmd, f)

		cmd.SetArgs([]string{})

		err := cmd.Execute()

		require.NoError(t, err)
	})

	t.Run("sync full - items", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1000000/rules_engine/request/rules"),
			httpmock.JSONFromFile("./fixtures/rules_results.json"),
		)

		mock.Register(
			httpmock.REST("GET", "edge_applications/1000000/cache_settings"),
			httpmock.JSONFromFile("./fixtures/cache_results.json"),
		)

		mock.Register(
			httpmock.REST("GET", "edge_applications/1000000/origins"),
			httpmock.JSONFromFile("./fixtures/origins_results.json"),
		)

		mock.Register(
			httpmock.REST("GET", "variables"),
			httpmock.JSONFromFile("./fixtures/variables.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		syncCmd := NewSyncCmd(f)
		syncCmd.GetAzionJsonContent = func(confPath string) (*contracts.AzionApplicationOptions, error) {
			return &contracts.AzionApplicationOptions{
				Application: contracts.AzionJsonDataApplication{
					ID:   1000000,
					Name: "testename",
				},
			}, nil
		}
		syncCmd.WriteManifest = func(manifest *contracts.Manifest, pathMan string) error {
			return nil
		}
		syncCmd.CommandRunInteractive = func(f *cmdutil.Factory, comm string) error {
			return nil
		}
		syncCmd.ReadEnv = func(filenames ...string) (envMap map[string]string, err error) {
			return nil, nil
		}

		syncCmd.WriteAzionJsonContent = func(conf *contracts.AzionApplicationOptions, confPath string) error {
			return nil
		}

		cmd := NewCobraCmd(syncCmd, f)

		cmd.SetArgs([]string{})

		err := cmd.Execute()

		require.NoError(t, err)
	})

	t.Run("sync full - failed to write", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1000000/rules_engine/request/rules"),
			httpmock.JSONFromFile("./fixtures/rules_results.json"),
		)

		mock.Register(
			httpmock.REST("GET", "edge_applications/1000000/cache_settings"),
			httpmock.JSONFromFile("./fixtures/cache_results.json"),
		)

		mock.Register(
			httpmock.REST("GET", "edge_applications/1000000/origins"),
			httpmock.JSONFromFile("./fixtures/origins_results.json"),
		)

		mock.Register(
			httpmock.REST("GET", "variables"),
			httpmock.JSONFromFile("./fixtures/variables.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		syncCmd := NewSyncCmd(f)
		syncCmd.GetAzionJsonContent = func(confPath string) (*contracts.AzionApplicationOptions, error) {
			return &contracts.AzionApplicationOptions{
				Application: contracts.AzionJsonDataApplication{
					ID:   1000000,
					Name: "testename",
				},
			}, nil
		}
		syncCmd.WriteManifest = func(manifest *contracts.Manifest, pathMan string) error {
			return nil
		}
		syncCmd.CommandRunInteractive = func(f *cmdutil.Factory, comm string) error {
			return nil
		}
		syncCmd.ReadEnv = func(filenames ...string) (envMap map[string]string, err error) {
			return nil, nil
		}

		syncCmd.WriteAzionJsonContent = func(conf *contracts.AzionApplicationOptions, confPath string) error {
			return utils.ErrorWritingAzionJsonFile
		}

		cmd := NewCobraCmd(syncCmd, f)

		cmd.SetArgs([]string{})

		err := cmd.Execute()

		require.Error(t, err)
	})
}
