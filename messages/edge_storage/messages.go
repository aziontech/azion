package edge_storage

const (
	USAGE         = "edge-storage"
	USAGE_BUCKET  = "bucket"
	USAGE_OBJECTS = "object"

	SHORT_DESCRIPTION               = "Manages Edge Storage buckets and objects directly through the command line"
	SHORT_DESCRIPTION_CREATE_BUCKET = "Creates a bucket in Edge Storage"
	SHORT_DESCRIPTION_LIST_BUCKET   = "List the buckets in Edge Storage"
	SHORT_DESCRIPTION_DELETE_BUCKET = "Deleted the bucket in Edge Storage"

	LONG_DESCRIPTION               = "Allows users to perform a wide range of operations, from creating and deleting buckets to adding, removing, and manipulating objects."
	LONG_DESCRIPTION_CREATE_BUCKET = "Allows users to create a bucket in Edge Storage"
	LONG_DESCRIPTION_LIST_BUCKET   = "Allows users to list their buckets in Edge Storage"
	LONG_DESCRIPTION_DELETE_BUCKET = "Allows users to delete their buckets in Edge Storage"

	EXAMPLE_CREATE        = "$ azion create edge-storage\n$ azion create edge-storage --help"
	EXAMPLE_CREATE_BUCKET = "$ azion create edge-storage bucket --name 'zorosola' --edge-access 'read_only'\n$ azion create edge-storage bucket --help"
	EXAMPLE_LIST          = "$ azion list edge-storage\n$ azion list edge-storage --help"
	EXAMPLE_LIST_BUCKET   = "$ azion list edge-storage bucket\n$ azion list edge-storage bucket --page 1 --page-size 3\n$ azion list edge-storage --help"
	EXAMPLE_DELETE        = "$ azion delete edge-storage\n$ azion delete edge-storage --help"
	EXAMPLE_DELETE_BUCKET = "$ azion delete edge-storage bucket --name 'bucket-name'\n$ azion delete edge-storage bucket --help"
	EXAMPLE_UPDATE        = "$ azion update edge-storage\n$ azion update edge-storage --help"
	EXAMPLE_UPDATE_BUCKET = "$ azion update edge-storage bucket --name 'zorosola' --edge-access 'read_only'\n$ azion update edge-storage bucket --help"
	EXAMPLE_UPDATE_OBJECT = "$ azion update edge-storage object --name 'zorosola' --edge-access 'read_only'\n$ azion update edge-storage object --help"

	FLAG_HELP                      = "Displays more information about the edge-storage command"
	FLAG_HELP_CREATE_BUCKET        = "Displays more information about the create edge-storege bucket command"
	FLAG_HELP_LIST_BUCKET          = "Displays more information about the list edge-storege bucket command"
	FLAG_HELP_DELETE_BUCKET        = "Displays more information about the delete edge-storege bucket command"
	FLAG_NAME_BUCKET               = "The name of the Edge Storage bucket"
	FLAG_EDGE_ACCESS_CREATE_BUCKET = "Indicates the type of permission for actions within the bucket. Possible values:	read_only, read_write or restricted"
	FLAG_SOURCE_UPDATE_OBJECT      = "file path field to update the edge storage bucket object"
	FLAG_FILE_JSON_CREATE_BUCKET   = "Path to a JSON file containing the attributes of the bucket that will be created; you can use - for reading from stdin"
	FLAG_NAME_OBJECT               = "The name of the Edge Storage object"

	ASK_NAME_CREATE_BUCKET         = "Enter your bucket's name"
	ASK_EDGE_ACCESSS_CREATE_BUCKET = "Enter your bucket's edge access type"
	ASK_NAME_DELETE_BUCKET         = "Enter the name of the Bucket you wish to delete"
	ASK_NAME_UPDATE_BUCKET         = "Enter the name of the Bucket you wish to update"
	ASK_NAME_UPDATE_OBJECT         = "Enter the name of the Object you wish to update"
	ASK_OBJECT_KEY_CREATE_OBJECT   = "Enter your object name"
	ASK_SOURCE_CREATE_OBJECT       = "Enter your path source"

	OUTPUT_CREATE_BUCKET = "Bucket created successfully"
	OUTPUT_DELETE_BUCKET = "Bucket %s was deleted successfully"
	OUTPUT_UPDATE_BUCKET = "Bucket updated successfully"
)
