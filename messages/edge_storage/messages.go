package edge_storage

const (
	USAGE                           = "edge-storage"
	USAGE_BUCKET                    = "bucket"
	SHORT_DESCRIPTION               = "Manages Edge Storage buckets and objects directly through the command line"
	SHORT_DESCRIPTION_CREATE_BUCKET = "Creates a bucket in Edge Storage"
	LONG_DESCRIPTION                = "Allows users to perform a wide range of operations, from creating and deleting buckets to adding, removing, and manipulating objects."
	LONG_DESCRIPTION_CREATE_BUCKET  = "Allows users to create a bucket in Edge Storage"
	EXAMPLE                         = "$ azion create edge-storage\n$ azion create edge-storage --help"
	EXAMPLE_CREATE_BUCKET           = "$ azion create edge-storage bucket --name 'zorosola' --edge-access 'read_only'"
	FLAG_HELP                       = "Displays more information about the edge-storage command"
	FLAG_HELP_CREATE_BUCKET         = "Displays more information about the create edge-storege bucket command"
	FLAG_NAME_CREATE_BUCKET         = "The name of the Edge Storage bucket"
	FLAG_EDGE_ACCESS_CREATE_BUCKET  = "Indicates the type of permission for actions within the bucket. Possible values:	read_only, read_write or restricted"
	FLAG_FILE_JSON_CREATE_BUCKET    = "Path to a JSON file containing the attributes of the bucket that will be created; you can use - for reading from stdin"
	SUCCESS_CREATE_BUCKET           = "Bucket created successfully"
	ASK_NAME_CREATE_BUCKET          = "Enter your bucket's name"
	ASK_EDGE_ACCESSS_CREATE_BUCKET  = "Enter your bucket's edge access type (Possible values: read_only, read_write or restricted)"
)
