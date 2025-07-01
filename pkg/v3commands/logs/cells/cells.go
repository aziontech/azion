package cells

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/logs/cells"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/v3api/graphql/cells"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	startTime = time.Now()
	utcTime   = startTime.UTC()
)

type LogsCmd struct {
	Io         *cmdutil.Factory
	FunctionId string
	Tail       bool
	Pretty     bool
	LogTime    time.Time
	Limit      string
	GetLogs    func(f *cmdutil.Factory, functionId string, logTime time.Time, limit string) (cells.CellsConsoleEventsResponse, error)
}

func NewLogsCmd(f *cmdutil.Factory) *LogsCmd {
	return &LogsCmd{
		Io:      f,
		LogTime: utcTime.Add(-5 * time.Minute),
		GetLogs: cells.CellsConsoleLogs,
	}
}

func NewCobraCmd(logs *LogsCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion logs cells
		$ azion logs cells --tail
		$ azion logs cells --function-id 1234 --limit 10
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := printLogs(logs, cmd)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	cmd.Flags().StringVar(&logs.FunctionId, "function-id", "", msg.FlagFunctionId)
	cmd.Flags().StringVar(&logs.Limit, "limit", "100", msg.LimitFlag)
	cmd.Flags().BoolVar(&logs.Tail, "tail", false, msg.FlagTail)
	cmd.Flags().BoolVar(&logs.Pretty, "pretty", false, msg.FlagPretty)
	return cmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewLogsCmd(f))
}

func printLogs(logs *LogsCmd, cmd *cobra.Command) error {
	resp, err := logs.GetLogs(logs.Io, logs.FunctionId, logs.LogTime, logs.Limit)
	if err != nil {
		return err
	}

	for _, event := range resp.CellsConsoleEvents {
		if logs.Tail && logs.LogTime.After(event.Ts) {
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

		if logs.Pretty {
			logger.FInfo(logs.Io.IOStreams.Out, "Function ID: ")
			logger.FInfo(logs.Io.IOStreams.Out, event.FunctionId)
			logger.FInfo(logs.Io.IOStreams.Out, "\n")
			logger.FInfo(logs.Io.IOStreams.Out, "Timestamp: ")
			logger.FInfo(logs.Io.IOStreams.Out, event.Ts.String())
			logger.FInfo(logs.Io.IOStreams.Out, "\n")
			logger.FInfo(logs.Io.IOStreams.Out, "Log: \n")
			color.New(colorLog).Fprintln(logs.Io.IOStreams.Out, event.Line)
			logger.FInfo(logs.Io.IOStreams.Out, "\n\n")
		} else {
			logger.FInfo(logs.Io.IOStreams.Out, fmt.Sprintf("Function ID: %s, Timestamp: %s, Log: %s \n", event.FunctionId, event.Ts.String(), event.Line))
		}

		logs.LogTime = event.Ts
	}

	if logs.Tail {
		logger.FInfo(logs.Io.IOStreams.Out, msg.NewLogs)
		logger.FInfo(logs.Io.IOStreams.Out, "\n\n")
		time.Sleep(10 * time.Second)
		return printLogs(logs, cmd)
	}
	return nil
}
