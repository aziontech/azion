package manifest

import (
	"fmt"
	"math"
	"strconv"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	edgesdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
)

// convertFirewallRuleToSDK converts a FirewallManifestRule (lenient JSON types)
// to an edgesdk.FirewallRuleRequest (strict SDK types).
func convertFirewallRuleToSDK(rule contracts.FirewallManifestRule) (edgesdk.FirewallRuleRequest, error) {
	criteria, err := convertFirewallCriteria(rule.Criteria)
	if err != nil {
		return edgesdk.FirewallRuleRequest{}, fmt.Errorf("error converting criteria for rule %q: %w", rule.Name, err)
	}

	behaviors, err := convertFirewallBehaviors(rule.Behaviors)
	if err != nil {
		return edgesdk.FirewallRuleRequest{}, fmt.Errorf("error converting behaviors for rule %q: %w", rule.Name, err)
	}

	sdkRule := edgesdk.FirewallRuleRequest{
		Name:      rule.Name,
		Active:    rule.Active,
		Criteria:  criteria,
		Behaviors: behaviors,
	}

	if rule.Description != nil {
		sdkRule.Description = rule.Description
	}

	return sdkRule, nil
}

// convertFirewallCriteria converts manifest criterion slices to SDK types.
func convertFirewallCriteria(criteria [][]contracts.FirewallManifestCriterion) ([][]edgesdk.FirewallCriterionFieldRequest, error) {
	result := make([][]edgesdk.FirewallCriterionFieldRequest, len(criteria))
	for i, group := range criteria {
		sdkGroup := make([]edgesdk.FirewallCriterionFieldRequest, len(group))
		for j, c := range group {
			sdkCriterion := edgesdk.FirewallCriterionFieldRequest{
				Variable:    c.Variable,
				Operator:    c.Operator,
				Conditional: c.Conditional,
			}

			if c.Argument != nil {
				arg := edgesdk.FirewallCriterionArgumentRequest{}
				switch v := c.Argument.(type) {
				case string:
					arg.String = &v
				case float64:
					intVal := int64(v)
					arg.Int64 = &intVal
				default:
					str := fmt.Sprintf("%v", v)
					arg.String = &str
				}
				nullableArg := edgesdk.NewNullableFirewallCriterionArgumentRequest(&arg)
				sdkCriterion.Argument = *nullableArg
			}

			sdkGroup[j] = sdkCriterion
		}
		result[i] = sdkGroup
	}
	return result, nil
}

// convertFirewallBehaviors converts manifest behaviors to SDK types.
func convertFirewallBehaviors(behaviors []contracts.FirewallManifestBehavior) ([]edgesdk.FirewallBehaviorRequest, error) {
	result := make([]edgesdk.FirewallBehaviorRequest, len(behaviors))
	for i, b := range behaviors {
		sdkBehavior, err := convertSingleFirewallBehavior(b)
		if err != nil {
			return nil, fmt.Errorf("error converting behavior %q: %w", b.Type, err)
		}
		result[i] = sdkBehavior
	}
	return result, nil
}

// convertSingleFirewallBehavior converts a single manifest behavior to an SDK behavior.
func convertSingleFirewallBehavior(b contracts.FirewallManifestBehavior) (edgesdk.FirewallBehaviorRequest, error) {
	switch b.Type {
	case "deny", "drop":
		// No-args behaviors
		return edgesdk.FirewallBehaviorNoArgsRequestAsFirewallBehaviorRequest(
			edgesdk.NewFirewallBehaviorNoArgsRequest(b.Type),
		), nil

	case "run_function":
		// Args behavior (simple argument)
		value, err := toInt64(b.Attributes["value"])
		if err != nil {
			return edgesdk.FirewallBehaviorRequest{}, fmt.Errorf("run_function behavior requires integer 'value' attribute: %w", err)
		}
		attrs := edgesdk.FirewallBehaviorRunFunctionAttributesRequest{
			Value: value,
		}
		return edgesdk.FirewallBehaviorArgsRequestAsFirewallBehaviorRequest(
			edgesdk.NewFirewallBehaviorArgsRequest(b.Type, attrs),
		), nil

	case "set_rate_limit":
		return convertRateLimitBehavior(b)

	case "set_custom_response":
		return convertCustomResponseBehavior(b)

	case "set_waf":
		return convertWafBehavior(b)

	default:
		logger.Debug(fmt.Sprintf("Unknown firewall behavior type: %s, attempting object args conversion", b.Type))
		return convertGenericObjectArgsBehavior(b)
	}
}

