package manifest

import (
	"context"
	"strconv"

	apiCache "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	apiWorkloads "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	edgesdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/aziontech/azionapi-v4-go-sdk/edge"
	"go.uber.org/zap"
)

//TODO: FIX HERE

// func transformEdgeConnectorRequest(connectorRequest edgesdk.EdgeConnectorPolymorphicRequest) *apiConnector.UpdateRequest {
// 	request := &apiConnector.UpdateRequest{}

// 	if connectorRequest.EdgeConnectorHTTPRequest != nil {
// 		bodyRequest := connectorRequest.EdgeConnectorHTTPRequest
// 		body := edgesdk.PatchedEdgeConnectorHTTPRequest{}
// 		if bodyRequest.Active != nil {
// 			body.SetActive(*bodyRequest.Active)
// 		}
// 		if len(bodyRequest.Addresses) > 0 {
// 			body.SetAddresses(bodyRequest.Addresses)
// 		}
// 		if len(bodyRequest.ConnectionPreference) > 0 {
// 			body.SetConnectionPreference(bodyRequest.ConnectionPreference)
// 		}
// 		if bodyRequest.ConnectionTimeout != nil {
// 			body.SetConnectionTimeout(*bodyRequest.ConnectionTimeout)
// 		}
// 		if bodyRequest.LoadBalanceMethod != nil {
// 			body.SetLoadBalanceMethod(*bodyRequest.LoadBalanceMethod)
// 		}
// 		if bodyRequest.MaxRetries != nil {
// 			body.SetMaxRetries(*bodyRequest.MaxRetries)
// 		}

// 		body.SetModules(bodyRequest.Modules)

// 		if bodyRequest.Name != "" {
// 			body.SetName(bodyRequest.Name)
// 		}
// 		if bodyRequest.ReadWriteTimeout != nil {
// 			body.SetReadWriteTimeout(*bodyRequest.ReadWriteTimeout)
// 		}
// 		if bodyRequest.Tls != nil {
// 			body.SetTls(*bodyRequest.Tls)
// 		}
// 		if bodyRequest.Type != "" {
// 			body.SetType(bodyRequest.Type)
// 		}
// 		body.SetTypeProperties(bodyRequest.TypeProperties)
// 		request.PatchedEdgeConnectorHTTPRequest = &body
// 		return request
// 	}

// 	if connectorRequest.EdgeConnectorLiveIngestRequest != nil {
// 		body := edgesdk.PatchedEdgeConnectorLiveIngestRequest{}
// 		bodyRequest := connectorRequest.EdgeConnectorLiveIngestRequest
// 		if bodyRequest.Active != nil {
// 			body.SetActive(*bodyRequest.Active)
// 		}
// 		if len(bodyRequest.ConnectionPreference) > 0 {
// 			body.SetConnectionPreference(bodyRequest.ConnectionPreference)
// 		}
// 		if bodyRequest.ConnectionTimeout != nil {
// 			body.SetConnectionTimeout(*bodyRequest.ConnectionTimeout)
// 		}
// 		if bodyRequest.LoadBalanceMethod != nil {
// 			body.SetLoadBalanceMethod(*bodyRequest.LoadBalanceMethod)
// 		}
// 		if bodyRequest.MaxRetries != nil {
// 			body.SetMaxRetries(*bodyRequest.MaxRetries)
// 		}

// 		body.SetModules(bodyRequest.Modules)

// 		if bodyRequest.Name != "" {
// 			body.SetName(bodyRequest.Name)
// 		}
// 		if bodyRequest.ReadWriteTimeout != nil {
// 			body.SetReadWriteTimeout(*bodyRequest.ReadWriteTimeout)
// 		}
// 		if bodyRequest.Tls != nil {
// 			body.SetTls(*bodyRequest.Tls)
// 		}
// 		if bodyRequest.Type != "" {
// 			body.SetType(bodyRequest.Type)
// 		}
// 		body.SetTypeProperties(bodyRequest.TypeProperties)
// 		request.PatchedEdgeConnectorLiveIngestRequest = &body
// 		return request
// 	}

