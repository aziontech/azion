package edgeapplication

var (
	Usage            = "edge-application"
	ShortDescription = "Removes an Edge Application"
	LongDescription  = "Removes an Edge Application from the Edge Applications library based on a given ID"
	OutputSuccess    = "Edge Application %d was successfully deleted"
	HelpFlag         = "Displays more information about the delete edge-application command"
	CascadeFlag      = "Deletes all resources created through the 'azion deploy' command"
	MissingFunction  = "Missing Edge Function ID in azion.json file. Skipping deletion\n"
	CascadeSuccess   = "Remote resources deleted successfully"
	FlagId           = "Unique identifier of the Edge Application"
	AskInput         = "Enter the ID of the Edge Application you wish to delete:"
)
