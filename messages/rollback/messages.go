package rollback

const (
	USAGE            = "rollback"
	SHORTDESCRIPTION = "Sets static files from a previous deploy"
	LONGDESCRIPTION  = "Sets static files from a previous deploy within the same bucket"
	FLAGHELP         = "Displays more information about the rollback command"
	FLAGORIGINKEY    = "Origin key of the origin used for storage of static files"
	CONFFLAG         = "Relative path to where your custom azion.json and args.json files are stored"
	ASKORIGIN        = "Enter the key of the Origin you wish to update:"
	SUCCESS          = "Static files rolled back successfully"
)
