package functioninstance

import "errors"

var (
	ErrorGetFunctionInstance       = "Error getting Function Instance: %s"
	ErrorConvertApplicationId      = errors.New("Invalid --application-id flag provided. The value must be an integer. Run the command 'azion describe function-instance --help' to display more information and try again")
	ErrorConvertFunctionInstanceId = errors.New("Invalid --instance-id flag provided. The value must be an integer. Run the command 'azion describe function-instance --help' to display more information and try again")
)
