package update

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/edge_functions"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	Id            int64
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
		Use:           edgefunctions.EdgeFunctionUpdateUsage,
		Short:         edgefunctions.EdgeFunctionUpdateShortDescription,
		Long:          edgefunctions.EdgeFunctionUpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion edge_functions update --function-id 1234 --name 'Hello'
		$ azion edge_functions update -f 4185 --code ./mycode/function.js --args ./mycode/myargs.json
		$ azion edge_functions update -f 9123 --active true
		$ azion edge_functions update -f 9123 --active false
		$ azion edge_functions update --in "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			// either function-id or in path should be passed
			if !cmd.Flags().Changed("function-id") && !cmd.Flags().Changed("in") {
				return edgefunctions.ErrorMissingFunctionIdArgument
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

				request.Id = fields.Id

				if cmd.Flags().Changed("active") {
					active, err := strconv.ParseBool(fields.Active)
					if err != nil {
						return fmt.Errorf("%w: %q", edgefunctions.ErrorActiveFlag, fields.Active)
					}
					request.SetActive(active)
				}

				if cmd.Flags().Changed("code") {
					code, err := os.ReadFile(fields.Code)
					if err != nil {
						return fmt.Errorf("%s: %w", edgefunctions.ErrorCodeFlag, err)
					}
					request.SetCode(string(code))
				}

				if cmd.Flags().Changed("args") {
					marshalledArgs, err := os.ReadFile(fields.Args)
					if err != nil {
						return fmt.Errorf("%s: %w", edgefunctions.ErrorArgsFlag, err)
					}
					args := make(map[string]interface{})
					if err := json.Unmarshal(marshalledArgs, &args); err != nil {
						return fmt.Errorf("%s: %w", edgefunctions.ErrorParseArgs, err)
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
				return fmt.Errorf(edgefunctions.ErrorUpdateFunction.Error(), err)
			}

			fmt.Fprintf(f.IOStreams.Out, "Updated Edge Function with ID %v\n", response.GetId())

			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&fields.Id, "function-id", "f", 0, edgefunctions.EdgeFunctionFlagId)
	flags.StringVar(&fields.Name, "name", "", edgefunctions.EdgeFunctionUpdateFlagName)
	flags.StringVar(&fields.Code, "code", "", edgefunctions.EdgeFunctionUpdateFlagCode)
	flags.StringVar(&fields.Args, "args", "", edgefunctions.EdgeFunctionUpdateFlagArgs)
	flags.StringVar(&fields.Active, "active", "", edgefunctions.EdgeFunctionUpdateFlagActive)
	flags.StringVar(&fields.InPath, "in", "", edgefunctions.EdgeFunctionUpdateFlagIn)
	flags.BoolP("help", "h", false, edgefunctions.EdgeFunctionUpdateHelpFlag)

	return cmd
}
