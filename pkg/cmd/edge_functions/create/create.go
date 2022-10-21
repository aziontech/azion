package create

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
		Use:           msg.EdgeFunctionCreateUsage,
		Short:         msg.EdgeFunctionCreateShortDescription,
		Long:          msg.EdgeFunctionCreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions create --name myjsfunc --code ./mycode/function.js --active false
        $ azioncli edge_functions create --name withargs --code ./mycode/function.js --args ./args.json --active true
        $ azioncli edge_functions create --in "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.NewCreateRequest()

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
						return fmt.Errorf("%s %s", utils.ErrorOpeningFile, fields.InPath)
					}
				}

				err = cmdutil.UnmarshallJsonFromReader(file, &request)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
			} else {
				if !cmd.Flags().Changed("active") || !cmd.Flags().Changed("code") || !cmd.Flags().Changed("name") {
					return msg.ErrorMandatoryCreateFlags
				}
				isActive, err := strconv.ParseBool(fields.Active)
				if err != nil {
					return fmt.Errorf("%w: %s", msg.ErrorActiveFlag, fields.Active)
				}
				request.SetActive(isActive)

				code, err := ioutil.ReadFile(fields.Code)
				if err != nil {
					return fmt.Errorf("%s: %w", msg.ErrorCodeFlag, err)
				}
				request.SetCode(string(code))

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

				request.SetName(fields.Name)

			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Create(ctx, request)

			if err != nil {
				return fmt.Errorf(msg.ErrorCreateFunction.Error(), err)
			}

			fmt.Fprintf(f.IOStreams.Out, msg.EdgeFunctionCreateOutputSuccess, response.GetId())

			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&fields.Name, "name", "", msg.EdgeFunctionCreateFlagName)
	flags.StringVar(&fields.Code, "code", "", msg.EdgeFunctionCreateFlagCode)
	flags.StringVar(&fields.Active, "active", "", msg.EdgeFunctionCreateFlagActive)
	flags.StringVar(&fields.Args, "args", "", msg.EdgeFunctionCreateFlagArgs)
	flags.StringVar(&fields.InPath, "in", "", msg.EdgeFunctionCreateFlagIn)
	flags.BoolP("help", "h", false, msg.EdgeFunctionCreateHelpFlag)

	return cmd
}
