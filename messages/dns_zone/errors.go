package dns_zone

import (
	"errors"
)

var (
	ErrorGetDNSZone    = errors.New("Failed to describe the DNS zone: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorCreateDNSZone = errors.New("Failed to create the DNS zone: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorUpdateDNSZone = errors.New("Failed to update the DNS zone: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorListDNSZone   = errors.New("Failed to list your DNS zones: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorFailToDelete  = errors.New("Failed to delete the DNS zone: %s. Check your settings and try again. If the error persists, contact Azion support.")

	ErrorConvertIdZone = errors.New("The DNS zone ID you provided is invalid. The value must be an integer. You may run the 'azion list dns-zone' command to check your DNS zones' IDs")
)
