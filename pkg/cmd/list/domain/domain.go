package domain

import (
	"context"
	"strings"

	"github.com/fatih/color"

	"github.com/MakeNowJust/heredoc"
	table "github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/list/domain"
	api "github.com/aziontech/azion-cli/pkg/api/domain"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
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
		$ azion list domain
		$ azion list domain --details
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := PrintTable(cmd, f, opts); err != nil {
				return msg.ErrorGetDomains
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions) error {
	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	for {
		resp, err := client.List(ctx, opts)
		if err != nil {
			return err
		}

		tbl := table.New("ID", "NAME")
		tbl.WithWriter(f.IOStreams.Out)

		if opts.Details {
			tbl = table.New("ID", "NAME", "EDGE DOMAIN", "DIGITAL CERTIFICATE ID", "EDGE APPLICATION ID", "CNAME ACCESS ONLY", "CNAMES", "ACTIVE")
		}

		headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgGreen).SprintfFunc()
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		for _, v := range resp.Results {
			tbl.AddRow(
				v.GetId(),
				utils.TruncateString(v.GetName()),
				v.GetDomainName(),
				v.GetDigitalCertificateId(),
				v.GetDigitalCertificateId(),
				v.GetCnameAccessOnly(),
				v.GetCnames(),
				v.GetIsActive(),
			)
		}

		format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
		tbl.CalculateWidths([]string{})

		// print the header only in the first flow
		if opts.Page == 1 {
			logger.PrintHeader(tbl, format)
		}

		for _, row := range tbl.GetRows() {
			logger.PrintRow(tbl, format, row)
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
