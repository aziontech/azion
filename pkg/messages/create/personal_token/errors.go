package personaltoken

import "errors"

var (
	ErrorCreate            = errors.New("Failed to create the Personal Token: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMissingExpiration = errors.New("Failed to create the Personal Token: You must provide an expiration value.")
)
