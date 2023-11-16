package personaltoken

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/create/personal_token"
	api "github.com/aziontech/azion-cli/pkg/api/personal_token"
	"github.com/aziontech/azion-cli/pkg/logger"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	Name        string `json:"name,omitempty"`
	ExpiresAt   string `json:"expires_at,omitempty"`
	Description string `json:"description,omitempty"`
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
        $ azion create personal-token --name "ranking of kings" --expiration "9m" 
        $ azion create personal-token --name "sakura" --expiration "9m" 
        $ azion create personal-token --name "strawhat" --expiration "9m" --description "gear five"
        $ azion create personal-token --file "create.json"
        $ "create.json" example: 
        {   
            "name": "One day token",
            "expires_at": "9m",
            "description": "example"
        }
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.Request{}

			if cmd.Flags().Changed("file") {
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

				err = cmdutil.UnmarshallJsonFromReader(file, &fields)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}

			} else {
				if !cmd.Flags().Changed("name") {
					answer, err := utils.AskInput(msg.AskInputName)
					if err != nil {
						return err
					}

					fields.Name = answer
				}
				if !cmd.Flags().Changed("expiration") {
					answer, err := utils.AskInput(msg.AskInputExpiration)
					if err != nil {
						return err
					}

					fields.ExpiresAt = answer
				}
			}

			date, err := ParseExpirationDate(time.Now(), fields.ExpiresAt)
			if err != nil {
				return err
			}

			request.SetName(fields.Name)
			request.SetExpiresAt(date)
			request.SetDescription(fields.Description)

			response, err := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token")).Create(context.Background(), &request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreate.Error(), err)
			}

			fmt.Fprintf(f.IOStreams.Out, msg.CreateOutputSuccess, response.GetKey())

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.Name, "name", "", msg.CreateFlagName)
	flags.StringVar(&fields.ExpiresAt, "expiration", "", msg.CreateFlagExpiresAt)
	flags.StringVar(&fields.Description, "description", "", msg.CreateFlagDescription)
	flags.StringVar(&fields.Path, "file", "", msg.CreateFlagFile)
	flags.BoolP("help", "h", false, msg.CreateHelpFlag)

	return cmd
}
