package init

import (
	"os"

	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/logger"
)

func InitNextjs(info *InitInfo, cmd *InitCmd) error {
	pathWorker := info.PathWorkingDir + "/worker"
	if err := cmd.Mkdir(pathWorker, os.ModePerm); err != nil {
		return msg.ErrorFailedCreatingWorkerDirectory
	}

	logger.FInfo(cmd.Io.Out, `	[ General Instructions ]
    - Requirements:
        - Tools: npm
    [ Usage ]
    	- Build Command: azion build
    	- Deploy Command: azion deploy`) //nolint:all

	return nil
}
