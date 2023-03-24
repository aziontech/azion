package rules_engine

var (
	RulesEngineUsage            = "rules_engine"
	RulesEngineShortDescription = "Create rules engines for edges on Azion's platform"
	RulesEngineLongDescription  = "Create rules engines for edges on Azion's platform"
	RulesEngineFlagHelp         = "Displays more information about the rules engines command"
	ApplicationFlagId           = "Unique identifier for an edge application used by this rules engine. The '--application-id' flag is required"
	RulesEngineFlagId           = "Unique identifier for a Rules Engine. The '--rule-id' flag is required"
	RulesEngineListHelpPhase    = "Rules Engine Phase <request|response>. The '--phase' flag is required"

	//list cmd
	RulesEngineListUsage            = "list [flags]"
	RulesEngineListShortDescription = "Displays your rules engines"
	RulesEngineListLongDescription  = "Displays all your rules engine references to your edge locations"
	RulesEngineListHelpFlag         = "Displays more information about the list subcommand"

	//describe cmd
	RulesEngineDescribeUsage            = "describe --application-id <application_id> --phase <phase> --rule-id <rule_id> [flags]"
	RulesEngineDescribeShortDescription = "Returns the rules engine data"
	RulesEngineDescribeLongDescription  = "Displays information about the rules engine via the given IDs and phase to show the rules engine’s attributes in detail"
	RulesEngineDescribeFlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	RulesEngineDescribeFlagFormat       = "Changes the output format passing the json value to the flag"
	RulesEngineDescribeHelpFlag         = "Displays more information about the describe command"
	RulesEngineFileWritten              = "File successfully written to: %s\n"
)
