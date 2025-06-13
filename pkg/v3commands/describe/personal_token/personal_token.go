package personal_token

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"

	msg "github.com/aziontech/azion-cli/messages/personal_token"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	api "github.com/aziontech/azion-cli/pkg/v3api/personal_token"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/personal_tokens"
	"github.com/spf13/cobra"
)

type DescribeCmd struct {
	Io       *iostreams.IOStreams
	AskInput func(string) (string, error)
	Get      func(context.Context, string) (*sdk.PersonalTokenResponseGet, error)
}

var (
	personalTokenID string
)

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		Get: func(ctx context.Context, id string) (*sdk.PersonalTokenResponseGet, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			return client.Get(ctx, id)
		},
	}
}

func NewCobraCmd(describe *DescribeCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.DescribeOptions{}
	cobraCmd := &cobra.Command{
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
				answer, err := describe.AskInput(msg.ASK_PERSONAL_TOKEN_ID)
				if err != nil {
					return err
				}

				personalTokenID = answer
			}

			ctx := context.Background()
			personalToken, err := describe.Get(ctx, personalTokenID)
			if err != nil {
				return fmt.Errorf(msg.ERROR_GET_PERSONAL_TOKEN, err)
			}

			fields := map[string]string{
				"Uuid":        "Uuid",
				"Name":        "Name",
				"Created":     "Created",
				"ExpiresAt":   "Expires At",
				"Description": "Description",
			}

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Out:   f.IOStreams.Out,
					Msg:   filepath.Clean(opts.OutPath),
					Flags: f.Flags,
				},
				Fields: fields,
				Values: personalToken,
			}

			return output.Print(&describeOut)
		},
	}

	cobraCmd.Flags().StringVar(&personalTokenID, "id", "", msg.FLAG_PERSONAL_TOKEN_ID)
	cobraCmd.Flags().BoolP("help", "h", false, msg.FLAG_HELP_DESCRIBE)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}
