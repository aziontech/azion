package digitalcertificate

var (
	Usage            = "digital-certificate"
	ShortDescription = "Returns the Digital Certificate data"
	LongDescription  = "Displays information about the Digital Certificate via a given ID to show the certificate's attributes in detail"
	FlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	FlagFormat       = "Changes the output format passing the json value to the flag"
	HelpFlag         = "Displays more information about the describe command"

	FlagId                        = "Unique identifier of the Digital Certificate"
	FileWritten                   = "File successfully written to: %s\n"
	AskInputDigitalCertificateID  = "Enter the Digital Certificate's ID:"
)
