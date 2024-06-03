package variables

import (
	"context"
	"fmt"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/variables"

	api "github.com/aziontech/azion-cli/pkg/api/variables"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var variableID string
	opts := &contracts.DescribeOptions{}
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.DescribeShortDescription,
		Long:          msg.DescribeLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
      $ azion describe variables --variable-id 7a187044-4a00-4a4a-93ed-d230900421f3
      $ azion describe variables --variable-id 7a187044-4a00-4a4a-93ed-d230900421f3 --format json
      $ azion describe variables --variable-id 7a187044-4a00-4a4a-93ed-d230900421f3 --out "./tmp/test.json" --format json
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("variable-id") {
				answers, err := utils.AskInput(msg.AskVariableID)

				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				variableID = answers
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			ctx := context.Background()
			variable, err := client.Get(ctx, variableID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetItem.Error(), err)
			}

			fields := make(map[string]string, 0)
			fields["Uuid"] = "Uuid"
			fields["Key"] = "Key"
			fields["Value"] = "Value"
			fields["Secret"] = "Secret"
			fields["LastEditor"] = "Last Editor"
			fields["CreatedAt"] = "Create At"
			fields["UpdatedAt"] = "Update At"

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Msg:   filepath.Clean(opts.OutPath),
					Out:   f.IOStreams.Out,
					Flags: f.Flags,
				},
				Fields: fields,
				Values: variable,
			}
			return output.Print(&describeOut)
		},
	}

	cmd.Flags().StringVar(&variableID, "variable-id", "", msg.FlagVariableID)
	cmd.Flags().BoolP("help", "h", false, msg.DescribeHelpFlag)

	return cmd
}
