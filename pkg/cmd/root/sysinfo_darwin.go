//go:build darwin

package root

import (
	msg "github.com/aziontech/azion-cli/messages/root"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
)

func showUpdadeMessageSystem(f *cmdutil.Factory) error {
	logger.FInfo(f.IOStreams.Out, msg.BrewUpdate)
	return nil
}
