package origin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/origin"
	api "github.com/aziontech/azion-cli/pkg/api/origin"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	applicationID int64
	originKey     string
)

type DeleteCmd struct {
	Io            *iostreams.IOStreams
	ReadInput     func(string) (string, error)
	DeleteOrigins func(context.Context, int64, string) error
	AskInput      func(string) (string, error)
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		DeleteOrigins: func(ctx context.Context, appID int64, key string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			return client.DeleteOrigins(ctx, appID, key)
		},
		AskInput: utils.AskInput,
	}
}

func NewCobraCmd(delete *DeleteCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.DeleteShortDescription,
		Long:          msg.DeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		  $ azion delete origin --application-id 1673635839 --origin-key 03a6e7bf-8e26-49c7-a66e-ab8eaa425086
		  $ azion delete origin
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if !cmd.Flags().Changed("application-id") {
				answer, err := delete.AskInput(msg.DeleteAskInputApp)
				if err != nil {
					return err
				}
				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApp
				}
				applicationID = num
			}

			if !cmd.Flags().Changed("origin-key") {
				answer, err := delete.AskInput(msg.DeleteAskInputOri)
				if err != nil {
					return err
				}
				originKey = answer
			}

			ctx := context.Background()

			err = delete.DeleteOrigins(ctx, applicationID, originKey)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDelete.Error(), err)
			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DeleteOutputSuccess, originKey),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&deleteOut)
		},
	}

	cobraCmd.Flags().Int64Var(&applicationID, "application-id", 0, msg.FlagEdgeApplicationID)
	cobraCmd.Flags().StringVar(&originKey, "origin-key", "", msg.FlagOriginKey)
	cobraCmd.Flags().BoolP("help", "h", false, msg.DeleteHelpFlag)
	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f), f)
}
