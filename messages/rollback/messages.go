package rollback

const (
	USAGE            = "rollback"
	SHORTDESCRIPTION = "Sets static files from a previous deploy"
	LONGDESCRIPTION  = "Sets static files from a previous deploy within the same bucket"
	FLAGHELP         = "Displays more information about the rollback command"
	FLAGCONNECTORID  = "Connector ID of the storage connector used for static files"
	CONFFLAG         = "Relative path to where your custom azion.json and args.json files are stored"
	ASKCONNECTOR     = "Enter the ID of the Connector you wish to update:"
	SUCCESS          = "Static files rolled back successfully"
)