// convertRateLimitBehavior converts a set_rate_limit behavior.
func convertRateLimitBehavior(b contracts.FirewallManifestBehavior) (edgesdk.FirewallBehaviorRequest, error) {
	attrs := b.Attributes

	limitBy, err := toString(attrs["limit_by"])
	if err != nil {
		return edgesdk.FirewallBehaviorRequest{}, fmt.Errorf("set_rate_limit requires 'limit_by' attribute: %w", err)
	}

	averageRateLimit, err := toInt64(attrs["average_rate_limit"])
	if err != nil {
		return edgesdk.FirewallBehaviorRequest{}, fmt.Errorf("set_rate_limit requires integer 'average_rate_limit' attribute: %w", err)
	}

	rateLimitAttrs := edgesdk.NewFirewallBehaviorSetRateLimitAttributesRequest(limitBy, averageRateLimit)

	if typeVal, ok := attrs["type"]; ok {
		typeStr, err := toString(typeVal)
		if err == nil {
			rateLimitAttrs.SetType(typeStr)
		}
	}

	if burstVal, ok := attrs["maximum_burst_size"]; ok {
		burstSize, err := toInt64(burstVal)
		if err == nil {
			rateLimitAttrs.SetMaximumBurstSize(burstSize)
		}
	}

	objAttrs := edgesdk.FirewallBehaviorSetRateLimitAttributesRequestAsFirewallBehaviorObjectArgsRequestAttributes(rateLimitAttrs)
	objArgs := edgesdk.NewFirewallBehaviorObjectArgsRequest(b.Type, objAttrs)

	return edgesdk.FirewallBehaviorObjectArgsRequestAsFirewallBehaviorRequest(objArgs), nil
}

// convertCustomResponseBehavior converts a set_custom_response behavior.
func convertCustomResponseBehavior(b contracts.FirewallManifestBehavior) (edgesdk.FirewallBehaviorRequest, error) {
	attrs := b.Attributes

	statusCode, err := toInt64(attrs["status_code"])
	if err != nil {
		return edgesdk.FirewallBehaviorRequest{}, fmt.Errorf("set_custom_response requires integer 'status_code' attribute: %w", err)
	}

	customResponseAttrs := edgesdk.NewFirewallBehaviorSetCustomResponseAttributesRequest(statusCode)

	if ct, ok := attrs["content_type"]; ok {
		ctStr, err := toString(ct)
		if err == nil {
			customResponseAttrs.SetContentType(ctStr)
		}
	}

	if cb, ok := attrs["content_body"]; ok {
		cbStr, err := toString(cb)
		if err == nil {
			customResponseAttrs.SetContentBody(cbStr)
		}
	}

	objAttrs := edgesdk.FirewallBehaviorSetCustomResponseAttributesRequestAsFirewallBehaviorObjectArgsRequestAttributes(customResponseAttrs)
	objArgs := edgesdk.NewFirewallBehaviorObjectArgsRequest(b.Type, objAttrs)

	return edgesdk.FirewallBehaviorObjectArgsRequestAsFirewallBehaviorRequest(objArgs), nil
}

// convertWafBehavior converts a set_waf behavior.
func convertWafBehavior(b contracts.FirewallManifestBehavior) (edgesdk.FirewallBehaviorRequest, error) {
	attrs := b.Attributes

	wafId, err := toInt64(attrs["waf_id"])
	if err != nil {
		return edgesdk.FirewallBehaviorRequest{}, fmt.Errorf("set_waf requires integer 'waf_id' attribute: %w", err)
	}

	mode, err := toString(attrs["mode"])
	if err != nil {
		return edgesdk.FirewallBehaviorRequest{}, fmt.Errorf("set_waf requires string 'mode' attribute: %w", err)
	}

	wafAttrs := edgesdk.NewFirewallBehaviorSetWafAttributesRequest(wafId, mode)

	objAttrs := edgesdk.FirewallBehaviorSetWafAttributesRequestAsFirewallBehaviorObjectArgsRequestAttributes(wafAttrs)
	objArgs := edgesdk.NewFirewallBehaviorObjectArgsRequest(b.Type, objAttrs)

	return edgesdk.FirewallBehaviorObjectArgsRequestAsFirewallBehaviorRequest(objArgs), nil
}

// convertGenericObjectArgsBehavior attempts to convert an unknown behavior type
// as a no-args behavior (fallback).
func convertGenericObjectArgsBehavior(b contracts.FirewallManifestBehavior) (edgesdk.FirewallBehaviorRequest, error) {
	if len(b.Attributes) == 0 {
		return edgesdk.FirewallBehaviorNoArgsRequestAsFirewallBehaviorRequest(
			edgesdk.NewFirewallBehaviorNoArgsRequest(b.Type),
		), nil
	}
	return edgesdk.FirewallBehaviorRequest{}, fmt.Errorf("unsupported firewall behavior type %q with attributes", b.Type)
}

// toInt64 converts an interface{} value to int64, handling strings and float64.
func toInt64(v interface{}) (int64, error) {
	if v == nil {
		return 0, fmt.Errorf("value is nil")
	}
	switch val := v.(type) {
	case float64:
		if val != math.Trunc(val) {
			return 0, fmt.Errorf("value %f is not an integer", val)
		}
		return int64(val), nil
	case int64:
		return val, nil
	case int:
		return int64(val), nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", v)
	}
}

// toString converts an interface{} value to string.
func toString(v interface{}) (string, error) {
	if v == nil {
		return "", fmt.Errorf("value is nil")
	}
	switch val := v.(type) {
	case string:
		return val, nil
	default:
		return fmt.Sprintf("%v", val), nil
	}
}
