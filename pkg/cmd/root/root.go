package root

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/root"
	buildCmd "github.com/aziontech/azion-cli/pkg/cmd/build"
	"github.com/aziontech/azion-cli/pkg/cmd/completion"
	"github.com/aziontech/azion-cli/pkg/cmd/create"
	"github.com/aziontech/azion-cli/pkg/cmd/delete"
	"github.com/aziontech/azion-cli/pkg/cmd/describe"
	"github.com/aziontech/azion-cli/pkg/cmd/list"
	"github.com/aziontech/azion-cli/pkg/cmd/login"
	"github.com/aziontech/azion-cli/pkg/cmd/logout"
	logcmd "github.com/aziontech/azion-cli/pkg/cmd/logs"
	"github.com/aziontech/azion-cli/pkg/cmd/purge"
	"github.com/aziontech/azion-cli/pkg/cmd/reset"
	"github.com/aziontech/azion-cli/pkg/cmd/sync"
	"github.com/aziontech/azion-cli/pkg/cmd/unlink"
	"github.com/aziontech/azion-cli/pkg/cmd/update"
	"github.com/aziontech/azion-cli/pkg/cmd/whoami"
	"github.com/aziontech/azion-cli/pkg/metric"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/schedule"

	deploycmd "github.com/aziontech/azion-cli/pkg/cmd/deploy"
	devcmd "github.com/aziontech/azion-cli/pkg/cmd/dev"
	initcmd "github.com/aziontech/azion-cli/pkg/cmd/init"
	linkcmd "github.com/aziontech/azion-cli/pkg/cmd/link"
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
}

type globals struct {
	commandName    string
	globalSettings *token.Settings
	startTime      time.Time
}

func (fact *factoryRoot) persistentPreRunE(cmd *cobra.Command, _ []string) error {
	fact.startTime = time.Now()
	logger.LogLevel(fact.factory.Logger)

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
	cobraCmd.Flags().BoolP("help", "h", false, msg.RootHelpFlag)
}

func (fact *factoryRoot) setCmds(cobraCmd *cobra.Command) {
	cobraCmd.AddCommand(initcmd.NewCmd(fact.factory))
	cobraCmd.AddCommand(logcmd.NewCmd(fact.factory))
	cobraCmd.AddCommand(deploycmd.NewCmd(fact.factory))
	cobraCmd.AddCommand(buildCmd.NewCmd(fact.factory))
	cobraCmd.AddCommand(devcmd.NewCmd(fact.factory))
	cobraCmd.AddCommand(linkcmd.NewCmd(fact.factory))
	cobraCmd.AddCommand(unlink.NewCmd(fact.factory))
	cobraCmd.AddCommand(completion.NewCmd(fact.factory))
	cobraCmd.AddCommand(describe.NewCmd(fact.factory))
	cobraCmd.AddCommand(login.New(fact.factory))
	cobraCmd.AddCommand(logout.NewCmd(fact.factory))
	cobraCmd.AddCommand(create.NewCmd(fact.factory))
	cobraCmd.AddCommand(list.NewCmd(fact.factory))
	cobraCmd.AddCommand(delete.NewCmd(fact.factory))
	cobraCmd.AddCommand(update.NewCmd(fact.factory))
	cobraCmd.AddCommand(version.NewCmd(fact.factory))
	cobraCmd.AddCommand(whoami.NewCmd(fact.factory))
	cobraCmd.AddCommand(purge.NewCmd(fact.factory))
	cobraCmd.AddCommand(reset.NewCmd(fact.factory))
	cobraCmd.AddCommand(sync.NewCmd(fact.factory))
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

	// set template for -v flag
	cobraCmd.SetVersionTemplate(color.New(color.Bold).Sprint("Azion CLI " + version.BinVersion + "\n"))

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
			errMetrics := metric.TotalCommandsCount(cmd, f.commandName, executionTime, err)
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
