package ruleengine

import (
	"context"
	"fmt"
	"io"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/list/rules_engine"
	api "github.com/aziontech/azion-cli/pkg/api/applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/spf13/cobra"
)

var phase string

type ListCmd struct {
	Io                      *iostreams.IOStreams
	ReadInput               func(string) (string, error)
	ListRulesEngineRequest  func(context.Context, *contracts.ListOptions, string) (*sdk.PaginatedApplicationRequestPhaseRuleEngineList, error)
	ListRulesEngineResponse func(context.Context, *contracts.ListOptions, string) (*sdk.PaginatedApplicationResponsePhaseRuleEngineList, error)
	AskInput                func(string) (string, error)
	EdgeApplicationID       string
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		ListRulesEngineRequest: func(ctx context.Context, opts *contracts.ListOptions, appID string) (*sdk.PaginatedApplicationRequestPhaseRuleEngineList, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.ListRulesEngineRequest(ctx, opts, appID)
		},
		ListRulesEngineResponse: func(ctx context.Context, opts *contracts.ListOptions, appID string) (*sdk.PaginatedApplicationResponsePhaseRuleEngineList, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.ListRulesEngineResponse(ctx, opts, appID)
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

			if !cmd.Flags().Changed("phase") {
				answer, err := list.AskInput(msg.AskInputPhase)
				if err != nil {
					return err
				}

				phase = answer
			}

			if err := PrintTable(cmd, f, opts, list); err != nil {
				return fmt.Errorf(msg.ErrorGetRulesEngines.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().StringVar(&list.EdgeApplicationID, "application-id", "", msg.ApplicationFlagId)
	cmd.Flags().StringVar(&phase, "phase", "request", msg.RulesEnginePhase)
	cmd.Flags().BoolP("help", "h", false, msg.RulesEngineListHelpFlag)

	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions, list *ListCmd) error {
	ctx := context.Background()

	switch phase {
	case "request":
		rules, err := list.ListRulesEngineRequest(ctx, opts, list.EdgeApplicationID)
		if err != nil {
			return err
		}

		extractor := func(rule sdk.ApplicationRequestPhaseRuleEngine, details bool) []string {
			if details {
				return []string{
					fmt.Sprintf("%d", rule.Id),
					rule.Name,
					fmt.Sprintf("%d", rule.Order),
					phase,
					fmt.Sprintf("%v", rule.Active),
				}
			}
			return []string{
				fmt.Sprintf("%d", rule.Id),
				rule.Name,
			}
		}

		// Ensure rules.Results is not nil before passing to RenderList
		if rules == nil {
			rules = &sdk.PaginatedApplicationRequestPhaseRuleEngineList{}
		}
		if rules.Results == nil {
			rules.Results = []sdk.ApplicationRequestPhaseRuleEngine{}
		}
		return RenderList(rules.Results, opts.Details, f.IOStreams.Out, f.Flags, extractor)

	case "response":
		rules, err := list.ListRulesEngineResponse(ctx, opts, list.EdgeApplicationID)
		if err != nil {
			return err
		}

		extractor := func(rule sdk.ApplicationResponsePhaseRuleEngine, details bool) []string {
			if details {
				return []string{
					fmt.Sprintf("%d", rule.Id),
					rule.Name,
					fmt.Sprintf("%d", rule.Order),
					phase,
					fmt.Sprintf("%v", rule.Active),
				}
			}
			return []string{
				fmt.Sprintf("%d", rule.Id),
				rule.Name,
			}
		}

		// Ensure rules.Results is not nil before passing to RenderList
		if rules == nil {
			rules = &sdk.PaginatedApplicationResponsePhaseRuleEngineList{}
		}
		if rules.Results == nil {
			rules.Results = []sdk.ApplicationResponsePhaseRuleEngine{}
		}
		return RenderList(rules.Results, opts.Details, f.IOStreams.Out, f.Flags, extractor)
	default:
		return msg.ErrorInvalidPhase

	}
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewListCmd(f), f)
}

type Rule interface{}

type ExtractListDataFunc[T Rule] func(T, bool) []string

func RenderList[T Rule](items []T, details bool, outWriter io.Writer, flags cmdutil.Flags, extract ExtractListDataFunc[T]) error {
	listOut := output.ListOutput{}
	listOut.Out = outWriter
	listOut.Flags = flags

	if details {
		listOut.Columns = []string{"ID", "NAME", "ORDER", "PHASE", "ACTIVE"}
	} else {
		listOut.Columns = []string{"ID", "NAME"}
	}

	if len(items) == 0 {
		return output.Print(&listOut)
	}

	for _, item := range items {
		line := extract(item, details)
		listOut.Lines = append(listOut.Lines, line)
	}

	return output.Print(&listOut)
}
