package cmdutil

import (
	msg "github.com/aziontech/azion-cli/messages/general"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/spf13/cobra"
)

func AddAzionApiFlags(cmd *cobra.Command, opts *contracts.ListOptions) {
	cmd.Flags().BoolVar(&opts.Details, "details", false, msg.ApiListFlagDetails)
	cmd.Flags().StringVar(&opts.OrderBy, "order_by", "", msg.ApiListFlagOrderBy)
	cmd.Flags().StringVar(&opts.Sort, "sort", "", msg.ApiListFlagSort)
	cmd.Flags().Int64Var(&opts.Page, "page", 1, msg.ApiListFlagPage)
	cmd.Flags().Int64Var(&opts.PageSize, "page_size", 10, msg.ApiListFlagPageSize)
	cmd.Flags().StringVar(&opts.Filter, "filter", "", msg.ApiListFlagFilter)
}
