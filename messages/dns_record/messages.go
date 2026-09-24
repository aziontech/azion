package dns_record

// Used only by the v4 command tree.
var (
	// [ dns record ]
	DNSRecordUsage = "dns-record"

	DNSRecordListShortDescription = "Displays the records of a DNS zone"
	DNSRecordListLongDescription  = "Displays all records related to a specific DNS zone"
	DNSRecordListHelpFlag         = "Displays more information about the list subcommand"

	DNSRecordDescribeShortDescription = "Returns the information related to a specific DNS record"
	DNSRecordDescribeLongDescription  = "Returns the information related to a specific DNS record, informed through the flag '--record-id', in detail"
	DNSRecordDescribeHelpFlag         = "Displays more information about the describe subcommand"

	DNSRecordCreateShortDescription = "Creates a new DNS record"
	DNSRecordCreateLongDescription  = "Creates a DNS record in a given DNS zone based on given attributes"
	DNSRecordCreateFlagName         = "The name (entry) of the DNS record"
	DNSRecordCreateFlagType         = "The type of the DNS record (e.g. A, AAAA, CNAME, MX, TXT, NS)"
	DNSRecordCreateFlagRdata        = "The value(s) of the DNS record; repeat the flag or pass a comma-separated list for multiple values"
	DNSRecordCreateFlagTTL          = "The time to live (TTL) of the DNS record, in seconds"
	DNSRecordCreateFlagPolicy       = "The routing policy of the DNS record. Must be 'simple' or 'weighted'"
	DNSRecordCreateFlagWeight       = "The weight of the DNS record; only used when '--policy' is 'weighted'"
	DNSRecordCreateFlagDescription  = "A description for the DNS record; only used when '--policy' is 'weighted'"
	DNSRecordCreateFlagIn           = "Path to a JSON file containing the attributes of the DNS record that will be created; you can use - for reading from stdin"
	DNSRecordCreateOutputSuccess    = "Created DNS record with ID %d\n"
	DNSRecordCreateHelpFlag         = "Displays more information about the create subcommand"

	DNSRecordUpdateShortDescription = "Updates a DNS record"
	DNSRecordUpdateLongDescription  = "Updates a DNS record based on given attributes"
	DNSRecordUpdateFlagName         = "The name (entry) of the DNS record"
	DNSRecordUpdateFlagType         = "The type of the DNS record (e.g. A, AAAA, CNAME, MX, TXT, NS)"
	DNSRecordUpdateFlagRdata        = "The value(s) of the DNS record; repeat the flag or pass a comma-separated list for multiple values"
	DNSRecordUpdateFlagTTL          = "The time to live (TTL) of the DNS record, in seconds"
	DNSRecordUpdateFlagPolicy       = "The routing policy of the DNS record. Must be 'simple' or 'weighted'"
	DNSRecordUpdateFlagWeight       = "The weight of the DNS record; only used when '--policy' is 'weighted'"
	DNSRecordUpdateFlagDescription  = "A description for the DNS record; only used when '--policy' is 'weighted'"
	DNSRecordUpdateFlagIn           = "Path to a JSON file containing the attributes of the DNS record that will be updated; you can use - for reading from stdin"
	DNSRecordUpdateOutputSuccess    = "DNS record %d was updated\n"
	DNSRecordUpdateHelpFlag         = "Displays more information about the update subcommand"

	DNSRecordDeleteShortDescription = "Deletes a DNS record"
	DNSRecordDeleteLongDescription  = "Deletes a DNS record based on the given '--zone-id' and '--record-id'"
	DNSRecordDeleteOutputSuccess    = "DNS record %d was successfully deleted\n"
	DNSRecordDeleteHelpFlag         = "Displays more information about the delete subcommand"

	// [ flags ]
	DNSRecordFlagZoneID   = "Unique identifier for the DNS zone that hosts this record. The '--zone-id' flag is required"
	DNSRecordFlagRecordID = "Unique identifier for a DNS record. The '--record-id' flag is required"

	// [ ask input prompts ]
	DNSRecordCreateAskInputZoneID     = "Enter the ID of the DNS zone the record will be created in:"
	DNSRecordCreateAskInputName       = "Enter the new DNS record's name:"
	DNSRecordCreateAskInputType       = "Enter the new DNS record's type (e.g. A, CNAME, TXT):"
	DNSRecordCreateAskInputRdata      = "Enter the new DNS record's value(s) (comma-separated for multiple):"
	DNSRecordListAskInputZoneID       = "Enter the ID of the DNS zone the records are linked to:"
	DNSRecordDescribeAskInputZoneID   = "Enter the ID of the DNS zone the record is linked to:"
	DNSRecordDescribeAskInputRecordID = "Enter the ID of the DNS record you wish to describe:"
	DNSRecordUpdateAskInputZoneID     = "Enter the ID of the DNS zone the record is linked to:"
	DNSRecordUpdateAskInputRecordID   = "Enter the ID of the DNS record you wish to update:"
	DNSRecordDeleteAskInputZoneID     = "Enter the ID of the DNS zone the record is linked to:"
	DNSRecordDeleteAskInputRecordID   = "Enter the ID of the DNS record you wish to delete:"
)
