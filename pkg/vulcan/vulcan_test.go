package vulcan

import "testing"

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
			want: "npx --yes  edge-functions@1.7.0 presets ls",
		},
		{
			name: "with flags",
			args: args{
				flags:  "--loglevel=error --no-update-notifier",
				params: "presets ls",
			},
			want: "npx --yes --loglevel=error --no-update-notifier edge-functions@1.7.0 presets ls",
		},
		{
			name: "no params",
			args: args{
				flags: "--loglevel=error --no-update-notifier",
			},
			want: "npx --yes --loglevel=error --no-update-notifier edge-functions@1.7.0 ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Command(tt.args.flags, tt.args.params); got != tt.want {
				t.Errorf("Command() = %v, want %v", got, tt.want)
			}
		})
	}
}
