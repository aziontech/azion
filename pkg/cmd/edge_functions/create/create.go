package create

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
		Use:           "create [flags]",
		Short:         "Create a new Edge Function",
		Long:          "Create a new Edge Function",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions create --name myjsfunc --code ./mycode/function.js --active true
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
					return errmsg.ErrorMandatoryCreateFlags
				}
				isActive, err := strconv.ParseBool(fields.Active)
				if err != nil {
					return fmt.Errorf("%w: %s", errmsg.ErrorActiveFlag, fields.Active)
				}
				request.SetActive(isActive)

				code, err := ioutil.ReadFile(fields.Code)
				if err != nil {
					return fmt.Errorf("%s: %w", errmsg.ErrorCodeFlag, err)
				}
				request.SetCode(string(code))

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

				request.SetName(fields.Name)

				if cmd.Flags().Changed("initiator-type") {
					request.SetInitiatorType(fields.InitiatorType)
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Create(ctx, request)

			if err != nil {
				return fmt.Errorf("%s: %w", errmsg.ErrorCreateFunction, err)
			}

			fmt.Fprintf(f.IOStreams.Out, "Created Edge Function with ID %d\n", response.GetId())

			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&fields.Name, "name", "", "Your Edge Function's name (Mandatory if --in is not sent)")
	flags.StringVar(&fields.Code, "code", "", "Path to the file containing your Edge Function's code (Mandatory if --in is not sent)")
	flags.StringVar(&fields.Active, "active", "", "Whether your Edge Function should be active or not: <true|false> (Mandatory if --in is not sent)")
	flags.StringVar(&fields.Args, "args", "", "Path to the file containing your Edge Function's JSON arguments")
	flags.StringVar(&fields.InPath, "in", "", "Uses provided file path to create an Edge Function. You can use - for reading from stdin")

	return cmd
}
