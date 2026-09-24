package dnssec

// Used only by the v4 command tree.
var (
	// [ dnssec ]
	DNSSECUsage = "dnssec"

	DNSSECDescribeShortDescription = "Returns the DNSSEC information of a specific DNS zone"
	DNSSECDescribeLongDescription  = "Returns the DNSSEC information of a specific DNS zone, informed through the flag '--zone-id', in detail"
	DNSSECDescribeHelpFlag         = "Displays more information about the describe subcommand"
	DNSSECDescribeAskInputZoneID   = "Enter the ID of the DNS zone whose DNSSEC you wish to describe:"

	DNSSECUpdateShortDescription = "Updates the DNSSEC of a DNS zone"
	DNSSECUpdateLongDescription  = "Enables or disables DNSSEC for a DNS zone based on the given attributes"
	DNSSECUpdateFlagEnabled      = "Whether DNSSEC should be enabled for the DNS zone"
	DNSSECUpdateFlagIn           = "Path to a JSON file containing the attributes of the DNSSEC that will be updated; you can use - for reading from stdin"
	DNSSECUpdateOutputSuccess    = "DNSSEC of DNS zone %d was updated\n"
	DNSSECUpdateHelpFlag         = "Displays more information about the update subcommand"
	DNSSECUpdateAskInputZoneID   = "Enter the ID of the DNS zone whose DNSSEC you wish to update:"
	DNSSECUpdateAskInputEnabled  = "Enter whether DNSSEC should be enabled (true/false):"

	// [ flags ]
	DNSSECFlagZoneID = "Unique identifier for a DNS zone. The '--zone-id' flag is required"
)
