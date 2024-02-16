package unlink

import (
	"context"

	app "github.com/aziontech/azion-cli/pkg/cmd/delete/edge_application"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	helpers "github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func shouldClean(f *cmdutil.Factory) bool {
	msg := "Do you want to unlink this project? (y/N)"
	return helpers.Confirm(f.GlobalFlagAll, msg, false)
}

func clean(f *cmdutil.Factory, cmd *UnlinkCmd) error {
	var err error
	var shouldCascade bool
	if empty, _ := cmd.IsDirEmpty("./azion"); !empty {
		if f.GlobalFlagAll {
			shouldCascade = true
		} else {
			answer := helpers.Confirm(f.GlobalFlagAll, "Would you like to delete remote resources as well? (y/N)", false)
			shouldCascade = answer
		}

		if shouldCascade {
			cmd := app.NewDeleteCmd(f)
			ctx := context.Background()
			err := cmd.Cascade(ctx)
			if err != nil {
				return err
			}

		}
		err = cmd.CleanDir("./azion")
		if err != nil {
			logger.Debug("Error while trying to clean azion directory", zap.Error(err))
			return err
		}
		return nil
	}
	return nil
}
