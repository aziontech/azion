package node

import (
	"errors"
	"os/exec"
	"testing"
)

func Test_checkNode(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		wantErr bool
		err     string
	}{
		{
			name:    "case 01",
			args:    "v21.7.3",
			wantErr: false,
		},
		{
			name:    "case 02",
			args:    "v15.3.8",
			wantErr: true,
			err:     NODE_OLDER_VERSION,
		},
		{
			name:    "case 03",
			args:    "v",
			wantErr: true,
			err:     NODE_OLDER_VERSION,
		},
		{
			name:    "case 04",
			args:    "",
			wantErr: true,
			err:     NODE_OLDER_VERSION,
		},
		{
			name:    "case 05",
			args:    "v18",
			wantErr: false,
		},
		{
			name:    "case 06",
			args:    "18.7.3",
			wantErr: false,
		},
		{
			name:    "case 07",
			args:    "v18.0.0beta",
			wantErr: false,
		},
		{
			name:    "case 08",
			args:    "vX.Y.Z",
			wantErr: true,
			err:     NODE_OLDER_VERSION,
		},
		{
			name:    "case 09",
			args:    "v16",
			wantErr: true,
			err:     NODE_OLDER_VERSION,
		},
		{
			name:    "case 10",
			args:    "v18.10",
			wantErr: false,
		},
		{
			name:    "case 11",
			args:    "v18.0.0-alpha",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodeManager := NewNode()
			err := nodeManager.CheckNode(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkNode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.err {
				t.Errorf("checkNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_nodeVersion(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		wantErr bool
		err     string
	}{
		{
			name:    "case 01",
			args:    "ls",
			wantErr: false,
		},
		{
			name:    "case 02",
			args:    "ls",
			wantErr: true,
			err:     NODE_NOT_INSTALLED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodeManager := NewNode()
			nodeManager.CmdBuilder = func(name string, arg ...string) *exec.Cmd {
				cmd := exec.Command(tt.args)
				return cmd
			}
			nodeManager.CheckNode = func(str string) error {
				if tt.wantErr {
					return errors.New(tt.err)
				} else {
					return nil
				}
			}
			err := nodeManager.NodeVer(nodeManager)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkNode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.err {
				t.Errorf("checkNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
