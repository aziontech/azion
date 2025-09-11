package functioninstance

import "errors"

var (
	ErrorFailToDeletInstance = errors.New("failed to delete function instance: %s")
)
