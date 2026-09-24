package domain

import "errors"

// Used only by the v3 command tree (bug-fix-only; see doc/plan.md).
var (
	ErrorConvertId          = errors.New("The Domain ID you provided is invalid. The value must be an integer. You may run the 'azion list domains' command to check your Domain ID")
	ErrorFailToDeleteDomain = errors.New("Failed to delete the Domain: %s. Check your settings and try again. If the error persists, contact Azion support")
)
