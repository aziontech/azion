package deploy

import (
	"context"
	"encoding/json"
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/deploy"
	apidom "github.com/aziontech/azion-cli/pkg/api/domains"
	apiapp "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
	apipurge "github.com/aziontech/azion-cli/pkg/api/realtime_purge"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

func (cmd *DeployCmd) doFunction(client *api.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) error {
	if conf.Function.Id == 0 {
		//Create New function
		PublishId, err := cmd.createFunction(client, ctx, conf)
		if err != nil {
			logger.Debug("Error while creating edge functions", zap.Error(err))
			return err
		}

		conf.Function.Id = PublishId
	} else {
		//Update existing function
		_, err := cmd.updateFunction(client, ctx, conf)
		if err != nil {
			logger.Debug("Error while updating edge functions", zap.Error(err))
			return err
		}
	}
	return nil
}

func (cmd *DeployCmd) doApplication(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) error {
	if conf.Application.Id == 0 {
		applicationId, _, err := cmd.createApplication(client, ctx, conf)
		if err != nil {
			logger.Debug("Error while creating edge application", zap.Error(err))
			return err
		}
		conf.Application.Id = applicationId

		err = cmd.WriteAzionJsonContent(conf)
		if err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
			return err
		}

		//TODO: Review what to do when user updates Function ID directly in azion.json
		err = cmd.updateRulesEngine(client, ctx, conf)
		if err != nil {
			logger.Debug("Error while updating rules engine", zap.Error(err))
			return err
		}
	} else {
		err := cmd.updateApplication(client, ctx, conf)
		if err != nil {
			logger.Debug("Error while updating edge application", zap.Error(err))
			return err
		}
	}
	return nil
}

func (cmd *DeployCmd) doDomain(client *apidom.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) error {
	var domain apidom.DomainResponse
	var err error

	newDomain := false
	if conf.Domain.Id == 0 {
		domain, err = cmd.createDomain(client, ctx, conf)
		if err != nil {
			logger.Debug("Error while creating domain", zap.Error(err))
			return err
		}
		conf.Domain.Id = domain.GetId()
		newDomain = true

	} else {
		domain, err = cmd.updateDomain(client, ctx, conf)
		if err != nil {
			logger.Debug("Error while updating domain", zap.Error(err))
			return err
		}
	}

	domainReturnedName := []string{domain.GetDomainName()}

	if conf.RtPurge.PurgeOnPublish && !newDomain {
		err := cmd.purgeDomains(cmd.F, domainReturnedName)
		if err != nil {
			logger.Debug("Error while purging domain", zap.Error(err))
			return err
		}
	}
	return nil
}

func (cmd *DeployCmd) doOrigin(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) error {
	if conf.Origin.Id == 0 {
		err := cmd.createAppRequirements(client, ctx, conf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cmd *DeployCmd) createFunction(client *api.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) (int64, error) {
	reqCre := api.CreateRequest{}

	if conf.Template == "static" {
		code, err := cmd.applyTemplate(conf)
		if err != nil {
			return 0, err
		}
		reqCre.SetCode(code)
	} else {
		//Read code to upload
		code, err := cmd.FileReader(conf.Function.File)
		if err != nil {
			logger.Debug("Error while reading edge function file <"+conf.Function.File+">", zap.Error(err))
			return 0, fmt.Errorf("%s: %w", msg.ErrorCodeFlag, err)
		}

		reqCre.SetCode(string(code))
	}

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
		logger.Debug("Error while creating edge function", zap.Error(err))
		return 0, fmt.Errorf(msg.ErrorCreateFunction.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputEdgeFunctionCreate, response.GetName(), response.GetId()))
	return response.GetId(), nil
}

func (cmd *DeployCmd) updateFunction(client *api.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) (int64, error) {
	reqUpd := api.UpdateRequest{}

	if conf.Template == "static" {
		code, err := cmd.applyTemplate(conf)
		if err != nil {
			return 0, err
		}
		reqUpd.SetCode(code)
	} else {
		//Read code to upload
		code, err := cmd.FileReader(conf.Function.File)
		if err != nil {
			logger.Debug("Error while reading edge function file <"+conf.Function.File+">", zap.Error(err))
			return 0, fmt.Errorf("%s: %w", msg.ErrorCodeFlag, err)
		}

		reqUpd.SetCode(string(code))
	}

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

	reqUpd.Id = conf.Function.Id
	reqUpd.SetJsonArgs(args)
	response, err := client.Update(ctx, &reqUpd)
	if err != nil {
		return 0, fmt.Errorf(msg.ErrorUpdateFunction.Error(), err)
	}

	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputEdgeFunctionUpdate, response.GetName(), conf.Function.Id))
	return response.GetId(), nil
}

