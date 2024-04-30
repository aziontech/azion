package cachesetting

import (
	"context"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/cache_setting"
	api "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var edgeApplicationID int64

func NewCmd(f *cmdutil.Factory) *cobra.Command {
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
				answer, err := utils.AskInput(msg.ListAskInputApplicationID)
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

			if err := PrintTable(cmd, f, opts); err != nil {
				return msg.ErrorGetCaches
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().Int64Var(&edgeApplicationID, "application-id", 0, msg.FlagEdgeApplicationID)
	cmd.Flags().BoolP("help", "h", false, msg.ListHelpFlag)
	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions) error {
	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	for {
		cache, err := client.List(ctx, opts, edgeApplicationID)
		if err != nil {
			return msg.ErrorGetCaches
		}

		listOut := output.ListOutput{}

		listOut.Columns = []string{"ID", "NAME", "BROWSER CACHE SETTINGS"}
		listOut.Out = f.IOStreams.Out

		if cmd.Flags().Changed("details") {
			listOut.Columns = []string{"ID", "NAME", "BROWSER CACHE SETTINGS", "CDN CACHE SETTINGS", "CACHE BY COOKIES", "ENABLE CACHING FOR POST"}
		}

		for _, v := range cache.Results {
			ln := []string{
				fmt.Sprintf("%d", v.Id),
				v.Name,
				v.BrowserCacheSettings,
				v.CdnCacheSettings,
				v.CacheByCookies,
				fmt.Sprintf("%v", v.EnableCachingForPost),
			}
			listOut.Lines = append(listOut.Lines, ln)
		}

		listOut.Page = opts.Page
		err = output.Print(&listOut)
		if err != nil {
			return err
		}

		if opts.Page >= cache.TotalPages {
			break
		}

		if cmd.Flags().Changed("page") || cmd.Flags().Changed("page-size") {
			break
		}

		opts.Page++
	}

	return nil
}
