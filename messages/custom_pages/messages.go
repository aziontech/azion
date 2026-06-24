package custompages

var (
	// general
	Usage            = "custom-pages"
	FileWritten      = "File successfully written to: %s\n"
	ShortDescription = "Manages your Azion account's Custom Pages"
	LongDescription  = "Manages custom error pages used to customize the responses returned to your users"

	// create cmd
	CreateShortDescription = "Creates a new Custom Page"
	CreateLongDescription  = "Creates a Custom Page based on given attributes"
	CreateOutputSuccess    = "Created Custom Page with ID %d"
	CreateFlagHelp         = "Displays more information about the create custom-pages command"

	// delete cmd
	DeleteShortDescription = "Deletes a Custom Page"
	DeleteLongDescription  = "Removes a Custom Page based on its given ID"
	DeleteOutputSuccess    = "Custom Page %d was successfully deleted"
	DeleteHelpFlag         = "Displays more information about the delete custom-pages command"

	// describe cmd
	DescribeShortDescription = "Returns the Custom Page data"
	DescribeLongDescription  = "Displays information about the Custom Page via a given ID to show its attributes in detail"
	DescribeFlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	DescribeFlagFormat       = "Changes the output format passing the json value to the flag"
	DescribeHelpFlag         = "Displays more information about the describe custom-pages command"

	// list cmd
	ListShortDescription = "Displays your account's Custom Pages"
	ListLongDescription  = "Displays all Custom Pages in the user account"
	ListHelpFlag         = "Displays more information about the list custom-pages command"

	// update cmd
	UpdateShortDescription  = "Updates a Custom Page"
	UpdateLongDescription   = "Modifies a Custom Page based on its ID to update its name, activity status, pages, and other attributes"
	UpdateFlagFile          = "Given path and JSON file to automatically update the Custom Page attributes; you can use - for reading from stdin"
	UpdateOutputSuccess     = "Updated Custom Page with ID %d"
	UpdateHelpFlag          = "Displays more information about the update custom-pages command"
	UpdateAskCustomPageID   = "Enter the ID of the Custom Page you wish to update:"
	UpdateAskCustomPageFile = "Enter the path of the json to update the Custom Page:"

	// flags
	FlagID = "Unique identifier of the Custom Page"
	FlagIn = "Given file path to create a Custom Page; you can use - for reading from stdin"

	// ask
	AskCustomPageID = "Enter the Custom Page's ID:"
	AskCreateFile   = "Enter the path of the json to create the Custom Page:"
)
