package personaltoken

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/create/personal_token"
	"os"
	"time"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

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
		Use:           personaltoken.CreateUsage,
		Short:         personaltoken.CreateShortDescription,
		Long:          personaltoken.CreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create personal-token --name "ranking of kings" --expiration "9m" 
        $ azion create personal-token --name "sakura" --expiration "9m" 
        $ azion create personal-token --name "strawhat" --expiration "9m" --description "gear five"
        $ azion create personal-token --in "create.json"
        $ "create.json" example: 
        {   
            "name": "One day token",
            "expires_at": "9m",
            "description": "example"
        }
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

				err = cmdutil.UnmarshallJsonFromReader(file, &fields)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}

			} else {
				if !cmd.Flags().Changed("name") {
					answer, err := utils.AskInput(personaltoken.AskInputName)
					if err != nil {
						return err
					}

					fields.Name = answer
				}
				if !cmd.Flags().Changed("expiration") {
					answer, err := utils.AskInput(personaltoken.AskInputExpiration)
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
				return fmt.Errorf(personaltoken.ErrorCreate.Error(), err)
			}

			fmt.Fprintf(f.IOStreams.Out, personaltoken.CreateOutputSuccess, response.GetKey())

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.Name, "name", "", personaltoken.CreateFlagName)
	flags.StringVar(&fields.ExpiresAt, "expiration", "", personaltoken.CreateFlagExpiresAt)
	flags.StringVar(&fields.Description, "description", "", personaltoken.CreateFlagDescription)
	flags.StringVar(&fields.Path, "in", "", personaltoken.CreateFlagIn)
	flags.BoolP("help", "h", false, personaltoken.CreateHelpFlag)

	return cmd
}
