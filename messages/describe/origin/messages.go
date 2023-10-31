package origin

import "github.com/MakeNowJust/heredoc"

var (
	Usage             = "origin --application-id <application_id> --origin-id <origin_id>"
	ShortDescription  = "Returns information about a specific origin"
	LongDescription   = "Returns information about a specific origin, based on a given ID, in details"
	FlagApplicationID = "Unique identifier for an edge application. The '--application-id' flag is mandatory"
	FlagOriginKey     = "Unique identifier for an origin. The '--origin-key' flag is mandatory"
	FlagOut           = "Exports the output to the given <file_path/file_name.ext>"
	FlagFormat        = "Changes the output format passing the json value to the flag"
	HelpFlag          = "Displays more information about the describe subcommand"

	AskAppID     = "What is the ID of the Edge Application this origin is linked to?"
	AskOriginKey = "What is the ID of the Origin?"

	Example string = heredoc.Doc(`
	$ azion origin describe --application-id 1673635839 --origin-key 0000000-00000000-00a0a00s0as0-000000
	$ azion origin describe --application-id 1673635839 --origin-key 0000000-00000000-00a0a00s0as0-000000 --format json
	$ azion origin describe --application-id 1673635839 --origin-key 0000000-00000000-00a0a00s0as0-000000 --out "./tmp/test.json" --format json
	`)
)
