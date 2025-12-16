package functioninstance

import "errors"

var (
	ErrorIsActiveFlag              = errors.New("Invalid --active flag provided. The value must be 'true' or 'false'. Run the command 'azion update function-instance --help' to display more information and try again")
	ErrorUpdate                    = errors.New("Failed to update the Function Instance: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorArgsFlag                  = errors.New("Failed to read the args file. Verify if the file name and its path are correct and the file's content has a valid JSON format. Run the command again using the --help flag to display more information and try again")
	ErrorParseArgs                 = errors.New("Failed to parse JSON args. Verify if the file's content has a valid JSON format. Run the command again using the --help flag to display more information and try again")
	ErrorConvertApplicationId      = errors.New("Invalid --application-id flag provided. The value must be an integer. Run the command 'azion update function-instance --help' to display more information and try again")
	ErrorConvertFunctionInstanceId = errors.New("Invalid --instance-id flag provided. The value must be an integer. Run the command 'azion update function-instance --help' to display more information and try again")
)
