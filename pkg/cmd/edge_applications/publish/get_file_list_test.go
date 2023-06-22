package publish

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aziontech/azion-cli/pkg/testutils"
)

func TestPublishGetFileList(t *testing.T) {
	t.Run("ignore specified files", func(t *testing.T) {
		files := []string{
			"index.html",
			"ignore.js",
			"some/",
			"some/file.html",
			"some/ignore.js",
			"ignore/",
			"ignore/file.html",
			"ignore/file.txt",
			".git/",
			".git/ignore.txt",
			".git/ignore_too.txt",
		}
		globs := []string{
			"ignore/*.html",
			".git",
			"*.js",
		}
		ignored_files := []string{
			"ignore.js",
			"some/ignore.js",
			"ignore/file.html",
			".git/ignore.txt",
			".git/ignore_too.txt",
		}

		f, _, _ := testutils.NewFactory(nil)

		cmd := NewPublishCmd(f)

		cmd.FilepathWalk = func(root string, fn filepath.WalkFunc) error {
			for _, file := range files {
				info := MockedFileInfo{isDir: strings.HasSuffix(file, "/")}
				file = strings.TrimSuffix(file, "/")

				_ = fn(file, info, nil)
			}
			return nil
		}

		filelist, _ := cmd.getFileList("", globs)

		for _, ignored_file := range ignored_files {
			for _, file := range filelist {
				if file == ignored_file {
					t.Errorf("file %s should be ignored", file)
				}
			}
		}
	})
}

type MockedFileInfo struct {
	isDir bool
}

func (m MockedFileInfo) IsDir() bool {
	return m.isDir
}

func (m MockedFileInfo) ModTime() time.Time {
	return time.Now()
}

func (m MockedFileInfo) Mode() os.FileMode {
	return 0
}

func (m MockedFileInfo) Name() string {
	return ""
}

func (m MockedFileInfo) Size() int64 {
	return 0
}

func (m MockedFileInfo) Sys() interface{} {
	var sys interface{}
	return sys
}
