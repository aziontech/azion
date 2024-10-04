package output

import (
	"errors"
	"testing"
)

func TestPrint(t *testing.T) {
	type args struct {
		out TypeOutputInterface
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				out: &MockTypeOutputInterface{
					formatBool: true,
					formatErr:  nil,
				},
			},
			wantErr: false,
		},
		{
			name: "success with no format",
			args: args{
				out: &MockTypeOutputInterface{
					formatBool: false,
					formatErr:  nil,
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				out: &MockTypeOutputInterface{
					formatBool: false,
					formatErr:  errors.New("error format"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Print(tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("Print() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
