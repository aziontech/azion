package init

import (
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
)

type ReferenceIter interface {
	ForEach(func(*plumbing.Reference) error) error
}

// sortTag return value in format refs/tags/v0.10.0
func sortTag(tags ReferenceIter, majorStr string) (tag string, err error) {
	major, _ := strconv.Atoi(majorStr)

	var previousMinor int = 0
	var previousPatch int = 0

	err = tags.ForEach(func(t *plumbing.Reference) error {
		tagCurrent := t.Name().String() // return this format "refs/tags/v0.10.0"
		if !strings.Contains(tagCurrent, "dev") {
			versionParts := strings.Split(tagCurrent, ".")

			majorCurrent, _ := strconv.Atoi(versionParts[0])
			minorCurrent, _ := strconv.Atoi(versionParts[1])
			patchCurrent, _ := strconv.Atoi(versionParts[2])

			if majorCurrent == major {
				if minorCurrent > previousMinor {
					previousMinor = minorCurrent
					previousPatch = patchCurrent
					tag = tagCurrent
				} else if minorCurrent == previousMinor && patchCurrent > previousPatch {
					previousPatch = patchCurrent
					tag = tagCurrent
				}
			}
		}

		return err
	})

	if err != nil {
		return tag, err
	}

	return tag, err
}
