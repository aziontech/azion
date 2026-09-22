package profiles

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/profile"
	"github.com/aziontech/azion-cli/pkg/apiversion"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/constants"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
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
        $ azion profiles --refresh
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if refresh, _ := cmd.Flags().GetBool("refresh"); refresh {
				return runRefresh(f)
			}

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
	flags.Bool("refresh", false, msg.ProfilesFlagRefresh)
}

// runRefresh re-checks the current profile's API generation against the SSO
// service and stores the result, without switching profiles or asking for
// credentials. It is the escape hatch for an account whose generation changed
// before the cached answer expires on its own.
func runRefresh(f *cmdutil.Factory) error {
	profile, _, err := token.ReadProfiles()
	if err != nil {
		return err
	}

	settings, err := token.ReadSettings(profile.Name)
	if err != nil {
		return err
	}

	// Resolve against the credential this invocation would actually use, which
	// is what root binds the cache to. Reading settings.Token instead would
	// write a hash that root never matches, re-resolving on every invocation.
	tok := f.Config.GetString("token")
	if tok == "" {
		tok = settings.Token
	}
	if tok == "" {
		return fmt.Errorf(msg.ErrorRefreshNoToken.Error(), profile.Name)
	}

	previous := settings.APIVersionCache().Version
	logger.Debug("Refreshing the cached API version on request; ignoring the TTL",
		zap.String("profile", profile.Name),
		zap.String("cached", previous.String()))

	current, err := apiversion.Resolve(f.HttpClient, constants.AuthURL, tok)
	if err != nil {
		if errors.Is(err, apiversion.ErrUnauthorized) {
			return fmt.Errorf(msg.ErrorRefreshUnauthorized.Error(), profile.Name)
		}
		return fmt.Errorf(msg.ErrorRefreshFailed.Error(), profile.Name, err)
	}

	settings.SetAPIVersionCache(apiversion.Refreshed(current, tok, time.Now()))
	if err := token.WriteSettings(settings, profile.Name); err != nil {
		return err
	}

	message := fmt.Sprintf(msg.RefreshUnchanged, profile.Name, current)
	if previous != "" && previous != current {
		message = fmt.Sprintf(msg.RefreshChanged, profile.Name, previous, current)
	}

	refreshOut := output.GeneralOutput{
		Msg:   message,
		Out:   f.IOStreams.Out,
		Flags: f.Flags,
	}
	return output.Print(&refreshOut)
}
