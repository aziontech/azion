package configure

import (
	"fmt"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var configureToken string

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	// configureCmd represents the configure command
	configureCmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure parameters and credentials",
		Long:  `This command configures cli parameters and credentials used for connecting to our services.`,
		Annotations: map[string]string{
			"IsAdditional": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.HttpClient()
			if err != nil {
				return fmt.Errorf("%s: %w", utils.ErrorGetHttpClient, err)
			}

			t, err := token.New(&token.Config{
				Client: client,
				Out:    f.IOStreams.Out,
			})
			if err != nil {
				return fmt.Errorf("%s: %w", utils.ErrorTokenManager, err)
			}

			if configureToken == "" {
				return utils.ErrorTokenNotProvided
			}

			valid, err := t.Validate(&configureToken)
			if err != nil {
				return err
			}

			if !valid {
				return utils.ErrorInvalidToken
			}

			if err := t.Save(); err != nil {
				return err
			}

			return nil
		},
	}

	configureCmd.SetIn(f.IOStreams.In)
	configureCmd.SetOut(f.IOStreams.Out)
	configureCmd.SetErr(f.IOStreams.Err)

	configureCmd.Flags().StringVarP(&configureToken, "token", "t", "", "Save provided token to use on any command")

	return configureCmd
}
