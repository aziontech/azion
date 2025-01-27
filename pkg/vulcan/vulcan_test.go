package vulcan

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/pkg/token"
	"go.uber.org/zap/zapcore"
)

func TestCommand(t *testing.T) {
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
			name: "no flags - debug off",
			args: args{
				params: "presets ls",
			},
			debug: false,
			want:  fmt.Sprintf("npx --yes  edge-functions%s presets ls", versionVulcan),
		},
		{
			name: "with flags - debug off",
			args: args{
				flags:  "--loglevel=error --no-update-notifier",
				params: "presets ls",
			},
			debug: false,
			want:  fmt.Sprintf("npx --yes --loglevel=error --no-update-notifier edge-functions%s presets ls", versionVulcan),
		},
		{
			name: "no params - debug off",
			args: args{
				flags: "--loglevel=error --no-update-notifier",
			},
			debug: false,
			want:  fmt.Sprintf("npx --yes --loglevel=error --no-update-notifier edge-functions%s ", versionVulcan),
		},
		{
			name: "no flags - debug on",
			args: args{
				params: "presets ls",
			},
			debug: true,
			want:  fmt.Sprintf("DEBUG=true npx --yes  edge-functions%s presets ls", versionVulcan),
		},
		{
			name: "with flags - debug on",
			args: args{
				flags:  "--loglevel=error --no-update-notifier",
				params: "presets ls",
			},
			debug: true,
			want:  fmt.Sprintf("DEBUG=true npx --yes --loglevel=error --no-update-notifier edge-functions%s presets ls", versionVulcan),
		},
		{
			name: "no params - debug on",
			args: args{
				flags: "--loglevel=error --no-update-notifier",
			},
			debug: true,
			want:  fmt.Sprintf("DEBUG=true npx --yes --loglevel=error --no-update-notifier edge-functions%s ", versionVulcan),
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

func TestCheckVulcanMajor(t *testing.T) {
	type args struct {
		currentVersion string
	}
	tests := []struct {
		name            string
		args            args
		lastVulcanVer   string
		expectedVersion string
		wantErr         bool
		err             string
	}{
		{
			name: "new major version without last version",
			args: args{
				currentVersion: "5.0.0",
			},
			lastVulcanVer:   "",
			expectedVersion: firstTimeExecuting,
			wantErr:         false,
		},
		{
			name: "new major version with last version",
			args: args{
				currentVersion: "5.0.0",
			},
			lastVulcanVer:   "4.4.2",
			expectedVersion: "@4.4.2",
			wantErr:         false,
		},
		{
			name: "same major version",
			args: args{
				currentVersion: "5.0.0",
			},
			lastVulcanVer:   "4.4.2",
			expectedVersion: "@4.4.2",
			wantErr:         false,
		},
		{
			name: "failed to parse version",
			args: args{
				currentVersion: "invalid",
			},
			expectedVersion: firstTimeExecuting,
			wantErr:         true,
			err:             "strconv.Atoi: parsing \"invalid\": invalid syntax",
		},
		{
			name: "empty version string",
			args: args{
				currentVersion: "",
			},
			lastVulcanVer:   "2.5.0",
			expectedVersion: firstTimeExecuting,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _, _ := testutils.NewFactory(nil)
			vul := NewVulcan()
			vul.ReadSettings = func() (token.Settings, error) {
				return token.Settings{}, nil
			}
			err := vul.CheckVulcanMajor(tt.args.currentVersion, f, vul)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckVulcanMajor() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.err {
				t.Errorf("CheckVulcanMajor() error = %v, wantErr %v", err, tt.wantErr)
			}
			if versionVulcan != tt.expectedVersion {
				t.Errorf("versionVulcan = %v, expectedVersion %v", versionVulcan, tt.expectedVersion)
			}
		})
	}
}
