package edgeapplication

var (
	Usage            = "edge-application"
	ShortDescription = "Removes an edge application"
	LongDescription  = "Removes an edge application from the Edge Applications library based on a given ID"
	OutputSuccess    = "Edge application %d was successfully deleted\n"
	HelpFlag         = "Displays more information about the delete edge-application command"
	CascadeFlag      = "Deletes all resources created through the 'azion deploy' command"
	MissingFunction  = "Missing Edge Function ID in azion.json file. Skipping deletion"
	CascadeSuccess   = "Cascade delete carried out successfully"
	FlagId           = "Unique identifier of the edge application"
)
