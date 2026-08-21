package manifest

import (
	apiApplications "github.com/aziontech/azion-cli/pkg/api/applications"
	apiCache "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	apiConnector "github.com/aziontech/azion-cli/pkg/api/connector"
	apiWorkloads "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/contracts"
	edgesdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
)

func transformEdgeConnectorRequest(connectorRequest edgesdk.ConnectorRequest) *apiConnector.UpdateRequest {
	if connectorRequest.ConnectorHTTPRequest != nil {
		request := &apiConnector.UpdateRequest{}
		bodyRequest := connectorRequest.ConnectorHTTPRequest
		atts := bodyRequest.Attributes
		body := edgesdk.PatchedConnectorHTTPRequest{}
		if bodyRequest.Active != nil {
			body.SetActive(*bodyRequest.Active)
		}

		if bodyRequest.Name != "" {
			body.SetName(bodyRequest.Name)
		}

		body.SetType(bodyRequest.Type)

		body.SetAttributes(atts)
		request.PatchedConnectorHTTPRequest = &body
		return request
	}

	if connectorRequest.ConnectorStorageRequest != nil {
		request := &apiConnector.UpdateRequest{}
		bodyRequest := connectorRequest.ConnectorStorageRequest
		body := edgesdk.PatchedConnectorRequest{}
		internalBody := edgesdk.PatchedConnectorStorageRequest{}
		if bodyRequest.Active != nil {
			internalBody.SetActive(*bodyRequest.Active)
		}

		if bodyRequest.Name != "" {
			internalBody.SetName(bodyRequest.Name)
		}

		internalBody.SetType(bodyRequest.Type)

		internalBody.SetAttributes(bodyRequest.Attributes)
		body.PatchedConnectorStorageRequest = &internalBody
		request.PatchedConnectorRequest = body
		return request
	}

	return &apiConnector.UpdateRequest{}
}

func transformWorkloadRequestUpdate(createRequest contracts.WorkloadManifest) *apiWorkloads.UpdateRequest {
	request := &apiWorkloads.UpdateRequest{}

	if createRequest.Name != "" {
		request.SetName(createRequest.Name)
	}
	if createRequest.Active != nil {
		request.SetActive(*createRequest.Active)
	}
	if len(createRequest.Domains) > 0 {
		request.SetDomains(createRequest.Domains)
	}
	if createRequest.Mtls != nil {
		request.SetMtls(*createRequest.Mtls)
	}
	if createRequest.Protocols != nil {
		request.SetProtocols(*createRequest.Protocols)
	}
	if createRequest.Tls != nil {
		request.SetTls(*createRequest.Tls)
	}

	return request
}

func transformWorkloadDeploymentRequestUpdate(updateRequest contracts.WorkloadDeployment, conf *contracts.AzionApplicationOptions) edgesdk.PatchedWorkloadDeploymentRequest {
	request := edgesdk.PatchedWorkloadDeploymentRequest{}

	if updateRequest.Name != "" {
		request.SetName(updateRequest.Name)
	}

	request.SetActive(updateRequest.Active)

	request.SetCurrent(updateRequest.Current)

	strategy := edgesdk.DeploymentStrategyDefaultDeploymentStrategyRequest{}
	attributes := edgesdk.DefaultDeploymentStrategyAttrsRequest{}

	if updateRequest.Strategy.Type != "" {
		strategy.SetType(updateRequest.Strategy.Type)
	}

	attributes.SetApplication(conf.Application.ID)
	strategy.SetAttributes(attributes)
	request.SetStrategy(strategy)

	return request
}

func transformWorkloadDeploymentRequestCreate(createRequest contracts.WorkloadDeployment, conf *contracts.AzionApplicationOptions) edgesdk.WorkloadDeploymentRequest {
	request := edgesdk.WorkloadDeploymentRequest{}

	if createRequest.Name != "" {
		request.SetName(createRequest.Name)
	}
	request.SetActive(createRequest.Active)

	request.SetCurrent(createRequest.Current)

	strategy := edgesdk.DeploymentStrategyDefaultDeploymentStrategyRequest{}
	attributes := edgesdk.DefaultDeploymentStrategyAttrsRequest{}

	if createRequest.Strategy.Type != "" {
		strategy.SetType(createRequest.Strategy.Type)
	}

	attributes.SetApplication(conf.Application.ID)
	strategy.SetAttributes(attributes)
	request.SetStrategy(strategy)

	return request
}

