package datastream

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/data_stream"
	streams "github.com/aziontech/azion-cli/pkg/cmd/create/data_stream/streams"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.CreateShortDescription,
		Long:          msg.CreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       heredoc.Doc(msg.CreateExample),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(streams.NewCmd(f))
	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
