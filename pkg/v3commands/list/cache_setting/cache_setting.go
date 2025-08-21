package cachesetting

import (
	"context"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/cache_setting"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	api "github.com/aziontech/azion-cli/pkg/v3api/cache_setting"
	"github.com/aziontech/azion-cli/utils"
	"github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"github.com/spf13/cobra"
)

type ListCmd struct {
	Io                *iostreams.IOStreams
	ReadInput         func(string) (string, error)
	ListCaches        func(context.Context, *contracts.ListOptions, int64) (*edgeapplications.ApplicationCacheGetResponse, error)
	AskInput          func(string) (string, error)
	EdgeApplicationID int64
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		ListCaches: func(ctx context.Context, opts *contracts.ListOptions, appID int64) (*edgeapplications.ApplicationCacheGetResponse, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			return client.List(ctx, opts, appID)
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
			$ azion list cache-setting --application-id 16736354321
			$ azion list cache-setting --application-id 16736354321 --details
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") {
				answer, err := list.AskInput(msg.ListAskInputApplicationID)
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
				return msg.ErrorGetCaches
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

	response, err := list.ListCaches(ctx, opts, list.EdgeApplicationID)
	if err != nil {
		return msg.ErrorGetCaches
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "NAME", "BROWSER CACHE SETTINGS"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	if opts.Details {
		listOut.Columns = []string{"ID", "NAME", "BROWSER CACHE SETTINGS", "CDN CACHE SETTINGS", "CACHE BY COOKIES", "ENABLE CACHING FOR POST"}
	}

	for _, v := range response.Results {
		var ln []string
		if opts.Details {
			ln = []string{
				fmt.Sprintf("%d", v.Id),
				v.Name,
				v.BrowserCacheSettings,
				v.CdnCacheSettings,
				v.CacheByCookies,
				fmt.Sprintf("%v", v.EnableCachingForPost),
			}
		} else {
			ln = []string{
				fmt.Sprintf("%d", v.Id),
				v.Name,
				v.BrowserCacheSettings,
			}
		}

		listOut.Lines = append(listOut.Lines, ln)
	}

	return output.Print(&listOut)
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewListCmd(f), f)
}
