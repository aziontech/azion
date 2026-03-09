package wafexceptions

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/create/waf_exceptions"
	api "github.com/aziontech/azion-cli/pkg/api/waf_exceptions"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"github.com/spf13/cobra"
)

type Fields struct {
	Name       string `json:"name"`
	RuleID     int64  `json:"rule_id,omitempty"`
	Path       string `json:"path,omitempty"`
	Conditions string `json:"conditions,omitempty"`
	Operator   string `json:"operator,omitempty"`
	Active     string `json:"active,omitempty"`
	InPath     string
	WafID      int64
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
		$ azion create waf-exceptions --name "My Exception" --waf-id 1234
		$ azion create waf-exceptions --name "My Exception" --waf-id 1234 --rule-id 1000
		$ azion create waf-exceptions --file "create.json" --waf-id 1234
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("waf-id") {
				answer, err := utils.AskInput(msg.AskInputWafID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertWafID
				}

				fields.WafID = num
			}

			var request sdk.WAFRuleRequest
			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.InPath+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				if !cmd.Flags().Changed("name") {
					answer, err := utils.AskInput(msg.AskInputName)
					if err != nil {
						return err
					}
					fields.Name = answer
				}

				request = *sdk.NewWAFRuleRequest(fields.Name, []sdk.WAFExceptionConditionRequest{})

				if cmd.Flags().Changed("rule-id") {
					request.SetRuleId(fields.RuleID)
				}

				if cmd.Flags().Changed("path") {
					request.SetPath(fields.Path)
				}

				if cmd.Flags().Changed("operator") {
					request.SetOperator(fields.Operator)
				}

				if cmd.Flags().Changed("active") {
					isActive, err := strconv.ParseBool(fields.Active)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorIsActiveFlag, fields.Active)
					}
					request.SetActive(isActive)
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			response, err := client.Create(context.Background(), fields.WafID, request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreate.Error(), err)
			}

			createOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OutputSuccess, response.GetId()),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&createOut)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.Int64Var(&fields.RuleID, "rule-id", 0, msg.FlagRuleID)
	flags.StringVar(&fields.Path, "path", "", msg.FlagPath)
	flags.StringVar(&fields.Conditions, "conditions", "", msg.FlagConditions)
	flags.StringVar(&fields.Operator, "operator", "", msg.FlagOperator)
	flags.StringVar(&fields.Active, "active", "true", msg.FlagActive)
	flags.StringVar(&fields.InPath, "file", "", msg.FlagFile)
	flags.Int64Var(&fields.WafID, "waf-id", 0, msg.FlagWafID)
	flags.BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}
