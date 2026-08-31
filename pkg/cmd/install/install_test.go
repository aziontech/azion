package install

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	f, _, _ := testutils.NewFactory(nil)
	cmd := NewCmd(f)

	assert.Equal(t, "install", cmd.Use)
	assert.Equal(t, "Install bundled resources", cmd.Short)
}

func TestInstallCmd_NoFlags(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	f, _, _ := testutils.NewFactory(nil)
	cmd := NewCmd(f)

	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	// Should show help when no flags provided - no error expected
}

func TestInstallCmd_HomeDirError(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	f, _, _ := testutils.NewFactory(nil)
	cmd := &installCmd{
		f:           f,
		skills:      true,
		userHomeDir: func() string { return "" },
	}

	err := cmd.installSkills()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to resolve home directory")
}

func TestInstallCmd_SourceDirNotFound(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	f, _, _ := testutils.NewFactory(nil)
	cmd := &installCmd{
		f:             f,
		skills:        true,
		userHomeDir:   func() string { return "/tmp" },
		executableDir: func() (string, error) { return "", errors.New("not found") },
		stat:          func(name string) (os.FileInfo, error) { return nil, os.ErrNotExist },
		getWorkingDir: func() (string, error) { return "/tmp", nil },
	}

	err := cmd.installSkills()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to locate bundled skills directory")
}

func TestInstallCmd_NoSkillsFound(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	// Create temp source directory (empty - no skill subdirectories)
	sourceDir, err := os.MkdirTemp("", "skills-empty")
	require.NoError(t, err)
	defer os.RemoveAll(sourceDir)

	// Create the expected skills directory structure
	skillsDir := sourceDir + "/skills"
	require.NoError(t, os.MkdirAll(skillsDir, 0755))

	f, _, _ := testutils.NewFactory(nil)
	cmd := &installCmd{
		f:           f,
		skills:      true,
		userHomeDir: func() string { return "/tmp" },
		executableDir: func() (string, error) {
			return sourceDir, nil
		},
		readDir:  os.ReadDir,
		mkdirAll: os.MkdirAll,
		stat: func(name string) (os.FileInfo, error) {
			// Make the findSkillsSourceDir find our skills directory
			cleanPath := name
			if len(name) > 3 && name[len(name)-3:] == "/.." {
				cleanPath = name[:len(name)-3]
			}
			if cleanPath == skillsDir || name == skillsDir {
				return mockDirInfo{}, nil
			}
			return os.Stat(name)
		},
		getWorkingDir: func() (string, error) { return sourceDir, nil },
	}

	// Empty source dir should result in no skills found (not an error)
	err = cmd.installSkills()
	require.NoError(t, err)
}

func TestGetUserHomeDir(t *testing.T) {
	home := getUserHomeDir()
	assert.NotEmpty(t, home, "getUserHomeDir should return a non-empty string")
}

func TestGetExecutableDir(t *testing.T) {
	dir, err := getExecutableDir()
	require.NoError(t, err)
	assert.NotEmpty(t, dir, "getExecutableDir should return a non-empty string")
}

// mockDirInfo implements os.FileInfo for testing
type mockDirInfo struct{}

func (m mockDirInfo) Name() string       { return "mock" }
func (m mockDirInfo) Size() int64        { return 0 }
func (m mockDirInfo) Mode() os.FileMode  { return os.ModeDir | 0755 }
func (m mockDirInfo) ModTime() time.Time { return time.Time{} }
func (m mockDirInfo) IsDir() bool        { return true }
func (m mockDirInfo) Sys() interface{}   { return nil }
