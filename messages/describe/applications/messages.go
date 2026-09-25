package applications

// Used by both the v3 and the v4 command trees.
var (
	ShortDescription      = "Returns the Application data"
	LongDescription       = "Displays information about the Application via a given ID to show the application’s attributes in detail"
	HelpFlag              = "Displays more information about the describe command"
	FlagId                = "Unique identifier of the Application"
	AskInputApplicationID = "Enter the Application's ID:"
)

// Used only by the v4 command tree.
var (
	Usage = "application"
)
