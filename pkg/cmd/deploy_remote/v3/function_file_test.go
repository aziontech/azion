package deploy

import "testing"

func TestFunctionFile(t *testing.T) {
	tests := []struct {
		preset string
		want   string
	}{
		{preset: "opennextjs", want: ".edge/functions/worker.js"},
		{preset: "OpenNextJS", want: ".edge/functions/worker.js"},
		{preset: "next", want: ".edge/functions/handler.js"},
		{preset: "astro", want: ".edge/functions/handler.js"},
		{preset: "", want: ".edge/functions/handler.js"},
	}
	for _, tt := range tests {
		t.Run(tt.preset, func(t *testing.T) {
			if got := functionFile(tt.preset); got != tt.want {
				t.Errorf("functionFile(%q) = %q, want %q", tt.preset, got, tt.want)
			}
		})
	}
}
