package edge_storage

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	msg "github.com/aziontech/azion-cli/messages/edge_storage"
	"github.com/aziontech/azion-cli/messages/general"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
)

func NewObject(f *cmdutil.Factory) *cobra.Command {
	object := &Objects{
		Factory: f,
		Options: &contracts.ListOptions{},
	}
	cmd := &cobra.Command{
		Use:           msg.USAGE_OBJECTS,
		Short:         msg.SHORT_DESCRIPTION_LIST_OBJECT,
		Long:          msg.LONG_DESCRIPTION_LIST_OBJECT,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          object.RunE,
		Example:       heredoc.Doc(msg.EXAMPLE_LIST_OBJECT),
	}
	object.AddFlags(cmd.Flags())
	return cmd
}

func (b *Objects) RunE(cmd *cobra.Command, args []string) error {
	if !cmd.Flags().Changed("bucket-name") {
		answer, err := utils.AskInput(msg.ASK_NAME_CREATE_BUCKET)
		if err != nil {
			return err
		}
		b.BucketName = answer
	}
	client := api.NewClient(
		b.Factory.HttpClient,
		b.Factory.Config.GetString("storage_url"),
		b.Factory.Config.GetString("token"))
	return b.PrintTable(client)
}

func (b *Objects) PrintTable(client *api.Client) error {
	c := context.Background()

	settings, err := token.ReadSettings()
	if err != nil {
		return err
	}

	if len(settings.ContinuationToken) > 0 && b.Options.NextPage {
		b.Options.ContinuationToken = settings.ContinuationToken
	}

	resp, err := client.ListObject(c, b.BucketName, b.Options)
	if err != nil {
		return fmt.Errorf(msg.ERROR_LIST_BUCKET, err)
	}

	settings.ContinuationToken = resp.GetContinuationToken()
	err = token.WriteSettings(settings)
	if err != nil {
		return err
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"KEY", "LAST MODIFIED"}
	listOut.Out = b.Factory.IOStreams.Out
	listOut.Flags = b.Factory.Flags

	if b.Options.Details {
		listOut.Columns = []string{"KEY", "LAST MODIFIED", "SIZE", "ETAG"}
	}

	for _, v := range resp.Results {
		var ln []string
		if b.Options.Details {
			ln = []string{
				v.GetKey(),
				fmt.Sprintf("%v", v.GetLastModified()),
				fmt.Sprintf("%v", v.GetSize()),
				v.GetEtag(),
			}
		} else {
			ln = []string{
				v.GetKey(),
				fmt.Sprintf("%v", v.GetLastModified()),
			}
		}
		listOut.Lines = append(listOut.Lines, ln)
	}
	return output.Print(&listOut)
}

func (b *Objects) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&b.BucketName, "bucket-name", "", msg.FLAG_NAME_BUCKET)
	flags.BoolVar(&b.Options.Details, "details", false, msg.FLAG_HELP_DETAILS_OBJECTS)
	flags.Int64Var(&b.Options.PageSize, "page-size", 50, general.ApiListFlagPageSize)
	flags.BoolVar(&b.Options.NextPage, "next-page", false, general.ApiListFlagNextPage)
	flags.BoolP("help", "h", false, msg.FLAG_HELP_LIST_OBJECT)
}
