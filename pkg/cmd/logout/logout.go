package logout

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/logout"
	api "github.com/aziontech/azion-cli/pkg/api/personal_token"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/spf13/cobra"
)

type LogoutCmd struct {
	Io            *iostreams.IOStreams
	ReadSettings  func() (token.Settings, error)
	WriteSettings func(token.Settings) error
	DeleteToken   func(context.Context, string) error
}

func NewLogoutCmd(f *cmdutil.Factory) *LogoutCmd {
	return &LogoutCmd{
		Io:            f.IOStreams,
		ReadSettings:  token.ReadSettings,
		WriteSettings: token.WriteSettings,
		DeleteToken: func(ctx context.Context, uuid string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			return client.Delete(ctx, uuid)
		},
	}
}

func NewCobraCmd(logout *LogoutCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   msg.USAGE,
		Short: msg.SHORTDESCRIPTION,
		Long:  msg.LONGDESCRIPTION,
		Example: heredoc.Doc(`
		$ azion logout --help
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logout.run()
		},
	}

	cobraCmd.Flags().BoolP("help", "h", false, msg.FLAGHELP)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewLogoutCmd(f), f)
}

func (cmd *LogoutCmd) run() error {
	settings, err := cmd.ReadSettings()
	if err != nil {
		return err
	}

	if settings.UUID != "" {
		err = cmd.DeleteToken(context.Background(), settings.UUID)
		if err != nil {
			return fmt.Errorf(msg.ErrorLogout, err.Error())
		}
	}

	settings.UUID = ""
	settings.Token = ""
	err = cmd.WriteSettings(settings)
	if err != nil {
		return err
	}

	logoutOut := output.GeneralOutput{
		Msg: msg.SUCCESS,
		Out: cmd.Io.Out,
	}

	return output.Print(&logoutOut)
}
