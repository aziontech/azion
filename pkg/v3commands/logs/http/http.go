package http

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/logs/http"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/v3api/graphql/http"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	startTime = time.Now()
	utcTime   = startTime.UTC()
)

type LogsCmd struct {
	Io        *cmdutil.Factory
	Tail      bool
	Pretty    bool
	LogTime   time.Time
	Limit     string
	GetEvents func(f *cmdutil.Factory, currentTime time.Time, limitFlag string) (http.HTTPEventsResponse, error)
}

func NewLogsCmd(f *cmdutil.Factory) *LogsCmd {
	return &LogsCmd{
		Io:        f,
		LogTime:   utcTime.Add(-5 * time.Minute),
		GetEvents: http.HttpEvents,
	}
}

func NewCobraCmd(logs *LogsCmd, f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion logs http
		$ azion logs http --tail
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := printLogs(logs, cmd, f)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	cmd.Flags().StringVar(&logs.Limit, "limit", "100", msg.LimitFlag)
	cmd.Flags().BoolVar(&logs.Tail, "tail", false, msg.FlagTail)
	cmd.Flags().BoolVar(&logs.Pretty, "pretty", false, msg.FlagPretty)
	return cmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewLogsCmd(f), f)
}

func printLogs(logs *LogsCmd, cmd *cobra.Command, f *cmdutil.Factory) error {
	resp, err := logs.GetEvents(f, logs.LogTime, logs.Limit)
	if err != nil {
		return err
	}

	for _, event := range resp.HTTPEvents {
		if logs.Tail && logs.LogTime.After(event.Ts) {
			continue
		}

		colorLog := color.FgGreen

		if logs.Pretty {
			color.New(colorLog).Fprint(f.IOStreams.Out, "Timestamp: ")
			logger.FInfo(f.IOStreams.Out, event.Ts.String())
			logger.FInfo(f.IOStreams.Out, "\n")
			color.New(colorLog).Fprint(f.IOStreams.Out, "Host: ")
			logger.FInfo(f.IOStreams.Out, event.Host)
			logger.FInfo(f.IOStreams.Out, "\n")
			color.New(colorLog).Fprint(f.IOStreams.Out, "Request URI: ")
			logger.FInfo(f.IOStreams.Out, event.RequestURI)
			logger.FInfo(f.IOStreams.Out, "\n")
			color.New(colorLog).Fprint(f.IOStreams.Out, "Status: ")
			logger.FInfo(f.IOStreams.Out, fmt.Sprint(event.Status))
			logger.FInfo(f.IOStreams.Out, "\n")
			color.New(colorLog).Fprint(f.IOStreams.Out, "User Agent: ")
			logger.FInfo(f.IOStreams.Out, event.HTTPUserAgent)
			logger.FInfo(f.IOStreams.Out, "\n")
			color.New(colorLog).Fprint(f.IOStreams.Out, "Bytes Sent: ")
			logger.FInfo(f.IOStreams.Out, fmt.Sprint(event.UpstreamBytesSent))
			logger.FInfo(f.IOStreams.Out, "\n")
			color.New(colorLog).Fprint(f.IOStreams.Out, "Request Time: ")
			logger.FInfo(f.IOStreams.Out, event.RequestTime)
			logger.FInfo(f.IOStreams.Out, "\n")
			color.New(colorLog).Fprint(f.IOStreams.Out, "Request Method: ")
			logger.FInfo(f.IOStreams.Out, event.RequestMethod)
			logger.FInfo(f.IOStreams.Out, "\n")
			color.New(colorLog).Fprint(f.IOStreams.Out, "Region Name: ")
			logger.FInfo(f.IOStreams.Out, event.GeolocRegion)
			logger.FInfo(f.IOStreams.Out, "\n\n")
		} else {
			logger.FInfo(f.IOStreams.Out, fmt.Sprintf("Timestamp: %s, Host: %s, Request URI: %s, Status: %s, User Agent: %s, Region Name: %s, Bytes Sent: %s, Request Time: %s, Request Method: %s \n\n",
				event.Ts.String(), event.Host, event.RequestURI, fmt.Sprint(event.Status), event.HTTPUserAgent, event.GeolocRegion, fmt.Sprint(event.UpstreamBytesSent), event.RequestTime, event.RequestMethod))
		}

		logs.LogTime = event.Ts
	}

	if logs.Tail {
		logger.FInfo(f.IOStreams.Out, msg.NewLogs)
		logger.FInfo(f.IOStreams.Out, "\n\n")
		time.Sleep(10 * time.Second)
		return printLogs(logs, cmd, f)
	}
	return nil
}
