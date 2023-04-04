package edgefunctionsinstances

import "errors"

var (
	ErrorGetFunctions = errors.New("failed to get the Edge Functions Instances: %s. Check your settings and try again. If the error persists, contact Azion support")
)
