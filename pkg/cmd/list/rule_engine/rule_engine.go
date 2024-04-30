package ruleengine

import (
	"context"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/list/rules_engine"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var edgeApplicationID int64
var phase string

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	cmd := &cobra.Command{
		Use:           msg.RulesEngineListUsage,
		Short:         msg.RulesEngineListShortDescription,
		Long:          msg.RulesEngineListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true, Example: heredoc.Doc(`
		$ azion list rules-engine --application-id 1673635839 --phase request
		$ azion list rules-engine --application-id 1673635839 --phase response --details
		`),
		RunE: func(cmd *cobra.Command, args []string) error {

			if !cmd.Flags().Changed("application-id") {

				answer, err := utils.AskInput(msg.AskInputApplicationId)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}

				edgeApplicationID = num
			}

			if !cmd.Flags().Changed("phase") {

				answer, err := utils.AskInput(msg.AskInputPhase)
				if err != nil {
					return err
				}

				phase = answer
			}

			if err := PrintTable(cmd, f, opts); err != nil {
				return fmt.Errorf(msg.ErrorGetRulesEngines.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().BoolP("help", "h", false, msg.RulesEngineListHelpFlag)
	cmd.Flags().Int64Var(&edgeApplicationID, "application-id", 0, msg.ApplicationFlagId)
	cmd.Flags().StringVar(&phase, "phase", "request", msg.RulesEnginePhase)
	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions) error {
	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	rules, err := client.ListRulesEngine(ctx, opts, edgeApplicationID, phase)
	if err != nil {
		return err
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "NAME"}
	listOut.Out = f.IOStreams.Out
	listOut.FlagOutPath = f.Out
	listOut.FlagFormat = f.Format

	if cmd.Flags().Changed("details") {
		listOut.Columns = []string{"ID", "NAME", "ORDER", "PHASE", "ACTIVE"}
	}

	for _, v := range rules.Results {
		ln := []string{
			fmt.Sprintf("%d", v.Id),
			v.Name,
			fmt.Sprintf("%d", v.Order),
			v.Phase,
			fmt.Sprintf("%v", v.IsActive),
		}
		listOut.Lines = append(listOut.Lines, ln)
	}
	return output.Print(&listOut)
}
