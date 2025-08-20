package edgeconnector

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_connector"
	api "github.com/aziontech/azion-cli/pkg/api/edge_connector"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/spf13/cobra"
)

var (
	connectorID string
)

type DescribeCmd struct {
	Io       *iostreams.IOStreams
	AskInput func(string) (string, error)
	Get      func(context.Context, string) (sdk.EdgeConnectorPolymorphic, error)
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		Get: func(ctx context.Context, connectorID string) (sdk.EdgeConnectorPolymorphic, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.Get(ctx, connectorID)
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
        $ azion describe edge-connector --connector-id 4312
        $ azion describe edge-connector --connector-id 1337 --out "./tmp/test.json" --format json
        $ azion describe edge-connector --connector-id 1337 --format json
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("connector-id") {
				answer, err := describe.AskInput(msg.AskEdgeConnectorID)
				if err != nil {
					return err
				}

				connectorID = answer
			}

			ctx := context.Background()
			resp, err := describe.Get(ctx, connectorID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetConnector.Error(), err)
			}

			fields := map[string]string{
				"Id":           "ID",
				"Name":         "Name",
				"Active":       "Active",
				"Type":         "Language",
				"LastModified": "Modified at",
				"LastEditor":   "Last Editor",
			}

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Out:   f.IOStreams.Out,
					Msg:   filepath.Clean(opts.OutPath),
					Flags: f.Flags,
				},
				Fields: fields,
				Values: &resp,
			}

			return output.Print(&describeOut)
		},
	}

	cobraCmd.Flags().StringVar(&connectorID, "connector-id", "", msg.FlagID)
	cobraCmd.Flags().BoolP("help", "h", false, msg.DescribeHelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}
