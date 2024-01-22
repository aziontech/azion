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

		_ = WriteAzionJsonContent(&azJsonData)

		require.NoError(t, err)
	})

	t.Run("read json content", func(t *testing.T) {
		path, _ := GetWorkingDir()

		jsonConf := path + "/azion/azion.json"

		_ = os.MkdirAll(filepath.Dir(jsonConf), os.ModePerm)

		azJsonData, err := GetAzionJsonContent()

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
		name string
		args args
		want bool
	}{
		{
			name: "string",
			args: args{
				value: "string",
			},
			want: false,
		},
		{
			name: "string empty",
			args: args{
				value: "",
			},
			want: true,
		},
		{
			name: "string empty pointer",
			args: args{
				value: str,
			},
			want: true,
		},
		{
			name: "int",
			args: args{
				value: 1,
			},
			want: false,
		},
		{
			name: "int number zero",
			args: args{
				value: 0,
			},
			want: false,
		},
		{
			name: "int pointer",
			args: args{
				value: num,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.args.value); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
