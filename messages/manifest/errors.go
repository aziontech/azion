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

// ============================================================================
// v6 — used only by the v6 deploy (pkg/manifest/v6).
// ============================================================================
var (
	ErrorV6DraftTaken            = errors.New("the Draft is no longer editable: another process (the Console or another CLI run) started its build during this deploy. Run the deploy again")
	ErrorV6UnknownRef            = errors.New("the manifest refers to a resource it does not declare")
	ErrorV6UnknownDependencyType = errors.New("a Function can only depend on a connector or a network_list")
	ErrorV6TooMany               = errors.New("a release serves at most one Application, one Firewall and one Custom Page")
	ErrorV6NoName                = errors.New("every resource in the manifest needs a name")
	ErrorV6DuplicateName         = errors.New("two resources of the same kind share a name in the manifest")
	ErrorV6FunctionCode          = errors.New("failed to read the code of a Function")
	ErrorV6InvalidRule           = errors.New("a rule in the manifest is not valid")
)

// v6, continued.
var (
	ErrorV6Internal       = errors.New("the deploy stopped on an internal error. Run it again with --debug and report the output")
	ErrorV6WrongVersion   = errors.New("azion.config is not a v6 config. This account uses API v6: set 'version: 6' in azion.config and use the v6 resource shapes, or create the project again with 'azion init'")
	ErrorV6UnknownBinding = errors.New("a workload binding refers to an Environment or Deployment the manifest does not declare")
)
