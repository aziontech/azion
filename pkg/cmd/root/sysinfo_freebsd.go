//go:build freebsd

package root

import (
	msg "github.com/aziontech/azion-cli/messages/root"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
)

func showUpdadeMessageSystem(f *cmdutil.Factory) error {
	logger.FInfo(f.IOStreams.Out, msg.DownloadRelease)
	logger.FInfo(f.IOStreams.Out, msg.PkgUpdate)
	return nil
}
