package cachesetting

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"
	"github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/cache_setting"

	api "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var (
	applicationID   int64
	cacheSettingsID int64
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.DescribeOptions{}
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.DescribeShortDescription,
		Long:          msg.DescribeLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion describe cache-setting --application-id 1673635839 --cache-setting-id 107313
        $ azion describe cache-setting --application-id 1673635839 --cache-setting-id 107313 --format json
        $ azion describe cache-setting --application-id 1673635839 --cache-setting-id 107313 --out "./tmp/test.json" 
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			if !cmd.Flags().Changed("application-id") {
				answer, err := utils.AskInput(msg.DescibeAskInputApplicationID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}

				applicationID = num
			}

			if !cmd.Flags().Changed("cache-setting-id") {
				answer, err := utils.AskInput(msg.DescribeAskInputCacheID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}

				cacheSettingsID = num
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			ctx := context.Background()
			resp, err := client.Get(ctx, applicationID, cacheSettingsID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetCache.Error(), err)
			}

			out := f.IOStreams.Out
			formattedFuction, err := format(cmd, resp)
			if err != nil {
				return utils.ErrorFormatOut
			}

			if cmd.Flags().Changed("out") {
				err := cmdutil.WriteDetailsToFile(formattedFuction, opts.OutPath, out)
				if err != nil {
					return fmt.Errorf("%s: %w", utils.ErrorWriteFile, err)
				}
				fmt.Fprintf(out, msg.FileWritten, opts.OutPath)
			} else {
				_, err := out.Write(formattedFuction[:])
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().Int64Var(&applicationID, "application-id", 0, msg.DescribeFlagApplicationID)
	cmd.Flags().Int64Var(&cacheSettingsID, "cache-setting-id", 0, msg.DescribeFlagCacheSettingsID)
	cmd.Flags().StringVar(&opts.OutPath, "out", "", msg.DescribeFlagOut)
	cmd.Flags().StringVar(&opts.Format, "format", "", msg.DescribeFlagFormat)
	cmd.Flags().BoolP("help", "h", false, msg.DescribeHelpFlag)
	return cmd
}

func format(cmd *cobra.Command, strResp api.GetResponse) ([]byte, error) {
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return nil, err
	}

	if format == "json" || cmd.Flags().Changed("out") {
		return json.MarshalIndent(strResp, "", " ")
	}

	tbl := tablecli.New("", "")
	tbl.WithFirstColumnFormatter(color.New(color.FgGreen).SprintfFunc())
	tbl.AddRow("Id: ", strResp.GetId())
	tbl.AddRow("Name: ", strResp.GetName())
	tbl.AddRow("Browser cache settings: ", strResp.GetBrowserCacheSettings())
	tbl.AddRow("Browser cache settings maximum TTL: ", strResp.GetBrowserCacheSettingsMaximumTtl())
	tbl.AddRow("Cdn cache settings: ", strResp.GetCdnCacheSettings())
	tbl.AddRow("Cdn cache settings maximum TTL: ", strResp.GetCdnCacheSettingsMaximumTtl())
	tbl.AddRow("Cache by query string: ", strResp.GetCacheByQueryString())
	tbl.AddRow("Query string fields: ", strResp.GetQueryStringFields())
	tbl.AddRow("Enable query string sort: ", strResp.GetEnableCachingForPost())
	tbl.AddRow("Cache by cookies: ", strResp.GetCacheByCookies())
	tbl.AddRow("Cookie Names: ", strResp.GetCookieNames())
	tbl.AddRow("Adaptive delivery action: ", strResp.GetAdaptiveDeliveryAction())
	tbl.AddRow("Device group: ", strResp.GetDeviceGroup())
	tbl.AddRow("EnableCachingForPost: ", strResp.GetEnableCachingForPost())
	tbl.AddRow("L2 caching enabled: ", strResp.GetL2CachingEnabled())
	return tbl.GetByteFormat(), nil
}
