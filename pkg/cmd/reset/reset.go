package reset

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/reset"
	api "github.com/aziontech/azion-cli/pkg/api/personal_token"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/spf13/cobra"
)

type ResetCmd struct {
	Io            *iostreams.IOStreams
	ReadSettings  func(string) (token.Settings, error)
	WriteSettings func(token.Settings, string) error
	DeleteToken   func(context.Context, string) error
}

func NewResetCmd(f *cmdutil.Factory) *ResetCmd {
	return &ResetCmd{
		Io:            f.IOStreams,
		ReadSettings:  token.ReadSettings,
		WriteSettings: token.WriteSettings,
		DeleteToken: func(ctx context.Context, uuid string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			return client.Delete(ctx, uuid)
		},
	}
}

func NewCobraCmd(reset *ResetCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   msg.USAGE,
		Short: msg.SHORTDESCRIPTION,
		Long:  msg.LONGDESCRIPTION,
		Example: heredoc.Doc(`
		$ azion reset --help
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return reset.run(f)
		},
	}
	cobraCmd.Flags().BoolP("help", "h", false, msg.FLAGHELP)
	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewResetCmd(f), f)
}

func (cmd *ResetCmd) run(f *cmdutil.Factory) error {
	activeProfile := f.GetActiveProfile()
	settings, err := cmd.ReadSettings(activeProfile)
	if err != nil {
		return err
	}

	if settings.UUID != "" {
		err = cmd.DeleteToken(context.Background(), settings.UUID)
		if err != nil {
			return fmt.Errorf(msg.ERRORLOGOUT, err.Error())
		}
	}

	settings = token.Settings{}
	err = cmd.WriteSettings(settings, activeProfile)
	if err != nil {
		return err
	}

	resetOut := output.GeneralOutput{
		Msg: msg.SUCCESS,
		Out: cmd.Io.Out,
	}

	return output.Print(&resetOut)
}
