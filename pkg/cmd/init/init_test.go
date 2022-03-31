package init

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Run("Init without package.json", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript"})

		err := cmd.Execute()

		require.ErrorIs(t, err, ErrorPackageJsonNotFound)
	})

	t.Run("Init with unsupported type", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--name", "BLEBLEBLE", "--type", "demeuamor"})

		err := cmd.Execute()

		require.ErrorIs(t, err, utils.ErrorUnsupportedType)
	})

	t.Run("Init success with javascript", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		err := ioutil.WriteFile("package.json", []byte(""), 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove("package.json")

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript"})

		in := bytes.NewBuffer(nil)
		in.WriteString("yes\n")
		f.IOStreams.In = io.NopCloser(in)

		err = cmd.Execute()

		require.NoError(t, err)
		require.Contains(t, stdout.String(), `Template successfully fetched and configured
`)
	})

	t.Run("Init does not overwrite contents", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		err := ioutil.WriteFile("package.json", []byte(""), 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove("package.json")

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript"})

		in := bytes.NewBuffer(nil)
		in.WriteString("no\n")
		f.IOStreams.In = io.NopCloser(in)

		err = cmd.Execute()

		require.NoError(t, err)
		require.Contains(t, stdout.String(), `Init command stopped
`)
	})

	t.Run("Init invalid option", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		err := ioutil.WriteFile("package.json", []byte(""), 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove("package.json")
		defer os.RemoveAll("./azion")

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript"})

		in := bytes.NewBuffer(nil)
		in.WriteString("pix\n")
		f.IOStreams.In = io.NopCloser(in)

		err = cmd.Execute()

		require.ErrorIs(t, err, utils.ErrorInvalidOption)
	})

	t.Run("Init valid but noop", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		err := ioutil.WriteFile("package.json", []byte(""), 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove("package.json")

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "nextjs"})

		err = cmd.Execute()

		require.NoError(t, err)
	})

}
