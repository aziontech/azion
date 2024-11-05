package variables

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/variables"
	api "github.com/aziontech/azion-cli/pkg/api/variables"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var (
	words = []string{"PASSWORD", "PWD", "SECRET", "HASH", "ENCRYPTED", "PASSCODE", "AUTH", "TOKEN", "SECRET"}
)

type Fields struct {
	Key      string
	Value    string
	Secret   string
	FileJSON string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.CreateShortDescription,
		Long:          msg.CreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion create variables --key "Content-Type" --value "string" --secret false
		$ azion create variables --file "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			request := api.Request{}

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.FileJSON, &request)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
			} else {
				err := createRequestFromFlags(cmd, fields, &request)
				if err != nil {
					return err
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.Create(context.Background(), request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateItem.Error(), err)
			}
			creatOut := output.GeneralOutput{
				Msg: fmt.Sprintf(msg.CreateOutputSuccess, response.GetUuid()),
				Out: f.IOStreams.Out,
			}
			return output.Print(&creatOut)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.Key, "key", "", msg.CreateFlagKey)
	flags.StringVar(&fields.Value, "value", "", msg.CreateFlagValue)
	flags.StringVar(&fields.Secret, "secret", "", msg.CreateFlagSecret)
	flags.StringVar(&fields.FileJSON, "file", "", msg.CreateFlagFileJSON)
	flags.BoolP("help", "h", false, msg.CreateHelpFlag)
	return cmd
}

func createRequestFromFlags(cmd *cobra.Command, fields *Fields, request *api.Request) error {
	if !cmd.Flags().Changed("key") {
		answers, err := utils.AskInput(msg.AskKey)

		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}

		fields.Key = answers
	}

	if !cmd.Flags().Changed("value") {
		answers, err := utils.AskInput(msg.AskValue)

		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}

		fields.Value = answers
	}

	if cmd.Flags().Changed("secret") {
		secret, err := strconv.ParseBool(fields.Secret)
		if err != nil {
			return fmt.Errorf("%w: %q", msg.ErrorSecretFlag, fields.Secret)
		}
		request.SetSecret(secret)
	} else {
		if utils.ContainSubstring(fields.Key, words) {
			request.SetSecret(true)
		}
	}

	request.SetKey(fields.Key)
	request.SetValue(fields.Value)

	return nil
}
