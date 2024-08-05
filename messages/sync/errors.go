package sync

import "errors"

var (
	ERRORSYNC        = "Failed to synchronize local resources with remote resources: %s"
	ERRORNOTDEPLOYED = errors.New("Failed to synchronize local resources with remote resources: You must deploy your project at least once before trying to synchronize with remote resources")
)
