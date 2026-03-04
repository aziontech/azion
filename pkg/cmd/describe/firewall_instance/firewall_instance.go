package firewallinstance

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe/firewall_instance"
	api "github.com/aziontech/azion-cli/pkg/api/firewall_instance"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	firewallID int64
	instanceID int64
)

type DescribeCmd struct {
	Io                          *iostreams.IOStreams
	AskInput                    func(string) (string, error)
	GetFirewallFunctionInstance func(ctx context.Context, firewallId, instanceId int64) (sdk.FirewallFunctionInstance, error)
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		GetFirewallFunctionInstance: func(ctx context.Context, firewallId, instanceId int64) (sdk.FirewallFunctionInstance, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.Get(ctx, firewallId, instanceId)
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
		$ azion describe firewall-instance --firewall-id 4312 --instance-id 42069
		$ azion describe firewall-instance --firewall-id 1337 --instance-id 42069 --out "./firewallinstance.json" --format json
		$ azion describe firewall-instance --firewall-id 1337 --instance-id 42069 --format json
		`),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !cmd.Flags().Changed("firewall-id") {
				answer, err := describe.AskInput(msg.AskInputFirewallID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertFirewallId
				}

				firewallID = num
			}

			if !cmd.Flags().Changed("instance-id") {
				answer, err := describe.AskInput(msg.AskInputFirewallFunctionInstanceID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertFirewallFunctionInstanceId
				}

				instanceID = num
			}

			ctx := context.Background()
			instance, err := describe.GetFirewallFunctionInstance(ctx, firewallID, instanceID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetFirewallFunctionInstance, err.Error())
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

	cobraCmd.Flags().Int64Var(&firewallID, "firewall-id", 0, msg.FlagFirewallID)
	cobraCmd.Flags().Int64Var(&instanceID, "instance-id", 0, msg.FlagFirewallFunctionInstanceID)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}
