package cachesettings

import (
	"context"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"
	table "github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/list/cache_settings"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
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
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion list cache-setting --application-id 16736354321
		$ azion list cache-setting --application-id 16736354321 --details
        `),

		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") {
				answer, err := utils.AskInput(msg.AskInputApplicationID)
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

			var numberPage int64 = opts.Page
			if !cmd.Flags().Changed("page") && !cmd.Flags().Changed("page-size") {
				for {
					pages, err := PrintTable(cmd, f, opts, &numberPage)
					if numberPage > pages && err == nil {
						return nil
					}
					if err != nil {
						return err
					}
				}
			}

			if _, err := PrintTable(cmd, f, opts, &numberPage); err != nil {
				return err
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().Int64Var(&edgeApplicationID, "application-id", 0, "")
	cmd.Flags().BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions, numberPage *int64) (int64, error) {
	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	applications, err := client.ListCacheSettings(ctx, opts, edgeApplicationID)
	if err != nil {
		return 0, msg.ErrorGetCaches
	}

	tbl := table.New("ID", "NAME", "BROWSER CACHE SETTINGS")
	tbl.WithWriter(f.IOStreams.Out)
	if cmd.Flags().Changed("details") {
		tbl = table.New("ID", "NAME", "BROWSER CACHE SETTINGS", "CDN CACHE SETTINGS", "CACHE BY COOKIES", "ENABLE CACHING FOR POST")
	}

	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range applications.Results {
		if cmd.Flags().Changed("details") {
			tbl.AddRow(v.Id, v.Name, v.BrowserCacheSettings, v.CdnCacheSettings, v.CacheByCookies, v.EnableCachingForPost)
		} else {
			tbl.AddRow(v.Id, v.Name, v.BrowserCacheSettings)
		}
	}

	format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
	tbl.CalculateWidths([]string{})
	if *numberPage == 1 {
		tbl.PrintHeader(format)
	}

	for _, row := range tbl.GetRows() {
		tbl.PrintRow(format, row)
	}

	*numberPage += 1
	opts.Page = *numberPage
	return applications.TotalPages, nil
}
