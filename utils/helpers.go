package utils

import (
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
