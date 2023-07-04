package create

import (
	"context"
	"fmt"
	"os"

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
	Secret bool
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
	$ azioncli variables create --key "Content-Type" --value "string" --secret false
	$ azioncli variables create --in "create.json"
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
				// flags requireds
				if !cmd.Flags().Changed("key") || !cmd.Flags().Changed("value") || !cmd.Flags().Changed("secret") {
					return msg.ErrorMandatoryCreateFlags
				}

				request.SetKey(fields.Key)
				request.SetValue(fields.Value)
				request.SetSecret(fields.Secret)
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
	flags.BoolVar(&fields.Secret, "secret", true, msg.CreateFlagSecret)
	flags.StringVar(&fields.Path, "in", "", msg.CreateFlagIn)
	flags.BoolP("help", "h", false, msg.CreateHelpFlag)
	return cmd
}