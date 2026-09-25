package functioninstance

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete/function_instance"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/registry"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DeleteCmd struct {
	Io                 *iostreams.IOStreams
	ReadInput          func(string) (string, error)
	AskInput           func(string) (string, error)
	EdgeApplicationID  int64
	FunctionInstanceID int64
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		AskInput: utils.AskInput,
	}
}

func NewCobraCmd(delete *DeleteCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion delete function-instance --application-id 1234 --instance-id 4321
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if !cmd.Flags().Changed("application-id") {
				answer, err := delete.AskInput(msg.AskDeleteApplicationInput)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertApplicationId
				}

				delete.EdgeApplicationID = num
			}

			if !cmd.Flags().Changed("instance-id") {
				answer, err := delete.AskInput(msg.AskDeleteInput)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertFunctionInstanceId
				}

				delete.FunctionInstanceID = num
			}

			client := registry.FunctionInstance(f)

			ctx := context.Background()

			err = client.Delete(ctx, delete.EdgeApplicationID, delete.FunctionInstanceID)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDeletInstance.Error(), err)
			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OutputSuccess, delete.FunctionInstanceID),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&deleteOut)
		},
	}

	cobraCmd.Flags().Int64Var(&delete.EdgeApplicationID, "application-id", 0, msg.FlagId)
	cobraCmd.Flags().Int64Var(&delete.FunctionInstanceID, "instance-id", 0, msg.FlagId)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f), f)
}
