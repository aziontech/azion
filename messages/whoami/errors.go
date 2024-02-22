package whoami

import "errors"

var (
	ErrorNotLoggedIn = errors.New("You must be logged in to use the 'whoami' command")
)
