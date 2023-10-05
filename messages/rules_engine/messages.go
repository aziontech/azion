package rules_engine

var (
	RulesEngineUsage            = "rules_engine"
	RulesEngineShortDescription = "Manage your edge applications' rules in Rules Engine on the Azion Edge platform"
	RulesEngineLongDescription  = `Manage your edge applications' rules in Rules Engine on the Azion Edge platform. 
	Rules Engine handles the conditional execution of behaviors through logical operators`
	RulesEngineFlagHelp = "Displays more information about the Rules Engine command"
	ApplicationFlagId   = "Unique identifier for the edge application that implements these rules. The '--application-id' flag is required"
	RulesEngineFlagId   = "Unique identifier for a rule in Rules Engine. The '--rule-id' flag is required"
	RulesEnginePhase    = "Rules Engine Phase <request|response>. The '--phase' flag is required"

	// [ Update ]
	RulesEngineUpdateUsage                 = "update [flags]"
	RulesEngineUpdateShortDescription      = "Updates a rule in Rules Engine"
	RulesEngineUpdateLongDescription       = "Updates a rule in Rules Engine based on given attributes to be used in edge applications"
	RulesEngineUpdateFlagEdgeApplicationId = "Unique identifier for an edge application"
	RulesEngineUpdateFlagIn                = "Path to a JSON file containing the attributes of the rule that will be updated; you can use - for reading from stdin"
	RulesEngineUpdateOutputSuccess         = "Updated rule engine with ID %d\n"
	RulesEngineUpdateHelpFlag              = "Displays more information about the update subcommand"

	// [ create ]
	RulesEngineCreateUsage                 = "create [flags]"
	RulesEngineCreateShortDescription      = "Creates a new rule in Rules Engine"
	RulesEngineCreateLongDescription       = "Creates a new rule in Rules Engine based on given attributes to be used in edge applications"
	RulesEngineCreateFlagEdgeApplicationID = "Unique identifier for an edge application"
	RulesEngineCreateFlagName              = "The rule name"
	RulesEngineCreateFlagPhase             = "The phase is either 'request' or 'response'"
	RulesEngineCreateFlagIn                = "Path to a JSON file containing the attributes of the rule that will be created; you can use - for reading from stdin"
	RulesEngineCreateOutputSuccess         = "Created Rules Engine with ID %d\n"
	RulesEngineCreateHelpFlag              = "Displays more information about the create subcommand"

	// [ template ]
	RulesEngineTemplateUsage            = "template [flags]"
	RulesEngineTemplateShortDescription = "Generates a default template for rules engine"
	RulesEngineTemplateLongDescription  = "Generates a default template to be used with rules engine's create and update commands"
	RulesEngineTemplateFlagOut          = "Exports the template to the given file path <file_path/file_name.ext>"
	RulesEngineTemplateHelpFlag         = "Displays more information about the template command"
)
