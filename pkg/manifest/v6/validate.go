package v6

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	msg "github.com/aziontech/azion-cli/messages/manifest"
	"github.com/aziontech/azion-cli/pkg/contracts"
)

// Read loads a v6 manifest and validates it. It is step 3 of the deploy:
// everything the CLI can check without the API is checked here, before any
// call.
func Read(path string) (*contracts.ManifestV6, error) {
	m, err := Parse(path)
	if err != nil {
		return nil, err
	}
	if err := Validate(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Parse reads the manifest without validating it. A first deploy reads the
// template's manifest once before the bundler replaces its placeholders
// ($DEPLOYMENT_NAME and so on), and at that point it is not meant to be valid.
func Parse(path string) (*contracts.ManifestV6, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read the manifest %s: %w", path, err)
	}
	m := &contracts.ManifestV6{}
	if err := json.Unmarshal(data, m); err != nil {
		return nil, fmt.Errorf("failed to parse the manifest %s: %w", path, err)
	}
	return m, nil
}

// Validate checks what a v6 manifest must hold before any of it is applied.
// The bundler validates the azion.config already; this catches what reaches
// the CLI anyway (an old manifest, one edited by hand) and what only the CLI
// resolves.
func Validate(m *contracts.ManifestV6) error {
	if m.Version != contracts.ManifestVersionV6 {
		return fmt.Errorf("version %d: %w", m.Version, msg.ErrorV6WrongVersion)
	}
	var errs []error

	if len(m.Applications) > 1 || len(m.Firewalls) > 1 || len(m.CustomPages) > 1 {
		errs = append(errs, msg.ErrorV6TooMany)
	}

	kinds := []struct {
		kind  string
		names []string
	}{
		{"function", namesOf(m.Functions)},
		{"network list", namesOf(m.NetworkLists)},
		{"waf", namesOf(m.WAFs)},
		{"connector", connectorNames(m.Connectors)},
		{"application", namesOf(m.Applications)},
		{"firewall", namesOf(m.Firewalls)},
		{"custom page", rawNames(m.CustomPages)},
		{"environment", namesOf(m.Environments)},
		{"deployment", namesOf(m.Deployments)},
		{"workload", namesOf(m.Workloads)},
	}
	declared := map[string]map[string]bool{}
	for _, k := range kinds {
		seen := map[string]bool{}
		for _, n := range k.names {
			switch {
			case n == "":
				errs = append(errs, fmt.Errorf("%s: %w", k.kind, msg.ErrorV6NoName))
			case seen[n]:
				errs = append(errs, fmt.Errorf("%s %q: %w", k.kind, n, msg.ErrorV6DuplicateName))
			}
			seen[n] = true
		}
		declared[k.kind] = seen
	}

	for _, fn := range m.Functions {
		if fn.Path == "" {
			errs = append(errs, fmt.Errorf("function %q has no path to its code: %w", fn.Name, msg.ErrorV6FunctionCode))
		}
	}

	checkInstances := func(owner string, in []contracts.ManifestV6FunctionInstance) {
		for _, inst := range in {
			if inst.Function.ID == 0 && !declared["function"][inst.Function.Name] {
				errs = append(errs, fmt.Errorf("%s: function instance %q runs function %q: %w", owner, inst.Name, inst.Function.Name, msg.ErrorV6UnknownRef))
			}
		}
	}
	for _, app := range m.Applications {
		checkInstances("application "+app.Name, app.FunctionInstances)
		if _, _, err := splitRules(app.Rules); err != nil {
			errs = append(errs, fmt.Errorf("application %s: %w", app.Name, err))
		}
	}
	for _, fw := range m.Firewalls {
		checkInstances("firewall "+fw.Name, fw.FunctionInstances)
	}

	for _, w := range m.Workloads {
		for _, b := range w.Bindings {
			if !declared["environment"][b.EnvironmentName] {
				errs = append(errs, fmt.Errorf("workload %q: environment %q: %w", w.Name, b.EnvironmentName, msg.ErrorV6UnknownBinding))
			}
			if !declared["deployment"][b.DeploymentName] {
				errs = append(errs, fmt.Errorf("workload %q: deployment %q: %w", w.Name, b.DeploymentName, msg.ErrorV6UnknownBinding))
			}
		}
	}

	return errors.Join(errs...)
}

func rawNames(raws []json.RawMessage) []string {
	out := make([]string, len(raws))
	for i, r := range raws {
		out[i] = nameOf(r)
	}
	return out
}
