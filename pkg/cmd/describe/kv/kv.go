package kv

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe/kv"
	api "github.com/aziontech/azion-cli/pkg/api/kv"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/kv-api"
	"github.com/spf13/cobra"
)

var (
	namespace string
)

type DescribeCmd struct {
	Io       *iostreams.IOStreams
	AskInput func(string) (string, error)
	Get      func(ctx context.Context, namespace string) (*sdk.Namespace, error)
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		Get: func(ctx context.Context, namespace string) (*sdk.Namespace, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.Get(ctx, namespace)
		},
	}
}

func NewCobraCmd(describe *DescribeCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.DescribeOptions{}
	cobraCmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion describe kv --namespace "my-namespace"
		$ azion describe kv --namespace "my-namespace" --format json
		`),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !cmd.Flags().Changed("namespace") {
				answer, err := describe.AskInput(msg.AskInputNamespace)
				if err != nil {
					return err
				}

				namespace = answer
			}

			ctx := context.Background()
			resp, err := describe.Get(ctx, namespace)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetNamespace, err)
			}

			fields := map[string]string{
				"Name":           "Name",
				"CreatedAt":      "Created At",
				"LastModifiedAt": "Last Modified At",
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
			return output.Print(&describeOut)

		},
	}

	cobraCmd.Flags().StringVar(&namespace, "namespace", "", msg.FlagNamespace)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}
