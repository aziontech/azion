package firewallrules

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe/firewall_rules"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/registry"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DescribeCmd struct {
	Io              *iostreams.IOStreams
	AskInput        func(string) (string, error)
	GetFirewallRule func(ctx context.Context, firewallId, ruleId int64) (sdk.FirewallRule, error)
	FirewallID      int64
	RuleID          int64
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		GetFirewallRule: func(ctx context.Context, firewallId, ruleId int64) (sdk.FirewallRule, error) {
			client := registry.FirewallRules(f)
			return client.Get(ctx, firewallId, ruleId)
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
		$ azion describe firewall-rule --firewall-id 4312 --rule-id 42069
		$ azion describe firewall-rule --firewall-id 1337 --rule-id 42069 --out "./firewallrule.json" --format json
		$ azion describe firewall-rule --firewall-id 1337 --rule-id 42069 --format json
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

				describe.FirewallID = num
			}

			if !cmd.Flags().Changed("rule-id") {
				answer, err := describe.AskInput(msg.AskInputRuleID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertRuleId
				}

				describe.RuleID = num
			}

			ctx := context.Background()
			rule, err := describe.GetFirewallRule(ctx, describe.FirewallID, describe.RuleID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetFirewallRule, err.Error())
			}

			fields := make(map[string]string)
			fields["Id"] = "ID"
			fields["Name"] = "Name"
			fields["Description"] = "Description"
			fields["Active"] = "Active"
			fields["LastEditor"] = "Last Editor"
			fields["LastModified"] = "Last Modified"

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Msg:   filepath.Clean(opts.OutPath),
					Flags: f.Flags,
					Out:   f.IOStreams.Out,
				},
				Fields: fields,
				Values: &rule,
			}
			return output.Print(&describeOut)
		},
	}

	cobraCmd.Flags().Int64Var(&describe.FirewallID, "firewall-id", 0, msg.FlagFirewallID)
	cobraCmd.Flags().Int64Var(&describe.RuleID, "rule-id", 0, msg.FlagRuleID)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}
