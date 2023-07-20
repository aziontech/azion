package upbin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReplaceBinary(t *testing.T) {
	filePath := "testfile.deb"

	downloadFileFunc = func(filePath string) error {
		return nil
	}

	replaceFileFunc = func(filePath string) error {
		return nil
	}

	err := replaceBinary(filePath)
	assert.NoError(t, err)
}
