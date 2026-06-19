package dnssec

import (
	"errors"
)

var (
	ErrorGetDNSSEC    = errors.New("Failed to describe the DNSSEC: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorUpdateDNSSEC = errors.New("Failed to update the DNSSEC: %s. Check your settings and try again. If the error persists, contact Azion support.")

	ErrorConvertIdZone  = errors.New("The DNS zone ID you provided is invalid. The value must be an integer. You may run the 'azion list dns-zone' command to check your DNS zones' IDs")
	ErrorConvertEnabled = errors.New("The value provided for '--enabled' is invalid. The value must be a boolean (true or false)")
)
