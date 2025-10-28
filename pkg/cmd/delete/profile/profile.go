package profile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/profile"
	api "github.com/aziontech/azion-cli/pkg/api/personal_token"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Fields struct {
	Name string
}

var confirmFn = utils.Confirm

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.UsageDelete,
		Short:         msg.DeleteShortDescription,
		Long:          msg.DeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion delete profile --name "my-profile"
		$ azion delete profile
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			var profileToDelete string

			if !cmd.Flags().Changed("name") {
				// List profiles and let user choose
				dir := config.Dir()
				entries, err := os.ReadDir(dir.Dir)
				if err != nil {
					return fmt.Errorf("%w", msg.ErrorReadDir)
				}

				var profileNames []string
				for _, entry := range entries {
					if entry.IsDir() && !strings.HasPrefix(entry.Name(), "tempclonesamples") {
						profileNames = append(profileNames, entry.Name())
					}
				}

				if len(profileNames) == 0 {
					return fmt.Errorf("No profiles found")
				}

				prompt := &survey.Select{
					Message:  msg.QuestionDeleteProfile,
					Options:  profileNames,
					PageSize: len(profileNames),
				}

				err = survey.AskOne(prompt, &profileToDelete)
				if err != nil {
					return err
				}
			} else {
				profileToDelete = fields.Name
			}

			if profileToDelete == "default" {
				return msg.ErrorCannotDeleteDefault
			}

			confirmDelete := confirmFn(false, fmt.Sprintf(msg.ConfirmDeleteProfile, profileToDelete), true)
			if !confirmDelete {
				return msg.ErrorDeleteCancelled
			}

			dir := config.Dir()
			profilePath := filepath.Join(dir.Dir, profileToDelete)
			if _, err := os.Stat(profilePath); os.IsNotExist(err) {
				return fmt.Errorf(msg.ErrorProfileNotFound.Error(), profileToDelete)
			}

			settings, err := token.ReadSettings(profileToDelete)
			if err == nil && settings.UUID != "" {
				client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), settings.Token)
				err = client.Delete(context.Background(), settings.UUID)
				if err != nil {
					fmt.Fprintf(f.IOStreams.Out, msg.WarningDeleteToken+"\n", err)
				}
			}

			err = os.RemoveAll(profilePath)
			if err != nil {
				return fmt.Errorf(msg.ErrorDeleteProfile.Error(), err)
			}

			currentProfile, _, err := token.ReadProfiles()
			if err == nil && currentProfile.Name == profileToDelete {
				defaultProfile := token.Profile{Name: "default"}
				err = token.WriteProfiles(defaultProfile)
				if err != nil {
					fmt.Fprintf(f.IOStreams.Out, msg.WarningSetActiveProfile+"\n", err)
				}
			}

			profileOut := output.GeneralOutput{
				Msg: fmt.Sprintf(msg.DeleteOutputSuccess, profileToDelete),
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
	flags.BoolP("help", "h", false, msg.DeleteFlagHelp)
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
}
