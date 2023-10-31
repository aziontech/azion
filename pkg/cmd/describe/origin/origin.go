package origin

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
	"go.uber.org/zap"

	"github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/describe/origin"

	api "github.com/aziontech/azion-cli/pkg/api/origin"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
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
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       msg.Example,
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

			out := f.IOStreams.Out
			formattedFuction, err := format(cmd, origin)
			if err != nil {
				return utils.ErrorFormatOut
			}

			if cmd.Flags().Changed("out") {
				err := cmdutil.WriteDetailsToFile(formattedFuction, opts.OutPath, out)
				if err != nil {
					return fmt.Errorf("%s: %w", utils.ErrorWriteFile, err)
				}
				fmt.Fprintf(out, msg.OriginsFileWritten, filepath.Clean(opts.OutPath))
			} else {
				_, err := out.Write(formattedFuction[:])
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().Int64Var(&applicationID, "application-id", 0, msg.FlagApplicationID)
	cmd.Flags().StringVar(&originKey, "origin-key", "", msg.FlagOriginKey)
	cmd.Flags().StringVar(&opts.OutPath, "out", "", msg.FlagOut)
	cmd.Flags().StringVar(&opts.Format, "format", "", msg.FlagFormat)
	cmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cmd
}

func format(cmd *cobra.Command, origin sdk.OriginsResultResponse) ([]byte, error) {
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return nil, err
	}

	if format == "json" || cmd.Flags().Changed("out") {
		return json.MarshalIndent(origin, "", " ")
	}

	tbl := tablecli.New("", "")
	tbl.WithFirstColumnFormatter(color.New(color.FgGreen).SprintfFunc())
	tbl.AddRow("Origin ID: ", origin.OriginId)
	tbl.AddRow("Origin Key: ", origin.OriginKey)
	tbl.AddRow("Name: ", origin.Name)
	tbl.AddRow("Origin Type: ", origin.OriginType)
	tbl.AddRow("Addresses: ", origin.Addresses)
	tbl.AddRow("Origin Protocol Policy: ", origin.OriginProtocolPolicy)
	tbl.AddRow("Is Origin Redirection Enable: ", origin.IsOriginRedirectionEnabled)
	tbl.AddRow("Host Header: ", origin.HostHeader)
	tbl.AddRow("Method: ", origin.Method)
	tbl.AddRow("Origin Path: ", origin.OriginPath)
	tbl.AddRow("Connection Timeout: ", origin.ConnectionTimeout)
	tbl.AddRow("Timeout Between Bytes: ", origin.TimeoutBetweenBytes)
	tbl.AddRow("Hmac Authentication: ", origin.HmacAuthentication)
	tbl.AddRow("Hmac Region Name: ", origin.HmacRegionName)
	tbl.AddRow("Hmac Secret Key: ", origin.HmacSecretKey)
	tbl.AddRow("Hmac Access Key: ", origin.HmacAccessKey)
	return tbl.GetByteFormat(), nil
}
