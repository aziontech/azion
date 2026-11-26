package deploy

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	msg "github.com/aziontech/azion-cli/messages/deploy-remote"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

// errFoundFile stops the walk at the first file found.
var errFoundFile = errors.New("found a file")

// ValidateStorageDirs fails when a storage entry of the manifest names a dir
// that does not exist or holds no file to upload. Entries without a dir are
// not checked: they declare a bucket and nothing to upload.
func ValidateStorageDirs(storages []contracts.StorageManifest) error {
	for _, st := range storages {
		if st.Dir == "" {
			continue
		}
		if !hasFiles(st.Dir) {
			return fmt.Errorf(msg.ErrorStorageDirEmpty, st.Dir)
		}
	}
	return nil
}

// hasFiles reports whether dir is a directory with at least one file in it or
// in its subdirectories. Symlinks are not counted, since the upload skips them.
func hasFiles(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		logger.Debug("Storage dir is missing or is not a directory", zap.String("dir", dir), zap.Error(err))
		return false
	}
	err = filepath.WalkDir(dir, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() {
			return errFoundFile
		}
		return nil
	})
	if errors.Is(err, errFoundFile) {
		return true
	}
	if err != nil {
		logger.Debug("Error while reading the storage dir", zap.String("dir", dir), zap.Error(err))
	}
	return false
}
