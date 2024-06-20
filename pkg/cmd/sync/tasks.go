package sync

import (
	"context"
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/sync"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

var (
	C    context.Context
	Opts *contracts.ListOptions
)

func (synch *SyncCmd) SyncResources(f *cmdutil.Factory, info contracts.SyncOpts) error {
	C = context.Background()
	Opts = &contracts.ListOptions{
		PageSize: 1000,
		Page:     1,
	}
	err := synch.syncRules(info, f)
	if err != nil {
		return fmt.Errorf(msg.ERRORSYNC, err.Error())
	}
	return nil
}

func (synch *SyncCmd) syncRules(info contracts.SyncOpts, f *cmdutil.Factory) error {

	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	resp, err := client.ListRulesEngine(C, Opts, info.Conf.Application.ID, "request")
	if err != nil {
		return err
	}

	for _, rule := range resp.Results {
		if r := info.RuleIds[rule.Name]; r.Id > 0 || rule.Name == "Default Rule" {
			// if remote rule is also on local environment, no action is needed
			continue
		}
		newRule := contracts.AzionJsonDataRules{
			Id:   rule.GetId(),
			Name: rule.GetName(),
		}
		info.Conf.RulesEngine.Rules = append(info.Conf.RulesEngine.Rules, newRule)
		err := synch.WriteAzionJsonContent(info.Conf, ProjectConf)
		if err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
			return err
		}
		logger.FInfo(synch.Io.Out, fmt.Sprintf(msg.SYNCMESSAGERULE, rule.Name))
	}
	return nil
}
