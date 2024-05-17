package edge_storage

const (
	ERROR_CREATE_BUCKET   = "Failed to create the bucket: %s. Check your settings and try again. If the error persists, contact Azion support."
	ERROR_LIST_BUCKET     = "Failed to list your buckets: %s. Check your settings and try again. If the error persists, contact Azion support."
	ERROR_DELETE_BUCKET   = "Failed to delete the Bucket: %s. Check your settings and try again. If the error persists, contact Azion support."
	ERROR_UPDATE_BUCKET   = "Failed to update the bucket: %s. Check your settings and try again. If the error persists, contact Azion support."
	ERROR_CREATE_OBJECT   = "Failed to create the object: %s. Check your settings and try again. If the error persists, contact Azion support."
	ERROR_DELETE_OBJECT   = "Failed to delete the Object: %s. Check your settings and try again. If the error persists, contact Azion support."
	ERROR_DESCRIBE_OBJECT = "Failed to describe the object: %s. Check your settings and try again. If the error persists, contact Azion support."
	ERROR_NO_EMPTY_BUCKET = "Unable to delete a non-empty bucket. Additionally, objects deleted within the last 24 hours are also taken into consideration."
)
