package domains

var (
	//domains cmd
	RulesEngineUsage            = "rules_engine"
	RulesEngineShortDescription = "Create rules engines for edges on Azion's platform"
	RulesEngineLongDescription  = "Create rules engines for edges on Azion's platform"
	RulesEngineFlagHelp         = "Displays more information about the rules engines command"
	ApplicationFlagId           = "Unique identifier for an edge application used by this rules engine.. The '--application-id' flag is mandatory"

	//list cmd
	RulesEngineListUsage            = "list [flags]"
	RulesEngineListShortDescription = "Displays your domains"
	RulesEngineListLongDescription  = "Displays all your domain references to your edge locations"
	RulesEngineListHelpFlag         = "Displays more information about the list subcommand"
	RulesEngineListHelpPhase        = "Rules Engine Phase <request|response>"
)
