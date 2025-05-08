package ruleengine

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/list/rules_engine"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk/edge"
	"github.com/spf13/cobra"
)

type ListCmd struct {
	Io                *iostreams.IOStreams
	ReadInput         func(string) (string, error)
	ListRulesEngine   func(context.Context, *contracts.ListOptions, string) (*sdk.PaginatedResponseListEdgeApplicationRuleEngineList, error)
	AskInput          func(string) (string, error)
	EdgeApplicationID string
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		ListRulesEngine: func(ctx context.Context, opts *contracts.ListOptions, appID string) (*sdk.PaginatedResponseListEdgeApplicationRuleEngineList, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.ListRulesEngine(ctx, opts, appID)
		},
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
	}
}

func NewCobraCmd(list *ListCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	cmd := &cobra.Command{
		Use:           msg.RulesEngineListUsage,
		Short:         msg.RulesEngineListShortDescription,
		Long:          msg.RulesEngineListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ azion list rules-engine --application-id 1673635839
			$ azion list rules-engine --application-id 1673635839 --details
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") {
				answer, err := list.AskInput(msg.AskInputApplicationId)
				if err != nil {
					return err
				}

				list.EdgeApplicationID = answer
			}

			if err := PrintTable(cmd, f, opts, list); err != nil {
				return msg.ErrorGetRulesEngines
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().StringVar(&list.EdgeApplicationID, "application-id", "", msg.ApplicationFlagId)
	cmd.Flags().BoolP("help", "h", false, msg.RulesEngineListHelpFlag)

	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions, list *ListCmd) error {
	ctx := context.Background()

	rules, err := list.ListRulesEngine(ctx, opts, list.EdgeApplicationID)
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
				fmt.Sprintf("%v", v.Active),
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
