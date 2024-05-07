package node

import (
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkNode(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkNode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if err.Error() != tt.err {
					t.Errorf("checkNode() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
