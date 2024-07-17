package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/stretchr/testify/require"
)

func TestCobraCmd(t *testing.T) {
	t.Run("clean directory", func(t *testing.T) {
		_ = os.MkdirAll("/tmp/ThisIsAzionCliTestDir", os.ModePerm)
		err := CleanDirectory("/tmp/ThisIsAzionCliTestDir")

		require.NoError(t, err)
	})

	t.Run("response to bool yes", func(t *testing.T) {
		resp, err := ResponseToBool("yes")
		require.True(t, resp)
		require.NoError(t, err)
	})

	t.Run("response to bool no", func(t *testing.T) {
		resp, err := ResponseToBool("no")
		require.False(t, resp)
		require.NoError(t, err)
	})

	t.Run("is directory empty", func(t *testing.T) {
		_ = os.MkdirAll("/tmp/ThisIsAzionCliTestDir", os.ModePerm)

		isEmpty, err := IsDirEmpty("/tmp/ThisIsAzionCliTestDir")
		require.True(t, isEmpty)

		require.NoError(t, err)
	})

	t.Run("load env from file vars", func(t *testing.T) {
		_ = os.MkdirAll("/tmp/ThisIsAzionCliFileVarTest", os.ModePerm)

		data := []byte("VAR1=test1\nVAR2=test2")
		_ = os.WriteFile("/tmp/ThisIsAzionCliFileVarTest/vars.txt", data, 0644)

		envs, err := LoadEnvVarsFromFile("/tmp/ThisIsAzionCliFileVarTest/vars.txt")
		require.Contains(t, envs[0], "test1")
		require.Contains(t, envs[1], "test2")

		require.NoError(t, err)
	})

	t.Run("write json content", func(t *testing.T) {
		path, _ := GetWorkingDir()

		jsonConf := path + "/azion/azion.json"

		err := os.MkdirAll(filepath.Dir(jsonConf), os.ModePerm)

		var azJsonData contracts.AzionApplicationOptions
		azJsonData.Name = "Test01"
		azJsonData.Function.Name = "MyFunc"
		azJsonData.Function.File = "myfile.js"
		azJsonData.Function.ID = 476

		_ = WriteAzionJsonContent(&azJsonData, "azion")

		require.NoError(t, err)
	})

	t.Run("read json content", func(t *testing.T) {
		path, _ := GetWorkingDir()

		jsonConf := path + "/azion/azion.json"

		_ = os.MkdirAll(filepath.Dir(jsonConf), os.ModePerm)

		azJsonData, err := GetAzionJsonContent("azion")

		require.NoError(t, err)
		require.Contains(t, azJsonData.Name, "Test01")
		require.Contains(t, azJsonData.Function.Name, "MyFunc")
		require.Contains(t, azJsonData.Function.File, "myfile.js")
		require.EqualValues(t, azJsonData.Function.ID, 476)
	})

	t.Run("returns invalid order_by", func(t *testing.T) {
		body := `{"invalid_order_field":"'edge_domain' is not a valid option for 'order_by'","available_order_fields":["id","name","cnames","cname_access_only","digital_certificate_id","edge_application_id","is_active"]}`
		err := checkOrderField(body)

		require.Equal(t, `'edge_domain' is not a valid option for 'order_by'`, err.Error())
	})
}

func TestIsEmpty(t *testing.T) {
	type args struct {
		value interface{}
	}

	var str *string
	var num *int

	tests := []struct {
		value interface{}
		want  bool
	}{
		{value: "string", want: false},
		{value: "", want: true},
		{value: str, want: true},
		{value: 1, want: false},
		{value: 0, want: false},
		{value: num, want: true},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := IsEmpty(tt.value); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCleanDirectory(t *testing.T) {
	t.Run("Successful cleaning", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "testdir")
		if err != nil {
			t.Fatalf("Error creating temporary directory: %v", err)
		}
		defer os.RemoveAll(dir)

		tempFile1 := filepath.Join(dir, "file1.txt")
		tempFile2 := filepath.Join(dir, "file2.txt")
		if err := os.WriteFile(tempFile1, []byte("content1"), 0666); err != nil {
			t.Fatalf("Error creating temporary file: %v", err)
		}
		if err := os.WriteFile(tempFile2, []byte("content2"), 0666); err != nil {
			t.Fatalf("Error creating temporary file: %v", err)
		}

		if _, err := os.Stat(tempFile1); os.IsNotExist(err) {
			t.Fatalf("file %s not created", tempFile1)
		}
		if _, err := os.Stat(tempFile2); os.IsNotExist(err) {
			t.Fatalf("file %s not created", tempFile2)
		}

		if err := CleanDirectory(dir); err != nil {
			t.Errorf("CleanDirectory failed: %v", err)
		}

		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			t.Errorf("Dir %s not removed", dir)
		}
	})

	t.Run("Error cleaning directory", func(t *testing.T) {
		nonExistentDir := "."
		err := CleanDirectory(nonExistentDir)
		errExpected := "Failed to clean the directory's contents because the directory is read-only and/or isn't accessible. Change the attributes of the directory to read/write and/or give access to it - ."
		if err != nil && err.Error() != errExpected {
			t.Errorf("Error not expected %q", err)
		}
	})
}
