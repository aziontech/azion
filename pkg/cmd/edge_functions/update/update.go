package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
	errmsg "github.com/aziontech/azion-cli/pkg/cmd/edge_functions/error_messages"
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
	InPath        string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           "update <edge_function_id> [flags]",
		Short:         "Updates an Edge Function",
		Long:          "Updates an Edge Function based on the id given",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions update 1234 --name 'Hello'
        $ azioncli edge_functions update 4185 --code ./mycode/function.js --args ./mycode/myargs.json
        $ azioncli edge_functions update 9123 --active false
        $ azioncli edge_functions update --in "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			// either id parameter or in path should be passed
			if len(args) < 1 && !cmd.Flags().Changed("in") {
				return errmsg.ErrorMissingArgumentUpdate
			}

			request := api.UpdateRequest{}

			if cmd.Flags().Changed("in") {
				var (
					file *os.File
					err  error
				)
				if fields.InPath == "-" {
					file = os.Stdin
				} else {
					file, err = os.Open(fields.InPath)
					if err != nil {
						return fmt.Errorf("%w: %s", utils.ErrorOpeningFile, fields.InPath)
					}
				}
				err = cmdutil.UnmarshallJsonFromReader(file, &request)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
			} else {
				ids, err := utils.ConvertIdsToInt(args[0])
				if err != nil {
					return utils.ErrorConvertingIdArgumentToInt
				}

				request.Id = ids[0]

				if cmd.Flags().Changed("active") {
					active, err := strconv.ParseBool(fields.Active)
					if err != nil {
						return fmt.Errorf("%w: %q", errmsg.ErrorActiveFlag, fields.Active)
					}
					request.SetActive(active)
				}

				if cmd.Flags().Changed("code") {
					code, err := ioutil.ReadFile(fields.Code)
					if err != nil {
						return fmt.Errorf("%s: %w", errmsg.ErrorCodeFlag, err)
					}
					request.SetCode(string(code))
				}

				if cmd.Flags().Changed("args") {
					marshalledArgs, err := ioutil.ReadFile(fields.Args)
					if err != nil {
						return fmt.Errorf("%s: %w", errmsg.ErrorArgsFlag, err)
					}
					args := make(map[string]interface{})
					if err := json.Unmarshal(marshalledArgs, &args); err != nil {
						return fmt.Errorf("%s: %w", errmsg.ErrorParseArgs, err)
					}
					request.SetJsonArgs(args)
				}

				if cmd.Flags().Changed("name") {
					request.SetName(fields.Name)
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Update(ctx, &request)

			if err != nil {
				return fmt.Errorf("%w: %s", errmsg.ErrorUpdateFunction, err)
			}

			fmt.Fprintf(f.IOStreams.Out, "Updated Edge Function with ID %d\n", response.GetId())

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.Name, "name", "", "Your Edge Function's name")
	flags.StringVar(&fields.Code, "code", "", "Path to the file containing your Edge Function's code")
	flags.StringVar(&fields.Args, "args", "", "Path to the file containing your Edge Function's JSON arguments")
	flags.StringVar(&fields.Active, "active", "", "Whether your Edge Function should be active or not: <true|false>")
	flags.StringVar(&fields.InPath, "in", "", "Uses provided file path to update the fields. You can use - for reading from stdin")

	return cmd
}
