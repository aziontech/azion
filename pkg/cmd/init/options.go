package init

import (
	"errors"
	"os"

	"github.com/aziontech/azion-cli/utils"
)

type TestFunc func() error

var types = map[string]TestFunc{
	"javascript": testJs,
	"nextjs":     noop,
	"flareact":   noop,
}

func testJs() error {
	if _, err := os.Stat("./package.json"); errors.Is(err, os.ErrNotExist) {
		return utils.ErrorPackageJsonNotFound
	}
	return nil
}

func noop() error {
	return nil
}
