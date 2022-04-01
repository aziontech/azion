package utils

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func ConvertIdsToInt(ids ...string) ([]int64, error) {
	converted_ids := make([]int64, len(ids))
	for index, id := range ids {
		converted_id, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		converted_ids[index] = int64(converted_id)
	}

	return converted_ids, nil

}

func CleanDirectory(dir string) error {

	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("%w - %s", ErrorCleaningDirectory, dir)
	}

	return nil
}

func IsDirEmpty(dir string) (bool, error) {
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
