package deploy

import (
	"github.com/aziontech/azion-cli/pkg/api/edge_applications"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"reflect"
	"testing"
)

func Test_readManifest(t *testing.T) {
	tests := []struct {
		name    string
		want    *Manifest
		path    string
		wantErr bool
	}{
		{
			name: "success simple manifest",
			path: "/fixtures/manifest1.json",
			want: &Manifest{
				Routes: Routes{
					Deliver: []Deliver{
						{
							Variable:   "/public",
							InputValue: "/.edge/storage/",
							Priority:   1,
						},
					},
					Compute: []Compute{
						{
							Variable:   "/",
							InputValue: "/.edge/worker.js",
							Priority:   2,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestFilePath = tt.path
			got, err := readManifest()
			if (err != nil) != tt.wantErr {
				t.Errorf("readManifest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readManifest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prepareRequestDeliverRulesEngine(t *testing.T) {
	type args struct {
		manifest Manifest
	}

	manifestFilePath = "/fixtures/manifest2.json"
	manf, _ := readManifest()

	tests := []struct {
		name string
		args args
		want edge_applications.RequestsRulesEngine
	}{
		{
			name: "success",
			args: args{
				manifest: *manf,
			},
			want: edge_applications.RequestsRulesEngine{
				Request: sdk.CreateRulesEngineRequest{
					Name: "",
				},
				Phase: "response",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prepareRequestDeliverRulesEngine(tt.args.manifest); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prepareRequestDeliverRulesEngine() = %v, want %v", got, tt.want)
			}
		})
	}
}
