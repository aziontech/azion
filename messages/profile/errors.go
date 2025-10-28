package profile

import "errors"

var (
	ErrorCreate        = errors.New("Failed to create the profile: %w")
	ErrorReadFile      = errors.New("Failed to read the file: %w")
	ErrorUnmarshalFile = errors.New("Failed to unmarshal the file: %w")
	ErrorReadDir       = errors.New("Failed to read the directory: %w")
	ErrorSwitchProfile = errors.New("Failed to switch profile: %w")
	ErrorDeleteProfile = errors.New("Failed to delete the profile: %w")
	ErrorProfileNotFound = errors.New("Profile '%s' not found")
	ErrorCannotDeleteDefault = errors.New("Cannot delete the 'default' profile")
	ErrorDeleteToken   = errors.New("Failed to delete token: %w")
	ErrorDeleteCancelled = errors.New("Profile deletion cancelled")
)
