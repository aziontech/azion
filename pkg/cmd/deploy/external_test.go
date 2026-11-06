package deploy

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/apiversion"
	v6deploy "github.com/aziontech/azion-cli/pkg/cmd/deploy/v6"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
)

// init and link deploy through NewExternalDeployer, so a v6 account must get
// the v6 pipeline there too, not the v4 one that creates v4 resources.
func TestExternalDeployerFollowsTheGeneration(t *testing.T) {
	for _, tt := range []struct {
		v    apiversion.Version
		isV6 bool
	}{{apiversion.V4, false}, {apiversion.V6, true}, {"", false}} {
		d := NewExternalDeployer(&cmdutil.Factory{APIVersion: tt.v})
		_, gotV6 := d.(*v6deploy.DeployCmd)
		if gotV6 != tt.isV6 {
			t.Errorf("generation %q: got the v6 deploy = %v, want %v", tt.v, gotV6, tt.isV6)
		}
	}
}
