package deploy

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func TestFunctionFile(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name      string
		functions []contracts.FunctionV3
		want      string
	}{
		{
			name:      "handler declared by most presets",
			functions: []contracts.FunctionV3{{Name: "handler", Target: "./functions/handler.js"}},
			want:      ".edge/functions/handler.js",
		},
		{
			name:      "worker declared by opennextjs",
			functions: []contracts.FunctionV3{{Name: "worker", Target: "./functions/worker.js"}},
			want:      ".edge/functions/worker.js",
		},
		{
			name: "only the first of several functions is used",
			functions: []contracts.FunctionV3{
				{Name: "first", Target: "functions/first.js"},
				{Name: "second", Target: "functions/second.js"},
			},
			want: ".edge/functions/first.js",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := functionFile(tt.functions); got != tt.want {
				t.Errorf("functionFile() = %q, want %q", got, tt.want)
			}
		})
	}
}
