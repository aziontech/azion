package upbin

import (
	"fmt"
	"testing"
)

func Test_getCurrentVersion(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "I want look tha version",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(GetCurrentVersion())
		})
	}
}

func Test_getLatestVersion(t *testing.T) {
	GetLatestVersion()
}

func TestWhich(t *testing.T) {
	type args struct {
		command string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			args: args{
				command: "git",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Which(tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("Which() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Which() = %v, want %v", got, tt.want)
			}
		})
	}
}
