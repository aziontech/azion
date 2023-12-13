package deploy

import (
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"testing"
)

func TestManifest_Interpreted(t *testing.T) {
	f, _, _ := testutils.NewFactory(nil)

	type args struct {
		f       *cmdutil.Factory
		cmd     func() *DeployCmd
		conf    *contracts.AzionApplicationOptions
		clients Clients
	}

	tests := []struct {
		name     string
		manifest *Manifest
		args     args
		wantErr  bool
	}{
		{
			name:     "case",
			manifest: &Manifest{},
			args: args{
				f: f,
				cmd: func() *DeployCmd {
					deployCmd := NewDeployCmd(f)
					return deployCmd
				},
				conf:    nil,
				clients: NewClients(f),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tt.manifest.Interpreted(tt.args.f, tt.args.cmd(), tt.args.conf, tt.args.clients); (err != nil) != tt.wantErr {
				t.Errorf("Interpreted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
