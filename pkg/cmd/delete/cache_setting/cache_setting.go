package cachesetting

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/cache_setting"
	api "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	applicationID   int64
	cacheSettingsID int64
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.DeleteShortDescription,
		Long:          msg.DeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion delete cache-setting --application-id 1673635839 --cache-setting-id 107313
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			if !cmd.Flags().Changed("application-id") {
				answer, err := utils.AskInput(msg.ListAskInputApplicationID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}

				applicationID = num
			}

			if !cmd.Flags().Changed("cache-setting-id") {
				answer, err := utils.AskInput(msg.DeleteAskInputCacheID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}

				cacheSettingsID = num
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err := client.Delete(ctx, applicationID, cacheSettingsID)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDelete.Error(), err)
			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DeleteOutputSuccess, cacheSettingsID),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&deleteOut)

		},
	}

	cmd.Flags().Int64Var(&applicationID, "application-id", 0, msg.DeleteFlagApplicationID)
	cmd.Flags().Int64Var(&cacheSettingsID, "cache-setting-id", 0, msg.DeleteFlagCacheSettingsID)
	cmd.Flags().BoolP("help", "h", false, msg.DeleteHelpFlag)
	return cmd
}
