package root

import (
	"fmt"
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
	deployremote "github.com/aziontech/azion-cli/pkg/cmd/deploy_remote"
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
	*cmdutil.Factory
	doPreCommandCheck func(cmd *cobra.Command, fact *factoryRoot) error //this package
	execSchedules     func(factory *cmdutil.Factory)                    //schedule.ExecShedules
	flags
	globals
}

func NewFactoryRoot(fact *cmdutil.Factory) *factoryRoot {
	return &factoryRoot{
		Factory:           fact,
		doPreCommandCheck: doPreCommandCheck,
		execSchedules:     schedule.ExecSchedules,
	}
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
	logger.LogLevel(fact.Factory.Logger)

	if strings.HasPrefix(fact.configFlag, PREFIX_FLAG) {
		return msg.ErrorPrefix
	}

	if err := fact.doPreCommandCheck(cmd, fact); err != nil {
		return err
	}

	fact.execSchedules(fact.Factory)
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
	cobraCmd.PersistentFlags().BoolVarP(&fact.Factory.Debug, "debug", "d", false, msg.RootLogDebug)
	cobraCmd.PersistentFlags().BoolVarP(&fact.Factory.Silent, "silent", "s", false, msg.RootLogSilent)
	cobraCmd.PersistentFlags().StringVarP(&fact.Factory.LogLevel, "log-level", "l", "info", msg.RootLogLevel)
	cobraCmd.PersistentFlags().BoolVarP(&fact.Factory.GlobalFlagAll, "yes", "y", false, msg.RootYesFlag)
	cobraCmd.PersistentFlags().StringVar(&fact.Factory.Out, "out", "", msg.RootFlagOut)
	cobraCmd.PersistentFlags().StringVar(&fact.Factory.Format, "format", "", msg.RootFlagFormat)
	cobraCmd.PersistentFlags().BoolVar(&fact.Factory.NoColor, "no-color", false, msg.RootFlagFormat)
	cobraCmd.Flags().BoolP("help", "h", false, msg.RootHelpFlag)
}

func (fact *factoryRoot) setCmds(cobraCmd *cobra.Command) {
	cobraCmd.AddCommand(initcmd.NewCmd(fact.Factory))
	cobraCmd.AddCommand(logcmd.NewCmd(fact.Factory))
	cobraCmd.AddCommand(deploycmd.NewCmd(fact.Factory))
	cobraCmd.AddCommand(buildCmd.NewCmd(fact.Factory))
	cobraCmd.AddCommand(devcmd.NewCmd(fact.Factory))
	cobraCmd.AddCommand(linkcmd.NewCmd(fact.Factory))
	cobraCmd.AddCommand(unlink.NewCmd(fact.Factory))
	cobraCmd.AddCommand(completion.NewCmd(fact.Factory))
	cobraCmd.AddCommand(describe.NewCmd(fact.Factory))
	cobraCmd.AddCommand(login.NewCmd(fact.Factory))
	cobraCmd.AddCommand(logout.NewCmd(fact.Factory))
	cobraCmd.AddCommand(create.NewCmd(fact.Factory))
	cobraCmd.AddCommand(list.NewCmd(fact.Factory))
	cobraCmd.AddCommand(delete.NewCmd(fact.Factory))
	cobraCmd.AddCommand(update.NewCmd(fact.Factory))
	cobraCmd.AddCommand(version.NewCmd(fact.Factory))
	cobraCmd.AddCommand(whoami.NewCmd(fact.Factory))
	cobraCmd.AddCommand(purge.NewCmd(fact.Factory))
	cobraCmd.AddCommand(reset.NewCmd(fact.Factory))
	cobraCmd.AddCommand(sync.NewCmd(fact.Factory))
	cobraCmd.AddCommand(deployremote.NewCmd(fact.Factory))
}

func CmdRoot(fact *factoryRoot) *cobra.Command {
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

	cobraCmd.SetIn(fact.Factory.IOStreams.In)
	cobraCmd.SetOut(fact.Factory.IOStreams.Out)
	cobraCmd.SetErr(fact.Factory.IOStreams.Err)

	cobraCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		rootHelpFunc(cmd, args)
	})

	fact.setFlags(cobraCmd)

	// set template for -v flag
	cobraCmd.SetVersionTemplate(color.New(color.Bold).Sprint("Azion CLI " + version.BinVersion + "\n"))

	fact.setCmds(cobraCmd)
	return cobraCmd
}

func Execute(factUtil *cmdutil.Factory) {
	logger.New(zapcore.InfoLevel)

	fact := NewFactoryRoot(factUtil)
	cmd := CmdRoot(fact)
	err := cmd.Execute()
	executionTime := time.Since(fact.startTime).Seconds()

	// 1 = authorize; anything different than 1 means that the user did not authorize metrics collection, or did not answer the question yet
	if fact.globalSettings != nil {
		if fact.globalSettings.AuthorizeMetricsCollection == 1 {
			errMetrics := metric.TotalCommandsCount(cmd, fact.commandName, executionTime, err)
			if errMetrics != nil {
				logger.Debug("Error while saving metrics", zap.Error(err))
			}
		}
	}

	if err != nil {
		output.Print(&output.ErrorOutput{
			GeneralOutput: output.GeneralOutput{
				Out:   factUtil.IOStreams.Out,
				Flags: factUtil.Flags,
			},
			Err: err,
		})
	}
}
