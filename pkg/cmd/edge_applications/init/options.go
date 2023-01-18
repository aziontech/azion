package init

import (
	"errors"
	"io/fs"
	"os"

	msg "github.com/aziontech/azion-cli/messages/edge_applications"
)

type TestFunc func(path string) error

type statFunc func(filename string) (fs.FileInfo, error)

func makeTestFuncMap(stat statFunc) map[string]TestFunc {
	return map[string]TestFunc{
		"javascript": testPackageJsonExists(stat),
		"flareact":   testPackageJsonExists(stat),
		"nextjs":     testPackageJsonExists(stat),
	}
}

func testPackageJsonExists(stat statFunc) TestFunc {
	return func(path string) error {
		if _, err := stat(path + "/package.json"); errors.Is(err, os.ErrNotExist) {
			return msg.ErrorPackageJsonNotFound
		}
		return nil
	}
}
