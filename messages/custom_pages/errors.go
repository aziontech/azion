package custompages

import "errors"

var (
	ErrorGetCustomPage          = errors.New("Failed to get the Custom Page: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorGetCustomPages         = errors.New("Failed to list the Custom Pages: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorFailToDeleteCustomPage = errors.New("Failed to delete the Custom Page: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateCustomPage       = errors.New("Failed to create Custom Page: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateCustomPage       = errors.New("Failed to update the Custom Page: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorConvertCustomPageId    = errors.New("Invalid --custom-page-id flag provided. The value must be an integer. Run the command 'azion describe custom-pages --help' to display more information and try again")
)
