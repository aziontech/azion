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
				return &contracts.AzionApplicationOptions{
					// Mock relevant fields
				}, nil
			},
			mockWriteFunc: func(conf *contracts.AzionApplicationOptions, confPath string) error {
				return nil
			},
			mockSyncResources: func(f *cmdutil.Factory, info contracts.SyncOpts, synch *SyncCmd) error {
				// Mock synchronization logic
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

			syncCmd := NewSync(f)

			// Mock GetAzionJsonContent and WriteAzionJsonContent functions
			syncCmd.GetAzionJsonContent = tt.mockGetContentFunc
			syncCmd.WriteAzionJsonContent = tt.mockWriteFunc

			// Replace syncResources function with mock
			syncCmd.SyncResources = tt.mockSyncResources

			err := Sync(syncCmd)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
