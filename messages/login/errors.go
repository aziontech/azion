package login

import "errors"

var (
	ErrorLogin              = "Failed to Login: %s. Check your settings and try again. If the error persists, contact Azion support."
	ErrorInvalidLogin       = errors.New("Invalid login method")
	ErrorTokenCreateInvalid = "Invalid token detected. The generated token appears to be corrupted or expired. Please check your authentication credentials and generate a new token to proceed."
)
