package upbin

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManagersPackages(t *testing.T) {
	packageManagerExistsFunc = func(manager string) bool {
		if manager == "brew" {
			return true
		} else {
			return false
		}
	}

	installPackageManagerFunc = func(manager string) error {
		if manager == "brew" {
			return nil
		}
		return errors.New("generic error")
	}

	downloadAndInstallPackageFunc = func(url string) error {
		return nil
	}

	install, err := managersPackages()
	assert.True(t, install)
	assert.Nil(t, err)
}

func TestGetFileExtension(t *testing.T) {
	fileName := "example.txt"
	expectedExtension := "txt"

	extension := getFileExtension(fileName)

	assert.Equal(t, expectedExtension, extension, "File extension should be 'txt'")
}
