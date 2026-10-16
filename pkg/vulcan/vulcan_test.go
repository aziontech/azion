//go:build !stage
// +build !stage

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
			want:  fmt.Sprintf("npx --yes  @aziontech/bundler%s presets ls --config-version 4", versionVulcan),
		},
		{
			name: "with flags - debug off",
			args: args{
				flags:  "--loglevel=error --no-update-notifier",
				params: "presets ls",
			},
			debug: false,
			want:  fmt.Sprintf("npx --yes --loglevel=error --no-update-notifier @aziontech/bundler%s presets ls --config-version 4", versionVulcan),
		},
		{
			name: "no params - debug off",
			args: args{
				flags: "--loglevel=error --no-update-notifier",
			},
			debug: false,
			want:  fmt.Sprintf("npx --yes --loglevel=error --no-update-notifier @aziontech/bundler%s ", versionVulcan),
		},
		{
			name: "no flags - debug on",
			args: args{
				params: "presets ls",
			},
			debug: true,
			want:  fmt.Sprintf("DEBUG=true npx --yes  @aziontech/bundler%s presets ls --config-version 4", versionVulcan),
		},
		{
			name: "with flags - debug on",
			args: args{
				flags:  "--loglevel=error --no-update-notifier",
				params: "presets ls",
			},
			debug: true,
			want:  fmt.Sprintf("DEBUG=true npx --yes --loglevel=error --no-update-notifier @aziontech/bundler%s presets ls --config-version 4", versionVulcan),
		},
		{
			name: "no params - debug on",
			args: args{
				flags: "--loglevel=error --no-update-notifier",
			},
			debug: true,
			want:  fmt.Sprintf("DEBUG=true npx --yes --loglevel=error --no-update-notifier @aziontech/bundler%s ", versionVulcan),
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

func TestCommandConfigVersion(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name   string
		params string
		want   string
	}{
		{name: "build", params: "build --preset react", want: "build --preset react --config-version 4"},
		{name: "build with a format verb", params: "build%s", want: "build%s --config-version 4"},
		{name: "dev", params: "dev", want: "dev --config-version 4"},
		{name: "manifest", params: "manifest transform --output %s --entry %s", want: "manifest transform --output %s --entry %s --config-version 4"},
		{name: "presets", params: "presets ls", want: "presets ls --config-version 4"},
		{name: "config does not accept config version", params: "config replace --replacements '{}'", want: "config replace --replacements '{}'"},
		{name: "store does not accept config version", params: "store init", want: "store init"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _, _ := testutils.NewFactory(&httpmock.Registry{})
			want := fmt.Sprintf("npx --yes  @aziontech/bundler%s %s", versionVulcan, tt.want)
			if got := NewVulcan().Command("", tt.params, f); got != want {
				t.Errorf("Command() = %v, want %v", got, want)
			}
		})
	}
}

// v3 still runs on edge-functions, which has no --config-version.
func TestCommandV3(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	// NewVulcanV3 switches the package to the edge-functions version.
	saved := []any{versionVulcan, currentMajor, firstTimeExecuting}
	defer func() {
		versionVulcan, currentMajor, firstTimeExecuting = saved[0].(string), saved[1].(int), saved[2].(string)
	}()

	f, _, _ := testutils.NewFactory(&httpmock.Registry{})
	got := NewVulcanV3().Command("", "build --preset react", f)
	want := "npx --yes  edge-functions@5.3.2 build --preset react"
	if got != want {
		t.Errorf("Command() = %v, want %v", got, want)
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
				currentVersion: "2.0.0",
			},
			lastVulcanVer:   "",
			expectedVersion: firstTimeExecuting,
			wantErr:         false,
		},
		{
			name: "new major version with last version",
			args: args{
				currentVersion: "2.0.0",
			},
			lastVulcanVer:   "1.2.0",
			expectedVersion: "@1.2.0",
			wantErr:         false,
		},
		{
			name: "same major version",
			args: args{
				currentVersion: "1.2.0",
			},
			lastVulcanVer:   "1.2.0",
			expectedVersion: "@1.2.0",
			wantErr:         false,
		},
		{
			name: "failed to parse bundler version",
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
			lastVulcanVer:   "0.5.0",
			expectedVersion: firstTimeExecuting,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _, _ := testutils.NewFactory(nil)
			vul := NewVulcan()
			vul.ReadSettings = func(string) (token.Settings, error) {
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