// 	if connectorRequest.EdgeConnectorS3Request != nil {
// 		body := edgesdk.PatchedEdgeConnectorS3Request{}
// 		bodyRequest := connectorRequest.EdgeConnectorS3Request
// 		if bodyRequest.Active != nil {
// 			body.SetActive(*bodyRequest.Active)
// 		}
// 		if len(bodyRequest.Addresses) > 0 {
// 			body.SetAddresses(bodyRequest.Addresses)
// 		}
// 		if len(bodyRequest.ConnectionPreference) > 0 {
// 			body.SetConnectionPreference(bodyRequest.ConnectionPreference)
// 		}
// 		if bodyRequest.ConnectionTimeout != nil {
// 			body.SetConnectionTimeout(*bodyRequest.ConnectionTimeout)
// 		}
// 		if bodyRequest.LoadBalanceMethod != nil {
// 			body.SetLoadBalanceMethod(*bodyRequest.LoadBalanceMethod)
// 		}
// 		if bodyRequest.MaxRetries != nil {
// 			body.SetMaxRetries(*bodyRequest.MaxRetries)
// 		}
// 		body.SetModules(bodyRequest.Modules)

// 		if bodyRequest.Name != "" {
// 			body.SetName(bodyRequest.Name)
// 		}
// 		if bodyRequest.ReadWriteTimeout != nil {
// 			body.SetReadWriteTimeout(*bodyRequest.ReadWriteTimeout)
// 		}
// 		if bodyRequest.Tls != nil {
// 			body.SetTls(*bodyRequest.Tls)
// 		}
// 		if bodyRequest.Type != "" {
// 			body.SetType(bodyRequest.Type)
// 		}
// 		body.SetTypeProperties(bodyRequest.TypeProperties)

// 		request.PatchedEdgeConnectorS3Request = &body
// 		return request
// 	}

// 	if connectorRequest.EdgeConnectorStorageRequest != nil {
// 		body := edgesdk.PatchedEdgeConnectorStorageRequest{}
// 		bodyRequest := connectorRequest.EdgeConnectorStorageRequest
// 		if bodyRequest.Active != nil {
// 			body.SetActive(*bodyRequest.Active)
// 		}
// 		if len(bodyRequest.ConnectionPreference) > 0 {
// 			body.SetConnectionPreference(bodyRequest.ConnectionPreference)
// 		}
// 		if bodyRequest.ConnectionTimeout != nil {
// 			body.SetConnectionTimeout(*bodyRequest.ConnectionTimeout)
// 		}
// 		if bodyRequest.LoadBalanceMethod != nil {
// 			body.SetLoadBalanceMethod(*bodyRequest.LoadBalanceMethod)
// 		}
// 		if bodyRequest.MaxRetries != nil {
// 			body.SetMaxRetries(*bodyRequest.MaxRetries)
// 		}
// 		body.SetModules(bodyRequest.Modules)
// 		if bodyRequest.Name != "" {
// 			body.SetName(bodyRequest.Name)
// 		}
// 		if bodyRequest.ReadWriteTimeout != nil {
// 			body.SetReadWriteTimeout(*bodyRequest.ReadWriteTimeout)
// 		}
// 		if bodyRequest.Tls != nil {
// 			body.SetTls(*bodyRequest.Tls)
// 		}
// 		if bodyRequest.Type != "" {
// 			body.SetType(bodyRequest.Type)
// 		}
// 		body.SetTypeProperties(bodyRequest.TypeProperties)

// 		request.PatchedEdgeConnectorStorageRequest = &body
// 		return request
// 	}

// 	return request
// }

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

