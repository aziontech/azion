package concat

import "strings"

func String(strs ...string) string {
	var sb strings.Builder

	for i := range strs {
		sb.WriteString(strs[i])
	}

	return sb.String()
}
