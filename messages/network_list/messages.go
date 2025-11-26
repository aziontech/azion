package networklist

var (
	Usage = "network-list"

	ListShortDescription = "Displays your Network Lists in a list"
	ListLongDescription  = "Displays all your Network Lists in a list"
	ListHelpFlag         = "Displays more information about the 'list network-list' command"

	DeleteShortDescription = "Deletes a Network List"
	DeleteLongDescription  = "Deletes a Network List based on a given ID"
	DeleteHelpFlag         = "Displays more information about the 'delete network-list' command"
	DeleteOutputSuccess    = "Network List %s was successfully deleted"

	DescribeShortDescription = "Displays a Network List"
	DescribeLongDescription  = "Displays a Network List based on a given ID"
	DescribeHelpFlag         = "Displays more information about the 'describe network-list' command"

	CreateShortDescription = "Creates a Network List"
	CreateLongDescription  = "Creates a Network List based on given attributes"
	CreateHelpFlag         = "Displays more information about the 'create network-list' command"
	CreateOutputSuccess    = "Created Network List with ID %d"

	UpdateShortDescription = "Updates a Network List"
	UpdateLongDescription  = "Updates a Network List based on a given ID"
	UpdateHelpFlag         = "Displays more information about the 'update network-list' command"
	UpdateOutputSuccess    = "Updated Network List with ID %d"
	UpdateAskNetworkListID = "Enter the Network List's ID:"

	AskNetworkListID = "Enter the Network List's ID:"
	AskName          = "Enter the Network List's name:"
	AskType          = "Select the Network List type:"
	AskItems         = "Enter the items (comma-separated):"
	AskActive        = "Should the Network List be active?"

	FlagID     = "Unique identifier of the Network List"
	FlagName   = "Name of the Network List"
	FlagType   = "Type of the Network List (asn, countries, ip_cidr)"
	FlagItems  = "Items for the Network List (comma-separated)"
	FlagActive = "Whether the Network List is active"
	FlagIn     = "Path to a JSON file containing the Network List attributes"
)
