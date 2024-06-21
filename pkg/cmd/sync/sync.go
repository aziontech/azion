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

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmdFactory := NewDevCmd(f)
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
		RunE: runE(cmdFactory),
	}
	syncCmd.Flags().BoolP("help", "h", false, msg.HELPFLAG)
	syncCmd.Flags().StringVar(&ProjectConf, "config-dir", "azion", msg.CONFDIRFLAG)
	return syncCmd
}

func runE(cmdFac *SyncCmd) func(cmd *cobra.Command, _ []string) error {
	return func(cmd *cobra.Command, _ []string) error {
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

		info := contracts.SyncOpts{
			RuleIds:   ruleIds,
			OriginIds: originIds,
			Conf:      conf,
		}

		err = cmdFac.SyncResources(cmdFac.F, info)
		if err != nil {
			return err
		}

		return nil
	}
}