func (cmd *DeployCmd) createApplication(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) (int64, int64, error) {
	reqApp := apiapp.CreateRequest{}
	if conf.Application.Name == "__DEFAULT__" {
		reqApp.SetName(conf.Name)
	} else {
		reqApp.SetName(conf.Application.Name)
	}
	reqApp.SetDeliveryProtocol("http,https")
	application, err := client.Create(ctx, &reqApp)
	if err != nil {
		return 0, 0, fmt.Errorf(msg.ErrorCreateApplication.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputEdgeApplicationCreate, application.GetName(), application.GetId()))
	reqUpApp := apiapp.UpdateRequest{}
	reqUpApp.SetEdgeFunctions(true)
	reqUpApp.SetApplicationAcceleration(true)
	reqUpApp.Id = application.GetId()
	application, err = client.Update(ctx, &reqUpApp)
	if err != nil {
		logger.Debug("Error while setting up edge application", zap.Error(err))
		return 0, 0, fmt.Errorf(msg.ErrorUpdateApplication.Error(), err)
	}
	reqIns := apiapp.CreateInstanceRequest{}
	reqIns.SetEdgeFunctionId(conf.Function.Id)
	reqIns.SetName(conf.Name)
	reqIns.ApplicationId = application.GetId()
	instance, err := client.CreateInstancePublish(ctx, &reqIns)
	if err != nil {
		logger.Debug("Error while creating edge function instance", zap.Error(err))
		return 0, 0, fmt.Errorf(msg.ErrorCreateInstance.Error(), err)
	}
	InstanceId = instance.GetId()
	return application.GetId(), instance.GetId(), nil
}

func (cmd *DeployCmd) updateApplication(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) error {
	reqApp := apiapp.UpdateRequest{}
	if conf.Application.Name == "__DEFAULT__" {
		reqApp.SetName(conf.Name)
	} else {
		reqApp.SetName(conf.Application.Name)
	}
	reqApp.Id = conf.Application.Id
	application, err := client.Update(ctx, &reqApp)
	if err != nil {
		return fmt.Errorf(msg.ErrorUpdateApplication.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputEdgeApplicationUpdate, application.GetName(), application.GetId()))
	return nil
}

func (cmd *DeployCmd) updateRulesEngine(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) error {
	reqRules := apiapp.UpdateRulesEngineRequest{}
	reqRules.IdApplication = conf.Application.Id

	_, err := client.UpdateRulesEnginePublish(ctx, &reqRules, InstanceId)
	if err != nil {
		return err
	}

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

	logger.FInfo(cmd.F.IOStreams.Out, msg.EdgeApplicationsPublishOutputCachePurge)
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
	reqDom.SetEdgeApplicationId(conf.Application.Id)
	domain, err := client.Create(ctx, &reqDom)
	if err != nil {
		return nil, fmt.Errorf(msg.ErrorCreateDomain.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputDomainCreate, conf.Name, domain.GetId()))
	return domain, nil
}

func (cmd *DeployCmd) updateDomain(client *apidom.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) (apidom.DomainResponse, error) {
	reqDom := apidom.UpdateRequest{}
	if conf.Domain.Name == "__DEFAULT__" {
		reqDom.SetName(conf.Name)
	} else {
		reqDom.SetName(conf.Domain.Name)
	}
	reqDom.SetEdgeApplicationId(conf.Application.Id)
	reqDom.Id = conf.Domain.Id
	domain, err := client.Update(ctx, &reqDom)
	if err != nil {
		return nil, fmt.Errorf(msg.ErrorUpdateDomain.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputDomainUpdate, conf.Name, domain.GetId()))
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

func (cmd *DeployCmd) createAppRequirements(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) error {
	reqOrigin := apiapp.CreateOriginsRequest{}
	var addresses []string
	if len(conf.Origin.Address) > 0 {
		address := prepareAddresses(conf.Origin.Address)
		addresses = conf.Origin.Address
		reqOrigin.SetAddresses(address)
	} else {
		addresses := prepareAddresses(DEFAULTORIGIN[:])
		reqOrigin.SetAddresses(addresses)
	}
	reqOrigin.SetName(conf.Name)
	reqOrigin.SetHostHeader("${host}")
	origin, err := client.CreateOrigins(ctx, conf.Application.Id, &reqOrigin)
	if err != nil {
		logger.Debug("Error while creating origin", zap.Error(err))
		return err
	}
	conf.Origin.Id = origin.GetOriginId()
	conf.Origin.Address = addresses
	conf.Origin.Name = origin.GetName()
	reqCache := apiapp.CreateCacheSettingsRequest{}
	reqCache.SetName(conf.Name)
	cache, err := client.CreateCacheSettingsNextApplication(ctx, &reqCache, conf.Application.Id)
	if err != nil {
		logger.Debug("Error while creating cache settings for Nextjs application", zap.Error(err))
		return err
	}
	logger.FInfo(cmd.F.IOStreams.Out, msg.EdgeApplicationsCacheSettingsSuccessful)
	err = client.CreateRulesEngineNextApplication(ctx, conf.Application.Id, cache.GetId(), conf.Template)
	if err != nil {
		logger.Debug("Error while creating rules engine for Nextjs application", zap.Error(err))
		return err
	}
	logger.FInfo(cmd.F.IOStreams.Out, msg.EdgeApplicationsRulesEngineSuccessful)

	return nil
}
