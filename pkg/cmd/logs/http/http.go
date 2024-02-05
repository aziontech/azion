package http

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/logs/http"
	"github.com/aziontech/azion-cli/pkg/api/graphql/http"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	tail      bool
	pretty    bool
	startTime = time.Now()
	utcTime   = startTime.UTC()
	logTime   time.Time
	limit     string
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	logTime = utcTime.Add(-5 * time.Minute)
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true, Example: heredoc.Doc(`
		$ azion logs http
		$ azion logs http --tail
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
	cmd.Flags().StringVar(&limit, "limit", "100", msg.LimitFlag)
	cmd.Flags().BoolVar(&tail, "tail", false, msg.FlagTail)
	cmd.Flags().BoolVar(&pretty, "pretty", false, msg.FlagPretty)
	return cmd
}

func printLogs(cmd *cobra.Command, f *cmdutil.Factory) error {

	resp, err := http.HttpEvents(f, logTime, limit)
	if err != nil {
		return err
	}

	for _, event := range resp.HTTPEvents {
		if tail && logTime.After(event.Ts) {
			continue
		}

		colorLog := color.FgGreen

		if pretty {
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
