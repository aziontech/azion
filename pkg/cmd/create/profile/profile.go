package profile

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/profile"
	api "github.com/aziontech/azion-cli/pkg/api/personal_token"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Fields struct {
	File  string
	Name  string
	Token string
}

var confirmFn = utils.Confirm

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.UsageCreate,
		Short:         msg.CreateShortDescription,
		Long:          msg.CreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create profile --file "create.toml"
		$ azion create profile"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			settings := &token.Settings{}
			if !cmd.Flags().Changed("name") {
				profileName, err := utils.AskInput(msg.FieldProfileName)
				if err != nil {
					return utils.ErrorParseResponse
				}
				fields.Name = profileName
			}

			if cmd.Flags().Changed("file") {
				fileData, err := os.ReadFile(fields.File)
				if err != nil {
					return fmt.Errorf(msg.ErrorReadFile.Error(), err)
				}

				err = toml.Unmarshal(fileData, settings)
				if err != nil {
					return fmt.Errorf(msg.ErrorUnmarshalFile.Error(), err)
				}
			} else {
				authorize := confirmFn(f.GlobalFlagAll, msg.QuestionCollectMetrics, true)
				if authorize {
					settings.AuthorizeMetricsCollection = 1
				} else {
					settings.AuthorizeMetricsCollection = 2
				}

				if !cmd.Flags().Changed("personal-token") {
					authorizeToken := confirmFn(f.GlobalFlagAll, msg.QuestionToken, true)
					if authorizeToken {
						provideToken := confirmFn(f.GlobalFlagAll, msg.QuestionProvideToken, true)
						if provideToken {
							tokenString, err := utils.AskInput(msg.FieldToken)
							if err != nil {
								return utils.ErrorParseResponse
							}
							fields.Token = tokenString
						} else {
							// Create new token via API
							request := api.Request{}
							request.SetName(fields.Name)
							request.SetExpiresAt(time.Now().Add(8760 * time.Hour))
							response, err := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token")).Create(context.Background(), &request)
							if err != nil {
								return fmt.Errorf(msg.ErrorCreate.Error(), err)
							}
							fields.Token = response.GetKey()
							settings.UUID = response.GetUuid()
						}
					}
				}

				// Validate token if one was provided or created
				if fields.Token != "" {
					t := token.New(&token.Config{
						Client: f.HttpClient,
						Out:    f.IOStreams.Out,
					})
					valid, user, err := t.Validate(&fields.Token)
					if err != nil {
						return err
					}

					if !valid {
						return utils.ErrorInvalidToken
					}
					settings.Token = fields.Token
					settings.ClientId = user.Results.ClientID
					settings.Email = user.Results.Email
				}
			}

			if err := token.WriteSettings(*settings, fields.Name); err != nil {
				return err
			}

			profileOut := output.GeneralOutput{
				Msg: fmt.Sprintf(msg.CreateOutputSuccess, fields.Name),
				Out: f.IOStreams.Out,
			}
			return output.Print(&profileOut)
		},
	}

	flags := cmd.Flags()
	addFlags(flags, fields)

	return cmd
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.BoolP("help", "h", false, msg.CreateFlagHelp)
	flags.StringVar(&fields.File, "file", "", msg.FlagFile)
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.Token, "personal-token", "", msg.FlagToken)
}
