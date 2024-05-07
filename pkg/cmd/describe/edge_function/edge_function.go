package edgefunction

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_function"
	api "github.com/aziontech/azion-cli/pkg/api/edge_function"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var function_id int64
	opts := &contracts.DescribeOptions{}
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.DescribeShortDescription,
		Long:          msg.DescribeLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion describe edge-function --function-id 4312
        $ azion describe edge-function --function-id 1337 --with-code
        $ azion describe edge-function --function-id 1337 --out "./tmp/test.json" --format json
        $ azion describe edge-function --function-id 1337 --format json
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
			resp, err := client.Get(ctx, function_id)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetFunction.Error(), err)
			}

			fields := make(map[string]string, 0)
			fields["Id"] = "ID"
			fields["Name"] = "Name"
			fields["Active"] = "Active"
			fields["Language"] = "Language"
			fields["ReferenceCount"] = "Reference Count"
			fields["Modified"] = "Modified at"
			fields["InitiatorType"] = "Initiator Type"
			fields["LastEditor"] = "Last Editor"
			fields["FunctionToRun"] = "Function to run"
			fields["JsonArgs"] = "JSON Args"

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Msg:         filepath.Clean(opts.OutPath),
					FlagOutPath: f.Out,
					FlagFormat:  f.Format,
					Out:         f.IOStreams.Out,
				},
				Fields: fields,
				Values: resp,
			}

			if cmd.Flags().Changed("with-code") {
				describeOut.Field = resp.GetCode()
			}

			return output.Print(&describeOut)
		},
	}

	cmd.Flags().Int64Var(&function_id, "function-id", 0, msg.FlagID)
	cmd.Flags().Bool("with-code", false, msg.DescribeFlagWithCode)
	cmd.Flags().BoolP("help", "h", false, msg.DescribeHelpFlag)

	return cmd
}
