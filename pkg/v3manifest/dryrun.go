package manifest

import (
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/pipeline"
)

// DryRunSteps reports, in order, the manifest steps this project would apply.
//
// It evaluates the real pipeline's predicates against a real ResourceContext
// and executes nothing: pipeline.Run in dry-run mode never calls a step's Run.
// Constructing the context builds API clients and seeds the id maps, but
// performs no request, so this is safe to call without touching the network.
func DryRunSteps(
	f *cmdutil.Factory,
	conf *contracts.AzionApplicationOptionsV3,
	manifest *contracts.Manifest,
	projectConf string,
) ([]string, error) {
	msgs := []string{}
	noWrite := func(*contracts.AzionApplicationOptionsV3, string) error { return nil }

	rc := newResourceContext(f, conf, manifest, projectConf, &msgs, noWrite)

	res, err := pipeline.Run(rc.ctx, V3Pipeline{}, rc, &pipeline.Plan{DryRun: true})
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(res.Steps))
	for _, s := range res.Steps {
		names = append(names, s.Name)
	}
	return names, nil
}
