package delete

const (
	Usage            = "delete"
	ShortDescription = "Delete all resources from azion.json"
	LongDescription  = "Deletes all resources found in the azion.json configuration file and resets it to a fresh state"
	FlagHelp         = "Displays more information about the delete command"
	FlagConfigDir    = "Relative path to where your azion.json file is stored (default: azion)"
	FlagForce        = "Force deletion without confirmation prompt"

	DeletingResources  = "Deleting all resources from azion.json\n"
	DeleteSuccess      = "All remote resources deleted successfully\n"
	ResettingConfig    = "Resetting azion.json to initial state\n"
	ConfigResetSuccess = "azion.json has been reset to initial state\n"

	// Resource deletion messages
	DeletingRulesEngineApp  = "Deleting Rules Engine (Application) rule '%s' (ID: %d)\n"
	DeletingRulesEngineFw   = "Deleting Rules Engine (Firewall) rule '%s' (ID: %d) from firewall %d\n"
	DeletingFuncInstanceApp = "Deleting Function Instance (Application) '%s' (ID: %d)\n"
	DeletingFuncInstanceFw  = "Deleting Function Instance (Firewall) '%s' (ID: %d) from firewall %d\n"
	DeletingCacheSetting    = "Deleting Cache Setting '%s' (ID: %d)\n"
	DeletingApplication     = "Deleting Application '%s' (ID: %d)\n"
	DeletingFirewall        = "Deleting Firewall '%s' (ID: %d)\n"
	DeletingFunction        = "Deleting Function '%s' (ID: %d)\n"
	DeletingWorkload        = "Deleting Workload '%s' (ID: %d)\n"
	DeletingBucket          = "Deleting Storage Bucket '%s'\n"
	DeletingConnector       = "Deleting Connector '%s' (ID: %d)\n"

	// Confirmation
	ConfirmDeletion = "This will delete all resources defined in azion.json. Continue? (y/N): "
	DeletionAborted = "Deletion aborted by user\n"

	// Summary
	DeletionSummary      = "\nDeletion Summary:\n"
	ResourcesDeleted     = "  Resources deleted successfully: %d\n"
	ResourcesFailed      = "  Resources failed to delete: %d\n"
	ErrorsDuringDeletion = "\nErrors occurred during deletion:\n"
)
