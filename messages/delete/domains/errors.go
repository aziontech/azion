package domains

import "errors"

var (
	ErrorConvertId          = errors.New("The domain ID you provided is invalid. The value must be an integer. You may run the 'azion list domains' command to check your domain ID")
	ErrorFailToDeleteDomain = errors.New("Failed to delete the Domain: %s. Check your settings and try again. If the error persists, contact Azion support")
)
