package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/MakeNowJust/heredoc"
	"github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/variables"

	api "github.com/aziontech/azion-cli/pkg/api/variables"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var variableID string
	opts := &contracts.DescribeOptions{}
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.DescribeShortDescription,
		Long:          msg.DescribeLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
      $ azion variables describe --variable-id 7a187044-4a00-4a4a-93ed-d230900421f3
      $ azion variables describe --variable-id 7a187044-4a00-4a4a-93ed-d230900421f3 --format json
      $ azion variables describe --variable-id 7a187044-4a00-4a4a-93ed-d230900421f3 --out "./tmp/test.json" --format json
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("variable-id") {
				return msg.ErrorMissingArguments
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			ctx := context.Background()
			variable, err := client.Get(ctx, variableID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetItem.Error(), err)
			}

			out := f.IOStreams.Out
			formattedFuction, err := format(cmd, variable)
			if err != nil {
				return utils.ErrorFormatOut
			}

			if cmd.Flags().Changed("out") {
				err := cmdutil.WriteDetailsToFile(formattedFuction, opts.OutPath, out)
				if err != nil {
					return fmt.Errorf("%s: %w", utils.ErrorWriteFile, err)
				}
				fmt.Fprintf(out, msg.FileWritten, filepath.Clean(opts.OutPath))
			} else {
				_, err := out.Write(formattedFuction[:])
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&variableID, "variable-id", "v", "", msg.FlagVariableID)
	cmd.Flags().StringVar(&opts.Format, "format", "", msg.DescribeFlagFormat)
	cmd.Flags().StringVar(&opts.OutPath, "out", "", msg.DescribeFlagOut)
	cmd.Flags().BoolP("help", "h", false, msg.DescribeHelpFlag)

	return cmd
}

func format(cmd *cobra.Command, variable api.Response) ([]byte, error) {
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return nil, err
	}

	if format == "json" || cmd.Flags().Changed("out") {
		return json.MarshalIndent(variable, "", " ")
	}

	tbl := tablecli.New("", "")
	tbl.WithFirstColumnFormatter(color.New(color.FgGreen).SprintfFunc())
	tbl.AddRow("Uuid: ", variable.GetUuid())
	tbl.AddRow("Key: ", variable.GetKey())
	tbl.AddRow("Value: ", variable.GetValue())
	tbl.AddRow("Secret: ", variable.GetSecret())
	tbl.AddRow("Last Editor: ", variable.GetLastEditor())
	tbl.AddRow("Create At: ", variable.GetCreatedAt())
	tbl.AddRow("Update At: ", variable.GetUpdatedAt())
	return tbl.GetByteFormat(), nil
}
