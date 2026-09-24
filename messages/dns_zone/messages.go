package dns_zone

// Used only by the v4 command tree.
var (
	// [ dns zone ]
	DNSZoneUsage = "dns-zone"

	DNSZoneListShortDescription = "Displays your DNS zones"
	DNSZoneListLongDescription  = "Displays all DNS zones in your account"
	DNSZoneListHelpFlag         = "Displays more information about the list subcommand"

	DNSZoneDescribeShortDescription = "Returns the information related to a specific DNS zone"
	DNSZoneDescribeLongDescription  = "Returns the information related to a specific DNS zone, informed through the flag '--zone-id', in detail"
	DNSZoneDescribeHelpFlag         = "Displays more information about the describe subcommand"

	DNSZoneCreateShortDescription = "Creates a new DNS zone"
	DNSZoneCreateLongDescription  = "Creates a DNS zone based on given attributes"
	DNSZoneCreateFlagName         = "The name of your DNS zone"
	DNSZoneCreateFlagDomain       = "The domain associated with the DNS zone"
	DNSZoneCreateFlagActive       = "Whether the DNS zone should be active"
	DNSZoneCreateFlagIn           = "Path to a JSON file containing the attributes of the DNS zone that will be created; you can use - for reading from stdin"
	DNSZoneCreateOutputSuccess    = "Created DNS zone with ID %d\n"
	DNSZoneCreateHelpFlag         = "Displays more information about the create subcommand"

	DNSZoneUpdateShortDescription = "Updates a DNS zone"
	DNSZoneUpdateLongDescription  = "Updates a DNS zone based on given attributes"
	DNSZoneUpdateFlagName         = "The name of your DNS zone"
	DNSZoneUpdateFlagActive       = "Whether the DNS zone should be active"
	DNSZoneUpdateFlagIn           = "Path to a JSON file containing the attributes of the DNS zone that will be updated; you can use - for reading from stdin"
	DNSZoneUpdateOutputSuccess    = "DNS zone %d was updated\n"
	DNSZoneUpdateHelpFlag         = "Displays more information about the update subcommand"

	DNSZoneDeleteShortDescription = "Deletes a DNS zone"
	DNSZoneDeleteLongDescription  = "Deletes a DNS zone based on the given '--zone-id'"
	DNSZoneDeleteOutputSuccess    = "DNS zone %d was successfully deleted\n"
	DNSZoneDeleteHelpFlag         = "Displays more information about the delete subcommand"

	// [ flags ]
	DNSZoneFlagId = "Unique identifier for a DNS zone. The '--zone-id' flag is required"

	// [ ask input prompts ]
	DNSZoneCreateAskInputName     = "Enter the new DNS zone's name:"
	DNSZoneCreateAskInputDomain   = "Enter the new DNS zone's domain:"
	DNSZoneDescribeAskInputZoneID = "Enter the ID of the DNS zone you wish to describe:"
	DNSZoneUpdateAskInputZoneID   = "Enter the ID of the DNS zone you wish to update:"
	DNSZoneDeleteAskInputZoneID   = "Enter the ID of the DNS zone you wish to delete:"
)
