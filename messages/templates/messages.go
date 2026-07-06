package templates

var (
	// general
	Usage            = "templates"
	FileWritten      = "File successfully written to: %s\n"
	ShortDescription = "Manages your Data Stream templates"
	LongDescription  = "Manages the Data Stream templates used to define the data set sent to your streams"

	// create cmd
	CreateShortDescription = "Creates a new template"
	CreateLongDescription  = "Creates a Data Stream template based on given attributes"
	CreateOutputSuccess    = "Created template with ID %d"
	CreateFlagHelp         = "Displays more information about the create templates command"

	// delete cmd
	DeleteShortDescription = "Deletes a template"
	DeleteLongDescription  = "Removes a template based on its given ID"
	DeleteOutputSuccess    = "Template %d was successfully deleted"
	DeleteHelpFlag         = "Displays more information about the delete templates command"

	// describe cmd
	DescribeShortDescription = "Returns the template data"
	DescribeLongDescription  = "Displays information about the template via a given ID to show its attributes in detail"
	DescribeFlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	DescribeFlagFormat       = "Changes the output format passing the json value to the flag"
	DescribeHelpFlag         = "Displays more information about the describe templates command"

	// list cmd
	ListShortDescription = "Displays your account's templates"
	ListLongDescription  = "Displays all Data Stream templates in the user account"
	ListHelpFlag         = "Displays more information about the list templates command"

	// update cmd
	UpdateShortDescription = "Updates a template"
	UpdateLongDescription  = "Modifies a template based on its ID to update its name, activity status, data set, and other attributes"
	UpdateFlagFile         = "Given path and JSON file to automatically update the template attributes; you can use - for reading from stdin"
	UpdateOutputSuccess    = "Updated template with ID %d"
	UpdateHelpFlag         = "Displays more information about the update templates command"
	UpdateAskTemplateID    = "Enter the ID of the template you wish to update:"
	UpdateAskTemplateFile  = "Enter the path of the json to update the template:"

	// flags
	FlagID      = "Unique identifier of the template"
	FlagIn      = "Given file path to create a template; you can use - for reading from stdin"
	FlagName    = "The template's name"
	FlagDataSet = "Path to a JSON file defining the set of variables sent to your streams; you can use - for reading from stdin"
	FlagActive  = "Whether the template is active or not"

	// ask
	AskTemplateID   = "Enter the template's ID:"
	AskCreateFile   = "Enter the path of the json to create the template:"
	AskInputName    = "Enter the new template's name:"
	AskInputDataSet = "Enter the path to the JSON file with the template's data set:"
)
