package create

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
	Key    string
	Value  string
	Secret string
	Path   string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.CreateUsage,
		Short:         msg.CreateShortDescription,
		Long:          msg.CreateLongDescription,
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
					return msg.ErrorMandatoryCreateFlags
				}

				if cmd.Flags().Changed("secret") {
					secret, err := strconv.ParseBool(fields.Secret)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorSecretFlag, fields.Secret)
					}
					request.SetSecret(secret)
				}

				request.SetKey(fields.Key)
				request.SetValue(fields.Value)
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.Create(context.Background(), request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateItem.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, msg.CreateOutputSuccess, response.GetUuid())
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.Key, "key", "", msg.CreateFlagKey)
	flags.StringVar(&fields.Value, "value", "", msg.CreateFlagValue)
	flags.StringVar(&fields.Secret, "secret", "", msg.CreateFlagSecret)
	flags.StringVar(&fields.Path, "in", "", msg.CreateFlagIn)
	flags.BoolP("help", "h", false, msg.CreateHelpFlag)
	return cmd
}
