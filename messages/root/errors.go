package root

import "errors"

var (
	ErrorCurrentUser       = errors.New("Failed to get current user's information.")
	ErrorMarshalUserInfo   = errors.New("Failed to marshal current user information.")
	ErrorUnmarshalUserInfo = errors.New("Failed to unmarshal current user information.")
)
