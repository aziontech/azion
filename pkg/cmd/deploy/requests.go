package deploy

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/deploy"
	apidom "github.com/aziontech/azion-cli/pkg/api/domain"
	apiapp "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	api "github.com/aziontech/azion-cli/pkg/api/edge_function"
	apiori "github.com/aziontech/azion-cli/pkg/api/origin"
	apipurge "github.com/aziontech/azion-cli/pkg/api/realtime_purge"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
)

// inject this code into worker.js
var injectIntoFunction = `
//---
//storages:
//   - name: assets
//     bucket: %s
//     prefix: %s
//---

`

func (cmd *DeployCmd) doFunction(clients *Clients, ctx context.Context, conf *contracts.AzionApplicationOptions) error {
	if conf.Function.ID == 0 {
		DeployID, err := cmd.createFunction(clients.EdgeFunction, ctx, conf)
		if err != nil {
			return err
		}
		conf.Function.ID = DeployID

		// create instance function
		reqIns := apiapp.CreateInstanceRequest{}
		reqIns.SetEdgeFunctionId(conf.Function.ID)
		reqIns.SetName(conf.Name)
		reqIns.ApplicationId = conf.Application.ID

		instance, err := clients.EdgeApplication.CreateInstancePublish(ctx, &reqIns)
		if err != nil {
			logger.Debug("Error while creating edge function instance", zap.Error(err))
			return fmt.Errorf(msg.ErrorCreateInstance.Error(), err)
		}
		conf.Function.InstanceID = instance.GetId()
	}

	_, err := cmd.updateFunction(clients.EdgeFunction, ctx, conf)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *DeployCmd) doApplication(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) error {
	if conf.Application.ID == 0 {
		applicationId, err := cmd.createApplication(client, ctx, conf)
		if err != nil {
			logger.Debug("Error while creating Edge Application", zap.Error(err))
			return err
		}
		conf.Application.ID = applicationId

		err = cmd.WriteAzionJsonContent(conf)
		if err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
			return err
		}
	} else {
		err := cmd.updateApplication(client, ctx, conf)
		if err != nil {
			logger.Debug("Error while updating Edge Application", zap.Error(err))
			return err
		}
	}
	return nil
}

func (cmd *DeployCmd) doDomain(client *apidom.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) (string, error) {
	var domain apidom.DomainResponse
	var err error

	newDomain := false
	if conf.Domain.Id == 0 {
		domain, err = cmd.createDomain(client, ctx, conf)
		if err != nil {
			logger.Debug("Error while creating domain", zap.Error(err))
			return "", err
		}
		conf.Domain.Id = domain.GetId()
		newDomain = true

		err = cmd.WriteAzionJsonContent(conf)
		if err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
			return "", err
		}

	} else {
		domain, err = cmd.updateDomain(client, ctx, conf)
		if err != nil {
			logger.Debug("Error while updating domain", zap.Error(err))
			return "", err
		}
	}

	domainReturnedName := []string{domain.GetDomainName()}

	if conf.RtPurge.PurgeOnPublish && !newDomain {
		err := cmd.purgeDomains(cmd.F, domainReturnedName)
		if err != nil {
			logger.Debug("Error while purging domain", zap.Error(err))
			return "", err
		}
	}

	return domainReturnedName[0], nil
}

