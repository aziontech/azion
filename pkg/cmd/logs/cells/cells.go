package cells

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/list/rules_engine"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/spf13/cobra"
)

var edgeApplicationID int64
var phase string
var watch bool

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	cmd := &cobra.Command{
		Use:           msg.RulesEngineListUsage,
		Short:         msg.RulesEngineListShortDescription,
		Long:          msg.RulesEngineListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true, Example: heredoc.Doc(`
		$ azion logs cells
		$ azion logs cells --tail
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			// startTime := time.Now()
			// resp, err := cells.Option1("wsdsds")
			// if err != nil {
			// 	return err
			// }

			// for _, event := range resp.HTTPEvents {
			// 	timeNow := time.Now()
			// 	if watch && (timeNow < startTime) {
			// 		continue
			// 	}
			// }
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().BoolP("help", "h", false, msg.RulesEngineListHelpFlag)
	cmd.Flags().Int64Var(&edgeApplicationID, "application-id", 0, msg.ApplicationFlagId)
	cmd.Flags().StringVar(&phase, "phase", "request", msg.RulesEnginePhase)
	return cmd
}
