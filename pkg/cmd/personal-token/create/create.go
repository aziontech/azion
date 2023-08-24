package create

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/MakeNowJust/heredoc"

	msg "github.com/aziontech/azion-cli/messages/personal-token"
	api "github.com/aziontech/azion-cli/pkg/api/personal_token"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	Name        string
	ExpiresAt   string
	Description string
	Path        string
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
        $ azion personal_token create --name "ranking of kings" --expiration "9m" 
        $ azion personal_token create --name "sakura" --expiration "9m" 
        $ azion personal_token create --name "luffy biruta" --expiration "9m" --description "gear five"
        $ azion personal_token create --in "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.Request{}

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

				if !cmd.Flags().Changed("name") || !cmd.Flags().Changed("expiration") {
					return msg.ErrorMandatoryCreateFlags
				}

				request.SetName(fields.Name)
				date, err := ParseExpirationDate(time.Now(), fields.ExpiresAt)
				if err != nil {
					return err
				}
				request.SetExpiresAt(date)
				request.SetDescription(fields.Description)
			}

			response, err := api.NewClient(
				f.HttpClient,
				f.Config.GetString("api_url"),
				f.Config.GetString("token"),
			).Create(context.Background(), &request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreate.Error(), err)
			}

			fmt.Fprintf(f.IOStreams.Out, msg.CreateOutputSuccess, response.GetUuid())

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.Name, "name", "", msg.CreateFlagName)
	flags.StringVar(&fields.ExpiresAt, "expiration", "", msg.CreateFlagExpiresAt)
	flags.StringVar(&fields.Description, "description", "", msg.CreateFlagDescription)
	flags.StringVar(&fields.Path, "in", "", msg.CreateFlagIn)
	flags.BoolP("help", "h", false, msg.CreateHelpFlag)

	return cmd
}
