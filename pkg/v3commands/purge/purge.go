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

type PurgeCmd struct {
	Io             *iostreams.IOStreams
	PurgeUrls      func([]string, *cmdutil.Factory) error
	PurgeWildcard  func([]string, *cmdutil.Factory) error
	PurgeCacheKeys func([]string, *cmdutil.Factory, string) error
	GetPurgeType   func() (string, error)
	AskForInput    func() ([]string, error)
	Layer          string
	Cachekeys      string
	Urls           string
	Wildcard       string
}

func NewPurgeCmd(f *cmdutil.Factory) *PurgeCmd {
	return &PurgeCmd{
		Layer:          "edge_caching",
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

	cobraCmd.Flags().StringVar(&purge.Urls, "urls", "", msg.FlagUrls)
	cobraCmd.Flags().StringVar(&purge.Wildcard, "wildcard", "", msg.FlagWildcard)
	cobraCmd.Flags().StringVar(&purge.Cachekeys, "cache-key", "", msg.FlagCacheKeys)
	cobraCmd.Flags().StringVar(&purge.Layer, "layer", "edge_caching", msg.FlagLayer)
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
			err := purge.PurgeCacheKeys(listOfUrls, f, purge.Layer)
			if err != nil {
				return err
			}
		}

		return nil
	}

	if cmd.Flags().Changed("urls") {
		err := purge.PurgeUrls(strings.Split(purge.Urls, ","), f)
		if err != nil {
			return err
		}
	}

	if cmd.Flags().Changed("wildcard") {
		err := purge.PurgeWildcard(strings.Split(purge.Wildcard, ","), f)
		if err != nil {
			return err
		}
	}

	if cmd.Flags().Changed("cache-key") {
		err := purge.PurgeCacheKeys(strings.Split(purge.Cachekeys, ","), f, purge.Layer)
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
