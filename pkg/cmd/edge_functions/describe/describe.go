package describe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_function"
	api "github.com/aziontech/azion-cli/pkg/api/edge_function"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var function_id int64
	opts := &contracts.DescribeOptions{}
	cmd := &cobra.Command{
		Use:           msg.DescribeUsage,
		Short:         msg.DescribeShortDescription,
		Long:          msg.DescribeLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion edge_functions describe --function-id 4312
        $ azion edge_functions describe --function-id 1337 --with-code
        $ azion edge_functions describe --function-id 1337 --out "./tmp/test.json" --format json
        $ azion edge_functions describe --function-id 1337 --format json
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("function-id") {
				return msg.ErrorMissingFunctionIdArgument
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			function, err := client.Get(ctx, function_id)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetFunction.Error(), err)
			}

			out := f.IOStreams.Out
			formattedFuction, err := format(cmd, function)
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

	cmd.Flags().Int64VarP(&function_id, "function-id", "f", 0, msg.FlagID)
	cmd.Flags().Bool("with-code", false, msg.DescribeFlagWithCode)
	cmd.Flags().StringVar(&opts.OutPath, "out", "", msg.DescribeFlagOut)
	cmd.Flags().StringVar(&opts.Format, "format", "", msg.DescribeFlagFormat)
	cmd.Flags().BoolP("help", "h", false, msg.DescribeHelpFlag)

	return cmd
}

func serializeToJson(data interface{}) string {
	// ignoring errors on purpose
	serialized, _ := json.Marshal(data)
	return string(serialized)
}

func format(cmd *cobra.Command, function api.EdgeFunctionResponse) ([]byte, error) {

	var b bytes.Buffer

	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return nil, err
	}

	if format == "json" || cmd.Flags().Changed("out") {
		file, err := json.MarshalIndent(function, "", " ")
		if err != nil {
			return nil, err
		}
		return file, nil
	} else {
		b.Write([]byte(fmt.Sprintf("ID: %d\n", uint64(function.GetId()))))
		b.Write([]byte(fmt.Sprintf("Name: %s\n", function.GetName())))
		b.Write([]byte(fmt.Sprintf("Active: %t\n", function.GetActive())))
		b.Write([]byte(fmt.Sprintf("Language: %s\n", function.GetLanguage())))
		b.Write([]byte(fmt.Sprintf("Reference Count: %d\n", uint64(function.GetReferenceCount()))))
		b.Write([]byte(fmt.Sprintf("Modified at: %s\n", function.GetModified())))
		b.Write([]byte(fmt.Sprintf("Initiator Type: %s\n", function.GetInitiatorType())))
		b.Write([]byte(fmt.Sprintf("Last Editor: %s\n", function.GetLastEditor())))
		b.Write([]byte(fmt.Sprintf("Function to run: %s\n", function.GetFunctionToRun())))
		b.Write([]byte(fmt.Sprintf("JSON Args: %s\n", serializeToJson(function.GetJsonArgs()))))
		if cmd.Flags().Changed("with-code") {
			b.Write([]byte(fmt.Sprintf("Code:\n%s\n", function.GetCode())))
		}

		return b.Bytes(), nil
	}
}
