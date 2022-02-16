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
	"github.com/aziontech/azion-cli/pkg/cmdutil"
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
						return fmt.Errorf("error while opening file %s", fields.InPath)
					}
				}

				err = cmdutil.UnmarshallJsonFromReader(file, &request)
				if err != nil {
					return fmt.Errorf("error while unmarshalling the file")
				}
			} else {
				if !cmd.Flags().Changed("active") || !cmd.Flags().Changed("code") || !cmd.Flags().Changed("name") {
					return fmt.Errorf("Error: flags --active, --code and --name are mandatory when flag --in is not sent")
				}
				isActive, err := strconv.ParseBool(fields.Active)
				if err != nil {
					return fmt.Errorf("invalid --active flag: %s", fields.Active)
				}
				request.SetActive(isActive)

				code, err := ioutil.ReadFile(fields.Code)
				if err != nil {
					return fmt.Errorf("failed to read code file: %w", err)
				}
				request.SetCode(string(code))

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

				request.SetName(fields.Name)

				if cmd.Flags().Changed("initiator-type") {
					request.SetInitiatorType(fields.InitiatorType)
				}
			}

			httpClient, err := f.HttpClient()
			if err != nil {
				return fmt.Errorf("failed to get http client: %w", err)
			}

			client := api.NewClient(httpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Create(ctx, request)

			if err != nil {
				return fmt.Errorf("failed to create edge function: %w", err)
			}

			fmt.Fprintf(f.IOStreams.Out, "Created Edge Function with ID %d\n", response.GetId())

			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&fields.Name, "name", "", "Name of your Edge Function (Mandatory)")
	flags.StringVar(&fields.Code, "code", "", "Path to the file containing your Edge Function code (Mandatory)")
	flags.StringVar(&fields.Active, "active", "", "Whether or not your Edge Function should be active: <true|false> (Mandatory)")
	flags.StringVar(&fields.Args, "args", "", "Path to the file containing the JSON arguments of your Edge Function")
	flags.StringVar(&fields.InPath, "in", "", "Use provided filepath to create an Edge Function")

	return cmd
}
