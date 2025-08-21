package purge

import (
	"context"
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/purge"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	urls      string
	wildcard  string
	cachekeys string
	Layer     string
)

type PurgeCmd struct {
	Io             *iostreams.IOStreams
	PurgeUrls      func([]string, *cmdutil.Factory) error
	PurgeWildcard  func([]string, *cmdutil.Factory) error
	PurgeCacheKeys func([]string, *cmdutil.Factory) error
	GetPurgeType   func() (string, error)
	AskForInput    func() ([]string, error)
}

func NewPurgeCmd(f *cmdutil.Factory) *PurgeCmd {
	return &PurgeCmd{
		Io:             f.IOStreams,
		PurgeUrls:      purgeUrls,
		PurgeWildcard:  purgeWildcard,
		PurgeCacheKeys: purgeCacheKeys,
		GetPurgeType:   getPurgeType,
		AskForInput:    askForInput,
	}
}

func NewCobraCmd(purge *PurgeCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
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
			ctx := context.Background()
			return purge.Run(ctx, cmd, f)
		},
	}

	cobraCmd.Flags().StringVar(&urls, "urls", "", msg.FlagUrls)
	cobraCmd.Flags().StringVar(&wildcard, "wildcard", "", msg.FlagWildcard)
	cobraCmd.Flags().StringVar(&cachekeys, "cache-key", "", msg.FlagCacheKeys)
	cobraCmd.Flags().StringVar(&Layer, "layer", "edge_caching", msg.FlagLayer)
	cobraCmd.Flags().BoolP("help", "h", false, msg.FlagHelp)

	return cobraCmd
}

func (purge *PurgeCmd) Run(ctx context.Context, cmd *cobra.Command, f *cmdutil.Factory) error {
	if !cmd.Flags().Changed("urls") && !cmd.Flags().Changed("wildcard") && !cmd.Flags().Changed("cache-key") {
		answer, err := purge.GetPurgeType()
		if err != nil {
			return err
		}

		listOfUrls, err := purge.AskForInput()
		if err != nil {
			return err
		}

		switch strings.ToLower(answer) {
		case "urls":
			err := purge.PurgeUrls(listOfUrls, f)
			if err != nil {
				return err
			}
		case "wildcard":
			err := purge.PurgeWildcard(listOfUrls, f)
			if err != nil {
				return err
			}
		case "cache-key":
			err := purge.PurgeCacheKeys(listOfUrls, f)
			if err != nil {
				return err
			}
		}

		return nil
	}

	if cmd.Flags().Changed("urls") {
		err := purge.PurgeUrls(strings.Split(urls, ","), f)
		if err != nil {
			return err
		}
	}

	if cmd.Flags().Changed("wildcard") {
		err := purge.PurgeWildcard(strings.Split(wildcard, ","), f)
		if err != nil {
			return err
		}
	}

	if cmd.Flags().Changed("cache-key") {
		err := purge.PurgeCacheKeys(strings.Split(cachekeys, ","), f)
		if err != nil {
			return err
		}
	}

	logger.FInfo(f.IOStreams.Out, msg.PurgeSuccessful)
	return nil
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewPurgeCmd(f), f)
}
