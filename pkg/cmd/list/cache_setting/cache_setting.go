package cachesetting

import (
	"context"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"
	table "github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/cache_setting"
	api "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
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

		tbl := table.New("ID", "NAME", "BROWSER CACHE SETTINGS")
		tbl.WithWriter(f.IOStreams.Out)

		if cmd.Flags().Changed("details") {
			tbl = table.New("ID", "NAME", "BROWSER CACHE SETTINGS", "CDN CACHE SETTINGS", "CACHE BY COOKIES", "ENABLE CACHING FOR POST")
		}

		headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgGreen).SprintfFunc()
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		for _, v := range cache.Results {
			tbl.AddRow(v.Id, v.Name, v.BrowserCacheSettings, v.CdnCacheSettings, v.CacheByCookies, v.EnableCachingForPost)
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
