package datastream

var (
	Usage = "data-stream <subcommand>"

	ShortDescription = "Manages your Azion Data Stream"
	LongDescription  = "Manages your Data Stream resources, such as streams, used to send your data to external endpoints and analytics tools"

	// create
	CreateShortDescription = "Creates Data Stream resources"
	CreateLongDescription  = "Creates Data Stream resources, such as streams"
	CreateExample          = "$ azion create data-stream streams --file \"create.json\"\n$ azion create data-stream --help"

	// list
	ListShortDescription = "Displays your Data Stream resources"
	ListLongDescription  = "Displays your Data Stream resources, such as streams"
	ListExample          = "$ azion list data-stream streams\n$ azion list data-stream --help"

	// describe
	DescribeShortDescription = "Returns Data Stream resource data"
	DescribeLongDescription  = "Displays information about a Data Stream resource, such as a stream"
	DescribeExample          = "$ azion describe data-stream streams --stream-id 1234\n$ azion describe data-stream --help"

	// update
	UpdateShortDescription = "Updates Data Stream resources"
	UpdateLongDescription  = "Updates Data Stream resources, such as streams"
	UpdateExample          = "$ azion update data-stream streams --stream-id 1234 --file \"update.json\"\n$ azion update data-stream --help"

	// delete
	DeleteShortDescription = "Deletes Data Stream resources"
	DeleteLongDescription  = "Deletes Data Stream resources, such as streams"
	DeleteExample          = "$ azion delete data-stream streams --stream-id 1234\n$ azion delete data-stream --help"

	FlagHelp = "Displays more information about the data-stream command"
)
