package manifest

// Used by both the v3 and the v4 command trees.
const (
	ManifestUpdateCache                    = "Cache Setting %s with id %d successfully updated\n"
	ManifestUpdateRule                     = "Rule Engine %s with id %d successfully updated\n"
	ManifestCreateCache                    = "Cache Setting %s with id %d successfully created\n"
	ManifestCreateRule                     = "Rule Engine %s with id %d successfully created\n"
	ManifestOrderRule                      = "Rules Engine of Application with id %d successfully ordered (%s phase)\n"
	ReadingManifest                        = "Reading manifest.json file\n"
	CreatingManifest                       = "Creating resources found in manifest.json file\n"
	SkipDeletion                           = "Skipping deletion of resources based on configuration found on azion.json file\n"
	AppInUse                               = "This Application's name is already in use, please try another one\n"
	FunctionInUse                          = "This Function's name is already in use, please try another one\n"
	WorkloadInUse                          = "This Workload's name is already in use, please try another one\n"
	FirewallInUse                          = "This Firewall's name is already in use, please try another one\n"
	AskInputName                           = "Type the new name:"
	ManifestCreateFirewall                 = "Firewall %s with id %d successfully created\n"
	ManifestUpdateFirewall                 = "Firewall %s with id %d successfully updated\n"
	ManifestCreateFirewallRule             = "Firewall Rule %s with id %d successfully created\n"
	ManifestUpdateFirewallRule             = "Firewall Rule %s with id %d successfully updated\n"
	ManifestOrderFirewallRule              = "Rules Engine of Firewall with id %d successfully ordered\n"
	ManifestCreateFirewallFunctionInstance = "Firewall Function Instance %s with id %d successfully created\n"
	ManifestUpdateFirewallFunctionInstance = "Firewall Function Instance %s with id %d successfully updated\n"
	ManifestCreateFunction                 = "Function %s with id %d successfully created\n"
	ManifestUpdateFunction                 = "Function %s with id %d successfully updated\n"
	ManifestCreateFunctionInstance         = "Function Instance %s with id %d successfully created\n"
	ManifestUpdateFunctionInstance         = "Function Instance %s with id %d successfully updated\n"
	ManifestCreateEdgeApplication          = "Edge Application %s with id %d successfully created\n"
	ManifestUpdateEdgeApplication          = "Edge Application %s with id %d successfully updated\n"
	ManifestCreateConnector                = "Connector %s with id %d successfully created\n"
	ManifestUpdateConnector                = "Connector %s with id %d successfully updated\n"
	ManifestCreateWorkload                 = "Workload %s with id %d successfully created\n"
	ManifestUpdateWorkload                 = "Workload %s with id %d successfully updated\n"
	ManifestCreateWorkloadDeployment       = "Workload Deployment %s with id %d successfully created\n"
	ManifestUpdateWorkloadDeployment       = "Workload Deployment %s with id %d successfully updated\n"
	ManifestCreateStorage                  = "Storage Bucket %s successfully created\n"
	ManifestUpdateStorage                  = "Storage Bucket %s successfully updated\n"
	ManifestPurgeSuccess                   = "Purge of type %s successfully executed\n"
	MessageDeleteResource                  = `It seems this resource was deleted from a previous version of the application.
One cause may be that the resource is not being used in any rule.
To avoid deleting resources that are not being used, you can add the field 'skip-deletion' to your azion.json file.`
	UpdateAzionConfig = "Updating azion.config file with new resource name\n"
)

// Used only by the v3 command tree (bug-fix-only; see doc/plan.md).
const (
	ManifestCreateOrigin = "Origin %s with id %d successfully created\n"
	ManifestUpdateOrigin = "Origin %s with key %s successfully updated\n"
)
