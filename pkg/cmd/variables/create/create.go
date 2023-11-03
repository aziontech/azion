package create

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/variables"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"

	api "github.com/aziontech/azion-cli/pkg/api/variables"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	Key    string
	Value  string
	Secret string
	Path   string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           variables.CreateUsage,
		Short:         variables.CreateShortDescription,
		Long:          variables.CreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
	$ azion variables create --key "Content-Type" --value "string" --secret false
	$ azion variables create --in "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			request := api.CreateRequest{}
			if cmd.Flags().Changed("in") {
				var (
					file *os.File
					err  error
				)
				if fields.Path == "-" {
					file = os.Stdin
				} else {
					file, err = os.Open(fields.Path)
					if err != nil {
						return fmt.Errorf("%w: %s", utils.ErrorOpeningFile, fields.Path)
					}
				}
				err = cmdutil.UnmarshallJsonFromReader(file, &request)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
			} else {
				// required flags
				if !cmd.Flags().Changed("key") || !cmd.Flags().Changed("value") {
					return variables.ErrorMandatoryCreateFlags
				}

				if cmd.Flags().Changed("secret") {
					secret, err := strconv.ParseBool(fields.Secret)
					if err != nil {
						return fmt.Errorf("%w: %q", variables.ErrorSecretFlag, fields.Secret)
					}
					request.SetSecret(secret)
				}

				request.SetKey(fields.Key)
				request.SetValue(fields.Value)
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.Create(context.Background(), request)
			if err != nil {
				return fmt.Errorf(variables.ErrorCreateItem.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, variables.CreateOutputSuccess, response.GetUuid())
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.Key, "key", "", variables.CreateFlagKey)
	flags.StringVar(&fields.Value, "value", "", variables.CreateFlagValue)
	flags.StringVar(&fields.Secret, "secret", "", variables.CreateFlagSecret)
	flags.StringVar(&fields.Path, "in", "", variables.CreateFlagIn)
	flags.BoolP("help", "h", false, variables.CreateHelpFlag)
	return cmd
}
