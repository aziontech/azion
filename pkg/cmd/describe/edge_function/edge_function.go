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
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	functionID int64
)

type DescribeCmd struct {
	Io       *iostreams.IOStreams
	AskInput func(string) (string, error)
	Get      func(context.Context, int64) (api.EdgeFunctionResponse, error)
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		Get: func(ctx context.Context, functionID int64) (api.EdgeFunctionResponse, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			return client.Get(ctx, functionID)
		},
	}
}

func NewCobraCmd(describe *DescribeCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.DescribeOptions{}
	cobraCmd := &cobra.Command{
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
				answer, err := describe.AskInput(msg.AskEdgeFunctionID)
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

			ctx := context.Background()
			resp, err := describe.Get(ctx, functionID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetFunction.Error(), err)
			}

			fields := map[string]string{
				"Id":             "ID",
				"Name":           "Name",
				"Active":         "Active",
				"Language":       "Language",
				"ReferenceCount": "Reference Count",
				"Modified":       "Modified at",
				"InitiatorType":  "Initiator Type",
				"LastEditor":     "Last Editor",
				"FunctionToRun":  "Function to run",
				"JsonArgs":       "JSON Args",
			}

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Out:   f.IOStreams.Out,
					Msg:   filepath.Clean(opts.OutPath),
					Flags: f.Flags,
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

	cobraCmd.Flags().Int64Var(&functionID, "function-id", 0, msg.FlagID)
	cobraCmd.Flags().Bool("with-code", false, msg.DescribeFlagWithCode)
	cobraCmd.Flags().BoolP("help", "h", false, msg.DescribeHelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}
