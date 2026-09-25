package purge

import (
	"context"
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/purge"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/registry"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type PurgeCmd struct {
	Io           *iostreams.IOStreams
	GetPurgeType func() (string, error)
	AskForInput  func() ([]string, error)
	Layer        string
	Cachekeys    string
	Urls         string
	Wildcard     string
}

func NewPurgeCmd(f *cmdutil.Factory) *PurgeCmd {
	return &PurgeCmd{
		Layer:        "cache",
		Io:           f.IOStreams,
		GetPurgeType: getPurgeType,
		AskForInput:  askForInput,
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
        $ azion purge --cachekey "www.domain.com/@@cookie_name=cookie_value,www.domain.com/test.js"
        `),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			return purge.Run(ctx, cmd, f)
		},
	}

	cobraCmd.Flags().StringVar(&purge.Urls, "urls", "", msg.FlagUrls)
	cobraCmd.Flags().StringVar(&purge.Wildcard, "wildcard", "", msg.FlagWildcard)
	cobraCmd.Flags().StringVar(&purge.Cachekeys, "cachekey", "", msg.FlagCacheKeys)
	cobraCmd.Flags().StringVar(&purge.Layer, "layer", "cache", msg.FlagLayer)
	cobraCmd.Flags().BoolP("help", "h", false, msg.FlagHelp)

	return cobraCmd
}

func (purge *PurgeCmd) Run(ctx context.Context, cmd *cobra.Command, f *cmdutil.Factory) error {
	clipurge := registry.RealtimePurge(f)
	//if none of the flags were sent
	if !cmd.Flags().Changed("urls") && !cmd.Flags().Changed("wildcard") && !cmd.Flags().Changed("cachekey") {
		answer, err := purge.GetPurgeType()
		if err != nil {
			return err
		}

		answer = strings.ReplaceAll(answer, " ", "")
		answer = strings.ToLower(answer)

		listOfUrls, err := purge.AskForInput()
		if err != nil {
			return err
		}

		err = clipurge.PurgeCache(ctx, listOfUrls, answer, purge.Layer)
		if err != nil {
			logger.Debug("Error while purging domains", zap.Error(err))
			return err
		}

		logger.FInfo(f.IOStreams.Out, msg.PurgeSuccessful)

		return nil
	}

	if cmd.Flags().Changed("urls") {
		err := clipurge.PurgeCache(ctx, strings.Split(purge.Urls, ","), "url", purge.Layer)
		if err != nil {
			logger.Debug("Error while purging domains", zap.Error(err))
			return err
		}
	}

	if cmd.Flags().Changed("wildcard") {
		splitWildcard := strings.Split(purge.Wildcard, ",")
		if len(splitWildcard) > 1 {
			logger.Debug("More than one URL for wildcard", zap.Any("Amount of URLs", len(splitWildcard)))
			return msg.ErrorTooManyUrls
		}
		err := clipurge.PurgeCache(ctx, splitWildcard, "wildcard", purge.Layer)
		if err != nil {
			logger.Debug("Error while purging domains", zap.Error(err))
			return err
		}
	}

	if cmd.Flags().Changed("cachekey") {
		err := clipurge.PurgeCache(ctx, strings.Split(purge.Cachekeys, ","), "cachekey", purge.Layer)
		if err != nil {
			logger.Debug("Error while purging domains", zap.Error(err))
			return err
		}
	}

	logger.FInfo(f.IOStreams.Out, msg.PurgeSuccessful)
	return nil
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewPurgeCmd(f), f)
}
