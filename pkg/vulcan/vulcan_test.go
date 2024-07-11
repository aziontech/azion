package vulcan

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/testutils"
)

func TestCommand(t *testing.T) {
	type args struct {
		flags  string
		params string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no flags",
			args: args{
				params: "presets ls",
			},
			want: fmt.Sprintf("npx --yes  edge-functions%s presets ls", versionVulcan),
		},
		{
			name: "with flags",
			args: args{
				flags:  "--loglevel=error --no-update-notifier",
				params: "presets ls",
			},
			want: fmt.Sprintf("npx --yes --loglevel=error --no-update-notifier edge-functions%s presets ls", versionVulcan),
		},
		{
			name: "no params",
			args: args{
				flags: "--loglevel=error --no-update-notifier",
			},
			want: fmt.Sprintf("npx --yes --loglevel=error --no-update-notifier edge-functions%s ", versionVulcan),
		},
	}
	f, _, _ := testutils.NewFactory(nil)
	for _, tt := range tests {
		vul := NewVulcan()
		t.Run(tt.name, func(t *testing.T) {
			if got := vul.Command(tt.args.flags, tt.args.params, f); got != tt.want {
				t.Errorf("Command() = %v, want %v", got, tt.want)
			}
		})
	}
}
