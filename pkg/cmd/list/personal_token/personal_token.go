package personaltoken

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/messages/general"
	msg "github.com/aziontech/azion-cli/messages/list/personal_token"
	api "github.com/aziontech/azion-cli/pkg/api/personal_token"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/constants"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var details bool

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true, Example: heredoc.Doc(`
        $ azion list personal-token  
        $ azion list personal-token --details
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient,
				f.Config.GetString("api_url"),
				f.Config.GetString("token"),
			)

			if err := PrintTable(client, f, details); err != nil {
				return fmt.Errorf(msg.ErrorList.Error(), err)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&details, "details", false, general.ApiListFlagDetails)
	flags.BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}

func PrintTable(client *api.Client, f *cmdutil.Factory, details bool) error {
	c := context.Background()

	resp, err := client.List(c)
	if err != nil {
		return err
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "NAME", "EXPIRES AT"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	if details {
		listOut.Columns = []string{"ID", "NAME", "EXPIRES AT", "CREATED AT", "DESCRIPTION"}
	}

	for _, v := range resp {

		ln := []string{}
		if details {
			var description string
			if v.Description.Get() != nil {
				description = *v.Description.Get()
			}
			ln = []string{
				*v.Uuid,
				utils.TruncateString(*v.Name),
				v.ExpiresAt.Format(constants.FORMAT_DATE),
				fmt.Sprintf("%v", *v.Created),
				utils.TruncateString(description),
			}
		} else {
			ln = []string{
				*v.Uuid,
				utils.TruncateString(*v.Name),
				v.ExpiresAt.Format(constants.FORMAT_DATE),
			}
		}
		listOut.Lines = append(listOut.Lines, ln)
	}
	return output.Print(&listOut)
}
