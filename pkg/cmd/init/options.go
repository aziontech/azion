package init

import (
	"errors"
	"os"
)

type TestFunc func(path string) error

var testFuncByType = map[string]TestFunc{
	"javascript": testJs,
	"nextjs":     nil,
	"flareact":   nil,
}

func testJs(path string) error {
	if _, err := os.Stat(path + "/package.json"); errors.Is(err, os.ErrNotExist) {
		return ErrorPackageJsonNotFound
	}
	return nil
}
