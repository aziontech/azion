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

var (
	ProjectConf string
)

type SyncCmd struct {
	Io                    *iostreams.IOStreams
	GetAzionJsonContent   func(confPath string) (*contracts.AzionApplicationOptions, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions, confPath string) error
	F                     *cmdutil.Factory
	SyncResources         func(f *cmdutil.Factory, info contracts.SyncOpts, synch *SyncCmd) error
	EnvPath               string
}

func NewSync(f *cmdutil.Factory) *SyncCmd {
	return &SyncCmd{
		F:                     f,
		Io:                    f.IOStreams,
		GetAzionJsonContent:   utils.GetAzionJsonContent,
		WriteAzionJsonContent: utils.WriteAzionJsonContent,
		SyncResources:         SyncLocalResources,
	}
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmdFactory := NewSync(f)
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
		RunE: func(_ *cobra.Command, _ []string) error {
			return Sync(cmdFactory)
		},
	}
	syncCmd.Flags().BoolP("help", "h", false, msg.HELPFLAG)
	syncCmd.Flags().StringVar(&ProjectConf, "config-dir", "azion", msg.CONFDIRFLAG)
	syncCmd.Flags().StringVar(&cmdFactory.EnvPath, "env", ".edge/.env", msg.ENVFLAG)
	return syncCmd
}

func Sync(cmdFac *SyncCmd) error {
	logger.Debug("Running sync command")
	conf, err := cmdFac.GetAzionJsonContent(ProjectConf)
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

	originIds := make(map[string]contracts.AzionJsonDataOrigin)
	for _, itemOrigin := range conf.Origin {
		originIds[itemOrigin.Name] = contracts.AzionJsonDataOrigin{
			OriginId:  itemOrigin.OriginId,
			OriginKey: itemOrigin.OriginKey,
			Name:      itemOrigin.Name,
			Address:   itemOrigin.Address,
		}
	}

	cacheIds := make(map[string]contracts.AzionJsonDataCacheSettings)
	for _, itemCache := range conf.CacheSettings {
		cacheIds[itemCache.Name] = contracts.AzionJsonDataCacheSettings{
			Id:   itemCache.Id,
			Name: itemCache.Name,
		}
	}

	info := contracts.SyncOpts{
		RuleIds:   ruleIds,
		OriginIds: originIds,
		CacheIds:  cacheIds,
		Conf:      conf,
	}

	err = cmdFac.SyncResources(cmdFac.F, info, cmdFac)
	if err != nil {
		return err
	}

	return nil
}
