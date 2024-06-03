package logout

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/logout"
	api "github.com/aziontech/azion-cli/pkg/api/personal_token"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   msg.USAGE,
		Short: msg.SHORTDESCRIPTION,
		Long:  msg.LONGDESCRIPTION,
		Example: heredoc.Doc(`
		$ azion logout --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			settings, err := token.ReadSettings()
			if err != nil {
				return err
			}

			if settings.UUID != "" {
				client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
				err = client.Delete(context.Background(), settings.UUID)
				if err != nil {
					return fmt.Errorf(msg.ErrorLogout, err.Error())
				}
			}

			settings.UUID = ""
			settings.Token = ""
			err = token.WriteSettings(settings)
			if err != nil {
				return err
			}

			logoutOut := output.GeneralOutput{
				Msg: msg.SUCCESS,
				Out: f.IOStreams.Out,
			}

			return output.Print(&logoutOut)
		},
	}

	flags := cmd.Flags()
	flags.BoolP("help", "h", false, msg.FLAGHELP)
	return cmd
}
