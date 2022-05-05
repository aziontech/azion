package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/stretchr/testify/require"
)

func TestCobraCmd(t *testing.T) {
	t.Run("Convert IDs to Int", func(t *testing.T) {
		ints, err := ConvertIdsToInt("10", "3")
		require.Equal(t, 10, int(ints[0]))
		require.Equal(t, 3, int(ints[1]))
		require.NoError(t, err)
	})

	t.Run("clean directory", func(t *testing.T) {
		os.MkdirAll("/tmp/ThisIsAzionCliTestDir", os.ModePerm)
		err := CleanDirectory("/tmp/ThisIsAzionCliTestDir")

		require.NoError(t, err)
	})

	t.Run("is directory empty", func(t *testing.T) {
		os.MkdirAll("/tmp/ThisIsAzionCliTestDir", os.ModePerm)

		isEmpty, err := IsDirEmpty("/tmp/ThisIsAzionCliTestDir")
		require.True(t, isEmpty)

		require.NoError(t, err)
	})

	t.Run("load env from file vars", func(t *testing.T) {
		os.MkdirAll("/tmp/ThisIsAzionCliFileVarTest", os.ModePerm)

		data := []byte("VAR1=test1\nVAR2=test2")
		os.WriteFile("/tmp/ThisIsAzionCliFileVarTest/vars.txt", data, 0644)

		envs, err := LoadEnvVarsFromFile("/tmp/ThisIsAzionCliFileVarTest/vars.txt")
		require.Contains(t, envs[0], "test1")
		require.Contains(t, envs[1], "test2")

		require.NoError(t, err)
	})

	t.Run("write json content", func(t *testing.T) {
		path, _ := GetWorkingDir()

		jsonConf := path + "/azion/azion.json"

		err := os.MkdirAll(filepath.Dir(jsonConf), os.ModePerm)

		var azJsonData contracts.AzionJsonData
		azJsonData.Name = "Test01"
		azJsonData.Function.Name = "MyFunc"
		azJsonData.Function.Active = true
		azJsonData.Function.File = "myfile.js"
		azJsonData.Function.Language = "javascript"
		azJsonData.Function.Id = 476

		WriteAzionJsonContent(&azJsonData)

		require.NoError(t, err)
	})

	t.Run("read json content", func(t *testing.T) {
		path, _ := GetWorkingDir()

		jsonConf := path + "/azion/azion.json"

		os.MkdirAll(filepath.Dir(jsonConf), os.ModePerm)

		azJsonData, err := GetAzionJsonContent()

		require.NoError(t, err)
		require.Contains(t, azJsonData.Name, "Test01")
		require.Contains(t, azJsonData.Function.Name, "MyFunc")
		require.Contains(t, azJsonData.Function.File, "myfile.js")
		require.Contains(t, azJsonData.Function.Language, "javascript")
		require.EqualValues(t, azJsonData.Function.Id, 476)
	})
}
