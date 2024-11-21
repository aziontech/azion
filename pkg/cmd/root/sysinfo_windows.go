//go:build windows

package root

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/root"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
)

func showUpdadeMessageSystem(f *cmdutil.Factory, vNumber string) error {
	logger.FInfo(f.IOStreams.Out, fmt.Sprintf(msg.NotSupported))
	return nil
}
