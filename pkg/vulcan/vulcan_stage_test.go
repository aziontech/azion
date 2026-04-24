//go:build stage
// +build stage

package vulcan

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"go.uber.org/zap/zapcore"
)

func TestCommandStage(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	type args struct {
		flags  string
		params string
	}
	tests := []struct {
		name  string
		args  args
		debug bool
		want  string
	}{
		{
			name: "stage - no flags - debug off",
			args: args{
				params: "presets ls",
			},
			debug: false,
			want:  fmt.Sprintf("npx --yes  %s presets ls", StagePkgURL),
		},
		{
			name: "stage - with flags - debug off",
			args: args{
				flags:  "--loglevel=error --no-update-notifier",
				params: "presets ls",
			},
			debug: false,
			want:  fmt.Sprintf("npx --yes --loglevel=error --no-update-notifier %s presets ls", StagePkgURL),
		},
		{
			name: "stage - no params - debug off",
			args: args{
				flags: "--loglevel=error --no-update-notifier",
			},
			debug: false,
			want:  fmt.Sprintf("npx --yes --loglevel=error --no-update-notifier %s ", StagePkgURL),
		},
		{
			name: "stage - no flags - debug on",
			args: args{
				params: "presets ls",
			},
			debug: true,
			want:  fmt.Sprintf("DEBUG=true npx --yes  %s presets ls", StagePkgURL),
		},
		{
			name: "stage - with flags - debug on",
			args: args{
				flags:  "--loglevel=error --no-update-notifier",
				params: "presets ls",
			},
			debug: true,
			want:  fmt.Sprintf("DEBUG=true npx --yes --loglevel=error --no-update-notifier %s presets ls", StagePkgURL),
		},
		{
			name: "stage - no params - debug on",
			args: args{
				flags: "--loglevel=error --no-update-notifier",
			},
			debug: true,
			want:  fmt.Sprintf("DEBUG=true npx --yes --loglevel=error --no-update-notifier %s ", StagePkgURL),
		},
	}
	for _, tt := range tests {
		logger.New(zapcore.DebugLevel)
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)
		if tt.debug {
			f.Logger.Debug = true
		}
		vul := NewVulcan()
		t.Run(tt.name, func(t *testing.T) {
			if got := vul.Command(tt.args.flags, tt.args.params, f); got != tt.want {
				t.Errorf("Command() = %v, want %v", got, tt.want)
			}
		})
	}
}
