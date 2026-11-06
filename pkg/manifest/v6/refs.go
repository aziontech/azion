package v6

import (
	"encoding/json"
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/manifest"
	"github.com/aziontech/azion-cli/pkg/contracts"
)

// target says where a behavior keeps the id of another resource, and which
// ids by name the run has for it.
type target struct {
	attribute string
	kind      string
	ids       map[string]int64
}

// refs is everything a rule may refer to by name.
type refs struct {
	// behaviors maps a behavior type to the resource it names.
	behaviors map[string]target
	// networkLists, when set, resolves the argument of is_in_list and
	// is_not_in_list criteria, which name a Network List (Firewall rules).
	networkLists map[string]int64
}

var listOperators = map[string]bool{"is_in_list": true, "is_not_in_list": true}

// resolveRule replaces, in one rule, every reference given by name with the id
// of the resource this run wrote. A reference already given as a number is
// left alone. It works on the raw JSON because the SDK types only accept ids,
// and the rule is decoded into them afterwards.
func resolveRule(raw json.RawMessage, r refs) (json.RawMessage, error) {
	var rule map[string]any
	if err := json.Unmarshal(raw, &rule); err != nil {
		return nil, fmt.Errorf("%w: %v", msg.ErrorV6InvalidRule, err)
	}

	behaviors, _ := rule["behaviors"].([]any)
	for _, b := range behaviors {
		behavior, ok := b.(map[string]any)
		if !ok {
			continue
		}
		kind, _ := behavior["type"].(string)
		t, ok := r.behaviors[kind]
		if !ok {
			continue
		}
		attrs, ok := behavior["attributes"].(map[string]any)
		if !ok {
			continue
		}
		name, isName := attrs[t.attribute].(string)
		if !isName {
			continue
		}
		id, found := t.ids[name]
		if !found {
			return nil, fmt.Errorf("rule %q: behavior %s refers to %s %q: %w", rule["name"], kind, t.kind, name, msg.ErrorV6UnknownRef)
		}
		attrs[t.attribute] = id
	}

	if r.networkLists != nil {
		groups, _ := rule["criteria"].([]any)
		for _, g := range groups {
			group, _ := g.([]any)
			for _, c := range group {
				criterion, ok := c.(map[string]any)
				if !ok {
					continue
				}
				op, _ := criterion["operator"].(string)
				name, isName := criterion["argument"].(string)
				if !listOperators[op] || !isName {
					continue
				}
				id, found := r.networkLists[name]
				if !found {
					return nil, fmt.Errorf("rule %q: criterion %s refers to network list %q: %w", rule["name"], op, name, msg.ErrorV6UnknownRef)
				}
				criterion["argument"] = id
			}
		}
	}
	return json.Marshal(rule)
}

// decodeRules resolves each rule's references and decodes it into the SDK
// request type, which also validates it.
func decodeRules[Req any](raws []json.RawMessage, r refs) ([]Req, error) {
	out := make([]Req, 0, len(raws))
	for _, raw := range raws {
		resolved, err := resolveRule(raw, r)
		if err != nil {
			return nil, err
		}
		var req Req
		if err := json.Unmarshal(resolved, &req); err != nil {
			return nil, fmt.Errorf("%w: %v", msg.ErrorV6InvalidRule, err)
		}
		out = append(out, req)
	}
	return out, nil
}

// splitRules separates an Application's rules by phase, keeping their order.
func splitRules(rules []contracts.ManifestV6Rule) (request, response []json.RawMessage, err error) {
	for _, r := range rules {
		switch r.Phase {
		case contracts.PhaseRequest:
			request = append(request, r.Rule)
		case contracts.PhaseResponse:
			response = append(response, r.Rule)
		default:
			return nil, nil, fmt.Errorf("rule phase %q (want request or response): %w", r.Phase, msg.ErrorV6InvalidRule)
		}
	}
	return request, response, nil
}

// runFunctionIDs maps what a run_function behavior may name to the instance
// it runs. The bundler writes the Function's name there, as the v4 manifest
// did, so each instance is reachable by the name of its Function (the first
// instance of a Function wins) and by its own name.
func runFunctionIDs(instances []contracts.ManifestV6FunctionInstance, instanceIDs map[string]int64) map[string]int64 {
	out := make(map[string]int64, 2*len(instances))
	for _, inst := range instances {
		id, ok := instanceIDs[inst.Name]
		if !ok {
			continue
		}
		if fn := inst.Function.Name; fn != "" {
			if _, taken := out[fn]; !taken {
				out[fn] = id
			}
		}
		out[inst.Name] = id
	}
	return out
}

// resolveCustomPage replaces the Connector names in a custom page's pages with
// their ids.
func resolveCustomPage(raw json.RawMessage, connectors map[string]int64) (json.RawMessage, error) {
	var cp map[string]any
	if err := json.Unmarshal(raw, &cp); err != nil {
		return nil, fmt.Errorf("custom page: %w", err)
	}
	pages, _ := cp["pages"].([]any)
	for _, p := range pages {
		entry, _ := p.(map[string]any)
		page, _ := entry["page"].(map[string]any)
		attrs, _ := page["attributes"].(map[string]any)
		name, isName := attrs["connector"].(string)
		if !isName {
			continue
		}
		id, found := connectors[name]
		if !found {
			return nil, fmt.Errorf("custom page %q refers to connector %q: %w", cp["name"], name, msg.ErrorV6UnknownRef)
		}
		attrs["connector"] = id
	}
	return json.Marshal(cp)
}
