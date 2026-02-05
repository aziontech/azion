package firewall

import "errors"

var (
	ErrorConvertId            = errors.New("The firewall ID you provided is invalid. The value must be an integer. You may run the 'azion list firewall' command to check your firewall ID")
	ErrorFailToDeleteFirewall = "Failed to delete the Firewall: %s. Check your settings and try again. If the error persists, contact Azion support"
)
