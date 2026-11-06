package v6

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	msg "github.com/aziontech/azion-cli/messages/manifest"
	apiApplications "github.com/aziontech/azion-cli/pkg/api/v6/applications"
	apiConnectors "github.com/aziontech/azion-cli/pkg/api/v6/connectors"
	apiCustomPages "github.com/aziontech/azion-cli/pkg/api/v6/custompages"
	apiFirewalls "github.com/aziontech/azion-cli/pkg/api/v6/firewalls"
	apiFunctions "github.com/aziontech/azion-cli/pkg/api/v6/functions"
	apiNetworkLists "github.com/aziontech/azion-cli/pkg/api/v6/networklists"
	apiWafs "github.com/aziontech/azion-cli/pkg/api/v6/wafs"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"
)

// Clients are the versioned-resource clients the applier writes through.
type Clients struct {
	Applications *apiApplications.Client
	Functions    *apiFunctions.Client
	Connectors   *apiConnectors.Client
	Firewalls    *apiFirewalls.Client
	CustomPages  *apiCustomPages.Client
	NetworkLists *apiNetworkLists.Client
	WAFs         *apiWafs.Client
}

// Applier writes a v6 manifest into Drafts. Each Apply method handles one
// resource kind and returns the Drafts it wrote, for the build step that
// follows. The kinds must be applied in dependency order — functions,
// network lists, WAFs and connectors before the Application, Firewall and
// Custom Page that refer to them — because references by name resolve against
// what was already applied.
type Applier struct {
	Clients  *Clients
	Conf     *contracts.AzionApplicationOptionsV6
	Manifest *contracts.ManifestV6

	// Say prints a line to the user.
	Say func(string)
	// WriteConfig saves azion.json; it runs after every Draft, so a run that
	// dies halfway keeps what it did.
	WriteConfig func() error
	// ReadFile reads a Function's code.
	ReadFile func(string) ([]byte, error)
}

func (a *Applier) say(line string) {
	if a.Say != nil {
		a.Say(line)
	}
}

// within reports a line about a sub-resource, prefixed with the resource whose
// Draft holds it: application "site": cache setting "site" created (id 7).
func (a *Applier) within(kind, name string) func(string) {
	return func(line string) {
		a.say(fmt.Sprintf("%s %q: %s", kind, name, line))
	}
}

func (a *Applier) writeConfig() {
	if a.WriteConfig == nil {
		return
	}
	if err := a.WriteConfig(); err != nil {
		logger.Debug("Error while writing azion.json", zap.Error(err))
	}
}

func (a *Applier) readFile(path string) ([]byte, error) {
	if a.ReadFile != nil {
		return a.ReadFile(path)
	}
	return os.ReadFile(path)
}

// ids returns the id of every resource of one kind in azion.json, by name.
func ids(list []contracts.AzionJsonV6Resource) map[string]int64 {
	out := make(map[string]int64, len(list))
	for _, r := range list {
		out[r.Name] = r.ID
	}
	return out
}

// functionCodePath is where a Function's code is on disk. The bundler writes
// the path relative to .edge (./functions/worker.js is
// .edge/functions/worker.js), as the v4 deploy reads it.
func functionCodePath(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(".edge", p)
}

