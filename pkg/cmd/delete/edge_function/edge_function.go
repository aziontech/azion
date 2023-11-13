package edgefunction

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_function"
	api "github.com/aziontech/azion-cli/pkg/api/edge_function"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var function_id int64
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.DeleteShortDescription,
		Long:          msg.DeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion delete edge-function --function-id 1234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("function-id") {
				answer, err := utils.AskInput(msg.AskEdgeFunctionID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdFunction
				}

				function_id = num
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err := client.Delete(ctx, function_id)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDeleteFunction.Error(), err)
			}

			out := f.IOStreams.Out
			fmt.Fprintf(out, msg.DeleteOutputSuccess, function_id)

			return nil
		},
	}

	cmd.Flags().Int64Var(&function_id, "function-id", 0, msg.FlagID)
	cmd.Flags().BoolP("help", "h", false, msg.DeleteHelpFlag)

	return cmd
}
