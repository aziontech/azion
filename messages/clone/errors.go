package clone

import "errors"

var (
	ErrorClone                = errors.New("Failed to clone Application: %s")
	ErrorConvertApplicationId = errors.New("The Application ID you provided is invalid. The value must be an integer. You may run the 'azion list application' command to check your Application ID")
)
