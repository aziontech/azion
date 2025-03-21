package deploy

import (
	"testing"

	"github.com/aziontech/azion-cli/utils"
)

func Test_replaceInvalidChars(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "case 01",
			args: "asdfadsf-asdfsdf_a",
			want: "asdfadsf-asdfsdfa",
		},
		{
			name: "case 02",
			args: "asdfad}}.sf-asdfsdf_a",
			want: "asdfadsf-asdfsdfa",
		},
		{
			name: "case 03",
			args: "asdfad..asdfsdf\\/",
			want: "asdfadasdfsdf",
		},
		{
			name: "case 04",
			args: "asdfad#!@asdfadsf-asdf/***((ˆ",
			want: "asdfadasdfadsf-asdf",
		},
		{
			name: "case 05",
			args: "azion-asdfad#!@asdfadsf-asdf/***((ˆ",
			want: "asdfadasdfadsf-asdf",
		},
		{
			name: "case 06",
			args: "b2-asdfad#!@asdfadsf-asdf/***((ˆ",
			want: "asdfadasdfadsf-asdf",
		},
		{
			name: "case 07",
			args: "azion-b2-asdfad#!@asdfadsf-asdf/***((ˆ",
			want: "asdfadasdfadsf-asdf",
		},
		{
			name: "case 08",
			args: "b2-azion-asdfad#!@asdfadsf-asdf/***((ˆ",
			want: "asdfadasdfadsf-asdf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.ReplaceInvalidCharsBucket(tt.args); got != tt.want {
				t.Errorf("replaceInvalidChars() = %v, want %v", got, tt.want)
			}
		})
	}
}
