package template

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/rules_engine"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestDescribe(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("export to a path", func(t *testing.T) {

		f, stdout, _ := testutils.NewFactory(nil)

		cmd := NewCmd(f)
		path := "./out.json"
		cmd.SetArgs([]string{"--out", path})

		err := cmd.Execute()
		if err != nil {
			log.Println("error executing cmd err: ", err.Error())
		}

		_, err = os.ReadFile(path)
		if err != nil {
			t.Fatalf("error reading `out.json`: %v", err)
		}
		defer func() {
			_ = os.Remove(path)
		}()

		require.NoError(t, err)

		require.Equal(t, `File successfully written to: out.json
`, stdout.String())
	})

	t.Run("failed to write file", func(t *testing.T) {

		f, _, _ := testutils.NewFactory(nil)

		cmdTemplate := newTemplateCmd(f)
		cmdTemplate.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return errors.New("failed")
		}

		cmd := NewCobraCmd(cmdTemplate)
		path := "./out.json"
		cmd.SetArgs([]string{"--out", path})

		err := cmd.Execute()
		if err != nil {
			log.Println("error executing cmd err: ", err.Error())
		}

		require.ErrorIs(t, err, msg.ErrorWriteTemplate)
	})
}
