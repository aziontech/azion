package manifest

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorCacheNotFound         = errors.New("Could not find this cache setting")
	ErrorConnectorNotFound     = errors.New("Could not find this connector")
	ErrorReadCodeFile          = errors.New("Failed to read target code file: %w")
	ErrorInvalidPhase          = errors.New("Invalid phase. Please use 'request' or 'response'")
	ErrorFuncNotFound          = errors.New("The Function Name informed does not exists. Please make sure to add this Function to your azion.config file")
	ErrorUnmarshalArgsFile     = errors.New("Failed to unmarshal args.json file: %w")
	ErrorApplicationIDRequired = errors.New("Application ID is required for this operation")
	ErrorWorkloadIDRequired    = errors.New("Workload ID is required for this operation")
	ErrorConnectorTypeNotFound = errors.New("Failed to determine connector type")
	ErrorReadManifest          = "Failed to read the manifest.json file: %w. Please remember to install dependencies and build your project before running the deploy command"
)

// Used only by the v3 command tree (bug-fix-only; see doc/plan.md).
var (
	ErrorOriginNotFound = errors.New("Could not find this origin")
	ErrorCreateOrigin   = errors.New("Failed to create the origin")
	ErrorCreateCache    = errors.New("Failed to create the cache setting")
	ErrorCreateRule     = errors.New("Failed to create the rule in Rules Engine")
	ErrorUpdateOrigin   = errors.New("Failed to update the origin")
	ErrorUpdateDomain   = errors.New("Failed to update the domain")
	ErrorUpdateCache    = errors.New("Failed to update the cache setting")
	ErrorUpdateRule     = errors.New("Failed to update the rule in Rules Engine")
)
