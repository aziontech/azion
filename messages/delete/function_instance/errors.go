package functioninstance

import "errors"

var (
	ErrorFailToDeletInstance       = errors.New("Failed to delete Function Instance: %s")
	ErrorConvertApplicationId      = errors.New("Invalid --application-id flag provided. The value must be an integer. Run the command 'azion delete function-instance --help' to display more information and try again")
	ErrorConvertFunctionInstanceId = errors.New("Invalid --instance-id flag provided. The value must be an integer. Run the command 'azion delete function-instance --help' to display more information and try again")
)
