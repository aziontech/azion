package warmup

import "errors"

// Used only by the v4 command tree.
var (
	ErrorInvalidUrl = errors.New("Invalid URL provided. URL must be a valid HTTP/HTTPS URL")
)
