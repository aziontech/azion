package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "describe <edge_function_id> [flags]",
		Short:         "Describe a given Edge Function",
		Long:          "Describe a given Edge Function",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions describe 1337 [--with-code]
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

			fmt.Fprintf(out, "ID: %d\n", uint64(function.GetId()))
			fmt.Fprintf(out, "Name: %s\n", function.GetName())
			fmt.Fprintf(out, "Language: %s\n", function.GetLanguage())
			fmt.Fprintf(out, "Reference Count: %d\n", uint64(function.GetReferenceCount()))
			fmt.Fprintf(out, "Modified at: %s\n", function.GetModified())
			fmt.Fprintf(out, "Initiator Type: %s\n", function.GetInitiatorType())
			fmt.Fprintf(out, "Last Editor: %s\n", function.GetLastEditor())
			fmt.Fprintf(out, "Function to run: %s\n", function.GetFunctionToRun())
			fmt.Fprintf(out, "JSON Args: %s\n", serializeToJson(function.GetJsonArgs())) // Show serialized JSON

			if cmd.Flags().Changed("with-code") {
				fmt.Fprintf(out, "Code:\n%s\n", function.GetCode())
			}

			return nil
		},
	}

	cmd.Flags().Bool("with-code", false, "Display the Edge Function code, disabled dy default")

	return cmd
}

func serializeToJson(data map[string]interface{}) string {
	// ignoring errors on purpose
	serialized, _ := json.Marshal(data)
	return string(serialized)
}