func transformWorkloadRequestCreate(createRequest contracts.WorkloadManifest, appid int64) *apiWorkloads.CreateRequest {
	request := &apiWorkloads.CreateRequest{}

	if createRequest.Name != "" {
		request.SetName(createRequest.Name)
	}
	if createRequest.Active != nil {
		request.SetActive(*createRequest.Active)
	}
	if len(createRequest.Domains) > 0 {
		request.SetDomains(createRequest.Domains)
	}
	if createRequest.Mtls != nil {
		request.SetMtls(*createRequest.Mtls)
	}
	if createRequest.Protocols != nil {
		request.SetProtocols(*createRequest.Protocols)
	}
	if createRequest.Tls != nil {
		request.SetTls(*createRequest.Tls)
	}

	return request
}

func transformEdgeApplicationRequestUpdate(edgeapprequest contracts.Applications) *apiApplications.UpdateRequest {
	request := &apiApplications.UpdateRequest{}

	if edgeapprequest.Active != nil {
		request.SetActive(*edgeapprequest.Active)
	}
	if edgeapprequest.Debug != nil {
		request.SetDebug(*edgeapprequest.Debug)
	}
	if edgeapprequest.Modules != nil {
		request.SetModules(*edgeapprequest.Modules)
	}

	if edgeapprequest.Name != "" {
		request.SetName(edgeapprequest.Name)
	}

	return request
}

func transformEdgeApplicationRequestCreate(edgeapprequest contracts.Applications) *apiApplications.CreateRequest {
	request := &apiApplications.CreateRequest{}

	if edgeapprequest.Active != nil {
		request.SetActive(*edgeapprequest.Active)
	}
	if edgeapprequest.Debug != nil {
		request.SetDebug(*edgeapprequest.Debug)
	}
	if edgeapprequest.Modules != nil {
		request.SetModules(*edgeapprequest.Modules)
	}

	if edgeapprequest.Name != "" {
		request.SetName(edgeapprequest.Name)
	}

	return request
}

func transformCacheRequest(cache contracts.ManifestCacheSetting) *apiCache.RequestUpdate {
	request := apiCache.RequestUpdate{}

	if cache.Name != "" {
		request.SetName(cache.Name)
	}

	if cache.BrowserCache != nil {
		request.SetBrowserCache(*cache.BrowserCache)
	}
	if cache.Modules != nil {
		request.SetModules(*cache.Modules)
	}

	return &request
}

func transformCacheRequestCreate(cache contracts.ManifestCacheSetting) *apiCache.Request {
	request := apiCache.Request{}

	if cache.Name != "" {
		request.SetName(cache.Name)
	}

	if cache.BrowserCache != nil {
		request.SetBrowserCache(*cache.BrowserCache)
	}
	if cache.Modules != nil {
		request.SetModules(*cache.Modules)
	}

	return &request
}

func transformRuleResponse(rule contracts.ManifestRule) *apiApplications.UpdateRulesEngineResponse {
	request := &apiApplications.UpdateRulesEngineResponse{}

	request.SetActive(rule.Active)
	if rule.Criteria != nil {
		request.SetCriteria(rule.Criteria)
	}

	request.SetDescription(rule.Description)
	request.SetName(rule.Name)

	return request
}

func transformRuleRequest(rule contracts.ManifestRule) *apiApplications.UpdateRulesEngineRequest {
	request := &apiApplications.UpdateRulesEngineRequest{}

	request.SetActive(rule.Active)
	if rule.Criteria != nil {
		request.SetCriteria(rule.Criteria)
	}

	request.SetDescription(rule.Description)
	request.SetName(rule.Name)

	return request
}

func getConnectorName(connector edgesdk.ConnectorRequest, defaultName string) (string, string) {
	if connector.ConnectorHTTPRequest != nil {
		return connector.ConnectorHTTPRequest.Name, "http"
	}

	if connector.ConnectorStorageRequest != nil {
		return connector.ConnectorStorageRequest.Name, "storage"
	}

	return defaultName, ""
}

func transformRuleRequestCreate(rule contracts.ManifestRule) edgesdk.RequestPhaseRuleRequest {
	request := edgesdk.RequestPhaseRuleRequest{}

	request.SetActive(rule.Active)
	if rule.Criteria != nil {
		request.SetCriteria(rule.Criteria)
	}
	request.SetDescription(rule.Description)
	request.SetName(rule.Name)

	return request
}

func transformRuleResponseCreate(rule contracts.ManifestRule) edgesdk.ResponsePhaseRuleRequest {
	request := edgesdk.ResponsePhaseRuleRequest{}

	request.SetActive(rule.Active)
	if rule.Criteria != nil {
		request.SetCriteria(rule.Criteria)
	}
	request.SetDescription(rule.Description)
	request.SetName(rule.Name)

	return request
}
