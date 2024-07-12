package schedule

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestNewSchedule(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name      string
		schedule  Schedule
		expectErr bool
	}{
		{
			name: "valid schedule",
			schedule: Schedule{
				Name: "TestSchedule",
				Kind: "TestKind",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewSchedule(tt.schedule.Name, tt.schedule.Kind)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateFileSchedule(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name      string
		schedules []Schedule
		expectErr bool
	}{
		{
			name: "valid schedules",
			schedules: []Schedule{
				{Name: "Schedule1", Time: time.Now(), Kind: "Kind1"},
				{Name: "Schedule2", Time: time.Now(), Kind: "Kind2"},
			},
			expectErr: false,
		},
		{
			name:      "empty schedules",
			schedules: []Schedule{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := createFileSchedule(tt.schedules)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReadFileSchedule(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	// Prepare test data
	configPath, _ := config.Dir()
	path := filepath.Join(configPath.Dir, configPath.Schedule)
	defer os.Remove(path)

	schedules := []Schedule{
		{Name: "TestSchedule", Time: time.Now(), Kind: "TestKind"},
	}
	data, _ := json.Marshal(&schedules)
	os.WriteFile(path, data, 0666)

	tests := []struct {
		name      string
		setup     func()
		expectErr bool
		expectLen int
	}{
		{
			name: "valid schedules file",
			setup: func() {
				os.WriteFile(path, data, 0666)
			},
			expectErr: false,
			expectLen: 1,
		},
		{
			name: "no schedules file",
			setup: func() {
				os.Remove(path)
			},
			expectErr: false,
			expectLen: 0,
		},
		{
			name: "empty schedules file",
			setup: func() {
				os.WriteFile(path, []byte(""), 0666)
			},
			expectErr: false,
			expectLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			schedules, err := readFileSchedule()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, schedules, tt.expectLen)
			}
		})
	}
}

func TestCheckIf24HoursPassed(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name     string
		time     time.Time
		expected bool
	}{
		{
			name:     "more than 24 hours passed",
			time:     time.Now().Add(-25 * time.Hour),
			expected: true,
		},
		{
			name:     "less than 24 hours passed",
			time:     time.Now().Add(-23 * time.Hour),
			expected: false,
		},
		{
			name:     "exactly 24 hours passed",
			time:     time.Now().Add(-24 * time.Hour),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckIf24HoursPassed(tt.time)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExecSchedules(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	mockFactory := &cmdutil.Factory{}

	tests := []struct {
		name      string
		schedules []Schedule
		expectErr bool
		Trigger   func(f *cmdutil.Factory, name string) error
	}{
		{
			name: "schedule with delete bucket",
			schedules: []Schedule{
				{Name: "BucketToDelete", Time: time.Now().Add(-25 * time.Hour), Kind: DELETE_BUCKET},
			},
			expectErr: false,
			Trigger: func(f *cmdutil.Factory, name string) error {
				return nil
			},
		},
		{
			name: "schedule without delete bucket",
			schedules: []Schedule{
				{Name: "NotToDelete", Time: time.Now().Add(-25 * time.Hour), Kind: "OtherKind"},
			},
			expectErr: false,
			Trigger: func(f *cmdutil.Factory, name string) error {
				return nil
			},
		},
		{
			name: "schedule failed to trigger",
			schedules: []Schedule{
				{Name: "BucketToDelete", Time: time.Now().Add(-25 * time.Hour), Kind: DELETE_BUCKET},
			},
			expectErr: true,
			Trigger: func(f *cmdutil.Factory, name string) error {
				return errors.New("mocked error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup schedules file
			configPath, _ := config.Dir()
			path := filepath.Join(configPath.Dir, configPath.Schedule)
			defer os.Remove(path)

			data, _ := json.Marshal(&tt.schedules)
			os.WriteFile(path, data, 0666)
			TriggerDelete = tt.Trigger

			ExecSchedules(mockFactory)
		})
	}
}
