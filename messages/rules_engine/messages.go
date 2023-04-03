package rules_engine

var (
	RulesEngineUsage            = "rules_engine"
	RulesEngineShortDescription = "Create rules engines for edges on Azion's platform"
	RulesEngineLongDescription  = "Create rules engines for edges on Azion's platform"
	RulesEngineFlagHelp         = "Displays more information about the rules engines command"
	ApplicationFlagId           = "Unique identifier for an edge application used by this rules engine. The '--application-id' flag is required"
	RulesEngineFlagId           = "Unique identifier for a Rules Engine. The '--rule-id' flag is required"
	RulesEnginePhase            = "Rules Engine Phase <request|response>. The '--phase' flag is required"

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

	// [ delete ]
	RulesEngineDeleteUsage             = "delete [flags]"
	RulesEngineDeleteShortDescription  = "Deletes an Rule Engine"
	RulesEngineDeleteLongDescription   = "Deletes an Rule Engine from the Edge Applications library based on its given ID"
	RulesEngineDeleteOutputSuccess     = "Rule Engine %d was successfully deleted\n"
	RulesEngineDeleteFlagApplicationID = "Unique identifier for an edge application"
	RulesEngineDeleteFlagRuleID        = "The Rule Engine's key unique identifier"
	RulesEngineDeleteFlagPhase         = "phase is request input or response output"
	RulesEngineDeleteHelpFlag          = "Displays more information about the delete subcommand"

	// [ Update ]
	RulesEngineUpdateUsage                 = "update [flags]"
	RulesEngineUpdateShortDescription      = "Updates a rule engine"
	RulesEngineUpdateLongDescription       = "Updates a rule engine based on given attributes to be used in edge applications"
	RulesEngineUpdateFlagEdgeApplicationId = "Unique identifier for an edge application"
	RulesEngineUpdateFlagIn                = "Path to a JSON file containing the attributes of the origin that will be updated; you can use - for reading from stdin"
	RulesEngineUpdateOutputSuccess         = "Updated rule engine with ID %d\n"
	RulesEngineUpdateHelpFlag              = "Displays more information about the update subcommand"

	// [ create ]
	RulesEngineCreateUsage                 = "create [flags]"
	RulesEngineCreateShortDescription      = "Creates a new Rules Engine"
	RulesEngineCreateLongDescription       = "Creates an Rules Engine based on given attributes to be used in edge applications"
	RulesEngineCreateFlagEdgeApplicationID = "Unique identifier for an edge application"
	RulesEngineCreateFlagName              = "The rule name"
	RulesEngineCreateFlagPhase             = "The rule phase"
	RulesEngineCreateFlagIn                = "Path to a JSON file containing the attributes of the rule that will be created; you can use - for reading from stdin"
	RulesEngineCreateOutputSuccess         = "Created Rules Engine with ID %d\n"
	RulesEngineCreateHelpFlag              = "Displays more information about the create subcommand"

	// [ template ]
	RulesEngineTemplateUsage            = "template [flags]"
	RulesEngineTemplateShortDescription = "Generates a default template for rules engine"
	RulesEngineTemplateLongDescription  = "Generates a default template to be used with rules engine's create and update commands"
	RulesEngineTemplateFlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	RulesEngineTemplateHelpFlag         = "Displays more information about the template command"
)
