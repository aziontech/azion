package origin

var (
	Usage             = "origin"
	ShortDescription  = "Deletes an Origin"
	LongDescription   = "Deletes an Origin from the Edge Applications library based on its given ID"
	OutputSuccess     = "Origin %s was successfully deleted\n"
	FlagApplicationID = "Unique identifier for an edge application"
	FlagOriginKey     = "The Origin's key unique identifier"
	HelpFlag          = "Displays more information about the delete origin command"
	AskInputApp       = "What is the id of the edge application linked to this origin?"
	AskInputOri       = "What is the key of the origin you wish to delete?"
)
