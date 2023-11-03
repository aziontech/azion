package list

import (
	"context"
	"github.com/aziontech/azion-cli/pkg/messages/cache_settings"
	"strings"

	"github.com/fatih/color"

	"github.com/MakeNowJust/heredoc"
	table "github.com/MaxwelMazur/tablecli"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/spf13/cobra"
)

var edgeApplicationID int64

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	cmd := &cobra.Command{
		Use:           cache_settings.CacheSettingsListUsage,
		Short:         cache_settings.CacheSettingsListShortDescription,
		Long:          cache_settings.CacheSettingsListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion cache_settings --application-id 16736354321 list --details
        $ azion cache_settings --application-id 16736354321 list --order_by "id"
        $ azion cache_settings --application-id 16736354321 list --page 1  
        $ azion cache_settings --application-id 16736354321 list --page_size 5
        $ azion cache_settings --application-id 16736354321 list --sort "asc" 
        `),

		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") {
				return cache_settings.ErrorMandatoryListFlags
			}

			var numberPage int64 = opts.Page
			if !cmd.Flags().Changed("page") && !cmd.Flags().Changed("page_size") {
				for {
					pages, err := PrintTable(cmd, f, opts, &numberPage)
					if numberPage > pages && err == nil {
						return nil
					}
					if err != nil {
						return cache_settings.ErrorGetCache
					}
				}
			}

			if _, err := PrintTable(cmd, f, opts, &numberPage); err != nil {
				return cache_settings.ErrorGetCache
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().Int64VarP(&edgeApplicationID, "application-id", "a", 0, "")
	cmd.Flags().BoolP("help", "h", false, cache_settings.CacheSettingsFlagHelp)
	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions, numberPage *int64) (int64, error) {
	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	applications, err := client.ListCacheSettings(ctx, opts, edgeApplicationID)
	if err != nil {
		return 0, cache_settings.ErrorGetCache
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
