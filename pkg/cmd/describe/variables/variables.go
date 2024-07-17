package variables

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/variables"

	api "github.com/aziontech/azion-cli/pkg/api/variables"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type DescribeCmd struct {
	Io       *iostreams.IOStreams
	AskInput func(string) (string, error)
	Get      func(context.Context, string) (api.Response, error)
}

var (
	variableID string
)

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		Get: func(ctx context.Context, id string) (api.Response, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			return client.Get(ctx, id)
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
      $ azion describe variables --variable-id 7a187044-4a00-4a4a-93ed-d230900421f3
      $ azion describe variables --variable-id 7a187044-4a00-4a4a-93ed-d230900421f3 --format json
      $ azion describe variables --variable-id 7a187044-4a00-4a4a-93ed-d230900421f3 --out "./tmp/test.json" --format json
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("variable-id") {
				answer, err := describe.AskInput(msg.AskVariableID)
				if err != nil {
					return err
				}

				variableID = answer
			}

			ctx := context.Background()
			variable, err := describe.Get(ctx, variableID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetItem.Error(), err)
			}

			fields := map[string]string{
				"Uuid":       "Uuid",
				"Key":        "Key",
				"Value":      "Value",
				"Secret":     "Secret",
				"LastEditor": "Last Editor",
				"CreatedAt":  "Create At",
				"UpdatedAt":  "Update At",
			}

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Out:   f.IOStreams.Out,
					Msg:   filepath.Clean(opts.OutPath),
					Flags: f.Flags,
				},
				Fields: fields,
				Values: variable,
			}

			return output.Print(&describeOut)
		},
	}

	cobraCmd.Flags().StringVar(&variableID, "variable-id", "", msg.FlagVariableID)
	cobraCmd.Flags().BoolP("help", "h", false, msg.DescribeHelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}
