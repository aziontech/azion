package variable

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/variables"

	api "github.com/aziontech/azion-cli/pkg/api/variables"
	"github.com/aziontech/azion-cli/pkg/logger"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
		Use:           msg.Usage,
		Short:         msg.CreateShortDescription,
		Long:          msg.CreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create variable --key "Content-Type" --value "string" --secret false
		$ azion create variable --in "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			request := api.CreateRequest{}

			if cmd.Flags().Changed("in") {
				err := utils.FlagINUnmarshalFileJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
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

				err := createRequestFromFlags(cmd, fields, &request)
				if err != nil {
					return err
				}
			}

			response, err := client.Create(context.Background(), request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateItem.Error(), err)
			}

			logger.LogSuccess(f.IOStreams.Out, fmt.Sprintf(msg.CreateOutputSuccess, response.GetUuid()))
			return nil
		},
	}

	flags := cmd.Flags()
	addFlags(flags, fields)
	return cmd
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.StringVar(&fields.Key, "key", "", msg.CreateFlagKey)
	flags.StringVar(&fields.Value, "value", "", msg.CreateFlagValue)
	flags.StringVar(&fields.Secret, "secret", "", msg.CreateFlagSecret)
	flags.StringVar(&fields.Path, "in", "", msg.CreateFlagIn)
	flags.BoolP("help", "h", false, msg.CreateHelpFlag)
}

func createRequestFromFlags(cmd *cobra.Command, fields *Fields, request *api.CreateRequest) error {
	request.SetKey(fields.Key)
	request.SetValue(fields.Value)

	if cmd.Flags().Changed("secret") {
		secret, err := strconv.ParseBool(fields.Secret)
		if err != nil {
			return fmt.Errorf("%w: %q", msg.ErrorSecretFlag, fields.Secret)
		}
		request.SetSecret(secret)
	}

	return nil
}
