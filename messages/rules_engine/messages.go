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

	//describe cmd
	RulesEngineDescribeUsage            = "describe --application-id <application_id> --phase <phase> --rule-id <rule_id> [flags]"
	RulesEngineDescribeShortDescription = "Returns the information related to the rule in Rules Engine"
	RulesEngineDescribeLongDescription  = "Returns the information related to the rule in Rules Engine, informed through the flag '--rule-id' in detail"
	RulesEngineDescribeFlagOut          = "Exports the output of the subcommand 'describe' to the given file path <file_path/file_name.ext>"
	RulesEngineDescribeFlagFormat       = "Changes the output format passing the json value to the flag. Example '--format json'"
	RulesEngineDescribeHelpFlag         = "Displays more information about the describe subcommand"
	RulesEngineFileWritten              = "File successfully written to: %s\n"

	// [ delete ]
	RulesEngineDeleteUsage             = "delete [flags]"
	RulesEngineDeleteShortDescription  = "Deletes a rule in Rules Engine"
	RulesEngineDeleteLongDescription   = "Deletes a rule in Rules Engine based on the given '--rule-id', '--application-id', and '--phase'"
	RulesEngineDeleteOutputSuccess     = "Rule %d was successfully deleted\n"
	RulesEngineDeleteFlagApplicationID = "Unique identifier for an edge application"
	RulesEngineDeleteFlagRuleID        = "The Rule Engine's rule unique identifier"
	RulesEngineDeleteFlagPhase         = "The phase is either 'request' or 'response'"
	RulesEngineDeleteHelpFlag          = "Displays more information about the delete subcommand"

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
