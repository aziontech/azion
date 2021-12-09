package utils

import "errors"

var (
	ErrorMissingServiceIdArgument = errors.New("You must provide a service_id as an argument. Use -h or --help for more information.")
)
