package list

import (
	"context"
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/printer"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
	"github.com/spf13/cobra"
)

type ListOptions struct {
	Limit int64
	Page  int64
	// FIXME: ENG-17161
	SortDesc bool
	Filter   string
	Details  bool
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{}

	// listCmd represents the list command
	listCmd := &cobra.Command{
		Use:           "list <service_id> [flags]",
		Short:         "Lists the Resources in a given Edge Service",
		Long:          `Lists the Resources in a given Edge Service`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_services resources list [--details]
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return utils.ErrorMissingServiceIdArgument
			}

			ids, err := utils.ConvertIdsToInt(args[0])
			if err != nil {
				return utils.ErrorConvertingIdArgumentToInt
			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			if err := listAllResources(client, f.IOStreams.Out, opts, ids[0]); err != nil {
				return err
			}
			return nil
		},
	}

	listCmd.Flags().Int64Var(&opts.Limit, "limit", 10, "Maximum number of items to fetch")
	listCmd.Flags().Int64Var(&opts.Page, "page", 1, "Select the page from results")
	listCmd.Flags().StringVar(&opts.Filter, "filter", "", "Filter results by their name")
	listCmd.Flags().BoolVar(&opts.Details, "details", false, "Show more fields when listing")

	return listCmd
}

func listAllResources(client *sdk.APIClient, out io.Writer, opts *ListOptions, service_id int64) error {
	c := context.Background()
	api := client.DefaultApi

	fields := []string{"Id", "Name"}
	headers := []string{"ID", "NAME"}

	resp, httpResp, err := api.GetResources(c, service_id).
		Page(opts.Page).
		Limit(opts.Limit).
		Filter(opts.Filter).
		Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}
		return err
	}

	resources := resp.Resources

	if len(resources) == 0 {
		return nil
	}

	tp := printer.NewTab(out)
	if opts.Details {
		fields = append(fields, "LastEditor", "UpdatedAt", "ContentType", "Type")
		headers = append(headers, "LAST EDITOR", "LAST MODIFIED", "CONTENT TYPE", "TRIGGER")
	}

	tp.PrintWithHeaders(resources, fields, headers)

	return nil
}
