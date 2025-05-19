package login

const (
	// general
	Usage            = "login"
	ShortDescription = "Logs in to your Azion account"
	LongDescription  = "Logs in to your Azion account and save a Personal Token locally to authorize CLI commands"
	Success          = "successfully logged in"

	// flags
	FlagUsername = "Your email address"
	FlagPassword = "Your password"
	FlagHelp     = "Displays more information about the login command"

	// Ask
	AskUsername = "Enter your email address:"
	AskPassword = "Enter your password:"

	//browser
	VisitMsg   = "Please visit https://console.azion.com/login?next=cli in case it did not open automatically\n"
	BrowserMsg = "You may now close this page and return to your terminal"
)
