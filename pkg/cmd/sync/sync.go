package sync

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/sync"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var ProjectConf string

type SyncCmd struct {
	Io                    *iostreams.IOStreams
	GetAzionJsonContent   func(confPath string) (*contracts.AzionApplicationOptions, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions, confPath string) error
	F                     *cmdutil.Factory
}

func NewDevCmd(f *cmdutil.Factory) *SyncCmd {
	return &SyncCmd{
		F:                     f,
		Io:                    f.IOStreams,
		GetAzionJsonContent:   utils.GetAzionJsonContent,
		WriteAzionJsonContent: utils.WriteAzionJsonContent,
	}
}

func NewCobraCmd(sync *SyncCmd) *cobra.Command {
	syncCmd := &cobra.Command{
		Use:           msg.USAGE,
		Short:         msg.SHORTDESCRIPTION,
		Long:          msg.LONGDESCRIPTION,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`       
        $ azion sync
        $ azion sync --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return sync.Run(sync.F)
		},
	}
	syncCmd.Flags().BoolP("help", "h", false, msg.HELPFLAG)
	syncCmd.Flags().StringVar(&ProjectConf, "--config-dir", "azion", msg.CONFDIRFLAG)
	return syncCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDevCmd(f))
}

func (cmd *SyncCmd) Run(f *cmdutil.Factory) error {
	logger.Debug("Running sync command")

	conf, err := cmd.GetAzionJsonContent(ProjectConf)
	if err != nil {
		logger.Debug("Failed to get Azion JSON content", zap.Error(err))
		return err
	}

	ruleIds := make(map[string]contracts.RuleIdsStruct)
	for _, ruleConf := range conf.RulesEngine.Rules {
		ruleIds[ruleConf.Name] = contracts.RuleIdsStruct{
			Id:    ruleConf.Id,
			Phase: ruleConf.Phase,
		}
	}

	info := contracts.SyncOpts{
		RuleIds: ruleIds,
		Conf:    conf,
	}

	err = cmd.SyncResources(f, info)
	if err != nil {
		return err
	}

	return nil
}
