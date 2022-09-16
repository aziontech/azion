package init

import (
	"errors"
	"io/fs"
	"os"

	msg "github.com/aziontech/azion-cli/messages/webapp"
)

type TestFunc func(path string) error

type statFunc func(filename string) (fs.FileInfo, error)

func makeTestFuncMap(stat statFunc) map[string]TestFunc {
	return map[string]TestFunc{
		"javascript": func(path string) error {
			if _, err := stat(path + "/package.json"); errors.Is(err, os.ErrNotExist) {
				return msg.ErrorPackageJsonNotFound
			}
			return nil
		},
		"flareact": func(path string) error {
			if _, err := stat(path + "/package.json"); errors.Is(err, os.ErrNotExist) {
				return msg.ErrorPackageJsonNotFound
			}
			return nil
		},
		"nextjs": func(path string) error {
			if _, err := stat(path + "/package.json"); errors.Is(err, os.ErrNotExist) {
				return msg.ErrorPackageJsonNotFound
			}
			return nil
		},
	}
}
