package create

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
		Use:           edgefunctions.EdgeFunctionCreateUsage,
		Short:         edgefunctions.EdgeFunctionCreateShortDescription,
		Long:          edgefunctions.EdgeFunctionCreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion edge_functions create --name myjsfunc --code ./mycode/function.js --active false
        $ azion edge_functions create --name withargs --code ./mycode/function.js --args ./args.json --active true
        $ azion edge_functions create --in "create.json"
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
					return edgefunctions.ErrorMandatoryCreateFlags
				}
				isActive, err := strconv.ParseBool(fields.Active)
				if err != nil {
					return fmt.Errorf("%w: %s", edgefunctions.ErrorActiveFlag, fields.Active)
				}
				request.SetActive(isActive)

				code, err := os.ReadFile(fields.Code)
				if err != nil {
					return fmt.Errorf("%s: %w", edgefunctions.ErrorCodeFlag, err)
				}
				request.SetCode(string(code))

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

				request.SetName(fields.Name)

			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Create(ctx, request)

			if err != nil {
				return fmt.Errorf(edgefunctions.ErrorCreateFunction.Error(), err)
			}

			fmt.Fprintf(f.IOStreams.Out, edgefunctions.EdgeFunctionCreateOutputSuccess, response.GetId())

			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&fields.Name, "name", "", edgefunctions.EdgeFunctionCreateFlagName)
	flags.StringVar(&fields.Code, "code", "", edgefunctions.EdgeFunctionCreateFlagCode)
	flags.StringVar(&fields.Active, "active", "", edgefunctions.EdgeFunctionCreateFlagActive)
	flags.StringVar(&fields.Args, "args", "", edgefunctions.EdgeFunctionCreateFlagArgs)
	flags.StringVar(&fields.InPath, "in", "", edgefunctions.EdgeFunctionCreateFlagIn)
	flags.BoolP("help", "h", false, edgefunctions.EdgeFunctionCreateHelpFlag)

	return cmd
}
