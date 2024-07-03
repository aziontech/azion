package edgefunction

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_function"
	api "github.com/aziontech/azion-cli/pkg/api/edge_function"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var functionID int64

type DeleteCmd struct {
	Io             *iostreams.IOStreams
	ReadInput      func(string) (string, error)
	DeleteFunction func(context.Context, int64) error
	AskInput       func(string) (string, error)
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		DeleteFunction: func(ctx context.Context, functionID int64) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			return client.Delete(ctx, functionID)
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
		$ azion delete edge-function --function-id 1234
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if !cmd.Flags().Changed("function-id") {
				answer, err := delete.AskInput(msg.AskEdgeFunctionID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdFunction
				}

				functionID = num
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err = client.Delete(ctx, functionID)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDeleteFunction.Error(), err)
			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DeleteOutputSuccess, functionID),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&deleteOut)
		},
	}

	cobraCmd.Flags().Int64Var(&functionID, "function-id", 0, msg.FlagID)
	cobraCmd.Flags().BoolP("help", "h", false, msg.DeleteHelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f), f)
}