func transformEdgeApplicationRequestUpdate(edgeapprequest contracts.EdgeApplications) *apiEdgeApplications.UpdateRequest {
	request := &apiEdgeApplications.UpdateRequest{}

	if edgeapprequest.Active != nil {
		request.SetActive(*edgeapprequest.Active)
	}
	if edgeapprequest.Debug != nil {
		request.SetDebug(*edgeapprequest.Debug)
	}
	type Modules struct {
		EdgeCacheEnabled              bool `json:"edge_cache_enabled"`
		EdgeFunctionsEnabled          bool `json:"edge_functions_enabled"`
		ApplicationAcceleratorEnabled bool `json:"application_accelerator_enabled"`
		ImageProcessorEnabled         bool `json:"image_processor_enabled"`
		TieredCacheEnabled            bool `json:"tiered_cache_enabled"`
	}
	// if edgeapprequest.Modules != nil {
	// 	modules := edgesdk.EdgeApplicationModulesRequest{}
	// 	if edgeapprequest.Modules.ApplicationAcceleratorEnabled != nil {
	// 		modules.SetApplicationAcceleratorEnabled(*edgeapprequest.Modules.ApplicationAcceleratorEnabled)
	// 	}
	// 	if edgeapprequest.Modules.EdgeCacheEnabled != nil {
	// 		modules.SetEdgeCacheEnabled(*edgeapprequest.Modules.EdgeCacheEnabled)
	// 	}
	// 	if edgeapprequest.Modules.EdgeFunctionsEnabled != nil {
	// 		modules.SetEdgeFunctionsEnabled(*edgeapprequest.Modules.EdgeFunctionsEnabled)
	// 	}
	// 	if edgeapprequest.Modules.ImageProcessorEnabled != nil {
	// 		modules.SetImageProcessorEnabled(*edgeapprequest.Modules.ImageProcessorEnabled)
	// 	}
	// 	if edgeapprequest.Modules.TieredCacheEnabled != nil {
	// 		modules.SetTieredCacheEnabled(*edgeapprequest.Modules.TieredCacheEnabled)
	// 	}

	// 	request.SetModules(modules)
	// }

	//TODO: Fix Here
	if edgeapprequest.Name != "" {
		request.SetName(edgeapprequest.Name)
	}

	return request
}

func transformEdgeApplicationRequestCreate(edgeapprequest contracts.EdgeApplications) *apiEdgeApplications.CreateRequest {
	request := &apiEdgeApplications.CreateRequest{}

	if edgeapprequest.Active != nil {
		request.SetActive(*edgeapprequest.Active)
	}
	if edgeapprequest.Debug != nil {
		request.SetDebug(*edgeapprequest.Debug)
	}
	type Modules struct {
		EdgeCacheEnabled              bool `json:"edge_cache_enabled"`
		EdgeFunctionsEnabled          bool `json:"edge_functions_enabled"`
		ApplicationAcceleratorEnabled bool `json:"application_accelerator_enabled"`
		ImageProcessorEnabled         bool `json:"image_processor_enabled"`
		TieredCacheEnabled            bool `json:"tiered_cache_enabled"`
	}
	// if edgeapprequest.Modules != nil {
	// 	modules := edgesdk.EdgeApplicationModulesRequest{}
	// 	if edgeapprequest.Modules.ApplicationAcceleratorEnabled != nil {
	// 		modules.SetApplicationAcceleratorEnabled(*edgeapprequest.Modules.ApplicationAcceleratorEnabled)
	// 	}
	// 	if edgeapprequest.Modules.EdgeCacheEnabled != nil {
	// 		modules.SetEdgeCacheEnabled(*edgeapprequest.Modules.EdgeCacheEnabled)
	// 	}
	// 	if edgeapprequest.Modules.EdgeFunctionsEnabled != nil {
	// 		modules.SetEdgeFunctionsEnabled(*edgeapprequest.Modules.EdgeFunctionsEnabled)
	// 	}
	// 	if edgeapprequest.Modules.ImageProcessorEnabled != nil {
	// 		modules.SetImageProcessorEnabled(*edgeapprequest.Modules.ImageProcessorEnabled)
	// 	}
	// 	if edgeapprequest.Modules.TieredCacheEnabled != nil {
	// 		modules.SetTieredCacheEnabled(*edgeapprequest.Modules.TieredCacheEnabled)
	// 	}

	// 	request.SetModules(modules)
	// }

	//TODO: Fix Here
	if edgeapprequest.Name != "" {
		request.SetName(edgeapprequest.Name)
	}

	return request
}

