package functioninstance

import "errors"

// Used only by the v4 command tree.
var (
	ErrorGetAll               = "Error getting Function Instances: %s. Check your settings and try again. If the error persists, contact Azion support."
	ErrorConvertApplicationId = errors.New("Invalid --application-id flag provided. The value must be an integer. Run the command 'azion list function-instance --help' to display more information and try again")
)
