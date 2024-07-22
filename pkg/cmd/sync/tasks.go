package sync

import (
	"context"
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/sync"
	edgeApp "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/api/origin"
	varApi "github.com/aziontech/azion-cli/pkg/api/variables"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	opts *contracts.ListOptions
	ctx  context.Context = context.Background()
)

func SyncLocalResources(f *cmdutil.Factory, info contracts.SyncOpts, synch *SyncCmd) error {
	opts = &contracts.ListOptions{
		PageSize: 1000,
		Page:     1,
	}

	var err error
	err = synch.syncRules(info, f)
	if err != nil {
		return fmt.Errorf(msg.ERRORSYNC, err.Error())
	}

	err = synch.syncCache(info, f)
	if err != nil {
		return fmt.Errorf(msg.ERRORSYNC, err.Error())
	}

	err = synch.syncOrigin(info, f)
	if err != nil {
		return fmt.Errorf(msg.ERRORSYNC, err.Error())
	}

	err = synch.syncEnv(f)
	if err != nil {
		return fmt.Errorf(msg.ERRORSYNC, err.Error())
	}

	return nil
}

func (synch *SyncCmd) syncOrigin(info contracts.SyncOpts, f *cmdutil.Factory) error {
	client := origin.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	resp, err := client.ListOrigins(ctx, opts, info.Conf.Application.ID)
	if err != nil {
		return err
	}
	for _, origin := range resp.Results {
		if r := info.OriginIds[origin.Name]; r.OriginId > 0 {
			continue
		}
		newOrigin := contracts.AzionJsonDataOrigin{
			OriginId:  origin.GetOriginId(),
			OriginKey: origin.GetOriginKey(),
			Name:      origin.GetName(),
		}
		info.Conf.Origin = append(info.Conf.Origin, newOrigin)
		err := synch.WriteAzionJsonContent(info.Conf, ProjectConf)
		if err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
			return err
		}
		logger.FInfoFlags(
			synch.Io.Out,
			fmt.Sprintf(msg.SYNCMESSAGEORIGIN, origin.Name),
			synch.F.Format,
			synch.F.Out)
	}
	return nil
}

func (synch *SyncCmd) syncCache(info contracts.SyncOpts, f *cmdutil.Factory) error {
	client := edgeApp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	resp, err := client.ListCacheEdgeApp(context.Background(), info.Conf.Application.ID)
	if err != nil {
		return err
	}
	for _, cache := range resp {
		if r := info.CacheIds[cache.Name]; r.Id > 0 {
			continue
		}
		newCache := contracts.AzionJsonDataCacheSettings{
			Id:   cache.GetId(),
			Name: cache.GetName(),
		}
		info.Conf.CacheSettings = append(info.Conf.CacheSettings, newCache)
		err := synch.WriteAzionJsonContent(info.Conf, ProjectConf)
		if err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
			return err
		}
		logger.FInfoFlags(synch.Io.Out, fmt.Sprintf(msg.SYNCMESSAGECACHE, cache.Name), synch.F.Format, synch.F.Out)
	}
	return nil
}

func (synch *SyncCmd) syncRules(info contracts.SyncOpts, f *cmdutil.Factory) error {
	client := edgeApp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	resp, err := client.ListRulesEngine(context.Background(), opts, info.Conf.Application.ID, "request")
	if err != nil {
		return err
	}

	for _, rule := range resp.Results {
		if r := info.RuleIds[rule.Name]; r.Id > 0 || rule.Name == "Default Rule" {
			// if remote rule is also on local environment, no action is needed
			continue
		}
		newRule := contracts.AzionJsonDataRules{
			Id:    rule.GetId(),
			Name:  rule.GetName(),
			Phase: rule.GetPhase(),
		}
		info.Conf.RulesEngine.Rules = append(info.Conf.RulesEngine.Rules, newRule)
		err := synch.WriteAzionJsonContent(info.Conf, ProjectConf)
		if err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
			return err
		}
		logger.FInfoFlags(
			synch.Io.Out, fmt.Sprintf(msg.SYNCMESSAGERULE, rule.Name), synch.F.Format, synch.F.Out)
	}
	return nil
}

func (synch *SyncCmd) syncEnv(f *cmdutil.Factory) error {

	client := varApi.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	resp, err := client.List(context.Background())
	if err != nil {
		return err
	}

	// Load the .env file
	envs, err := godotenv.Read(synch.EnvPath)
	if err != nil {
		logger.Debug("Error while loading .env file", zap.Error(err))
		return nil // not every project has a .env file... this should not stop the execution
	}

	for _, variable := range resp {
		if v := envs[variable.GetKey()]; v != "" {
			delete(envs, variable.GetKey())
		}
	}

	for key, value := range envs {
		createReq := &varApi.Request{}
		createReq.Key = key
		createReq.Value = value
		_, err := client.Create(ctx, *createReq)
		if err != nil {
			logger.Debug("Error while creating variables during sync process", zap.Error(err))
			return err
		}
		logger.FInfoFlags(
			synch.Io.Out, fmt.Sprintf(msg.SYNCMESSAGEENV, key), synch.F.Format, synch.F.Out)
	}
	return nil
}
