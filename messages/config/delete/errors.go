package delete

import "errors"

// Used only by the v4 command tree.
var (
	ErrorMissingAzionJson = errors.New("azion.json file is missing. Please initialize and deploy your project before using config delete")
	ErrorPartialDeletion  = errors.New("deletion completed with %d error(s). See output above for details")
)
