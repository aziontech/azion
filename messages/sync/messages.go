package sync

const (
	USAGE             = "sync"
	SHORTDESCRIPTION  = "Synchronizes the local azion.json file with remote resources"
	LONGDESCRIPTION   = "Synchronizes your local file containing your existing application resources configuration with remote resources"
	SYNCMESSAGERULE   = "Adding out of sync rule '%s' to your azion.json file\n"
	SYNCMESSAGECACHE  = "Adding out of sync cache '%s' to your azion.json file\n"
	SYNCMESSAGEORIGIN = "Adding out of sync origin '%s' to your azion.json file\n"
	SYNCMESSAGEENV    = "Adding out of sync variable '%s' to your azion account\n"
	HELPFLAG          = "Displays more information about the sync command"
	CONFDIRFLAG       = "Relative path to where your custom azion.json and args.json files are stored"
	ENVFLAG           = "Relative path to where your custom .env file is stored"
)
