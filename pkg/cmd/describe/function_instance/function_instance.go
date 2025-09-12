package functioninstance

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe/function_instance"
	api "github.com/aziontech/azion-cli/pkg/api/function_instance"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/spf13/cobra"
)

var (
	applicationID string
	instanceID    string
)

type DescribeCmd struct {
	Io                  *iostreams.IOStreams
	AskInput            func(string) (string, error)
	GetFunctionInstance func(ctx context.Context, applicationId, instanceId string) (sdk.ApplicationFunctionInstance, error)
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		GetFunctionInstance: func(ctx context.Context, applicationId, instanceId string) (sdk.ApplicationFunctionInstance, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.Get(ctx, applicationId, instanceId)
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
		$ azion describe function-instance --application-id 4312 --instance-id 42069
		$ azion describe function-instance --application-id 1337 --instance-id 42069 --out "./functioninstance.json" --format json
		$ azion describe function-instance --application-id 1337 --instance-id 42069 --format json
		`),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !cmd.Flags().Changed("application-id") {
				answer, err := describe.AskInput(msg.AskInputApplicationID)
				if err != nil {
					return err
				}

				applicationID = answer
			}

			if !cmd.Flags().Changed("instance-id") {
				answer, err := describe.AskInput(msg.AskInputFunctionInstanceID)
				if err != nil {
					return err
				}

				instanceID = answer
			}

			ctx := context.Background()
			instance, err := describe.GetFunctionInstance(ctx, applicationID, instanceID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetFunctionInstance, err.Error())
			}

			fields := make(map[string]string)
			fields["Id"] = "ID"
			fields["Name"] = "Name"
			fields["Function"] = "Function"
			fields["LastEditor"] = "Last Editor"
			fields["LastModified"] = "Last Modified"

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Msg:   filepath.Clean(opts.OutPath),
					Flags: f.Flags,
					Out:   f.IOStreams.Out,
				},
				Fields: fields,
				Values: &instance,
			}
			return output.Print(&describeOut)
		},
	}

	cobraCmd.Flags().StringVar(&applicationID, "application-id", "", msg.FlagApplicationID)
	cobraCmd.Flags().StringVar(&instanceID, "instance-id", "", msg.FlagFunctionInstanceID)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}
