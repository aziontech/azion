package domain

// Used only by the v3 command tree (bug-fix-only; see doc/plan.md).
var (
	Usage            = "domain"
	ShortDescription = "Returns the Domain data"
	LongDescription  = "Displays information about the Domain via a given ID to show the application’s attributes in detail"
	HelpFlag         = "Displays more information about the describe command"
	FlagDomainID     = "Unique identifier of the Domain"
	AskInputDomainID = "Enter the Domain's ID:"
)
