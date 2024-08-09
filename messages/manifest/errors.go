package manifest

import "errors"

var (
	ErrorUnmarshalAzionJsonFile = errors.New("Failed to parse the given 'azion.json' file. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorCacheNotFound          = errors.New("Could not find this cache setting")
	ErrorFunctionNotFound       = errors.New("Could not find this edge function")
	ErrorOriginNotFound         = errors.New("Could not find this origin")
	ErrorCreateOrigin           = errors.New("Failed to create the origin")
	ErrorCreateCache            = errors.New("Failed to create the cache setting")
	ErrorCreateRule             = errors.New("Failed to create the rule in Rules Engine")
	ErrorUpdateOrigin           = errors.New("Failed to update the origin")
	ErrorUpdateDomain           = errors.New("Failed to update the domain")
	ErrorCreateDomain           = errors.New("Failed to create the domain")
	ErrorUpdateCache            = errors.New("Failed to update the cache setting")
	ErrorUpdateRule             = errors.New("Failed to update the rule in Rules Engine")
)
