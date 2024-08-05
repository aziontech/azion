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
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"github.com/spf13/cobra"
)

type ListCmd struct {
	Io                *iostreams.IOStreams
	ReadInput         func(string) (string, error)
	ListRulesEngine   func(context.Context, *contracts.ListOptions, int64, string) (*edgeapplications.RulesEngineResponse, error)
	AskInput          func(string) (string, error)
	EdgeApplicationID int64
	Phase             string
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		ListRulesEngine: func(ctx context.Context, opts *contracts.ListOptions, appID int64, phase string) (*edgeapplications.RulesEngineResponse, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			return client.ListRulesEngine(ctx, opts, appID, phase)
		},
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
	}
}

func NewCobraCmd(list *ListCmd, f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           msg.RulesEngineListUsage,
		Short:         msg.RulesEngineListShortDescription,
		Long:          msg.RulesEngineListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ azion list rules-engine --application-id 1673635839 --phase request
			$ azion list rules-engine --application-id 1673635839 --phase response --details
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") {
				answer, err := list.AskInput(msg.AskInputApplicationId)
				if err != nil {
					return err
				}
				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}
				list.EdgeApplicationID = num
			}

			if !cmd.Flags().Changed("phase") {
				answer, err := list.AskInput(msg.AskInputPhase)
				if err != nil {
					return err
				}
				list.Phase = answer
			}

			opts := &contracts.ListOptions{}
			if err := PrintTable(cmd, f, opts, list); err != nil {
				return msg.ErrorGetRulesEngines
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, &contracts.ListOptions{})
	cmd.Flags().Int64Var(&list.EdgeApplicationID, "application-id", 0, msg.ApplicationFlagId)
	cmd.Flags().StringVar(&list.Phase, "phase", "request", msg.RulesEnginePhase)
	cmd.Flags().BoolP("help", "h", false, msg.RulesEngineListHelpFlag)

	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions, list *ListCmd) error {
	ctx := context.Background()

	rules, err := list.ListRulesEngine(ctx, opts, list.EdgeApplicationID, list.Phase)
	if err != nil {
		return err
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "NAME"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	if opts.Details {
		listOut.Columns = []string{"ID", "NAME", "ORDER", "PHASE", "ACTIVE"}
	}

	for _, v := range rules.Results {
		var ln []string
		if opts.Details {
			ln = []string{
				fmt.Sprintf("%d", v.Id),
				v.Name,
				fmt.Sprintf("%d", v.Order),
				v.Phase,
				fmt.Sprintf("%v", v.IsActive),
			}
		} else {
			ln = []string{
				fmt.Sprintf("%d", v.Id),
				v.Name,
			}
		}
		listOut.Lines = append(listOut.Lines, ln)
	}
	return output.Print(&listOut)
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewListCmd(f), f)
}