func (a *Applier) ApplyFunctions(ctx context.Context) ([]Pending, error) {
	var out []Pending
	for _, fn := range a.Manifest.Functions {
		codePath := functionCodePath(fn.Path)
		code, err := a.readFile(codePath)
		if err != nil {
			logger.Debug("Error while reading function code", zap.String("path", codePath), zap.Error(err))
			return nil, fmt.Errorf("%w %q from %s: %v", msg.ErrorV6FunctionCode, fn.Name, codePath, err)
		}
		req := sdk.FunctionsRequest{Name: fn.Name, Code: string(code), Active: fn.Active}
		if fn.Runtime != "" {
			req.SetRuntime(fn.Runtime)
		}
		if fn.ExecutionEnvironment != "" {
			req.SetExecutionEnvironment(fn.ExecutionEnvironment)
		}
		if fn.DefaultArgs != nil {
			req.DefaultArgs = fn.DefaultArgs
		}
		r := a.Clients.Functions
		p, err := draft(ctx, a, "function", r, &a.Conf.Functions, fn.Name, req)
		if err != nil {
			return nil, err
		}
		if err := verify(ctx, r, p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func (a *Applier) ApplyNetworkLists(ctx context.Context) ([]Pending, error) {
	var out []Pending
	r := a.Clients.NetworkLists
	for _, nl := range a.Manifest.NetworkLists {
		p, err := draft(ctx, a, "network list", r, &a.Conf.NetworkLists, nl.Name, nl)
		if err != nil {
			return nil, err
		}
		if err := verify(ctx, r, p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func (a *Applier) ApplyWAFs(ctx context.Context) ([]Pending, error) {
	var out []Pending
	r := a.Clients.WAFs
	for _, w := range a.Manifest.WAFs {
		req := sdk.WAFRequest{Name: w.Name, Active: w.Active, EngineSettings: w.EngineSettings}
		p, err := draft(ctx, a, "waf", r, &a.Conf.WAFs, w.Name, req)
		if err != nil {
			return nil, err
		}
		if _, err := reconcile(ctx, a.within("waf", w.Name), wafExceptions(r), p.ID, p.VersionID, w.Exceptions); err != nil {
			return nil, err
		}
		if err := verify(ctx, r, p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func (a *Applier) ApplyConnectors(ctx context.Context) ([]Pending, error) {
	var out []Pending
	r := a.Clients.Connectors
	for _, c := range a.Manifest.Connectors {
		p, err := draft(ctx, a, "connector", r, &a.Conf.Connectors, connectorName(c), c)
		if err != nil {
			return nil, err
		}
		if err := verify(ctx, r, p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func (a *Applier) ApplyApplications(ctx context.Context) ([]Pending, error) {
	var out []Pending
	c := a.Clients.Applications
	for _, app := range a.Manifest.Applications {
		req := sdk.ApplicationRequest{Name: app.Name, Modules: app.Modules, Active: app.Active, Debug: app.Debug}
		p, err := draft(ctx, a, "application", c, &a.Conf.Applications, app.Name, req)
		if err != nil {
			return nil, err
		}
		inApp := a.within("application", app.Name)

		if _, err := reconcile(ctx, inApp, deviceGroups(c), p.ID, p.VersionID, app.DeviceGroups); err != nil {
			return nil, err
		}
		caches, err := reconcile(ctx, inApp, cacheSettings(c), p.ID, p.VersionID, app.CacheSettings)
		if err != nil {
			return nil, err
		}
		instances, err := a.functionInstances(app.FunctionInstances)
		if err != nil {
			return nil, err
		}
		instanceIDs, err := reconcile(ctx, inApp, appFunctionInstances(c), p.ID, p.VersionID, instances)
		if err != nil {
			return nil, err
		}

		targets := refs{behaviors: map[string]target{
			"run_function":     {attribute: "value", kind: "function", ids: runFunctionIDs(app.FunctionInstances, instanceIDs)},
			"set_cache_policy": {attribute: "value", kind: "cache setting", ids: caches},
			"set_connector":    {attribute: "value", kind: "connector", ids: ids(a.Conf.Connectors)},
		}}
		requestRaw, responseRaw, err := splitRules(app.Rules)
		if err != nil {
			return nil, err
		}
		requestRules, err := decodeRules[sdk.RequestPhaseRuleRequest](requestRaw, targets)
		if err != nil {
			return nil, err
		}
		if _, err := reconcile(ctx, inApp, appRequestRules(c), p.ID, p.VersionID, requestRules); err != nil {
			return nil, err
		}
		responseRules, err := decodeRules[sdk.ResponsePhaseRuleRequest](responseRaw, targets)
		if err != nil {
			return nil, err
		}
		if _, err := reconcile(ctx, inApp, appResponseRules(c), p.ID, p.VersionID, responseRules); err != nil {
			return nil, err
		}

		if err := verify(ctx, c, p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func (a *Applier) ApplyFirewalls(ctx context.Context) ([]Pending, error) {
	var out []Pending
	c := a.Clients.Firewalls
	for _, fw := range a.Manifest.Firewalls {
		req := sdk.FirewallRequest{Name: fw.Name, Modules: fw.Modules, Active: fw.Active, Debug: fw.Debug}
		p, err := draft(ctx, a, "firewall", c, &a.Conf.Firewalls, fw.Name, req)
		if err != nil {
			return nil, err
		}
		inFw := a.within("firewall", fw.Name)

		appInstances, err := a.functionInstances(fw.FunctionInstances)
		if err != nil {
			return nil, err
		}
		// The firewall's instance request has the same fields as the
		// application's.
		instances := make([]sdk.FirewallFunctionInstanceRequest, len(appInstances))
		for i, in := range appInstances {
			instances[i] = sdk.FirewallFunctionInstanceRequest{Name: in.Name, Args: in.Args, AzionForm: in.AzionForm, Active: in.Active, Function: in.Function}
		}
		instanceIDs, err := reconcile(ctx, inFw, firewallFunctionInstances(c), p.ID, p.VersionID, instances)
		if err != nil {
			return nil, err
		}

		targets := refs{
			behaviors: map[string]target{
				"run_function": {attribute: "value", kind: "function", ids: runFunctionIDs(fw.FunctionInstances, instanceIDs)},
				"set_waf":      {attribute: "waf_id", kind: "waf", ids: ids(a.Conf.WAFs)},
			},
			networkLists: ids(a.Conf.NetworkLists),
		}
		rules, err := decodeRules[sdk.FirewallRuleRequest](fw.Rules, targets)
		if err != nil {
			return nil, err
		}
		if _, err := reconcile(ctx, inFw, firewallRules(c), p.ID, p.VersionID, rules); err != nil {
			return nil, err
		}

		if err := verify(ctx, c, p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func (a *Applier) ApplyCustomPages(ctx context.Context) ([]Pending, error) {
	var out []Pending
	r := a.Clients.CustomPages
	for _, raw := range a.Manifest.CustomPages {
		resolved, err := resolveCustomPage(raw, ids(a.Conf.Connectors))
		if err != nil {
			return nil, err
		}
		var cp sdk.CustomPageRequest
		if err := json.Unmarshal(resolved, &cp); err != nil {
			return nil, fmt.Errorf("custom page: %w", err)
		}
		p, err := draft(ctx, a, "custom page", r, &a.Conf.CustomPages, cp.Name, cp)
		if err != nil {
			return nil, err
		}
		if err := verify(ctx, r, p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

// functionInstances turns the manifest's instances into requests, resolving
// the Function each one runs.
func (a *Applier) functionInstances(in []contracts.ManifestV6FunctionInstance) ([]sdk.FunctionInstanceRequest, error) {
	functions := ids(a.Conf.Functions)
	out := make([]sdk.FunctionInstanceRequest, 0, len(in))
	for _, inst := range in {
		fnID := inst.Function.ID
		if fnID == 0 {
			id, ok := functions[inst.Function.Name]
			if !ok {
				return nil, fmt.Errorf("function instance %q runs function %q: %w", inst.Name, inst.Function.Name, msg.ErrorV6UnknownRef)
			}
			fnID = id
		}
		req := sdk.FunctionInstanceRequest{Name: inst.Name, Active: inst.Active, Function: fnID}
		if inst.Args != nil {
			req.Args = inst.Args
		}
		out = append(out, req)
	}
	return out, nil
}

// Orphans lists the resources azion.json has and the manifest no longer
// declares. They are reported, not deleted: another release may still use
// them (doc/V6-DEPLOY-PLAN.md, D4).
func (a *Applier) Orphans() []string {
	m := a.Manifest
	declared := func(names ...string) map[string]bool {
		set := make(map[string]bool, len(names))
		for _, n := range names {
			set[n] = true
		}
		return set
	}
	kinds := []struct {
		kind     string
		recorded []contracts.AzionJsonV6Resource
		names    map[string]bool
	}{
		{"function", a.Conf.Functions, declared(namesOf(m.Functions)...)},
		{"network list", a.Conf.NetworkLists, declared(namesOf(m.NetworkLists)...)},
		{"waf", a.Conf.WAFs, declared(namesOf(m.WAFs)...)},
		{"connector", a.Conf.Connectors, declared(connectorNames(m.Connectors)...)},
		{"application", a.Conf.Applications, declared(namesOf(m.Applications)...)},
		{"firewall", a.Conf.Firewalls, declared(namesOf(m.Firewalls)...)},
		{"custom page", a.Conf.CustomPages, declared(rawNames(m.CustomPages)...)},
	}
	var out []string
	for _, k := range kinds {
		for _, r := range k.recorded {
			if !k.names[r.Name] {
				out = append(out, fmt.Sprintf(msg.V6OrphanWarning, k.kind, r.Name))
			}
		}
	}
	return out
}

func namesOf[T any](items []T) []string {
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = nameOf(it)
	}
	return out
}

// connectorName reads the name of a connector, whichever of the oneOf
// variants it is.
func connectorName(c sdk.ConnectorRequest) string {
	return nameOf(c)
}

func connectorNames(cs []sdk.ConnectorRequest) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = connectorName(c)
	}
	return out
}
