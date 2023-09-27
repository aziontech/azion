package update

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/variables"
	api "github.com/aziontech/azion-cli/pkg/api/variables"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	Id     string
	Key    string
	Value  string
	Secret string
	InPath string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.UpdateUsage,
		Short:         msg.UpdateShortDescription,
		Long:          msg.UpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion variables update --variable-id 7a187044-4a00-4a4a-93ed-d230900421f3 --key 'Content-Type' --value 'json' --secret false
		$ azion variables update -v 7a187044-4a00-4a4a-93ed-d230900421f3 --key 'Content-Type' --value 'json' --secret false
		$ azion variables update --in variables.json
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			// either function-id or in path should be passed
			if (!cmd.Flags().Changed("variable-id") || !cmd.Flags().Changed("key") || !cmd.Flags().Changed("value") || !cmd.Flags().Changed("secret")) && !cmd.Flags().Changed("in") {
				return msg.ErrorMissingVariableIdArgument
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
				secret, err := strconv.ParseBool(fields.Secret)
				if err != nil {
					return fmt.Errorf("%w: %q", msg.ErrorSecretFlag, fields.Secret)
				}
				request.SetSecret(secret)
				request.SetKey(fields.Key)
				request.SetValue(fields.Value)
				request.Uuid = fields.Id

			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Update(ctx, &request)

			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateVariable.Error(), err)
			}

			fmt.Fprintf(f.IOStreams.Out, "Updated Variable with ID %s\n", response.GetUuid())

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&fields.Id, "variable-id", "v", "0", msg.FlagVariableID)
	flags.StringVar(&fields.Key, "key", "", msg.UpdateFlagKey)
	flags.StringVar(&fields.Value, "value", "", msg.UpdateFlagValue)
	flags.StringVar(&fields.Secret, "secret", "", msg.UpdateFlagSecret)
	flags.StringVar(&fields.InPath, "in", "", msg.UpdateFlagIn)
	flags.BoolP("help", "h", false, msg.UpdateHelpFlag)

	return cmd
}
