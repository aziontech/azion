package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/MakeNowJust/heredoc"
	table "github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/domains"
	api "github.com/aziontech/azion-cli/pkg/api/domains"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	cmd := &cobra.Command{
		Use:           msg.DomainsListUsage,
		Short:         msg.DomainsListShortDescription,
		Long:          msg.DomainsListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true, Example: heredoc.Doc(`
		$ azion domains list
		$ azion domains list --details
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := PrintTable(cmd, f, opts); err != nil {
				return fmt.Errorf(msg.ErrorGetDomain.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().BoolP("help", "h", false, msg.DomainsListHelpFlag)
	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions) error {
	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	domains, err := client.List(ctx, opts)
	if err != nil {
		return err
	}

	tbl := table.New("ID", "NAME")
	table.DefaultWriter = f.IOStreams.Out
	if cmd.Flags().Changed("details") {
		tbl = table.New("ID", "NAME", "EDGE DOMAIN", "DIGITAL CERTIFICATE ID", "EDGE APPLICATION ID", "CNAME ACCESS ONLY", "CNAMES", "ACTIVE")
	}

	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	if cmd.Flags().Changed("details") {
		for _, v := range domains.Results {
			tbl.AddRow(v.Id, v.Name, v.DomainName, v.DigitalCertificateId.Get(), v.EdgeApplicationId, v.CnameAccessOnly, v.Cnames, v.IsActive)
		}
	} else {
		for _, v := range domains.Results {
			tbl.AddRow(v.Id, v.Name)
		}
	}

	format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
	tbl.CalculateWidths([]string{})
	tbl.PrintHeader(format)
	for _, row := range tbl.GetRows() {
		tbl.PrintRow(format, row)
	}

	f.IOStreams.Out = table.DefaultWriter
	return nil
}
