package cells

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/logs/cells"
	"github.com/aziontech/azion-cli/pkg/api/graphql/cells"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	functionId string
	tail       bool
	pretty     bool
	startTime  = time.Now()
	utcTime    = startTime.UTC()
	logTime    time.Time
	limit      string
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	logTime = utcTime.Add(-5 * time.Minute)
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true, Example: heredoc.Doc(`
		$ azion logs cells
		$ azion logs cells --tail
		$ azion logs cells --function-id 1234 --limit 10
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := printLogs(cmd, f)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	cmd.Flags().StringVar(&functionId, "function-id", "", msg.FlagFunctionId)
	cmd.Flags().StringVar(&limit, "limit", "100", msg.LimitFlag)
	cmd.Flags().BoolVar(&tail, "tail", false, msg.FlagTail)
	cmd.Flags().BoolVar(&pretty, "pretty", false, msg.FlagPretty)
	return cmd
}

func printLogs(cmd *cobra.Command, f *cmdutil.Factory) error {

	resp, err := cells.CellsConsoleLogs(f, functionId, logTime, limit)
	if err != nil {
		return err
	}

	for _, event := range resp.CellsConsoleEvents {
		if tail && logTime.After(event.Ts) {
			continue
		}

		var colorLog color.Attribute

		switch event.Level {
		case "LOG":
			colorLog = color.FgGreen
		case "ERROR":
			colorLog = color.FgRed
		default:
			colorLog = color.FgWhite
		}

		if pretty {
			logger.FInfo(f.IOStreams.Out, "Function ID: ")
			logger.FInfo(f.IOStreams.Out, event.FunctionId)
			logger.FInfo(f.IOStreams.Out, "\n")
			logger.FInfo(f.IOStreams.Out, "Timestamp: ")
			logger.FInfo(f.IOStreams.Out, event.Ts.String())
			logger.FInfo(f.IOStreams.Out, "\n")
			logger.FInfo(f.IOStreams.Out, "Log: \n")
			color.New(colorLog).Fprintln(f.IOStreams.Out, event.Line)
			logger.FInfo(f.IOStreams.Out, "\n\n")
		} else {
			logger.FInfo(f.IOStreams.Out, fmt.Sprintf("Function ID: %s, Timestamp: %s, Log: %s \n", event.FunctionId, event.Ts.String(), event.Line))
		}

		logTime = event.Ts

	}

	if tail {
		logger.FInfo(f.IOStreams.Out, msg.NewLogs)
		logger.FInfo(f.IOStreams.Out, "\n\n")
		time.Sleep(10 * time.Second)
		return printLogs(cmd, f)
	}
	return nil
}