func (cmd *DeployCmd) doOrigin(client *apiapp.Client, clientOrigin *apiori.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) error {
	if conf.Template == "javascript" || conf.Template == "typescript" {
		return nil
	}

	var addresses []string
	var DefaultOrigin = [1]string{"httpbin.org"}

	if strings.ToLower(conf.Mode) == "compute" {
		if conf.Origin.SingleOriginID == 0 {
			reqSingleOrigin := apiori.CreateRequest{}

			if len(conf.Origin.Address) > 0 {
				address := prepareAddresses(conf.Origin.Address)
				reqSingleOrigin.SetAddresses(address)
			} else {
				addresses := prepareAddresses(DefaultOrigin[:])
				reqSingleOrigin.SetAddresses(addresses)
			}

			reqSingleOrigin.SetName(utils.Concat(conf.Name, "_single"))
			reqSingleOrigin.SetHostHeader("${host}")

			origin, err := clientOrigin.Create(ctx, conf.Application.ID, &reqSingleOrigin)
			if err != nil {
				logger.Debug("Error while creating default origin ", zap.Any("Error", err))
				return err
			}
			logger.FInfo(cmd.F.IOStreams.Out, msg.OriginsSuccessful)

			conf.Origin.SingleOriginID = origin.GetOriginId()
		}
	}

	if conf.Origin.StorageOriginID == 0 {
		reqObjectStorageOrigin := apiori.CreateRequest{}
		reqObjectStorageOrigin.SetName(utils.Concat(conf.Name, "_object"))
		reqObjectStorageOrigin.SetOriginType("object_storage")
		reqObjectStorageOrigin.Bucket = &conf.Bucket
		reqObjectStorageOrigin.Prefix = &conf.Prefix

		origin, err := clientOrigin.Create(ctx, conf.Application.ID, &reqObjectStorageOrigin)
		if err != nil {
			logger.Debug("Error while creating origin of type object storage", zap.Any("Error", err))
			return err
		}
		logger.FInfo(cmd.F.IOStreams.Out, msg.OriginsSuccessful)

		conf.Origin.StorageOriginID = origin.GetOriginId()
		conf.Origin.StorageOriginKey = origin.GetOriginKey()
		conf.Origin.Address = addresses
		conf.Origin.Name = origin.GetName()

		var reqCache apiapp.CreateCacheSettingsRequest
		reqCache.SetName(conf.Name)

		// create cache settings
		cache, err := client.CreateCacheSettingsNextApplication(ctx, &reqCache, conf.Application.ID)
		if err != nil {
			logger.Debug("Error while creating cache settings", zap.Error(err))
			return err
		}
		logger.FInfo(cmd.F.IOStreams.Out, msg.CacheSettingsSuccessful)

		// creates gzip and cache rules
		err = client.CreateRulesEngineNextApplication(ctx, conf.Application.ID, cache.GetId(), conf.Template, conf.Mode)
		if err != nil {
			logger.Debug("Error while creating rules engine", zap.Error(err))
			return err
		}
		logger.FInfo(cmd.F.IOStreams.Out, msg.RulesEngineSuccessful)
	} else {
		reqObjectStorageOrigin := apiori.UpdateRequest{}
		reqObjectStorageOrigin.Prefix = &conf.Prefix

		_, err := clientOrigin.Update(ctx, conf.Application.ID, conf.Origin.StorageOriginKey, &reqObjectStorageOrigin)
		if err != nil {
			logger.Debug("Error while updating origin of type object storage", zap.Any("Error", err))
			return err
		}
		logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.OriginsUpdateSuccessful, conf.Application.ID, conf.Origin.StorageOriginID))
	}
	return nil
}

func (cmd *DeployCmd) createFunction(client *api.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) (int64, error) {
	reqCre := api.CreateRequest{}

	code, err := cmd.FileReader(conf.Function.File)
	if err != nil {
		logger.Debug("Error while reading edge function file <"+conf.Function.File+">", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorCodeFlag, err)
	}

	prependText := fmt.Sprintf(injectIntoFunction, conf.Bucket, conf.Prefix)
	newCode := append([]byte(prependText), code...)

	reqCre.SetCode(string(newCode))

	reqCre.SetActive(true)
	if conf.Function.Name == "__DEFAULT__" {
		reqCre.SetName(conf.Name)
	} else {
		reqCre.SetName(conf.Function.Name)
	}

	//Read args
	marshalledArgs, err := cmd.FileReader(conf.Function.Args)
	if err != nil {
		logger.Debug("Error while reding args.json file <"+conf.Function.Args+">", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorArgsFlag, err)
	}
	args := make(map[string]interface{})
	if err := json.Unmarshal(marshalledArgs, &args); err != nil {
		logger.Debug("Error while unmarshling args.json file <"+conf.Function.Args+">", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorParseArgs, err)
	}

	reqCre.SetJsonArgs(args)
	response, err := client.Create(ctx, &reqCre)
	if err != nil {
		logger.Debug("Error while creating Edge Function", zap.Error(err))
		return 0, fmt.Errorf(msg.ErrorCreateFunction.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.DeployOutputEdgeFunctionCreate, response.GetName(), response.GetId()))
	return response.GetId(), nil
}

func (cmd *DeployCmd) updateFunction(client *api.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) (int64, error) {
	reqUpd := api.UpdateRequest{}

	code, err := cmd.FileReader(conf.Function.File)
	if err != nil {
		logger.Debug("Error while reading edge function file <"+conf.Function.File+">", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorCodeFlag, err)
	}

	prependText := fmt.Sprintf(injectIntoFunction, conf.Bucket, conf.Prefix)
	newCode := append([]byte(prependText), code...)

	reqUpd.SetCode(string(newCode))

	reqUpd.SetActive(true)
	if conf.Function.Name == "__DEFAULT__" {
		reqUpd.SetName(conf.Name)
	} else {
		reqUpd.SetName(conf.Function.Name)
	}

	//Read args
	marshalledArgs, err := cmd.FileReader(conf.Function.Args)
	if err != nil {
		logger.Debug("Error while reading args.json file <"+conf.Function.Args+">", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorArgsFlag, err)
	}
	args := make(map[string]interface{})
	if err := json.Unmarshal(marshalledArgs, &args); err != nil {
		logger.Debug("Error while unmarshling args.json file <"+conf.Function.Args+">", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorParseArgs, err)
	}

	reqUpd.SetJsonArgs(args)
	response, err := client.Update(ctx, &reqUpd, conf.Function.ID)
	if err != nil {
		return 0, fmt.Errorf(msg.ErrorUpdateFunction.Error(), err)
	}

	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.DeployOutputEdgeFunctionUpdate, response.GetName(), conf.Function.ID))
	return response.GetId(), nil
}

