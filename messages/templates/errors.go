package templates

import "errors"

var (
	ErrorGetTemplate          = errors.New("Failed to get the template: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorGetTemplates         = errors.New("Failed to list the templates: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorFailToDeleteTemplate = errors.New("Failed to delete the template: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateTemplate       = errors.New("Failed to create template: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateTemplate       = errors.New("Failed to update the template: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorConvertTemplateId    = errors.New("Invalid --template-id flag provided. The value must be an integer. Run the command 'azion describe templates --help' to display more information and try again")
	ErrorActiveFlag           = errors.New("Invalid --active flag provided. The value must be either 'true' or 'false'")
	ErrorDataSetFlag          = errors.New("Failed to read the data set file. Verify if the file name and its path are correct and the file's content has a valid JSON format")
	ErrorParseDataSet         = errors.New("Failed to parse the data set. Verify if the file's content has a valid JSON format")
	ErrorUpdateNoFields       = errors.New("No fields to update were provided. Inform at least one of the flags --name, --data-set, --active, or --file. Nothing was updated")
)
