package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_functions"
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
		Use:           msg.EdgeFunctionUpdateUsage,
		Short:         msg.EdgeFunctionUpdateShortDescription,
		Long:          msg.EdgeFunctionUpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions update --function-id 1234 --name 'Hello'
        $ azioncli edge_functions update -f 4185 --code ./mycode/function.js --args ./mycode/myargs.json
        $ azioncli edge_functions update -f 9123 --active false
        $ azioncli edge_functions update --in "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			// either function-id or in path should be passed
			if !cmd.Flags().Changed("function-id") && !cmd.Flags().Changed("in") {
				return msg.ErrorMissingArgumentUpdate
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
						return fmt.Errorf("%w: %q", msg.ErrorActiveFlag, fields.Active)
					}
					request.SetActive(active)
				}

				if cmd.Flags().Changed("code") {
					code, err := ioutil.ReadFile(fields.Code)
					if err != nil {
						return fmt.Errorf("%s: %w", msg.ErrorCodeFlag, err)
					}
					request.SetCode(string(code))
				}

				if cmd.Flags().Changed("args") {
					marshalledArgs, err := ioutil.ReadFile(fields.Args)
					if err != nil {
						return fmt.Errorf("%s: %w", msg.ErrorArgsFlag, err)
					}
					args := make(map[string]interface{})
					if err := json.Unmarshal(marshalledArgs, &args); err != nil {
						return fmt.Errorf("%s: %w", msg.ErrorParseArgs, err)
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
				return fmt.Errorf("%w: %s", msg.ErrorUpdateFunction, err)
			}

			fmt.Fprintf(f.IOStreams.Out, "Updated Edge Function with ID %d\n", response.GetId())

			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&fields.Id, "function-id", "f", 0, msg.EdgeFunctionFlagId)
	flags.StringVar(&fields.Name, "name", "", msg.EdgeFunctionUpdateFlagName)
	flags.StringVar(&fields.Code, "code", "", msg.EdgeFunctionUpdateFlagCode)
	flags.StringVar(&fields.Args, "args", "", msg.EdgeFunctionUpdateFlagArgs)
	flags.StringVar(&fields.Active, "active", "", msg.EdgeFunctionUpdateFlagActive)
	flags.StringVar(&fields.InPath, "in", "", msg.EdgeFunctionUpdateFlagIn)

	return cmd
}
