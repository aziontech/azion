package build

import (
	"bytes"
	"os"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/edge_applications"

	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	t.Run("in build.cmd to run, type not informed", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		command := newBuildCmd(f)

		command.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"build": {}}`), nil
		}

		err := command.run()
		require.ErrorContains(t, err, "Error executing Vulcan")
	})

	t.Run("missing config file", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		command := newBuildCmd(f)

		command.FileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		err := command.run()
		require.ErrorIs(t, err, msg.ErrorOpeningAzionFile)
	})

	t.Run("invalid json", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		jsonContent := bytes.NewBufferString(`
	   {
	       "build": {
	           "cmd": "rm -rm *",
			   "output-ctrl": "on-error"
	       }
	   }
	   `)

		command := newBuildCmd(f)

		command.FileReader = func(path string) ([]byte, error) {
			return jsonContent.Bytes(), nil
		}

		err := command.run()
		require.ErrorContains(t, err, "Error executing Vulcan")
	})
}
