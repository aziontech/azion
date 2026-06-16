package dns_record

import (
	"errors"
)

var (
	ErrorGetDNSRecord    = errors.New("Failed to describe the DNS record: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorCreateDNSRecord = errors.New("Failed to create the DNS record: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorUpdateDNSRecord = errors.New("Failed to update the DNS record: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorListDNSRecord   = errors.New("Failed to list your DNS records: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorFailToDelete    = errors.New("Failed to delete the DNS record: %s. Check your settings and try again. If the error persists, contact Azion support.")

	ErrorConvertIdZone   = errors.New("The DNS zone ID you provided is invalid. The value must be an integer. You may run the 'azion list dns-zone' command to check your DNS zones' IDs")
	ErrorConvertIdRecord = errors.New("The DNS record ID you provided is invalid. The value must be an integer. You may run the 'azion list dns-record' command to check your DNS records' IDs")
)
