package v6

import (
	"context"

	apiApplications "github.com/aziontech/azion-cli/pkg/api/v6/applications"
	apiFirewalls "github.com/aziontech/azion-cli/pkg/api/v6/firewalls"
	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	apiWafs "github.com/aziontech/azion-cli/pkg/api/v6/wafs"
	"github.com/aziontech/azion-cli/pkg/contracts"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
)

// The sub-resource kinds reconcile works on, each built from its client's
// own methods.

func appRequestRules(c *apiApplications.Client) sub[sdk.RequestPhaseRuleRequest] {
	return sub[sdk.RequestPhaseRuleRequest]{
		kind: "request rule",
		list: func(ctx context.Context, id int64, vid string) ([]version.Item, error) {
			return listAll[sdk.RequestPhaseRule](ctx, func(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedRequestPhaseRuleList, error) {
				return c.ListRequestRules(ctx, id, vid, opts)
			})
		},
		create: c.CreateRequestRule,
		update: c.UpdateRequestRule,
		delete: c.DeleteRequestRule,
		order:  c.OrderRequestRules,
	}
}

func appResponseRules(c *apiApplications.Client) sub[sdk.ResponsePhaseRuleRequest] {
	return sub[sdk.ResponsePhaseRuleRequest]{
		kind: "response rule",
		list: func(ctx context.Context, id int64, vid string) ([]version.Item, error) {
			return listAll[sdk.ResponsePhaseRule](ctx, func(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedResponsePhaseRuleList, error) {
				return c.ListResponseRules(ctx, id, vid, opts)
			})
		},
		create: c.CreateResponseRule,
		update: c.UpdateResponseRule,
		delete: c.DeleteResponseRule,
		order:  c.OrderResponseRules,
	}
}

func cacheSettings(c *apiApplications.Client) sub[sdk.CacheSettingRequest] {
	return sub[sdk.CacheSettingRequest]{
		kind: "cache setting",
		list: func(ctx context.Context, id int64, vid string) ([]version.Item, error) {
			return listAll[sdk.CacheSetting](ctx, func(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedCacheSettingList, error) {
				return c.ListCacheSettings(ctx, id, vid, opts)
			})
		},
		create: c.CreateCacheSetting,
		update: c.UpdateCacheSetting,
		delete: c.DeleteCacheSetting,
	}
}

func deviceGroups(c *apiApplications.Client) sub[sdk.DeviceGroupRequest] {
	return sub[sdk.DeviceGroupRequest]{
		kind: "device group",
		list: func(ctx context.Context, id int64, vid string) ([]version.Item, error) {
			return listAll[sdk.DeviceGroup](ctx, func(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedDeviceGroupList, error) {
				return c.ListDeviceGroups(ctx, id, vid, opts)
			})
		},
		create: c.CreateDeviceGroup,
		update: c.UpdateDeviceGroup,
		delete: c.DeleteDeviceGroup,
	}
}

func appFunctionInstances(c *apiApplications.Client) sub[sdk.FunctionInstanceRequest] {
	return sub[sdk.FunctionInstanceRequest]{
		kind: "function instance",
		list: func(ctx context.Context, id int64, vid string) ([]version.Item, error) {
			return listAll[sdk.FunctionInstance](ctx, func(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedFunctionInstanceList, error) {
				return c.ListFunctionInstances(ctx, id, vid, opts)
			})
		},
		create: c.CreateFunctionInstance,
		update: c.UpdateFunctionInstance,
		delete: c.DeleteFunctionInstance,
	}
}

func firewallRules(c *apiFirewalls.Client) sub[sdk.FirewallRuleRequest] {
	return sub[sdk.FirewallRuleRequest]{
		kind: "rule",
		list: func(ctx context.Context, id int64, vid string) ([]version.Item, error) {
			return listAll[sdk.FirewallRule](ctx, func(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedFirewallRuleList, error) {
				return c.ListRules(ctx, id, vid, opts)
			})
		},
		create: c.CreateRule,
		update: c.UpdateRule,
		delete: c.DeleteRule,
		order:  c.OrderRules,
	}
}

func firewallFunctionInstances(c *apiFirewalls.Client) sub[sdk.FirewallFunctionInstanceRequest] {
	return sub[sdk.FirewallFunctionInstanceRequest]{
		kind: "function instance",
		list: func(ctx context.Context, id int64, vid string) ([]version.Item, error) {
			return listAll[sdk.FirewallFunctionInstance](ctx, func(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedFirewallFunctionInstanceList, error) {
				return c.ListFunctionInstances(ctx, id, vid, opts)
			})
		},
		create: c.CreateFunctionInstance,
		update: c.UpdateFunctionInstance,
		delete: c.DeleteFunctionInstance,
	}
}

func wafExceptions(c *apiWafs.Client) sub[sdk.WAFRuleRequest] {
	return sub[sdk.WAFRuleRequest]{
		kind: "exception",
		list: func(ctx context.Context, id int64, vid string) ([]version.Item, error) {
			return listAll[sdk.WAFRule](ctx, func(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedWAFRuleList, error) {
				return c.ListExceptions(ctx, id, vid, opts)
			})
		},
		create: c.CreateException,
		update: c.UpdateException,
		delete: c.DeleteException,
	}
}
