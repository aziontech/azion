package edge_applications

import (
	"context"
	"fmt"

	"github.com/aziontech/azion-cli/pkg/output"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/messages/general"
	msg "github.com/aziontech/azion-cli/messages/list/edge_applications"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true, Example: heredoc.Doc(`
		$ azion list edge-application
		$ azion list edge-application --details
		$ azion list edge-application --page 1 
		$ azion list edge-application --page-size 5
		`),
		RunE: func(cmd *cobra.Command, _ []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			if err := PrintTable(cmd, client, f, opts); err != nil {
				return fmt.Errorf(msg.ErrorGetAll.Error(), err)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64Var(&opts.Page, "page", 1, general.ApiListFlagPage)
	flags.Int64Var(&opts.PageSize, "page-size", 10, general.ApiListFlagPageSize)
	flags.BoolVar(&opts.Details, "details", false, general.ApiListFlagDetails)
	flags.BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}

func PrintTable(cmd *cobra.Command, client *api.Client, f *cmdutil.Factory, opts *contracts.ListOptions) error {
	c := context.Background()

	for {
		resp, err := client.List(c, opts)
		if err != nil {
			return err
		}

		listOut := output.ListOutput{}
		listOut.Columns = []string{"ID", "NAME", "ACTIVE"}
		listOut.Out = f.IOStreams.Out

		if opts.Details {
			listOut.Columns = []string{"ID", "NAME", "ACTIVE", "LAST EDITOR", "LAST MODIFIED", "DEBUG RULES"}
		}

		for _, v := range resp.Results {
			ln := []string{
				fmt.Sprintf("%d", v.Id),
				fmt.Sprintf("%s", utils.TruncateString(v.Name)),
				fmt.Sprintf("%v", v.Active),
				fmt.Sprintf("%s", v.LastEditor),
				fmt.Sprintf("%s", v.LastModified),
				fmt.Sprintf("%v", v.DebugRules),
			}
			listOut.Lines = append(listOut.Lines, ln)
		}

		listOut.Page = opts.Page
		err = output.Print(&listOut)
		if err != nil {
			return err
		}

		if opts.Page >= resp.TotalPages {
			break
		}

		if cmd.Flags().Changed("page") || cmd.Flags().Changed("page-size") {
			break
		}

		opts.Page++
	}

	return nil
}
