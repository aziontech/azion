package describe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.DescribeOptions{}
	cmd := &cobra.Command{
		Use:           "describe <edge_function_id> [flags]",
		Short:         "Describe a given Edge Function",
		Long:          "Describe a given Edge Function",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions describe 4312
        $ azioncli edge_functions describe 1337 --with-code
        $ azioncli edge_functions describe 1337 --out "./tmp/test.json" --format json
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("missing edge function id argument")
			}

			ids, err := utils.ConvertIdsToInt(args[0])
			if err != nil {
				return err
			}

			httpClient, err := f.HttpClient()
			if err != nil {
				return fmt.Errorf("failed to get http client: %w", err)
			}

			client := api.NewClient(httpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			function, err := client.Get(ctx, ids[0])
			if err != nil {
				return err
			}

			out := f.IOStreams.Out
			formattedFuction, err := format(cmd, function)
			if err != nil {
				return err
			}

			if cmd.Flags().Changed("out") {
				err := cmdutil.WriteDetailsToScreenOrFile(formattedFuction, true, opts.OutPath, out)
				if err != nil {
					return fmt.Errorf("failed to write to file: %w", err)
				}
				fmt.Fprintf(out, "File successfuly written to: %s\n", filepath.Clean(opts.OutPath))
			} else {
				err := cmdutil.WriteDetailsToScreenOrFile(formattedFuction, false, "", out)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().Bool("with-code", false, "Display the Edge Function code, disabled dy default")
	cmd.Flags().StringVar(&opts.OutPath, "out", "", "Enables exporting result to defined file path")
	cmd.Flags().StringVar(&opts.Format, "format", "", "You can change the format of the results by sending json value to this flag")

	return cmd
}

func serializeToJson(data map[string]interface{}) string {
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
