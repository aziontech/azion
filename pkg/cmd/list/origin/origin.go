package origin

import (
	"context"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/origin"
	api "github.com/aziontech/azion-cli/pkg/api/origin"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var edgeApplicationID int64 = 0

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ListShortDescription,
		Long:          msg.ListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true, Example: heredoc.Doc(`
        $ azion list origin  --application-id 16736354321
        $ azion list origin  --application-id 16736354321 --details
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") {

				answer, err := utils.AskInput(msg.AskAppID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}

				edgeApplicationID = num
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			if err := PrintTable(client, f, opts); err != nil {
				return fmt.Errorf(msg.ErrorGetOrigins.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	flags := cmd.Flags()
	flags.Int64Var(&edgeApplicationID, "application-id", 0, msg.FlagEdgeApplicationID)
	flags.BoolP("help", "h", false, msg.ListHelpFlag)
	return cmd
}

func PrintTable(client *api.Client, f *cmdutil.Factory, opts *contracts.ListOptions) error {
	c := context.Background()

	resp, err := client.ListOrigins(c, opts, edgeApplicationID)
	if err != nil {
		return err
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ORIGIN KEY", "NAME"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	if opts.Details {
		listOut.Columns = []string{"ORIGIN KEY", "NAME", "ID", "ORIGIN TYPE", "ORIGIN PATH", "ADDRESSES", "CONNECTION TIMEOUT"}
	}

	for _, v := range resp.Results {
		ln := []string{}
		if opts.Details {
			ln = []string{
				*v.OriginKey,
				utils.TruncateString(v.Name),
				fmt.Sprintf("%d", *v.OriginId),
				*v.OriginType,
				*v.OriginPath,
				fmt.Sprintf("%v", v.Addresses),
				fmt.Sprintf("%d", *v.ConnectionTimeout),
			}
		} else {
			ln = []string{
				*v.OriginKey,
				utils.TruncateString(v.Name),
			}
		}
		listOut.Lines = append(listOut.Lines, ln)
	}

	return output.Print(&listOut)
}
