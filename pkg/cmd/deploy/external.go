package deploy

import (
	"github.com/aziontech/azion-cli/pkg/apiversion"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
)

// ExternalDeployer is the deploy init and link start when the user asks to
// deploy right away.
type ExternalDeployer interface {
	ExternalRun(f *cmdutil.Factory, configPath string, sync, local, skipFramework, aliasEnv bool) error
}

// NewExternalDeployer picks the deploy for the account's API generation, so a
// v6 account never runs the v4 pipeline from init or link.
func NewExternalDeployer(f *cmdutil.Factory) ExternalDeployer {
	if f.APIVersion == apiversion.V6 {
		return newV6Deploy(f)
	}
	return NewDeployCmd(f)
}
