package firewallrules

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/update/firewall_rules"
	api "github.com/aziontech/azion-cli/pkg/api/firewall_rules"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	Path       string
	FirewallID int64
	RuleID     int64
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
        $ azion update firewall-rule --firewall-id 1234 --rule-id 4321 --file "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateUserInput(cmd, fields); err != nil {
				return err
			}

			request := api.NewUpdateRequest()
			err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
			if err != nil {
				logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
				return utils.ErrorUnmarshalReader
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			response, err := client.Update(context.Background(), fields.FirewallID, fields.RuleID, request)
			if err != nil {
				return fmt.Errorf(msg.ErrorUpdate.Error(), err)
			}

			updateOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OutputSuccess, response.GetId()),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&updateOut)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.Int64Var(&fields.FirewallID, "firewall-id", 0, msg.FlagFirewallID)
	flags.Int64Var(&fields.RuleID, "rule-id", 0, msg.FlagRuleID)
	flags.BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}

func validateUserInput(cmd *cobra.Command, fields *Fields) error {
	if !cmd.Flags().Changed("firewall-id") {
		answer, err := utils.AskInput(msg.AskInputFirewallID)
		if err != nil {
			return err
		}

		num, err := strconv.ParseInt(answer, 10, 64)
		if err != nil {
			logger.Debug("Error while converting answer to int64", zap.Error(err))
			return msg.ErrorConvertFirewallId
		}

		fields.FirewallID = num
	}

	if !cmd.Flags().Changed("rule-id") {
		answer, err := utils.AskInput(msg.AskInputRuleID)
		if err != nil {
			return err
		}

		num, err := strconv.ParseInt(answer, 10, 64)
		if err != nil {
			logger.Debug("Error while converting answer to int64", zap.Error(err))
			return msg.ErrorConvertRuleId
		}

		fields.RuleID = num
	}

	if !cmd.Flags().Changed("file") {
		answer, err := utils.AskInput(msg.AskInputPathFile)
		if err != nil {
			return err
		}

		fields.Path = answer
	}

	return nil
}
