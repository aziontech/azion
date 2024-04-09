package purge

import "errors"

var (
	ErrorTooManyUrls = errors.New("Only one item is allowed for the Wildcard option")
)
