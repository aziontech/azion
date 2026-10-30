package application

// Used by the v3, v4 and v6 command trees. The cascade messages are v3 and v4
// only: v6 has no cascade delete until its deploy model is known.
var (
	ShortDescription = "Deletes an Application"
	LongDescription  = "Removes an Application from the Applications library based on a given ID"
	OutputSuccess    = "Application %d was successfully deleted"
	HelpFlag         = "Displays more information about the delete application command"
	CascadeFlag      = "Deletes all resources created through the 'azion deploy' command"
	MissingFunction  = "Missing Function ID in azion.json file. Skipping deletion\n"
	CascadeSuccess   = "Remote resources deleted successfully\n"
	FlagId           = "Unique identifier of the Application"
	AskInput         = "Enter the ID of the Application you wish to delete:"
	CONFDIRFLAG      = "Relative path to where your custom azion.json and args.json files are stored"
)

// Used by the v4 and v6 command trees.
var (
	Usage = "application"
)
