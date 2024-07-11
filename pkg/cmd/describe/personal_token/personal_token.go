package personal_token

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/personal_token"

	api "github.com/aziontech/azion-cli/pkg/api/personal_token"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var personalTokenID string

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.DescribeOptions{}
	cmd := &cobra.Command{
		Use:           msg.USAGE,
		Short:         msg.SHORT_DESCRIPTION_DESCRIBE,
		Long:          msg.LONG_DESCRIPTION_DESCRIBE,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion describe personal-token --help
		$ azion describe personal-token --personal-token-id 12345  
		$ azion describe personal-token --personal-token-id 12345 --format json
		$ azion describe personal-token --personal-token-id 12345 --out "./tmp/test.json"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("id") {
				answers, err := utils.AskInput(msg.ASK_PERSONAL_TOKEN_ID)
				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				personalTokenID = answers
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			ctx := context.Background()
			personalToken, err := client.Get(ctx, personalTokenID)
			if err != nil {
				return fmt.Errorf(msg.ERROR_GET_PERSONAL_TOKEN, err)
			}

			fields := make(map[string]string, 0)
			fields["Uuid"] = "Uuid"
			fields["Name"] = "Name"
			fields["Created"] = "Created"
			fields["ExpiresAt"] = "Expires At"
			fields["Description"] = "Description"

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Msg:   filepath.Clean(opts.OutPath),
					Flags: f.Flags,
					Out:   f.IOStreams.Out,
				},
				Fields: fields,
				Values: personalToken,
			}
			return output.Print(&describeOut)
		},
	}

	cmd.Flags().StringVar(&personalTokenID, "id", "", msg.FLAG_PERSONAL_TOKEN_ID)
	cmd.Flags().BoolP("help", "h", false, msg.FLAG_HELP_DESCRIBE)

	return cmd
}
