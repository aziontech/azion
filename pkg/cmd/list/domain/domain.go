package domain

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/list/domain"
	api "github.com/aziontech/azion-cli/pkg/api/domain"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/output"
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

	resp, err := client.List(ctx, opts)
	if err != nil {
		return err
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "NAME"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	if opts.Details {
		listOut.Columns = []string{"ID", "NAME", "EDGE DOMAIN", "DIGITAL CERTIFICATE ID", "EDGE APPLICATION ID", "CNAME ACCESS ONLY", "CNAMES", "ACTIVE"}
	}

	for _, v := range resp.Results {
		ln := []string{}
		if opts.Details {
			ln = []string{
				fmt.Sprintf("%d", v.GetId()),
				utils.TruncateString(v.GetName()),
				v.GetDomainName(),
				fmt.Sprintf("%d", v.GetDigitalCertificateId()),
				fmt.Sprintf("%d", v.GetDigitalCertificateId()),
				fmt.Sprintf("%v", v.GetCnameAccessOnly()),
				fmt.Sprintf("%v", v.GetCnames()),
				fmt.Sprintf("%v", v.GetIsActive()),
			}
		} else {
			ln = []string{
				fmt.Sprintf("%d", v.GetId()),
				utils.TruncateString(v.GetName()),
			}
		}

		listOut.Lines = append(listOut.Lines, ln)
	}

	return output.Print(&listOut)
}
