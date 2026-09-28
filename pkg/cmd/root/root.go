package root

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/root"
	"github.com/aziontech/azion-cli/pkg/cmdregistry"
	"github.com/aziontech/azion-cli/pkg/metric"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/schedule"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const PREFIX_FLAG = "--"

type factoryRoot struct {
	factory           *cmdutil.Factory
	doPreCommandCheck func(cmd *cobra.Command, fact *factoryRoot) error //this package
	execSchedules     func(factory *cmdutil.Factory)                    //schedule.ExecShedules
	command           cmdutil.Command
	osExit            func(code int)
	flags
	globals
}

type flags struct {
	tokenFlag  string
	configFlag string
	timeout    int
}

type globals struct {
	commandName    string
	globalSettings *token.Settings
	startTime      time.Time
	// how this invocation's API version was obtained; logged once the
	// requested log level is in effect
	apiVersionSource apiVersionSource
}

func (fact *factoryRoot) persistentPreRunE(cmd *cobra.Command, _ []string) error {
	fact.startTime = time.Now()
	logger.LogLevel(fact.factory.Logger)
	fact.logAPIVersionSource()

	if strings.HasPrefix(fact.configFlag, PREFIX_FLAG) {
		return msg.ErrorPrefix
	}

	if err := fact.doPreCommandCheck(cmd, fact); err != nil {
		return err
	}

	fact.execSchedules(fact.factory)
	return nil
}

func (fact *factoryRoot) runE(cmd *cobra.Command, _ []string) error {
	if cmd.Flags().Changed("token") {
		return nil
	}
	return cmd.Help()
}

func (fact *factoryRoot) setFlags(cobraCmd *cobra.Command) {
	cobraCmd.PersistentFlags().StringVarP(&fact.tokenFlag, "token", "t", "", msg.RootTokenFlag)
	cobraCmd.PersistentFlags().StringVarP(&fact.configFlag, "config", "c", "", msg.RootConfigFlag)
	cobraCmd.PersistentFlags().BoolVarP(&fact.factory.Debug, "debug", "d", false, msg.RootLogDebug)
	cobraCmd.PersistentFlags().BoolVarP(&fact.factory.Silent, "silent", "s", false, msg.RootLogSilent)
	cobraCmd.PersistentFlags().StringVarP(&fact.factory.LogLevel, "log-level", "l", "info", msg.RootLogLevel)
	cobraCmd.PersistentFlags().BoolVarP(&fact.factory.GlobalFlagAll, "yes", "y", false, msg.RootYesFlag)
	cobraCmd.PersistentFlags().StringVar(&fact.factory.Out, "out", "", msg.RootFlagOut)
	cobraCmd.PersistentFlags().StringVar(&fact.factory.Format, "format", "", msg.RootFlagFormat)
	cobraCmd.PersistentFlags().BoolVar(&fact.factory.NoColor, "no-color", false, msg.RootFlagFormat)
	cobraCmd.PersistentFlags().IntVar(&fact.timeout, "timeout", 50, msg.RootFlagTimeout)
	cobraCmd.Flags().BoolP("help", "h", false, msg.RootHelpFlag)
}

// setCmds registers the top-level commands the account's generation is
// entitled to. The list itself lives in table.go.
func (fact *factoryRoot) setCmds(cobraCmd *cobra.Command) {
	cobraCmd.AddCommand(cmdregistry.Children(children, fact.factory)...)
}

func (fact *factoryRoot) CmdRoot() cmdutil.Command {
	cobraCmd := &cobra.Command{
		Use:               msg.RootUsage,
		Long:              msg.RootDescription,
		Short:             color.New(color.Bold).Sprint(fmt.Sprintf(msg.RootDescription, version.BinVersion)),
		Version:           version.BinVersion,
		PersistentPreRunE: fact.persistentPreRunE,
		Example:           heredoc.Doc(msg.EXAMPLE),
		RunE:              fact.runE,
		SilenceErrors:     true, // Silence errors, so the help message won't be shown on flag error
		SilenceUsage:      true, // Silence usage on error
	}

	cobraCmd.SetIn(fact.factory.IOStreams.In)
	cobraCmd.SetOut(fact.factory.IOStreams.Out)
	cobraCmd.SetErr(fact.factory.IOStreams.Err)

	cobraCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		rootHelpFunc(cmd, args)
	})

	fact.setFlags(cobraCmd)

	// The tree below is chosen from the account's API generation, which depends
	// on --token and --config. Cobra has not parsed them yet, so read them now.
	fact.preParseGlobalFlags(os.Args[1:])
	fact.applyConfigFlag()

	// set template for -v flag
	cobraCmd.SetVersionTemplate(color.New(color.Bold).Sprint("Azion CLI " + version.BinVersion + "\n"))

	fact.factory.APIVersion = fact.resolveAPIVersion()

	fact.setCmds(cobraCmd)

	return cobraCmd
}

func NewFactoryRoot(fact *cmdutil.Factory) *factoryRoot {
	return &factoryRoot{
		factory:           fact,
		doPreCommandCheck: doPreCommandCheck,
		execSchedules:     schedule.ExecSchedules,
		command:           &cobra.Command{},
		osExit:            os.Exit,
	}
}

func Execute(f *factoryRoot) {
	logger.New(zapcore.InfoLevel)

	cmd := f.CmdRoot()
	err := cmd.Execute()
	executionTime := time.Since(f.startTime).Seconds()

	// 1 = authorize; anything different than 1 means that the user did not authorize metrics collection, or did not answer the question yet
	if f.globalSettings != nil {
		if f.globalSettings.AuthorizeMetricsCollection == 1 {
			activeProfile := f.factory.GetActiveProfile()
			errMetrics := metric.TotalCommandsCount(cmd, f.commandName, executionTime, err, activeProfile, f.factory.APIVersion.String())
			if errMetrics != nil {
				logger.Debug("Error while saving metrics", zap.Error(err))
			}
		}
	}

	if err != nil {
		output.Print(&output.ErrorOutput{
			GeneralOutput: output.GeneralOutput{
				Out:   f.factory.IOStreams.Out,
				Flags: f.factory.Flags,
			},
			Err: err,
		})
		f.osExit(1)
	}
}
