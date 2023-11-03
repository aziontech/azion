package origin

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/delete/origin"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/origin"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var applicationID int64
	var originKey string
	cmd := &cobra.Command{
		Use:           origin.Usage,
		Short:         origin.ShortDescription,
		Long:          origin.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		  $ azion delete origin --application-id 1673635839 --origin-key 03a6e7bf-8e26-49c7-a66e-ab8eaa425086
		  $ azion delete origin
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") {
				answer, err := utils.AskInput(origin.AskInputApp)
				if err != nil {
					return err
				}
				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return origin.ErrorConvertIdApp
				}
				applicationID = num
			}

			if !cmd.Flags().Changed("origin-key") {
				answer, err := utils.AskInput(origin.AskInputOri)
				if err != nil {
					return err
				}
				originKey = answer
			}

			ctx := context.Background()
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			err := client.DeleteOrigins(ctx, applicationID, originKey)
			if err != nil {
				return fmt.Errorf(origin.ErrorFailToDelete.Error(), err)
			}

			logger.LogSuccess(f.IOStreams.Out, fmt.Sprintf(origin.OutputSuccess, originKey))
			return nil
		},
	}

	cmd.Flags().Int64Var(&applicationID, "application-id", 0, origin.FlagApplicationID)
	cmd.Flags().StringVar(&originKey, "origin-key", "", origin.FlagOriginKey)
	cmd.Flags().BoolP("help", "h", false, origin.HelpFlag)
	return cmd
}
