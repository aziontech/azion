package firewall_rule_order

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/update/firewall_rule_order"
	api "github.com/aziontech/azion-cli/pkg/api/firewall_rules"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	FirewallID int64
	RuleIDs    string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion update firewall-rule-order --firewall-id 1234 --rule-ids "123,456,789"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("firewall-id") {
				answer, err := utils.AskInput(msg.AskInputFirewallID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertFirewallID
				}

				fields.FirewallID = num
			}

			if !cmd.Flags().Changed("rule-ids") {
				answer, err := utils.AskInput(msg.AskInputRuleIDs)
				if err != nil {
					return err
				}

				fields.RuleIDs = answer
			}

			order, err := utils.ParseInt64Slice(fields.RuleIDs)
			if err != nil {
				logger.Debug("Error while parsing rule ids", zap.Error(err))
				return msg.ErrorConvertRuleIDs
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			err = client.Order(context.Background(), fields.FirewallID, order)
			if err != nil {
				return fmt.Errorf(msg.ErrorOrder.Error(), err)
			}

			orderOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OutputSuccess, fields.FirewallID),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&orderOut)
		},
	}

	flags := cmd.Flags()
	flags.Int64Var(&fields.FirewallID, "firewall-id", 0, msg.FlagFirewallID)
	flags.StringVar(&fields.RuleIDs, "rule-ids", "", msg.FlagRuleIDs)
	flags.BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}
