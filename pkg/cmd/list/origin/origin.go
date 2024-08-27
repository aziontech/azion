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
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"github.com/spf13/cobra"
)

type ListCmd struct {
	Io                *iostreams.IOStreams
	ReadInput         func(string) (string, error)
	ListOrigins       func(context.Context, *contracts.ListOptions, int64) (*edgeapplications.OriginsResponse, error)
	AskInput          func(string) (string, error)
	EdgeApplicationID int64
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		ListOrigins: func(ctx context.Context, opts *contracts.ListOptions, appID int64) (*edgeapplications.OriginsResponse, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			return client.ListOrigins(ctx, opts, appID)
		},
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
	}
}

func NewCobraCmd(list *ListCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ListShortDescription,
		Long:          msg.ListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ azion list origin --application-id 16736354321
			$ azion list origin --application-id 16736354321 --details
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") {
				answer, err := list.AskInput(msg.AskAppID)
				if err != nil {
					return err
				}
				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}
				list.EdgeApplicationID = num
			}

			if err := PrintTable(cmd, f, opts, list); err != nil {
				return fmt.Errorf(msg.ErrorGetOrigins.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().Int64Var(&list.EdgeApplicationID, "application-id", 0, msg.FlagEdgeApplicationID)
	cmd.Flags().BoolP("help", "h", false, msg.ListHelpFlag)

	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions, list *ListCmd) error {
	ctx := context.Background()

	resp, err := list.ListOrigins(ctx, opts, list.EdgeApplicationID)
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
		var ln []string
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

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewListCmd(f), f)
}
