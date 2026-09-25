package whoami

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorNotLoggedIn = errors.New("You must be logged in to use the 'whoami' command")
)
