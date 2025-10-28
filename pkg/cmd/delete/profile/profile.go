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

			// Prevent deletion of default profile
			if profileToDelete == "default" {
				return msg.ErrorCannotDeleteDefault
			}

			// Confirm deletion
			confirmDelete := confirmFn(false, fmt.Sprintf(msg.ConfirmDeleteProfile, profileToDelete), true)
			if !confirmDelete {
				return fmt.Errorf("Profile deletion cancelled")
			}

			// Check if profile exists
			dir := config.Dir()
			profilePath := filepath.Join(dir.Dir, profileToDelete)
			if _, err := os.Stat(profilePath); os.IsNotExist(err) {
				return fmt.Errorf("Profile '%s' not found", profileToDelete)
			}

			// Try to delete token if UUID exists
			settings, err := token.ReadSettings(profileToDelete)
			if err == nil && settings.UUID != "" {
				client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), settings.Token)
				err = client.Delete(context.Background(), settings.UUID)
				if err != nil {
					// Log warning but don't fail the profile deletion
					fmt.Fprintf(f.IOStreams.Out, "Warning: Failed to delete token from server: %v\n", err)
				}
			}

			// Delete the profile directory and all its contents
			err = os.RemoveAll(profilePath)
			if err != nil {
				return fmt.Errorf("Failed to delete the profile: %w", err)
			}

			// Check if deleted profile was the active one
			currentProfile, _, err := token.ReadProfiles()
			if err == nil && currentProfile.Name == profileToDelete {
				// Set active profile to "default"
				defaultProfile := token.Profile{Name: "default"}
				err = token.WriteProfiles(defaultProfile)
				if err != nil {
					fmt.Fprintf(f.IOStreams.Out, "Warning: Failed to set active profile to default: %v\n", err)
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
