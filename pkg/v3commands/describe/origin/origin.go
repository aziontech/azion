package origin

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/origin"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	api "github.com/aziontech/azion-cli/pkg/v3api/origin"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var (
	applicationID int64
	originKey     string
)

type DescribeCmd struct {
	Io       *iostreams.IOStreams
	AskInput func(string) (string, error)
	Get      func(context.Context, int64, string) (api.GetResponse, error)
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		Get: func(ctx context.Context, appID int64, key string) (api.GetResponse, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			return client.Get(ctx, appID, key)
		},
	}
}

func NewCobraCmd(describe *DescribeCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.DescribeOptions{}
	cobraCmd := &cobra.Command{
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
				answer, err := describe.AskInput(msg.AskAppID)
				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return utils.ErrorConvertingStringToInt
				}

				applicationID = num
			}

			if !cmd.Flags().Changed("origin-key") {
				answer, err := describe.AskInput(msg.AskOriginKey)
				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				originKey = answer
			}

			ctx := context.Background()
			origin, err := describe.Get(ctx, applicationID, originKey)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetOrigin.Error(), err)
			}

			fields := make(map[string]string)
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
					Out:   f.IOStreams.Out,
					Msg:   filepath.Clean(opts.OutPath),
					Flags: f.Flags,
				},
				Fields: fields,
				Values: origin,
			}
			return output.Print(&describeOut)
		},
	}

	cobraCmd.Flags().Int64Var(&applicationID, "application-id", 0, msg.FlagEdgeApplicationID)
	cobraCmd.Flags().StringVar(&originKey, "origin-key", "", msg.FlagOriginKey)
	cobraCmd.Flags().BoolP("help", "h", false, msg.DescribeHelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}
