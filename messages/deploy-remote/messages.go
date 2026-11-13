package deployremote

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	// deploy cmd
	DeployUsage                          = "deploy-remote"
	DeployShortDescription               = "Deploys an Application"
	DeployLongDescription                = "Deploys an Application"
	DeploySuccessful                     = "Your Application was deployed successfully\n"
	DeployOutputWorkloadSuccess          = "\nTo visualize your application access the Domain: %s\n"
	EdgeApplicationDeployPathFlag        = "Path to where your static files are stored"
	EdgeApplicationDeployProjectConfFlag = "Relative path to where your custom azion.json and args.json files are stored"
	EdgeApplicationDeploySync            = "Synchronizes the local azion.json file with remote resources"
	EnvFlag                              = "Relative path to where your custom .env file is stored"
	GeneratingManifest                   = "manifest.json not found. Generating azion.config and manifest.json before deploying\n"
	AliasEnvFlag                         = "If sent, the --alias-env option is forwarded to Bundler during the build phase"
	DeployFlagHelp                       = "Displays more information about the deploy command"
	DeployFlagAuto                       = "If sent, the entire flow of the command will be run without interruptions"
	DeployFlagNoPrompt                   = "If sent, whenever the CLI would display an interactive prompt due to an error, it instead just returns the error"
	DeployPropagation                    = "Your application is being deployed to all Azion Edge Locations and it might take a few minutes.\n"
	SkipBucket                           = "Your project does not contain a '.edge/storage' folder. Skipping creation of bucket"
	SkipUpload                           = "Your project does not contain a '.edge/storage' folder. Skipping upload of static files"
	SkipUploadBuild                      = "Skipping upload of static files due to project not being built and no new static files being generated"
)

// Used only by the v3 command tree (bug-fix-only; see doc/plan.md).
var (
	DeployOutputDomainSuccess = "\nTo visualize your application access the Domain: %s\n"
	DeployFlagSkipBuild       = "If sent, the build command will not be called during the deploy process"
)

// ============================================================================
// v6 — used only by the v6 deploy (pkg/cmd/deploy_remote/v6).
// ============================================================================
var (
	V6RunningDeploy       = "Running deploy command (API v6)\n"
	V6FlagEnvironment     = "Deploys only to this Environment: releases the Deployments bound to it and updates only its bindings. On a first run it also names the preset's Environment (default: Production)"
	V6FlagNoActivate      = "If sent, the release is built but not activated: traffic keeps going to the current release"
	V6FlagTimeout         = "How long to wait for each build or activation, for example 90s or 5m"
	V6FlagDryRun          = "If sent, shows what the deploy would do without calling the API"
	V6BuildTransition     = "%s %q (%s): %s → %s (%s)\n"
	V6BuildsStarting      = "Building %d version(s): %s\n"
	V6BuildSucceeded      = "%s %q built successfully: version %s is ready (%s)\n"
	V6EnvironmentResolved = "Deploying to Environment %q (%s, %s)\n"
	V6DeploymentCreated   = "Deployment %q created (%s)\n"
	V6DeploymentRecreate  = "Deployment %s no longer exists; creating a new one\n"
	V6ReleaseCreated      = "Release %s created on Deployment %q\n"
	V6ReleaseReady        = "Release %s is ready\n"
	V6ReleaseNotActivated = "Release %s was built but not activated (--no-activate). Traffic still goes to the previous release\n"
	V6ReleaseActivated    = "Release %s is now %s\n"
	V6WorkloadSaved       = "Workload %q (id %d) now serves Deployment %q in Environment %q\n"
	V6DeploySuccessful    = "Your Application was deployed successfully\n"
	V6DeployOutputURL     = "\nTo visualize your application access: %s\n"
	V6Interrupted         = "Deploy interrupted %s activation. Builds started by this run were asked to stop\n"
	V6InterruptedBefore   = "before"
	V6InterruptedAfter    = "after"
	V6ActiveButIncomplete = "The new release is live; what failed was %s. Run the deploy again to finish it\n"
	V6DryRunStaleManifest = "The current manifest cannot be used (%v). The deploy builds it again from azion.config, so this report has no resource counts\n"
	V6DryRunHeader        = "Dry run: nothing was sent to the API. The deploy would:\n"
	V6DryRunStep          = "  - %s\n"
	V6DryRunSkipped       = "Skipped: %s\n"
	V6TimingHeader        = "\n=== Deploy timing (API v6) ===\n"
	V6RemoteNotAvailable  = "The remote deploy is not available for API v6 accounts yet; running the deploy locally instead\n"
)

var (
	ErrorV6EnvironmentNotFound = errors.New("the Environment was not found in this account. Check its name, or create it in the Console")
	ErrorV6ReleaseFailed       = errors.New("the release build did not produce a usable release; nothing changed in traffic")
	ErrorV6ActivationFailed    = errors.New("the release could not be activated; traffic still goes to the previous release")
	ErrorV6DeploymentFailed    = errors.New("the Deployment did not become ready")
	ErrorV6VerifyFailed        = errors.New("the deploy is live, but the check after activation found differences")
)

// v6, continued.
var (
	V6EnvironmentCreated   = "Environment %q created (%s)\n"
	V6WaitingToBind        = "Waiting for %s %q (%s) before binding it to the Workload: it is %s\n"
	ErrorV6NotBindable     = errors.New("it cannot be bound to a Workload")
	ErrorV6LocalOnlyFlags  = errors.New("--environment, --no-activate and --timeout apply only to a deploy that runs here. Add --local to use them")
	ErrorV6PolicyImmutable = errors.New("an Environment's deployment policy cannot change. Use the policy it has, or another Environment")
	V6WorkloadBuilding     = "Waiting for the Workload %q (id %d) to build its new version: it is %s\n"
	V6WorkloadReady        = "Workload %q (id %d) is ready (%s)\n"
	ErrorV6WorkloadFailed  = errors.New("the Workload did not build its new version, so its domains do not serve the new bindings. Check the Workload in the Console, then run the deploy again")
)
