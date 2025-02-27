package sync

const (
	USAGE             = "sync"
	SHORTDESCRIPTION  = "Synchronizes the local azion.json file with remote resources"
	LONGDESCRIPTION   = "Synchronizes your local file containing your existing application resources configuration with remote resources"
	SYNCMESSAGEENV    = "Adding out of sync variable '%s' to your azion account\n"
	SYNCSUCCESS       = "Synched successfully"
	SYNCUPDATEENV     = "Updating remote variable '%s' with local details\n"
	HELPFLAG          = "Displays more information about the sync command"
	CONFDIRFLAG       = "Relative path to where your custom azion.json and args.json files are stored"
	ENVFLAG           = "Relative path to where your custom .env file is stored"
	IACFLAG           = "If sent, sync command will also generate azion.config file"
	IACFORMATFLAG     = "Extension used to generate the azion.config file. Possible options: mjs, cjs, ts, js"
	VARIABLESETSECRET = "Setting secret to true due to the variable key using one of the following words: 'PASSWORD', 'PWD', 'SECRET', 'HASH', 'ENCRYPTED', 'PASSCODE', 'AUTH', 'TOKEN', 'SECRET'\n"
)
