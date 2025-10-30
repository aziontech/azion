package profiles

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/profile"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var confirmFn = utils.Confirm

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           msg.UsageProfiles,
		Short:         msg.ProfilesShortDescription,
		Long:          msg.ProfilesLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion profiles
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := config.Dir()
			entries, err := os.ReadDir(dir.Dir)
			if err != nil {
				return fmt.Errorf(msg.ErrorReadDir.Error(), err)
			}

			var profileNames []string
			for _, entry := range entries {
				if entry.IsDir() && !strings.HasPrefix(entry.Name(), "tempclonesamples") {
					profileNames = append(profileNames, entry.Name())
				}
			}

			prompt := &survey.Select{
				Message:  "Choose a profile:",
				Options:  profileNames,
				PageSize: len(profileNames),
			}

			var answer string
			err = survey.AskOne(prompt, &answer)
			if err != nil {
				return err
			}

			profile, _, err := token.ReadProfiles()
			if err != nil {
				return err
			}

			profile.Name = answer

			err = token.WriteProfiles(profile)
			if err != nil {
				return err
			}

			return nil
		},
	}

	flags := cmd.Flags()
	addFlags(flags)

	return cmd
}

func addFlags(flags *pflag.FlagSet) {
	flags.BoolP("help", "h", false, msg.ProfilesFlagHelp)
}
