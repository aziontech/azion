package describe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/MakeNowJust/heredoc"
	table "github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/origins"

	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
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
      $ azioncli origins describe --application-id 4312 --origin-id 31223
      $ azioncli origins describe --application-id 1337 --origin-id 31223--format json
      $ azioncli origins describe --application-id 1337 --origin-id 31223--out "./tmp/test.json" --format json
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("origin-id") {
				return msg.ErrorMissingArguments
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			ctx := context.Background()
			origin, err := client.GetOrigin(ctx, applicationID, originID)
			if err != nil {
				return msg.ErrorGetOrigin
			}

      // var mapsStr map[string]interface{}
      // bOrigin, err := json.Marshal(origin)
      // json.Unmarshal(bOrigin, &mapsStr)

			formattedFuction, err := Describe(cmd, origin, resultsSetFieldsOutput)
			if err != nil {
				return utils.ErrorFormatOut
			}

			if cmd.Flags().Changed("out") {
				err := cmdutil.WriteDetailsToFile(formattedFuction, opts.OutPath, f.IOStreams.Out)
				if err != nil {
					return fmt.Errorf("%s: %w", utils.ErrorWriteFile, err)
				}
				fmt.Fprintf(f.IOStreams.Out, msg.OriginsFileWritten, filepath.Clean(opts.OutPath))
			} else {
				_, err := f.IOStreams.Out.Write(formattedFuction[:])
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

type Fields struct {
  key   string 
  value interface{}
}

func Describe(cmd *cobra.Command, str interface{}, funcFields func(fields map[string]interface{}) []Fields ) ([]byte, error) {
  var mapsStr map[string]interface{}
  byteStr, err := json.Marshal(str)
  if err != nil {
    return nil, err
  }

  err = json.Unmarshal(byteStr, &mapsStr)
  if err != nil {
    return nil, err
  }


	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return nil, err
	}

	if format == "json" || cmd.Flags().Changed("out") {
		return json.MarshalIndent(mapsStr, "", " ")
	}

	tbl := table.New("", "")
	tbl.WithFirstColumnFormatter(color.New(color.FgGreen).SprintfFunc())

  for _, p := range funcFields(mapsStr) {
  	tbl.AddRow(p.key, p.value)
  }

	tbl.Print()
	var b bytes.Buffer
	b.WriteTo(tbl.GetWriter())
	return b.Bytes(), nil
}

func resultsSetFieldsOutput(fields map[string]interface{}) []Fields {
  return []Fields{
    {"Origin ID: ", fields["origin_id"]},
    {"Name: ", fields["name"]},
    {"Origin Type: ", fields["origin_type"]},
    {"Addresses: ", fields["addresses"]},
    {"Origin Protocol Policy: ", fields["origin_protocol_policy"]},
    {"Is Origin Redirection Enable: ", fields["is_origin_redirection_enabled"]},
    {"Host Header: ", fields["host_header"]},
    {"Method: ", fields["method"]},
    {"Origin Path: ", fields["origin_path"]},
    {"Connection Timeout: ", fields["connection_timeout"]},
    {"Timeout Between Bytes: ", fields["timeout_between_bytes"]},
    {"Hmac Authentication: ", fields["hmac_authentication"]},
    {"Hmac Region Name: ", fields["hmac_region_name"]},
    {"Hmac Secret Key: ", fields["hmac_secret_key"]},
    {"Hmac Access Key: ", fields["hmac_access_key"]},
  }
}

