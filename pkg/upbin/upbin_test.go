package upbin

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateBin(t *testing.T) {
	notifyFunc = func() (bool, error) {
		return true, nil
	}

	needToUpdateFunc = func() bool {
		return true
	}

	wantToUpdateFunc = func() bool {
		return true
	}

	prepareURLFunc = func() (string, error) {
		return "http://example.com/file", nil
	}

	managersPackagesFunc = func() (bool, error) {
		return false, nil
	}

	whichFunc = func(string) (string, error) {
		return "/path/to/azioncli", nil
	}

	downloadFileFunc = func(string, string) error {
		return nil
	}

	replaceFileFunc = func(string) error {
		return nil
	}

	err := UpdateBin()
	assert.NoError(t, err)
}

func TestGetLastActivity(t *testing.T) {
	openConfigFunc = func() (map[string]interface{}, error) {
		return map[string]interface{}{
			"LAST_ACTIVITY": "2021-10-01",
		}, nil
	}

	expectedLastActivity := "2021-10-01"

	lastActivity, err := getLastActivity()
	assert.NoError(t, err)
	assert.Equal(t, expectedLastActivity, lastActivity)
}

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
		return errors.New("Package manager installation failed")
	}

	install, err := managersPackages()
	assert.True(t, install)
	assert.Nil(t, err)
}
