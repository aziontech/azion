package build

import "time"

const VERSION_ID_FORMAT string = "20060102150405"

func createVersionID(dir string) string {
	return time.Now().Format(VERSION_ID_FORMAT)
}
