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
	api "github.com/aziontech/azion-cli/pkg/api/variables"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/spf13/cobra"
)

var dump bool

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}

	listCmd := &cobra.Command{
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

			if err := listAllVariables(client, f, opts); err != nil {
				return err
			}
			return nil
		},
	}

	listCmd.Flags().BoolVar(&opts.Details, "details", false, general.ApiListFlagDetails)
	listCmd.Flags().BoolVar(&dump, "dump", false, "")
	listCmd.Flags().BoolP("help", "h", false, msg.VariablesListHelpFlag)
	return listCmd
}

func listAllVariables(client *api.Client, f *cmdutil.Factory, opts *contracts.ListOptions) error {
	c := context.Background()

	resp, err := client.List(c)
	if err != nil {
		return err
	}

	if dump {
		err := dumpVariables(resp)
		if err != nil {
			return err
		}
		dumpOut := output.GeneralOutput{
			Msg:         msg.VariablesDump,
			Out:         f.IOStreams.Out,
			FlagOutPath: f.Out,
			FlagFormat:  f.Format,
		}
		return output.Print(&dumpOut)
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "KEY", "VALUE"}
	listOut.Out = f.IOStreams.Out
	listOut.FlagOutPath = f.Out
	listOut.FlagFormat = f.Format

	if opts.Details {
		listOut.Columns = []string{"ID", "KEY", "VALUE", "SECRET", "LAST EDITOR"}
	}

	for _, v := range resp {
		ln := []string{
			v.GetUuid(),
			v.GetKey(),
			utils.TruncateString(v.GetValue()),
			fmt.Sprintf("%v", v.GetSecret()),
			v.GetLastEditor(),
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
