package delete

import "errors"

var (
	ErrorMissingAzionJson         = errors.New("azion.json file is missing. Please initialize and deploy your project before using config delete")
	ErrorMissingApplicationId     = errors.New("application ID is missing from azion.json. Please deploy your project before using config delete")
	ErrorFailedDeleteRulesEngine  = errors.New("failed to delete Rules Engine rule: %w")
	ErrorFailedDeleteCacheSetting = errors.New("failed to delete Cache Setting: %w")
	ErrorFailedDeleteFuncInstance = errors.New("failed to delete Function Instance: %w")
	ErrorFailedDeleteApplication  = errors.New("failed to delete Application: %w")
	ErrorFailedDeleteFirewall     = errors.New("failed to delete Firewall: %w")
	ErrorFailedDeleteFunction     = errors.New("failed to delete Function: %w")
	ErrorFailedDeleteWorkload     = errors.New("failed to delete Workload: %w")
	ErrorFailedDeleteBucket       = errors.New("failed to delete Storage Bucket: %w")
	ErrorFailedDeleteConnector    = errors.New("failed to delete Connector: %w")
	ErrorFailedDeleteFwRule       = errors.New("failed to delete Firewall Rule: %w")
	ErrorFailedUpdateAzionJson    = errors.New("failed to reset azion.json file after deletion")
	ErrorDeletionAborted          = errors.New("deletion aborted by user")
	ErrorPartialDeletion          = errors.New("deletion completed with %d error(s). See output above for details")
)
