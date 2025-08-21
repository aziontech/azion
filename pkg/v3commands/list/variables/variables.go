package variables

import (
	"context"
	"fmt"
	"os"

	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/messages/general"
	msg "github.com/aziontech/azion-cli/messages/variables"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	api "github.com/aziontech/azion-cli/pkg/v3api/variables"
	"github.com/spf13/cobra"
)

var dump bool

type ListCmd struct {
	Io               *iostreams.IOStreams
	ListAllVariables func(context.Context, *api.Client, *cmdutil.Factory, *contracts.ListOptions) error
	Dump             bool
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		ListAllVariables: func(ctx context.Context, client *api.Client, f *cmdutil.Factory, opts *contracts.ListOptions) error {
			return listAllVariables(ctx, client, f, opts)
		},
	}
}

func NewCobraCmd(list *ListCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.VariablesListShortDescription,
		Long:          msg.VariablesListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ azion list variables -h
			$ azion list variables --details
			$ azion list variables
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			if err := list.ListAllVariables(context.Background(), client, f, opts); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&opts.Details, "details", false, general.ApiListFlagDetails)
	cmd.Flags().BoolVar(&dump, "dump", false, "")
	cmd.Flags().BoolP("help", "h", false, msg.VariablesListHelpFlag)

	return cmd
}

func listAllVariables(ctx context.Context, client *api.Client, f *cmdutil.Factory, opts *contracts.ListOptions) error {
	resp, err := client.List(ctx)
	if err != nil {
		return err
	}

	if dump {
		err := dumpVariables(resp)
		if err != nil {
			return err
		}
		dumpOut := output.GeneralOutput{
			Msg:   msg.VariablesDump,
			Out:   f.IOStreams.Out,
			Flags: f.Flags,
		}
		return output.Print(&dumpOut)
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "KEY", "VALUE"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	if opts.Details {
		listOut.Columns = []string{"ID", "KEY", "VALUE", "SECRET", "LAST EDITOR"}
	}

	for _, v := range resp {
		var ln []string
		if opts.Details {
			ln = []string{
				v.GetUuid(),
				v.GetKey(),
				utils.TruncateString(v.GetValue()),
				fmt.Sprintf("%v", v.GetSecret()),
				v.GetLastEditor(),
			}
		} else {
			ln = []string{
				v.GetUuid(),
				v.GetKey(),
				utils.TruncateString(v.GetValue()),
			}
		}

		listOut.Lines = append(listOut.Lines, ln)
	}
	return output.Print(&listOut)
}

func dumpVariables(resp []api.Response) error {
	file, err := os.Create(".env")
	if err != nil {
		logger.Debug("Error creating .env file", zap.Error(err))
		return err
	}
	defer file.Close()

	for _, v := range resp {
		envLine := fmt.Sprintf("%s=%s\n", v.GetKey(), v.GetValue())
		_, err := file.WriteString(envLine)
		if err != nil {
			logger.Debug("Error writing to .env file", zap.Error(err))
			return err
		}
	}

	return nil
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewListCmd(f), f)
}
