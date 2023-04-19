package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/MakeNowJust/heredoc"
	"github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/origins"

	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"github.com/spf13/cobra"
)

var (
	applicationID int64
	originID      int64
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.DescribeOptions{}
	cmd := &cobra.Command{
		Use:           msg.OriginsDescribeUsage,
		Short:         msg.OriginsDescribeShortDescription,
		Long:          msg.OriginsDescribeLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
      $ azioncli origins describe --application-id 1673635839 --origin-id 31223
      $ azioncli origins describe --application-id 1673635839 --origin-id 31223--format json
      $ azioncli origins describe --application-id 1673635839 --origin-id 31223--out "./tmp/test.json" --format json
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("origin-id") {
				return msg.ErrorMissingArguments
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			ctx := context.Background()
			origin, err := client.GetOrigin(ctx, applicationID, originID)
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

	cmd.Flags().Int64VarP(&applicationID, "application-id", "a", 0, msg.OriginsDescribeFlagApplicationID)
	cmd.Flags().Int64VarP(&originID, "origin-id", "o", 0, msg.OriginsDescribeFlagOriginID)
	cmd.Flags().StringVar(&opts.OutPath, "out", "", msg.OriginsDescribeFlagOut)
	cmd.Flags().StringVar(&opts.Format, "format", "", msg.OriginsDescribeFlagFormat)
	cmd.Flags().BoolP("help", "h", false, msg.OriginsDescribeHelpFlag)

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
