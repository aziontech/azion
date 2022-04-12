package init

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/contracts"
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

	t.Run("Init with -y and -n flags", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--name", "BLEBLEBLE", "--type", "demeuamor", "-y", "-n"})

		err := cmd.Execute()

		require.ErrorIs(t, err, ErrorYesAndNoOptions)
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

	t.Run("Init success with javascript using flag -y", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		err := ioutil.WriteFile("package.json", []byte(""), 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove("package.json")

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript", "-y"})

		err = cmd.Execute()

		require.NoError(t, err)
		require.Contains(t, stdout.String(), `Template successfully fetched and configured
`)
	})

	t.Run("Init does not overwrite contents", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

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
	})

	t.Run("Init does not overwrite contents using flag -n", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		err := ioutil.WriteFile("package.json", []byte(""), 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove("package.json")

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript", "-n"})

		err = cmd.Execute()

		require.NoError(t, err)
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

	t.Run("runInitCmdLine without config.json", func(t *testing.T) {
		config := &contracts.AzionApplicationConfig{}
		confDir, _ := os.Getwd()
		confDir = confDir + "/azion/"

		var err error
		_ = os.Remove(confDir + "config.json")

		err = runInitCmdLine(config)
		require.EqualError(t, err, "Failed to open config.json file")
	})

	t.Run("runInitCmdLine without init.env", func(t *testing.T) {
		var err error
		config := &contracts.AzionApplicationConfig{}
		confDir, _ := os.Getwd()
		confDir = confDir + "/azion/"
		_ = os.Remove("/tmp/ls-test.txt")
		_ = os.MkdirAll(confDir, os.ModePerm)

		file, err := os.Create(confDir + "config.json")
		if err == nil {
			_, err = file.WriteString("{\n	\"init\": {\n	\"cmd\": \"ls -1 $VAR1 $VAR2 > /tmp/ls-test.txt\",\n		\"env\": \"./azion/init.env\"\n		}\n	}\n")
			if err != nil {
				require.NoError(t, err)
			}
		}
		file.Close()

		// User specified envfile but it cannot be read correctly
		err = runInitCmdLine(config)
		require.Error(t, err)
	})

	t.Run("runInitCmdLine full", func(t *testing.T) {
		var err error
		config := &contracts.AzionApplicationConfig{}
		confDir, _ := os.Getwd()
		confDir = confDir + "/azion/"
		_ = os.Remove("/tmp/ls-test.txt")
		_ = os.MkdirAll(confDir, os.ModePerm)
		defer os.RemoveAll(confDir)
		jsonConf := confDir + "config.json"
		file, err := os.Create(jsonConf)
		if err == nil {
			_, err = file.WriteString("{\n	\"init\": {\n	\"cmd\": \"ls -1 $VAR1 $VAR2 > /tmp/ls-test.txt\",\n		\"env\": \"./azion/init.env\"\n		}\n	}\n")
			if err != nil {
				require.NoError(t, err)
			}

		}
		file.Close()

		file, err = os.Create(confDir + "init.env")
		if err == nil {
			_, err = file.WriteString("VAR1=/\nVAR2=/bin\n")
			if err != nil {
				require.NoError(t, err)
			}
		}
		file.Close()

		err = runInitCmdLine(config)
		if err != nil {
			require.NoError(t, err)
		}

		if _, err := os.Stat("/tmp/ls-test.txt"); errors.Is(err, os.ErrNotExist) {
			require.NoError(t, err)
		}

		fileContent, err := ioutil.ReadFile("/tmp/ls-test.txt")
		if err != nil {
			require.NoError(t, err)
		}
		strFromFile := string(fileContent)

		require.NoError(t, err)
		//As stated in VAR2, /bin should have 'bash'
		require.Contains(t, strFromFile, "bash")
	})

}
