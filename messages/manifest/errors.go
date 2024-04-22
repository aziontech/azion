package manifest

import "errors"

var (
	ErrorUnmarshalAzionJsonFile = errors.New("Failed to parse the given 'azion.json' file. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorCacheNotFound          = errors.New("Could not Find this cache setting")
)