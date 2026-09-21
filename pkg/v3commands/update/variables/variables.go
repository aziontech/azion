package variables

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/spf13/pflag"
	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"

	msg "github.com/aziontech/azion-cli/messages/variables"
	api "github.com/aziontech/azion-cli/pkg/api/variables"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	ID       string
	Key      string
	Value    string
	Secret   string
	FileJSON string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.UpdateShortDescription,
		Long:          msg.UpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion update variables --variable-id 7a187044-4a00-4a4a-93ed-d230900421f3 --key 'Content-Type' --value 'json' --secret false
		$ azion update variables --file variables.json
		$ Example JSON: {
		    "uuid": "32e8ffca-4021-49a4-971f-330935566af4",
		    "key": "Content-Type",
		    "value": "json",
		    "secret": false,
		    "last_editor": "hunter@hunter.com",
		    "created_at": "2023-06-13T13:17:13.145625Z",
		    "updated_at": "2023-06-13T13:17:13.145666Z"
		}`),
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

			response, err := client.Update(context.Background(), &request)
			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateVariable.Error(), err)
			}

			updateOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.UpdateOutputSuccess, response.GetUuid()),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&updateOut)
		},
	}

	flags := cmd.Flags()
	addFlags(flags, fields)

	return cmd
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.StringVar(&fields.ID, "variable-id", "0", msg.FlagVariableID)
	flags.StringVar(&fields.Key, "key", "", msg.UpdateFlagKey)
	flags.StringVar(&fields.Value, "value", "", msg.UpdateFlagValue)
	flags.StringVar(&fields.Secret, "secret", "", msg.UpdateFlagSecret)
	flags.StringVar(&fields.FileJSON, "file", "", msg.UpdateFlagIn)
	flags.BoolP("help", "h", false, msg.UpdateHelpFlag)
}

func createRequestFromFlags(cmd *cobra.Command, fields *Fields, request *api.Request) error {
	if !cmd.Flags().Changed("variable-id") {
		answers, err := utils.AskInput(msg.AskVariableID)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		fields.ID = answers
	}

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

	if !cmd.Flags().Changed("secret") {
		answers, err := utils.AskInput(msg.AskSecret)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		fields.Secret = answers
	}

	secret, err := strconv.ParseBool(fields.Secret)
	if err != nil {
		return fmt.Errorf("%w: %q", msg.ErrorSecretFlag, fields.Secret)
	}
	request.SetSecret(secret)
	request.SetKey(fields.Key)
	request.SetValue(fields.Value)
	request.Uuid = fields.ID

	return nil
}
