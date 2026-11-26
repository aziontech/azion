package deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/deploy-remote"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestValidateStorageDirs(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	root := t.TempDir()

	withFile := filepath.Join(root, "with-file")
	require.NoError(t, os.MkdirAll(withFile, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(withFile, "index.html"), []byte("ok"), 0o644))

	nested := filepath.Join(root, "nested")
	require.NoError(t, os.MkdirAll(filepath.Join(nested, "a", "b"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(nested, "a", "b", "app.js"), []byte("ok"), 0o644))

	empty := filepath.Join(root, "empty")
	require.NoError(t, os.MkdirAll(empty, 0o755))

	onlyDirs := filepath.Join(root, "only-dirs")
	require.NoError(t, os.MkdirAll(filepath.Join(onlyDirs, "a", "b"), 0o755))

	notDir := filepath.Join(root, "file.txt")
	require.NoError(t, os.WriteFile(notDir, []byte("ok"), 0o644))

	missing := filepath.Join(root, "missing")

	tests := []struct {
		name     string
		storages []contracts.StorageManifest
		wantDir  string
	}{
		{name: "no storage", storages: nil},
		{name: "storage without dir", storages: []contracts.StorageManifest{{Name: "b"}}},
		{name: "dir with a file", storages: []contracts.StorageManifest{{Name: "b", Dir: withFile}}},
		{name: "file in a subdirectory", storages: []contracts.StorageManifest{{Name: "b", Dir: nested}}},
		{name: "missing dir", storages: []contracts.StorageManifest{{Name: "b", Dir: missing}}, wantDir: missing},
		{name: "empty dir", storages: []contracts.StorageManifest{{Name: "b", Dir: empty}}, wantDir: empty},
		{name: "only empty subdirectories", storages: []contracts.StorageManifest{{Name: "b", Dir: onlyDirs}}, wantDir: onlyDirs},
		{name: "dir is a file", storages: []contracts.StorageManifest{{Name: "b", Dir: notDir}}, wantDir: notDir},
		{
			name:     "second entry empty",
			storages: []contracts.StorageManifest{{Name: "a", Dir: withFile}, {Name: "b", Dir: empty}},
			wantDir:  empty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStorageDirs(tt.storages)
			if tt.wantDir == "" {
				assert.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Equal(t, fmt.Sprintf(msg.ErrorStorageDirEmpty, tt.wantDir), err.Error())
			assert.Contains(t, err.Error(), "Storage folder '"+tt.wantDir+"' is missing or empty.")
		})
	}
}
