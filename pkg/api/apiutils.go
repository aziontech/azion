package api

import (
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/spf13/cobra"
)

func AddAzionApiFlags(cmd *cobra.Command, opts *contracts.ListOptions) {

	cmd.Flags().BoolVar(&opts.Details, "details", false, "Show all relevant fields when listing")
	cmd.Flags().StringVar(&opts.Order_by, "order_by", "", "Identifies which field the return should be sorted by")
	cmd.Flags().StringVar(&opts.Sort, "sort", "", "Defines which ordering to be used: <asc|desc>")
	cmd.Flags().Int64Var(&opts.Page, "page", 1, "Identifies which page should be returned")
	cmd.Flags().Int64Var(&opts.Page_size, "page_size", 10, "Identifies how many items should be returned per page")

}
