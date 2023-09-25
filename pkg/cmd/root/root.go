package root

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/root"
	buildCmd "github.com/aziontech/azion-cli/pkg/cmd/build"
	"github.com/aziontech/azion-cli/pkg/cmd/completion"
	"github.com/aziontech/azion-cli/pkg/cmd/describe"

	// "github.com/aziontech/azion-cli/pkg/cmd/create"
	deploycmd "github.com/aziontech/azion-cli/pkg/cmd/deploy"
	devcmd "github.com/aziontech/azion-cli/pkg/cmd/dev"
	initcmd "github.com/aziontech/azion-cli/pkg/cmd/init"
	linkcmd "github.com/aziontech/azion-cli/pkg/cmd/link"
	personal_token "github.com/aziontech/azion-cli/pkg/cmd/personal-token"
	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/constants"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type RootCmd struct {
	F       *cmdutil.Factory
	InitCmd func(f *cmdutil.Factory) *initcmd.InitCmd
}

func NewRootCmd(f *cmdutil.Factory) *RootCmd {
	return &RootCmd{
		F:       f,
		InitCmd: initcmd.NewInitCmd,
	}
}

var (
	tokenFlag  string
	configFlag string
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewRootCmd(f), f)
}

func NewCobraCmd(rootCmd *RootCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:     msg.RootUsage,
		Long:    msg.RootDescription,
		Short:   color.New(color.Bold).Sprint(fmt.Sprintf(msg.RootDescription, version.BinVersion)),
		Version: version.BinVersion,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			logger.LogLevel(f.Logger)
			err := doPreCommandCheck(cmd, f, PreCmd{
				config: configFlag,
				token:  tokenFlag,
			})
			if err != nil {
				return err
			}
			return nil
		},
		Example: heredoc.Doc(`
		$ azion
		$ azion -t azionb43a9554776zeg05b11cb1declkbabcc9la
		$ azion --debug
		$ azion -h
		`),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if cmd.Flags().Changed("token") {
				return nil
			}
			return rootCmd.Run()
		},
		SilenceErrors: true, // Silence errors, so the help message won't be shown on flag error
		SilenceUsage:  true, // Silence usage on error
	}

	cobraCmd.SetIn(f.IOStreams.In)
	cobraCmd.SetOut(f.IOStreams.Out)
	cobraCmd.SetErr(f.IOStreams.Err)

	cobraCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		rootHelpFunc(f, cmd, args)
	})

	// Global flags
	cobraCmd.PersistentFlags().StringVarP(&tokenFlag, "token", "t", "", msg.RootTokenFlag)
	cobraCmd.PersistentFlags().StringVarP(&configFlag, "config", "c", "", msg.RootConfigFlag)
	cobraCmd.PersistentFlags().BoolVarP(&f.GlobalFlagAll, "yes", "y", false, msg.RootYesFlag)
	cobraCmd.PersistentFlags().BoolVarP(&f.Debug, "debug", "d", false, msg.RootLogDebug)
	cobraCmd.PersistentFlags().BoolVarP(&f.Silent, "silent", "s", false, msg.RootLogSilent)
	cobraCmd.PersistentFlags().StringVarP(&f.LogLevel, "log-level", "l", "info", msg.RootLogDebug)

	// other flags
	cobraCmd.Flags().BoolP("help", "h", false, msg.RootHelpFlag)

	// set template for -v flag
	cobraCmd.SetVersionTemplate(color.New(color.Bold).Sprint("Azion CLI " + version.BinVersion + "\n"))

	cobraCmd.AddCommand(initcmd.NewCmd(f))
	cobraCmd.AddCommand(deploycmd.NewCmd(f))
	cobraCmd.AddCommand(buildCmd.NewCmd(f))
	cobraCmd.AddCommand(devcmd.NewCmd(f))
	cobraCmd.AddCommand(linkcmd.NewCmd(f))
	cobraCmd.AddCommand(personal_token.NewCmd(f))
	cobraCmd.AddCommand(completion.NewCmd(f))
	// cobraCmd.AddCommand(create.NewCmd(f))
	cobraCmd.AddCommand(describe.NewCmd(f))

	return cobraCmd
}

func (cmd *RootCmd) Run() error {
	logger.Debug("Running root command")
	info := &initcmd.InitInfo{}
	init := cmd.InitCmd(cmd.F)
	err := init.Run(info)
	if err != nil {
		logger.Debug("Error while running init command called by root command", zap.Error(err))
		return err
	}

	return nil
}

func Execute() {
	streams := iostreams.System()
	httpClient := &http.Client{
		Timeout: 10 * time.Second, // TODO: Configure this somewhere
	}

	// TODO: Ignoring errors since the file might not exist, maybe warn the user?
	tok, _ := token.ReadFromDisk()
	viper.SetEnvPrefix("AZIONCLI")
	viper.AutomaticEnv()
	viper.SetDefault("token", tok)
	viper.SetDefault("api_url", constants.ApiURL)
	viper.SetDefault("storage_url", constants.StorageApiURL)

	factory := &cmdutil.Factory{
		HttpClient: httpClient,
		IOStreams:  streams,
		Config:     viper.GetViper(),
	}

	cmd := NewCmd(factory)

	cobra.CheckErr(cmd.Execute())
}
