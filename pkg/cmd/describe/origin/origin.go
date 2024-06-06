package origin

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/origin"

	api "github.com/aziontech/azion-cli/pkg/api/origin"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var (
	applicationID int64
	originKey     string
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
		$ azion describe origin --application-id 1673635839 --origin-key 0000000-00000000-00a0a00s0as0-000000
		$ azion describe origin --application-id 1673635839 --origin-key 0000000-00000000-00a0a00s0as0-000000 --format json
		$ azion describe origin --application-id 1673635839 --origin-key 0000000-00000000-00a0a00s0as0-000000 --out "./tmp/test.json"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") {
				answers, err := utils.AskInput(msg.AskAppID)
				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				appID, err := strconv.Atoi(answers)
				if err != nil {
					logger.Debug("Error while parsing string to integer", zap.Error(err))
					return utils.ErrorConvertingStringToInt
				}

				applicationID = int64(appID)
			}

			if !cmd.Flags().Changed("origin-key") {
				answers, err := utils.AskInput(msg.AskOriginKey)
				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				originKey = answers
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			ctx := context.Background()
			origin, err := client.Get(ctx, applicationID, originKey)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetOrigin.Error(), err)
			}

			fields := make(map[string]string, 0)
			fields["OriginId"] = "Origin ID"
			fields["OriginKey"] = "Origin Key"
			fields["Name"] = "Name"
			fields["OriginType"] = "Origin Type"
			fields["Addresses"] = "Addresses"
			fields["OriginProtocolPolicy"] = "Origin Protocol Policy"
			fields["IsOriginRedirectionEnabled"] = "Is Origin Redirection Enable"
			fields["HostHeader"] = "Host Header"
			fields["Method"] = "Method"
			fields["OriginPath"] = "Origin Path"
			fields["ConnectionTimeout"] = "Connection Timeout"
			fields["TimeoutBetweenBytes"] = "Timeout Between Bytes"
			fields["HmacAuthentication"] = "Hmac Authentication"
			fields["HmacRegionName"] = "Hmac Region Name"
			fields["HmacAccessKey"] = "Hmac Secret Key"
			fields["HmacSecretKey"] = "Hmac Access Key"

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Msg:   filepath.Clean(opts.OutPath),
					Flags: f.Flags,
					Out:   f.IOStreams.Out,
				},
				Fields: fields,
				Values: origin,
			}
			return output.Print(&describeOut)
		},
	}

	cmd.Flags().Int64Var(&applicationID, "application-id", 0, msg.FlagEdgeApplicationID)
	cmd.Flags().StringVar(&originKey, "origin-key", "", msg.FlagOriginKey)
	cmd.Flags().BoolP("help", "h", false, msg.DescribeHelpFlag)

	return cmd
}
