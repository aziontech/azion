package init

import (
	"errors"
	"os"
)

type TestFunc func() error

var testFuncByType = map[string]TestFunc{
	"javascript": testJs,
	"nextjs":     nil,
	"flareact":   nil,
}

func testJs() error {
	if _, err := os.Stat("./package.json"); errors.Is(err, os.ErrNotExist) {
		return ErrorPackageJsonNotFound
	}
	return nil
}