func (cmd *DeployCmd) createApplication(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) (int64, error) {
	reqApp := apiapp.CreateRequest{}
	if conf.Application.Name == "__DEFAULT__" {
		reqApp.SetName(conf.Name)
	} else {
		reqApp.SetName(conf.Application.Name)
	}
	reqApp.SetDeliveryProtocol("http,https")

	application, err := client.Create(ctx, &reqApp)
	if err != nil {
		return 0, fmt.Errorf(msg.ErrorCreateApplication.Error(), err)
	}

	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.DeployOutputEdgeApplicationCreate, application.GetName(), application.GetId()))

	reqUpApp := apiapp.UpdateRequest{}
	reqUpApp.SetEdgeFunctions(true)
	reqUpApp.SetApplicationAcceleration(true)
	reqUpApp.Id = application.GetId()

	application, err = client.Update(ctx, &reqUpApp)
	if err != nil {
		logger.Debug("Error while setting up edge application", zap.Error(err))
		return 0, fmt.Errorf(msg.ErrorUpdateApplication.Error(), err)
	}

	return application.GetId(), nil
}

func (cmd *DeployCmd) updateApplication(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) error {
	reqApp := apiapp.UpdateRequest{}
	if conf.Application.Name == "__DEFAULT__" {
		reqApp.SetName(conf.Name)
	} else {
		reqApp.SetName(conf.Application.Name)
	}
	reqApp.Id = conf.Application.ID
	application, err := client.Update(ctx, &reqApp)
	if err != nil {
		return fmt.Errorf(msg.ErrorUpdateApplication.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.DeployOutputEdgeApplicationUpdate, application.GetName(), application.GetId()))
	return nil
}

func (cmd *DeployCmd) purgeDomains(f *cmdutil.Factory, domainNames []string) error {
	ctx := context.Background()
	clipurge := apipurge.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	err := clipurge.Purge(ctx, domainNames)
	if err != nil {
		logger.Debug("Error while purging domain", zap.Error(err))
		return err
	}

	logger.FInfo(cmd.F.IOStreams.Out, msg.DeployOutputCachePurge)
	return nil
}

func (cmd *DeployCmd) createDomain(client *apidom.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) (apidom.DomainResponse, error) {
	reqDom := apidom.CreateRequest{}
	if conf.Domain.Name == "__DEFAULT__" {
		reqDom.SetName(conf.Name)
	} else {
		reqDom.SetName(conf.Domain.Name)
	}
	reqDom.SetCnames([]string{})
	reqDom.SetCnameAccessOnly(false)
	reqDom.SetIsActive(true)
	reqDom.SetEdgeApplicationId(conf.Application.ID)
	domain, err := client.Create(ctx, &reqDom)
	if err != nil {
		return nil, fmt.Errorf(msg.ErrorCreateDomain.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.DeployOutputDomainCreate, conf.Name, domain.GetId()))
	return domain, nil
}

func (cmd *DeployCmd) updateDomain(client *apidom.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) (apidom.DomainResponse, error) {
	reqDom := apidom.UpdateRequest{}
	if conf.Domain.Name == "__DEFAULT__" {
		reqDom.SetName(conf.Name)
	} else {
		reqDom.SetName(conf.Domain.Name)
	}
	reqDom.SetEdgeApplicationId(conf.Application.ID)
	reqDom.Id = conf.Domain.Id
	domain, err := client.Update(ctx, &reqDom)
	if err != nil {
		return nil, fmt.Errorf(msg.ErrorUpdateDomain.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.DeployOutputDomainUpdate, conf.Name, domain.GetId()))
	return domain, nil
}

func prepareAddresses(addrs []string) (addresses []sdk.CreateOriginsRequestAddresses) {
	var addr sdk.CreateOriginsRequestAddresses
	for _, v := range addrs {
		addr.Address = v
		addresses = append(addresses, addr)
	}
	return
}
