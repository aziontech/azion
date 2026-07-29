package dnszone

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"

	msg "github.com/aziontech/azion-cli/messages/dns_zone"
	api "github.com/aziontech/azion-cli/pkg/api/dns_zone"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"github.com/spf13/cobra"
)

type ListCmd struct {
	Io          *iostreams.IOStreams
	ListDNSZone func(context.Context, *contracts.ListOptions) (*sdk.PaginatedZoneList, error)
	AskInput    func(string) (string, error)
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		ListDNSZone: func(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedZoneList, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.List(ctx, opts)
		},
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
	}
}

func NewCobraCmd(list *ListCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	cmd := &cobra.Command{
		Use:           msg.DNSZoneUsage,
		Short:         msg.DNSZoneListShortDescription,
		Long:          msg.DNSZoneListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ azion list dns-zone
			$ azion list dns-zone --details
			$ azion list dns-zone --order-by "name"
			$ azion list dns-zone --order-by "-name"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := PrintTable(cmd, f, opts, list); err != nil {
				return fmt.Errorf(msg.ErrorListDNSZone.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().BoolP("help", "h", false, msg.DNSZoneListHelpFlag)

	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions, list *ListCmd) error {
	ctx := context.Background()

	response, err := list.ListDNSZone(ctx, opts)
	if err != nil {
		return err
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "NAME", "DOMAIN", "ACTIVE"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	for _, v := range response.GetResults() {
		ln := []string{
			fmt.Sprintf("%d", v.GetId()),
			v.GetName(),
			v.GetDomain(),
			fmt.Sprintf("%t", v.GetActive()),
		}
		listOut.Lines = append(listOut.Lines, ln)
	}

	return output.Print(&listOut)
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewListCmd(f), f)
}
