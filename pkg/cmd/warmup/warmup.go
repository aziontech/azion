package warmup

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/warmup"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/spf13/cobra"
)

var (
	baseUrl       string
	maxUrls       int
	maxConcurrent int
	timeout       int
)

// WarmupCmd defines the command structure
type WarmupCmd struct {
	Io          *iostreams.IOStreams
	WarmupCache func(ctx context.Context, baseUrl string, maxUrls int, maxConcurrent int, timeout int, f *cmdutil.Factory) error
	AskForUrl   func() (string, error)
}

// NewWarmupCmd creates a new WarmupCmd instance
func NewWarmupCmd(f *cmdutil.Factory) *WarmupCmd {
	return &WarmupCmd{
		Io:          f.IOStreams,
		WarmupCache: warmupCache,
		AskForUrl:   askForUrl,
	}
}

// NewCobraCmd creates a new cobra command for warmup
func NewCobraCmd(warmup *WarmupCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion warmup --url "https://example.com"
        $ azion warmup --url "https://example.com/products"
        $ azion warmup --url "https://example.com/blog" --max-urls 500 --max-concurrent 5 --timeout 10000
        `),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			return warmup.Run(ctx, cmd, f)
		},
	}

	cobraCmd.Flags().StringVar(&baseUrl, "url", "", msg.FlagUrl)
	cobraCmd.Flags().IntVar(&maxUrls, "max-urls", 1500, msg.FlagMaxUrls)
	cobraCmd.Flags().IntVar(&maxConcurrent, "max-concurrent", 2, msg.FlagMaxConcurrent)
	cobraCmd.Flags().IntVar(&timeout, "timeout", 8000, msg.FlagTimeout)
	cobraCmd.Flags().BoolP("help", "h", false, msg.FlagHelp)

	return cobraCmd
}

// Run executes the warmup command
func (warmup *WarmupCmd) Run(ctx context.Context, cmd *cobra.Command, f *cmdutil.Factory) error {
	if !cmd.Flags().Changed("url") {
		url, err := warmup.AskForUrl()
		if err != nil {
			return err
		}
		baseUrl = url
	}

	err := warmup.WarmupCache(ctx, baseUrl, maxUrls, maxConcurrent, timeout, f)
	if err != nil {
		return err
	}

	warmupOut := output.GeneralOutput{
		Msg:   msg.WarmupSuccessful,
		Out:   f.IOStreams.Out,
		Flags: f.Flags,
	}
	return output.Print(&warmupOut)
}

// NewCmd creates a new cobra command for warmup
func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewWarmupCmd(f), f)
}
