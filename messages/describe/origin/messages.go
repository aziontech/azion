package origin

var (
	Usage             = "origin --application-id <application_id> --origin-id <origin_id>"
	ShortDescription  = "Returns information about a specific origin"
	LongDescription   = "Returns information about a specific origin, based on a given ID, in details"
	FlagApplicationID = "Unique identifier for an edge application. The '--application-id' flag is mandatory"
	FlagOriginKey     = "Unique identifier for an origin. The '--origin-key' flag is mandatory"
	FlagOut           = "Exports the output to the given <file_path/file_name.ext>"
	FlagFormat        = "Changes the output format passing the json value to the flag"
	HelpFlag          = "Displays more information about the describe subcommand"
)
