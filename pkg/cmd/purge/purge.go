package purge

import (
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/purge"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	urls      string
	wildcard  string
	cachekeys string
	Layer     string
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion purge --wildcard "www.example.com/*"
        $ azion purge --urls "www.example.com,www.pudim.com"
        $ azion purge --cache-key "www.domain.com/@@cookie_name=cookie_value,www.domain.com/test.js"
        `),
		RunE: func(cmd *cobra.Command, _ []string) error {

			if !cmd.Flags().Changed("urls") && !cmd.Flags().Changed("wildcard") && !cmd.Flags().Changed("cache-key") {
				answer, err := getPurgeType()
				if err != nil {
					return err
				}

				listOfUrls, err := askForInput()
				if err != nil {
					return err
				}

				switch strings.ToLower(answer) {
				case "urls":
					err := purgeUrls(listOfUrls, f)
					if err != nil {
						return err
					}
				case "wildcard":
					err := purgeWildcard(listOfUrls, f)
					if err != nil {
						return err
					}
				case "cache-key":
					err := purgeCacheKeys(listOfUrls, f)
					if err != nil {
						return err
					}
				}

				return nil
			}

			if cmd.Flags().Changed("urls") {
				err := purgeUrls(strings.Split(urls, ","), f)
				if err != nil {
					return err
				}
			}

			if cmd.Flags().Changed("wildcard") {
				err := purgeWildcard(strings.Split(wildcard, ","), f)
				if err != nil {
					return err
				}
			}

			if cmd.Flags().Changed("cache-key") {
				err := purgeCacheKeys(strings.Split(cachekeys, ","), f)
				if err != nil {
					return err
				}
			}

			logger.FInfo(f.IOStreams.Out, msg.PurgeSuccessful)
			return nil
		},
	}

	cmd.Flags().StringVar(&urls, "urls", "", msg.FlagUrls)
	cmd.Flags().StringVar(&wildcard, "wildcard", "", msg.FlagWildcard)
	cmd.Flags().StringVar(&cachekeys, "cache-key", "", msg.FlagCacheKeys)
	cmd.Flags().StringVar(&Layer, "layer", "edge_caching", msg.FlagLayer)
	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)

	return cmd
}
