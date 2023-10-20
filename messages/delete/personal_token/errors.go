package personaltoken

import "errors"

var (
	ErrorFailToDelete = errors.New("Failed to delete the personal token: %s. Check your settings and try again. If the error persists, contact Azion support")
)
