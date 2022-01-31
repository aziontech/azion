package update

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	Name          string
	Language      string
	Code          string
	Active        string
	InitiatorType string
	Args          string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           "update <edge_function_id> [flags]",
		Short:         "Update an Edge Function",
		Long:          "Update an Edge Function",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions update 1234 -–name 'Hello'
        $ azioncli edge_functions update 4185 -–code ./mycode/function.js -–args ./mycode/myargs.json
        $ azioncli edge_functions update 9123 -–active false
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing edge function id argument")
			}

			ids, err := utils.ConvertIdsToInt(args[0])
			if err != nil {
				return fmt.Errorf("invalid edge function id: %q", args[0])
			}

			request := api.NewUpdateRequest(ids[0])

			if cmd.Flags().Changed("active") {
				active, err := strconv.ParseBool(fields.Active)
				if err != nil {
					return fmt.Errorf("invalid --active flag: %q", fields.Active)
				}
				request.SetActive(active)
			}

			if cmd.Flags().Changed("code") {
				code, err := ioutil.ReadFile(fields.Code)
				if err != nil {
					return fmt.Errorf("failed to read code file: %w", err)
				}
				request.SetCode(string(code))
			}

			if cmd.Flags().Changed("args") {
				marshalledArgs, err := ioutil.ReadFile(fields.Args)
				if err != nil {
					return fmt.Errorf("failed to read args file: %w", err)
				}
				args := make(map[string]interface{})
				if err := json.Unmarshal(marshalledArgs, &args); err != nil {
					return fmt.Errorf("failed to parse json args: %w", err)
				}
				request.SetJsonArgs(args)
			}

			if cmd.Flags().Changed("name") {
				request.SetName(fields.Name)
			}

			httpClient, err := f.HttpClient()
			if err != nil {
				return fmt.Errorf("failed to get http client: %w", err)
			}

			client := api.NewClient(httpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Update(ctx, request)

			if err != nil {
				return fmt.Errorf("failed to create edge function: %w", err)
			}

			fmt.Fprintf(f.IOStreams.Out, "Updated Edge Function with ID %d\n", response.GetId())

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.Name, "name", "", "Name of your Edge Function.")
	flags.StringVar(&fields.Code, "code", "", "Path to the file containing your Edge Function code.")
	flags.StringVar(&fields.Args, "args", "", "Path to the file containing the JSON arguments of your Edge Function")
	flags.StringVar(&fields.Active, "active", "", "Whether or not your Edge Function should be active: <true|false>")

	return cmd
}
