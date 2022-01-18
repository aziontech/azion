package list

import (
	"context"
	"io"

	"github.com/aziontech/azion-cli/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/printer"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/edgeservices-go-sdk"
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
		Use:           "list [flags]",
		Short:         "Lists services of an Azion account",
		Long:          `Lists services of an Azion account`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			if err := listAllServices(client, f.IOStreams.Out, opts); err != nil {
				return err
			}
			return nil
		},
	}

	listCmd.Flags().Int64Var(&opts.Limit, "limit", 10, "Maximum number of items to fetch (default 10)")
	listCmd.Flags().Int64Var(&opts.Page, "page", 1, "Select the page from results (default 1)")
	listCmd.Flags().StringVar(&opts.Filter, "filter", "", "Filter results by their name")
	listCmd.Flags().BoolVar(&opts.Details, "details", false, "Show all relevant fields when listing")

	return listCmd
}

func listAllServices(client *sdk.APIClient, out io.Writer, opts *ListOptions) error {
	c := context.Background()
	api := client.DefaultApi

	resp, httpResp, err := api.GetServices(c).
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

	services := resp.Services

	if len(services) == 0 {
		return nil
	}

	tp := printer.NewTab(out)
	if opts.Details {
		tp.PrintWithHeaders(services, []string{"Id", "Name", "LastEditor", "UpdatedAt", "Active", "BoundNodes"}, []string{"ID", "NAME", "LAST EDITOR", "LAST MODIFIED", "ACTIVE", "BOUND NODES"})
	} else {
		tp.PrintWithHeaders(services, []string{"Id", "Name"}, []string{"ID", "NAME"})
	}

	return nil
}
