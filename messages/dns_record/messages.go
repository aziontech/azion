package dns_record

var (
	// [ dns record ]
	DNSRecordUsage            = "dns-record"
	DNSRecordShortDescription = "Manages records of an Intelligent DNS zone"
	DNSRecordLongDescription  = "Manages the DNS records hosted in a specific Intelligent DNS zone on Azion's edge network."
	DNSRecordFlagHelp         = "Displays more information about the dns-record command"

	// [ list ]
	DNSRecordListUsage            = "list [flags]"
	DNSRecordListShortDescription = "Displays the records of a DNS zone"
	DNSRecordListLongDescription  = "Displays all records related to a specific DNS zone"
	DNSRecordListHelpFlag         = "Displays more information about the list subcommand"

	// [ describe ]
	DNSRecordDescribeUsage            = "describe --zone-id <zone_id> --record-id <record_id> [flags]"
	DNSRecordDescribeShortDescription = "Returns the information related to a specific DNS record"
	DNSRecordDescribeLongDescription  = "Returns the information related to a specific DNS record, informed through the flag '--record-id', in detail"
	DNSRecordDescribeFlagOut          = "Exports the output of the subcommand 'describe' to the given file path <file_path/file_name.ext>"
	DNSRecordDescribeFlagFormat       = "Changes the output format passing the json value to the flag. Example '--format json'"
	DNSRecordDescribeHelpFlag         = "Displays more information about the describe subcommand"
	DNSRecordFileWritten              = "File successfully written to: %s\n"

	// [ create ]
	DNSRecordCreateUsage            = "create [flags]"
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

	// [ update ]
	DNSRecordUpdateUsage            = "update [flags]"
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

	// [ delete ]
	DNSRecordDeleteUsage            = "delete [flags]"
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