func transformCacheRequest(cache edgesdk.CacheSettingRequest) *apiCache.RequestUpdate {
	request := apiCache.RequestUpdate{}

	if cache.Name != "" {
		request.SetName(cache.Name)
	}
	// request.SetApplicationControls(cache.ApplicationControls)
	// request.SetBrowserCache(cache.BrowserCache)
	// request.SetEdgeCache(cache.EdgeCache)
	// request.SetSliceControls(cache.SliceControls)

	//TODO: Fix Here

	return &request
}

func transformRuleRequest(rule edgesdk.PatchedEdgeApplicationResponsePhaseRuleEngineRequest) *apiEdgeApplications.UpdateRulesEngineRequest {
	request := &apiEdgeApplications.UpdateRulesEngineRequest{}

	if rule.Active != nil {
		request.SetActive(*rule.Active)
	}
	// if rule.Behaviors != nil {
	// 	request.SetBehaviors(rule.Behaviors)
	// }
	// if rule.Criteria != nil {
	// 	request.SetCriteria(rule.Criteria)
	// }

	//TODO: Fix Here
	if rule.Description != nil {
		request.SetDescription(*rule.Description)
	}
	if rule.Name != nil {
		request.SetName(*rule.Name)
	}

	return request
}

func doCacheForRule(ctx context.Context, client *apiEdgeApplications.Client, conf *contracts.AzionApplicationOptions, function contracts.AzionJsonDataFunction) (int64, error) {
	if function.CacheId > 0 {
		return function.CacheId, nil
	}
	var reqCache apiEdgeApplications.CreateCacheSettingsRequest
	reqCache.SetName("function-policy")
	// reqCache.BrowserCache = edgesdk.BrowserCacheModuleRequest{
	// 	Behavior: "honor",
	// }
	// reqCache.EdgeCache = edgesdk.EdgeCacheModuleRequest{
	// 	Behavior: "honor",
	// 	MaxAge:   0,
	// }
	// reqCache.ApplicationControls = edgesdk.ApplicationControlsModuleRequest{
	// 	CacheByQueryString: "all",
	// 	CacheByCookies:     "all",
	// }

	//TODO: Fix Here

	// create cache to function next
	str := strconv.FormatInt(conf.Application.ID, 10)
	cache, err := client.CreateCacheEdgeApplication(ctx, &reqCache, str)
	if err != nil {
		logger.Debug("Error while creating Cache Settings", zap.Error(err))
		return 0, err
	}

	// conf.Function.CacheId = cache.GetId()

	return cache.GetId(), nil
}

func getConnectorName(connector edge.EdgeConnectorPolymorphicRequest, defaultName string) string {
	if connector.EdgeConnectorHTTPTypedRequest != nil {
		return connector.EdgeConnectorHTTPTypedRequest.Name
	}

	if connector.EdgeConnectorLiveIngestTypedRequest != nil {
		return connector.EdgeConnectorLiveIngestTypedRequest.Name
	}

	if connector.EdgeConnectorS3TypedRequest != nil {
		return connector.EdgeConnectorS3TypedRequest.Name
	}

	if connector.EdgeConnectorStorageTypedRequest != nil {
		return connector.EdgeConnectorStorageTypedRequest.Name
	}
	return defaultName
}
