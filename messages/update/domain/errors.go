package domain

import "errors"

// Used only by the v3 command tree (bug-fix-only; see doc/plan.md).
var (
	ErrorUpdateDomain    = errors.New("Failed to update the Domain: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorActiveFlag      = errors.New("Invalid --active flag provided. The flag must have  'true' or 'false' values. Run the command 'azion update domains --help' to display more information and try again.")
	ErrorConvertDomainID = errors.New("The domain ID you provided is invalid. The value must be an integer. You may run the 'azion list domains' command to check your domain ID")
)
