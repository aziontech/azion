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
// Constructing the context builds API clients but performs no request, so this
// is safe to call without touching the network.
//
// It exists so that `azion deploy --dry-run` can describe the manifest from the
// same step table the deploy uses, instead of from a second copy of the rules.
func DryRunSteps(
	f *cmdutil.Factory,
	conf *contracts.AzionApplicationOptions,
	manifest *contracts.ManifestV4,
	projectConf string,
) ([]string, error) {
	msgs := []string{}
	noWrite := func(*contracts.AzionApplicationOptions, string) error { return nil }

	rc := NewResourceContext(f, conf, manifest, projectConf, &msgs, noWrite)

	res, err := pipeline.Run(rc.Ctx, V4Pipeline{}, rc, &pipeline.Plan{DryRun: true})
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(res.Steps))
	for _, s := range res.Steps {
		names = append(names, s.Name)
	}
	return names, nil
}
